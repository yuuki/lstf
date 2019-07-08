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
	GO111MODULE=off go get -v \
	golang.org/x/tools/cmd/cover \
	github.com/mattn/goveralls \
	github.com/motemen/gobump/cmd/gobump \
	github.com/Songmu/ghch/cmd/ghch \
	github.com/Songmu/goxz/cmd/goxz \
	github.com/tcnksm/ghr \
	github.com/Songmu/gocredits/cmd/gocredits

.PHONY: credits
credits:
	GO111MODULE=off go get github.com/go-bindata/go-bindata/...
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

.PHONY: lint
lint:
	go vet ./...
	golint -set_exit_status ./...
