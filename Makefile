#
# MAKEFILE for stevedore

#
# binary output name
BINARY=stevedore

#
# project name
PROJECT=github.com/gostevedore/stevedore

#
# Values Version and Commit
VERSION=`cat version || echo "unknown"`
COMMIT=`git rev-parse --short HEAD || echo "unknown"`
BUILD_DATE=`date +"%c"`

#
# folder to store artifacts generated
ARTIFACTS_DIR=dist

#
# working dir
WORKING_DIR=`pwd`

#
# Go options
GO_TEST_OPTS=-count=1 -parallel=4 -v

#
# checksum
CHECKSUM=md5sum
CHECKSUM_EXT=md5

#
# Setup the -ldflags option for go build here, interpolate the variable values
#  -s: Omit the symbol table and debug information.
#  -w: Omit the DWARF symbol table
LDFLAGS=-ldflags "-s -w -X '${PROJECT}/internal/release.BuildDate=${BUILD_DATE}' -X ${PROJECT}/internal/release.Version=${VERSION} -X ${PROJECT}/internal/release.Commit=${COMMIT}"

#
# dafault target
.DEFAULT_GOAL: all

# define phony targets
.PHONY: build checksum clean dependencies install notes snapshot tag tar test vet


help:
	@grep -E '^[a-zA-Z1-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

all: vet test snapshot

hi:
	echo $(MAKEFILE_LIST)

#
# BINARY TARGETS
# 

#
# build the binary
#
# "...We’re disabling cgo which gives us a static binary. 
# We’re also setting the OS to Linux (in case someone builds this on a Mac or Windows) 
# and the -a flag means to rebuild all the packages we’re using, 
# which means all the imports will be rebuilt with cgo disabled..."
#
# reference: https://blog.codeship.com/building-minimal-docker-containers-for-go-applications/
#

build: clean ## run a golang build (it is recommend to use 'snapshot' target)
	CGO_ENABLED=0 GOOS=linux go build ${LDFLAGS} -a -o bin/${BINARY} cmd/${BINARY}.go

checksum: build ## generate binary checksum
	${CHECKSUM} bin/${BINARY} > bin/${BINARY}.${CHECKSUM_EXT}

clean: ## clear binaries generated by install or build targets
	if [ -f bin/${BINARY} ] ; then rm -f bin/${BINARY} ; rm -f $$GOPATH/bin/${BINARY} ; fi
	if [ -f bin/${BINARY}.${CHECKSUM_EXT} ] ; then rm -f bin/${BINARY}.${CHECKSUM_EXT} ; fi
	if [ -f ${ARTIFACTS_DIR}/${BINARY}-${VERSION}.tar.gz ] ; then rm -f ${ARTIFACTS_DIR}/${BINARY}-${VERSION}.tar.gz ; fi

dependencies: ## download dependencies
	go mod download

install: clean test ## install compiled dependencies in $GOPATH/pkg and put the binary in $GOPATH/bin
	go install ${LDFLAGS} ./...

notes: ## generate release notes from commits since last tag
	echo "# RELEASE NOTES" > aux && \
	echo "" >> aux & \
	echo "## ${VERSION}" >> aux && \
	echo "" >> aux & \
	git log --format="- %h: %s" `git describe --tags --abbrev=0 @^`..@ >> aux && \
	LINES=`cat RELEASE_NOTES.md | wc -l` && \
	tail -n `expr $$LINES - 1` RELEASE_NOTES.md >> aux && \
	mv aux RELEASE_NOTES.md

snapshot: ## create a goreleaser snapshot
	goreleaser --snapshot --skip-publish --rm-dist --release-notes RELEASE_NOTES.md

tag: ## generate a tag on main branch based on the Version file content
	git checkout main
	git pull origin main
	git tag -a v${VERSION} -m "Version v${VERSION}"
	git push origin v${VERSION}

tar: checksum ## generate and artifact (it is recommend to use 'snapshot' target)
	mkdir -p ${ARTIFACTS_DIR}
	tar cvzf ${ARTIFACTS_DIR}/${BINARY}-${VERSION}.tar.gz bin/${BINARY} bin/${BINARY}.${CHECKSUM_EXT}
	rm -rf bin/${BINARY} bin/${BINARY}.${CHECKSUM_EXT}

test: ## execute all tests
	go test ${GO_TEST_OPTS} ./...

vet: ## execute go vet
	go vet ${LDFLAGS} ./...
