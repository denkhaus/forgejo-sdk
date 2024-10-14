GO ?= go

WORK_DIR   := $(shell pwd)

FORGEJO_SDK_TEST_URL ?= http://localhost:3000
FORGEJO_SDK_TEST_USERNAME ?= test01
FORGEJO_SDK_TEST_PASSWORD ?= test01

PACKAGE := codeberg.org/mvdkleijn/forgejo-sdk/forgejo

GOFUMPT_PACKAGE ?= mvdan.cc/gofumpt@v0.7.0
GOLANGCI_LINT_PACKAGE ?= github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

FORGEJO_VERSION := 7.0.4
FORGEJO_DL := https://codeberg.org/forgejo/forgejo/releases/download/v$(FORGEJO_VERSION)/forgejo-$(FORGEJO_VERSION)-
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
  FORGEJO_DL := $(FORGEJO_DL)linux-

  UNAME_P := $(shell uname -p)
  ifeq ($(UNAME_P),unknown)
   FORGEJO_DL := $(FORGEJO_DL)amd64
  endif
  ifeq ($(UNAME_P),x86_64)
   FORGEJO_DL := $(FORGEJO_DL)amd64
  endif
  ifneq ($(filter %86,$(UNAME_P)),)
   FORGEJO_DL := $(FORGEJO_DL)386
  endif
  ifneq ($(filter arm%,$(UNAME_P)),)
    FORGEJO_DL := $(FORGEJO_DL)arm-5
  endif
endif
ifeq ($(UNAME_S),Darwin)
  FORGEJO_DL := $(FORGEJO_DL)darwin-10.12-amd64
endif

.PHONY: all
all: clean test build

.PHONY: help
help:
	@echo "Make Routines:"
	@echo " - \"\"              run \"make clean test build\""
	@echo " - build             build sdk"
	@echo " - clean             clean"
	@echo " - fmt               format the code"
	@echo " - lint              run golint"
	@echo " - vet               examines Go source code and reports"
	@echo " - test              run unit tests (need a running forgejo)"
	@echo " - test-instance     start a forgejo instance for test"


.PHONY: clean
clean:
	rm -r -f test
	cd forgejo && $(GO) clean -i ./...

.PHONY: fmt
fmt:
	find . -name "*.go" -type f | xargs gofmt -s -w; \
	$(GO) run $(GOFUMPT_PACKAGE) -extra -w ./forgejo

.PHONY: vet
vet:
	# Default vet
	cd forgejo && $(GO) vet $(PACKAGE)

.PHONY: ci-lint
ci-lint: 
	@cd forgejo/; echo -n "gofumpt ...";\
	diff=$$($(GO) run $(GOFUMPT_PACKAGE) -extra -l .); \
	if [ -n "$$diff" ]; then \
		echo; echo "Not gofumpt-ed"; \
		exit 1; \
	fi; echo " done"; echo -n "golangci-lint ...";\
	$(GO) run $(GOLANGCI_LINT_PACKAGE) run --timeout 5m; \
	if [ $$? -eq 1 ]; then \
		echo; echo "Doesn't pass golangci-lint"; \
		exit 1; \
	fi; echo " done"; \
	cd -; \

.PHONY: test
test:
	@export FORGEJO_SDK_TEST_URL=${FORGEJO_SDK_TEST_URL}; export FORGEJO_SDK_TEST_USERNAME=${FORGEJO_SDK_TEST_USERNAME}; export FORGEJO_SDK_TEST_PASSWORD=${FORGEJO_SDK_TEST_PASSWORD}; \
	if [ -z "$(shell curl --noproxy "*" "${FORGEJO_SDK_TEST_URL}/api/v1/version" 2> /dev/null)" ]; then \echo "No test-instance detected!"; exit 1; else \
	    cd forgejo && $(GO) test -race -cover -coverprofile coverage.out; \
	fi

.PHONY: test-instance
test-instance:
	rm -f -r ${WORK_DIR}/test 2> /dev/null; \
	mkdir -p ${WORK_DIR}/test/conf/ ${WORK_DIR}/test/data/
	wget ${FORGEJO_DL} -O ${WORK_DIR}/test/forgejo-main; \
	chmod +x ${WORK_DIR}/test/forgejo-main; \
	echo "[security]" > ${WORK_DIR}/test/conf/app.ini; \
	echo "INTERNAL_TOKEN = eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYmYiOjE1NTg4MzY4ODB9.LoKQyK5TN_0kMJFVHWUW0uDAyoGjDP6Mkup4ps2VJN4" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "INSTALL_LOCK   = true" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "SECRET_KEY     = 2crAW4UANgvLipDS6U5obRcFosjSJHQANll6MNfX7P0G3se3fKcCwwK3szPyGcbo" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "PASSWORD_COMPLEXITY = off" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "[database]" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "DB_TYPE = sqlite3" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "[repository]" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "ROOT = ${WORK_DIR}/test/data/" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "[server]" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "ROOT_URL = ${FORGEJO_SDK_TEST_URL}" >> ${WORK_DIR}/test/conf/app.ini; \
	${WORK_DIR}/test/forgejo-main migrate -c ${WORK_DIR}/test/conf/app.ini; \
	${WORK_DIR}/test/forgejo-main admin user create --username=${FORGEJO_SDK_TEST_USERNAME} --password=${FORGEJO_SDK_TEST_PASSWORD} --email=test01@forgejo.org --admin=true --must-change-password=false --access-token -c ${WORK_DIR}/test/conf/app.ini; \
	${WORK_DIR}/test/forgejo-main web -c ${WORK_DIR}/test/conf/app.ini

.PHONY: bench
bench:
	cd forgejo && $(GO) test -run=XXXXXX -benchtime=10s -bench=. || exit 1

.PHONY: build
build:
	cd forgejo && $(GO) build

