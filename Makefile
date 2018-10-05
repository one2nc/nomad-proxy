CLIENT_REPO := "tsl8/nomad-client-proxy"
SERVER_REPO := "tsl8/nomad-server-proxy"

.PHONY: nomad-proxy

deps:
	dep version || (curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh)
	dep ensure -v

test: deps
	go test -v ./...

build_client_proxy: deps
	env GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o client-proxy/nomad-client-proxy -a -installsuffix cgo \
		github.com/tsocial/nomad-proxy/client-proxy

build_server_proxy: deps
	env GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o server-proxy/nomad-server-proxy -a -installsuffix cgo \
		github.com/tsocial/nomad-proxy/server-proxy

build_linux: build_client_proxy build_server_proxy

build_mac: build_deps
	env GOOS=darwin GARCH=amd64 CGO_ENABLED=0 go build -o nomad-client-proxy -a -installsuffix \
    		cgo github.com/tsocial/nomad-proxy/client-proxy
	env GOOS=darwin GARCH=amd64 CGO_ENABLED=0 go build -o nomad-server-proxy -a -installsuffix \
                cgo github.com/tsocial/nomad-proxy/server-proxy

build_images: build_linux
	docker-compose -f docker-compose.yaml build client-proxy server-proxy

upload_image: docker_login
	docker tag $(CLIENT_REPO):latest $(CLIENT_REPO):$(TRAVIS_BRANCH)-latest
	docker tag $(CLIENT_REPO):latest $(CLIENT_REPO):$(TRAVIS_BRANCH)-$(TRAVIS_BUILD_NUMBER)
	docker push $(CLIENT_REPO):latest
	docker push $(CLIENT_REPO):$(TRAVIS_BRANCH)-latest
	docker push $(CLIENT_REPO):$(TRAVIS_BRANCH)-$(TRAVIS_BUILD_NUMBER)
	docker tag $(SERVER_REPO):latest $(SERVER_REPO):$(TRAVIS_BRANCH)-latest
	docker tag $(SERVER_REPO):latest $(SERVER_REPO):$(TRAVIS_BRANCH)-$(TRAVIS_BUILD_NUMBER)
	docker push $(SERVER_REPO):latest
	docker push $(SERVER_REPO):$(TRAVIS_BRANCH)-latest
	docker push $(SERVER_REPO):$(TRAVIS_BRANCH)-$(TRAVIS_BUILD_NUMBER)

docker_login:
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin
