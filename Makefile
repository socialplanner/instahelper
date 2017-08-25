default:
	debug

build:
	go get -u github.com/jteeuwen/go-bindata/...
	go-bindata -nometadata -pkg="assets" -ignore=\\.DS_Store -prefix "assets" -o app/assets/assets.go assets/...

debug:
	go-bindata -debug -nometadata -pkg="assets" -ignore=\\.DS_Store -prefix "assets" -o app/assets/assets.go assets/...
	go build -o test
	test

deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

test:
	deps
	go test -v ./app/...

version:
	go run app/update/gen_version.go $(v)