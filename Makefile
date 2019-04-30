export GO111MODULE=on

PKG = github.com/yuuki/lstf
COMMIT = $$(git describe --tags --always)
DATE = $$(date --utc '+%Y-%m-%d_%H:%M:%S')
BUILD_LDFLAGS = -X $(PKG).commit=$(COMMIT) -X $(PKG).date=$(DATE)
RELEASE_BUILD_LDFLAGS = -s -w $(BUILD_LDFLAGS)

.PHONY: build
build:
	go build -ldflags="$(BUILD_LDFLAGS)"

.PHONY: test
test:
	go test -v ./...

.PHONY: cover
cover: devel-deps
	goveralls -service=travis-ci

.PHONY: devel-deps
devel-deps:
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls
	go get github.com/motemen/gobump
	go get github.com/Songmu/ghch
	go get github.com/Songmu/goxz
	go get github.com/tcnksm/ghr
	go get github.com/Songmu/gocredits/cmd/gocredits

.PHONY: credits
credits:
	@go get github.com/go-bindata/go-bindata/...
	gocredits -w .
	go generate -x .

.PHONY: crossbuild
crossbuild: devel-deps credits
	$(eval ver = $(shell gobump show -r))
	goxz -pv=v$(ver) -os=linux -arch=386,amd64 -build-ldflags="$(RELEASE_BUILD_LDFLAGS)" \
	  -d=./dist/v$(ver)

.PHONY: release
release: devel-deps
	_tools/release
	_tools/upload_artifacts
