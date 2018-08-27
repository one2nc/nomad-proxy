package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	uuid "github.com/satori/go.uuid"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// Version of this Proxy.
const Version = "0.1.0"

var (
	prefix     = uuid.NewV4().String()
	port       = kingpin.Flag("port", "Port no.").Short('p').Default("9988").String()
	jobPrefix  = kingpin.Flag("job", "Job Prefix").Short('j').Default(prefix).OverrideDefaultFromEnvar("JOB_PREFIX").String()
	serverAddr = kingpin.Flag("server-addr", "Server Addr").
			Short('s').Default("http://127.0.0.1:8080").OverrideDefaultFromEnvar("SERVER_ADDR").String()
)

// Transformer accepts a Request and make in-place changes.
type Transformer interface {
	Transform(*http.Request) error
}

// Transformation contains Transformer attached to a path that proxy should interfere with.
// Example: If /v1/jobs is a path that needs some magic to be performed.
// Path is /v1/jobs & magic would be in tx.
type Transformation struct {
	path string
	tx   Transformer
}

// Sugar so that I don't have to write structs to implement single-func Transformer.
type ruleTransformer func(*http.Request) error

func (f ruleTransformer) Transform(r *http.Request) error {
	return f(r)
}

func newBody(r *http.Request, i interface{}) error {
	nb, err := json.Marshal(i)
	if err != nil {
		return err
	}

	r.ContentLength = int64(len(nb))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(nb))
	return nil
}

var rules = []Transformation{
	Transformation{path: "/v1/search", tx: ruleTransformer(search)},
	Transformation{path: "/v1/jobs", tx: ruleTransformer(jobs)},
	Transformation{path: "/v1/job", tx: ruleTransformer(job)},
	Transformation{path: "/v1/system", tx: ruleTransformer(system)},
}

// Accept a Request. Walk through the rules.
// The first path that matches apply the corresponding transformer.
func modifyRequest(r *http.Request) error {
	for _, t := range rules {
		if !strings.HasPrefix(r.URL.Path, t.path) {
			continue
		}

		if err := t.tx.Transform(r); err != nil {
			return err
		}
		return nil
	}
	return nil
}

// Main Engine function.
func main() {

	// Setup and Parse Kingpin.
	kingpin.Version(Version)
	kingpin.Parse()

	// Parse the Backend URL, ensure it works, panic if it doesnt.
	origin, err := url.Parse(*serverAddr)
	if err != nil {
		panic(err)
	}

	// Create a new SingleHost Proxy
	reverseProxy := httputil.NewSingleHostReverseProxy(origin)
	reverseProxy.ModifyResponse = func(res *http.Response) error {
		return nil
	}

	// Director accepts the incoming request and modifies it, if needed.
	reverseProxy.Director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = origin.Scheme
		req.URL.Host = origin.Host

		if err := modifyRequest(req); err != nil {
			log.Println("Cannot render", req.URL)
			panic(err)
		}
		log.Println("Hitting: ", req.Method, " ", req.URL)
	}

	// Start the Server. Listen to the specified Port.
	if err := http.ListenAndServe(fmt.Sprintf(":%v", *port), reverseProxy); err != nil {
		panic(err)
	}
}
