# Go related settings
GO := go
GODOC := godoc
GOMETALINTER := gometalinter
GO_VERBOSE := -v

BUILD_PATH := $(subst $(GOPATH)/src/,,$(CURDIR))
LIB_PATH := $(BUILD_PATH)/...

ifneq "$(VERBOSE)" "1"
GO_VERBOSE=
.SILENT:
endif

iterate = \
	for i in $(shell find . -maxdepth 1 -mindepth 1 -not -path '*/\.*' -not -path './ci' -type d -printf '%f\n');do \
		cd $$i && eval $(1) && cd - ;\
	done

# Go build targets
all: get
	$(call iterate, "$(GO) build $(GO_VERBOSE)")

utest: get
	$(GO) test \
		-v \
		-parallel 8 \
		-timeout 30s \
		-covermode atomic \
		$(LIB_PATH) \
		${GOTEST_COLORING}

get:
	$(GO) get $(GO_VERBOSE) $(LIB_PATH)

clean:
	$(GO) clean
	git clean -fdx

# Run source and package documentation server and open it in browser
doc: all
	-pkill $(GODOC)
	$(GODOC) -http=:6060 -links=true -index -play &
	xdg-open http://localhost:6060/pkg/$(BUILD_PATH)

lint: all
	$(GOMETALINTER) \
		--enable-all \
		--line-length=120 \
		--deadline=60s \
		--cyclo-over=25 \
		--exclude=_test.go \
		./...
