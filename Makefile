.PHONY: docker-builder vic-oskr-server

docker-builder:
	docker build -t armbuilder docker-builder/.

all: vic-oskr-server

go_deps:
	echo `/usr/local/go/bin/go version` && cd $(PWD) && /usr/local/go/bin/go mod download

vic-oskr-server: go_deps
	mkdir -p build
	docker container run  \
	-v "$(PWD)":/go/src/digital-dream-labs/vector-cloud \
	-v $(GOPATH)/pkg/mod:/go/pkg/mod \
	-w /go/src/digital-dream-labs/vector-cloud \
	--user $(UID):$(GID) \
	armbuilder \
	go build  \
	-tags nolibopusfile,vicos \
	--trimpath \
	-ldflags '-w -s -linkmode internal -extldflags "-static" -r /anki/lib' \
	-o build/vic-oskr-server \
	cmd/server.go

	docker container run \
	-v "$(PWD)":/go/src/digital-dream-labs/vector-cloud \
	-v $(GOPATH)/pkg/mod:/go/pkg/mod \
	-w /go/src/digital-dream-labs/vector-cloud \
	--user $(UID):$(GID) \
	armbuilder \
	upx build/vic-oskr-server
