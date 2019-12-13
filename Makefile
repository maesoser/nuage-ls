all: dependencies nuage-nsg

nuage-nsg:
	mkdir -p bin
	for pkgname in show fetch config util ; do \
		go vet github.com/maesoser/nuage-nsg/pkg/$$pkgname; \
		go fmt github.com/maesoser/nuage-nsg/pkg/$$pkgname; \
	done
	env GOOS=windows GOARCH=amd64 go build -o bin/nuage-nsg_amd64.exe cmd/main.go
	for arch in amd64 arm64 ; do \
		env GOOS=linux GOARCH=$$arch go build -o bin/nuage-nsg_$$arch cmd/main.go; \
	done

dependencies:
	go get github.com/nuagenetworks/vspk-go/vspk

.PHONY: clean

clean:
	rm -fr bin
