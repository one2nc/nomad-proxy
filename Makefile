REPO := "tsl8/nomad-proxy"

.PHONY: nomad-proxy

deps:
	dep version || (curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh)
	dep ensure -v

test: deps
	go test -v ./...

build: deps
	env GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o nomad-proxy -a -installsuffix cgo \
		github.com/tsocial/nomad-proxy

build_linux:
	go build -o nomad-proxy github.com/tsocial/nomad-proxy

build_mac: build_deps
	env GOOS=darwin GARCH=amd64 CGO_ENABLED=0 go build -o nomad-proxy -a -installsuffix \
    		cgo github.com/tsocial/nomad-proxy

build_image: build
	docker-compose -f docker-compose.yaml build proxy

upload_image: docker_login
	docker tag $(REPO):latest $(REPO):$(TRAVIS_BRANCH)-latest
	docker tag $(REPO):latest $(REPO):$(TRAVIS_BRANCH)-$(TRAVIS_BUILD_NUMBER)
	docker push $(REPO):latest
	docker push $(REPO):$(TRAVIS_BRANCH)-latest
	docker push $(REPO):$(TRAVIS_BRANCH)-$(TRAVIS_BUILD_NUMBER)

docker_login:
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin
