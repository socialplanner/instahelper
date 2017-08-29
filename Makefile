.PHONY: assets test

GOARCHES = 386 amd64
GOOSES = linux darwin windows

RUNCMD = ./instahelper
RUNCMDWIN = ./instahelper.exe 

default: debug

build: version deps
	go get -u github.com/jteeuwen/go-bindata/...
	go-bindata -nometadata -pkg="assets" -ignore=\\.DS_Store -prefix "assets" -o app/assets/assets.go assets/...
	rm -rf dist
	mkdir -p dist

	# Building each binary
	@$(foreach arch,$(GOARCHES),$(foreach os,  $(GOOSES), GOOS=$(os) GOARCH=$(arch) go build -o instahelper-$(v)-$(os)-$(arch);))

	@$(foreach arch,$(GOARCHES),$(foreach os,  $(GOOSES), mkdir -p dist/instahelper-$(v)-$(os)-$(arch);))

	@$(foreach arch,$(GOARCHES),$(foreach os,  $(GOOSES), mv instahelper-$(v)-$(os)-$(arch) dist/instahelper-$(v)-$(os)-$(arch)/instahelper;))

# Checks OS and adds the respective run command
	@$(foreach arch,$(GOARCHES),$(foreach os,  $(GOOSES), \
	$(if $(filter windows,$(os)), echo $(RUNCMDWIN) > dist/instahelper-$(v)-$(os)-$(arch)/run.bat;) \
	$(if $(filter darwin,$(os)), echo $(RUNCMD) > dist/instahelper-$(v)-$(os)-$(arch)/run.command;) \
	$(if $(filter linux,$(os)), echo $(RUNCMD) > dist/instahelper-$(v)-$(os)-$(arch)/run.sh;) \
	))
	
# Appends .exe to windows
	@$(foreach arch,$(GOARCHES), mv dist/instahelper-$(v)-windows-$(arch)/instahelper dist/instahelper-$(v)-windows-$(arch)/instahelper.exe;)

# Creates a zip archive of each folder
	@$(foreach arch,$(GOARCHES),$(foreach os,  $(GOOSES), zip -rj dist/instahelper-$(v)-$(os)-$(arch).zip dist/instahelper-$(v)-$(os)-$(arch);))

# Deletes the original folders
	@$(foreach arch,$(GOARCHES),$(foreach os,  $(GOOSES), rm -rf dist/instahelper-$(v)-$(os)-$(arch);))

debug: assets
	go run main.go

test: deps assets 
	go test -v ./app/...

deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

version:
ifeq ($(v),)
		$(error Set version with v={{VERSION}})
endif
	go run app/update/gen_version.go $(v)

assets-release:
	go get -u github.com/jteeuwen/go-bindata/...
	go-bindata -nometadata -pkg="assets" -ignore=\\.DS_Store -prefix "assets" -o app/assets/assets.go assets/...

assets:
	go get github.com/jteeuwen/go-bindata/...
	go-bindata -debug -nometadata -pkg="assets" -ignore=\\.DS_Store -prefix "assets" -o app/assets/assets.go assets/...