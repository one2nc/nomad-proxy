CLIENT_REPO := "tsl8/nomad-client-proxy"
SERVER_REPO := "tsl8/nomad-server-proxy"
SHELL := /bin/bash

.PHONY: nomad-proxy

deps:
	dep version || (curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh)
	dep ensure -v

test: deps
	env GOCACHE=/tmp go test -v ./...

build_client_proxy: deps
	env GOOS=linux GARCH=amd64 CGO_ENABLED=0 GOCACHE=/tmp go build -o client-proxy/nomad-client-proxy -a -installsuffix cgo \
		github.com/tsocial/nomad-proxy/client-proxy

build_server_proxy: deps
	env GOOS=linux GARCH=amd64 CGO_ENABLED=0 GOCACHE=/tmp go build -o server-proxy/nomad-server-proxy -a -installsuffix cgo \
		github.com/tsocial/nomad-proxy/server-proxy

build_linux: build_client_proxy build_server_proxy

build_mac: build_deps
	env GOOS=darwin GARCH=amd64 CGO_ENABLED=0 GOCACHE=/tmp go build -o nomad-client-proxy -a -installsuffix \
    		cgo github.com/tsocial/nomad-proxy/client-proxy
	env GOOS=darwin GARCH=amd64 CGO_ENABLED=0 GOCACHE=/tmp go build -o nomad-server-proxy -a -installsuffix \
                cgo github.com/tsocial/nomad-proxy/server-proxy

build_images: build_linux
	docker-compose -f docker-compose.yaml build client-proxy server-proxy

docker_login:
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin

upload_image: docker_login
	docker tag $(CLIENT_REPO):latest $(CLIENT_REPO):$(TRAVIS_BRANCH)-latest
	docker tag $(CLIENT_REPO):latest $(CLIENT_REPO):$(TRAVIS_BRANCH)-$(TRAVIS_BUILD_NUMBER)
	docker push $(CLIENT_REPO):latest
	docker push $(CLIENT_REPO):$(TRAVIS_BRANCH)-$(TRAVIS_BUILD_NUMBER)

	docker tag $(SERVER_REPO):latest $(SERVER_REPO):$(TRAVIS_BRANCH)-latest
	docker tag $(SERVER_REPO):latest $(SERVER_REPO):$(TRAVIS_BRANCH)-$(TRAVIS_BUILD_NUMBER)
	docker push $(SERVER_REPO):latest
	docker push $(SERVER_REPO):$(TRAVIS_BRANCH)-$(TRAVIS_BUILD_NUMBER)

docker_registry_login:
	echo "$(DOCKER_REGISTRY_PASSWORD)" | docker login -u "$(DOCKER_REGISTRY_USERNAME)" ${DOMAIN} --password-stdin

registry_upload_image: docker_registry_login
	docker tag $(CLIENT_REPO):latest $(REGISTRY_CLIENT_REPO):latest
	docker tag $(CLIENT_REPO):latest $(REGISTRY_CLIENT_REPO):$(TRAVIS_BRANCH)-$(TRAVIS_BUILD_NUMBER)
	docker push $(REGISTRY_CLIENT_REPO):latest
	docker push $(REGISTRY_CLIENT_REPO):$(TRAVIS_BRANCH)-$(TRAVIS_BUILD_NUMBER)

	docker tag $(SERVER_REPO):latest $(REGISTRY_SERVER_REPO):latest
	docker tag $(SERVER_REPO):latest $(REGISTRY_SERVER_REPO):$(TRAVIS_BRANCH)-$(TRAVIS_BUILD_NUMBER)
	docker push $(REGISTRY_SERVER_REPO):latest
	docker push $(REGISTRY_SERVER_REPO):$(TRAVIS_BRANCH)-$(TRAVIS_BUILD_NUMBER)

copy_certs:
	cp -r client-proxy/testdata/* /tmp/

run_nomad: copy_certs
	bash client-proxy/scripts/nomad.sh

validate_nomad_server_tls: run_nomad
	sleep 5
	curl -k https://localhost:4646; [[ $$? -eq "35" ]] && /bin/true

run_client_proxy: build_client_proxy
	./client-proxy/nomad-client-proxy --root-ca-file=/tmp/cert-chain.pem --cert-file=/tmp/client.pem --key-file=/tmp/client-key.pem --server-addr=https://localhost:4646 --dc=dc1 2>&1 &

validate_client_proxy: validate_nomad_server_tls run_client_proxy
	sleep 5
	curl -k http://localhost:9988; [[ $$? -eq "0" ]] && /bin/true
