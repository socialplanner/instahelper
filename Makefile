.PHONY: assets test

GOARCHES = 386 amd64
BITS = 32 64
GOOSES = linux darwin windows

RUNCMD = ./instahelper
RUNCMDWIN = ./instahelper.exe 

default: debug

build: version deps assets-release
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
	
# Renames 386 > 32, amd64 > 64
	@$(foreach arch,$(GOARCHES),$(foreach os,  $(GOOSES), \
	$(if $(filter 386,$(arch)), mv dist/instahelper-$(v)-$(os)-$(arch) dist/instahelper-$(v)-$(os)-32; ) \
	$(if $(filter amd64,$(arch)), mv dist/instahelper-$(v)-$(os)-$(arch) dist/instahelper-$(v)-$(os)-64; ) \
	))

# Appends .exe to windows
	@$(foreach bit,$(BITS), mv dist/instahelper-$(v)-windows-$(bit)/instahelper dist/instahelper-$(v)-windows-$(bit)/instahelper.exe;)

# Creates a zip archive of each folder
	@$(foreach bit,$(BITS),$(foreach os,  $(GOOSES), zip -rj dist/instahelper-$(v)-$(os)-$(bit).zip dist/instahelper-$(v)-$(os)-$(bit);))

# Remove 32 bit MacOS, all are 64 bit anyways
	rm dist/instahelper-$(v)-darwin-32.zip

# Renames darwin > macos
	mv dist/instahelper-$(v)-darwin-64.zip dist/instahelper-$(v)-macos-64.zip

# Deletes the original folders
	@$(foreach bit,$(BITS),$(foreach os,  $(GOOSES), rm -rf dist/instahelper-$(v)-$(os)-$(bit);))

debug: assets
	go run main.go -debug -noopen

test: deps assets 
	go test -v ./app/...

deps:
	go get github.com/golang/dep/cmd/dep
	dep ensure

version:
ifeq ($(v),)
		$(error Set version with v={{VERSION}})
endif
	go run app/update/gen_version.go $(v)

assets-release:
	go get github.com/jteeuwen/go-bindata/...
	go-bindata -nometadata -pkg="assets" -ignore=\\.DS_Store -prefix "assets" -o app/assets/assets.go assets/...

assets:
	go get github.com/jteeuwen/go-bindata/...
	go-bindata -debug -nometadata -pkg="assets" -ignore=\\.DS_Store -prefix "assets" -o app/assets/assets.go assets/...