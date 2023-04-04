package gofigure

import "fmt"

// Build identifiers. These are given sensible defaults, but can be overridden
// when an application is compiled using -ldflags. Within a Makefile this may
// look like:
//
//	 PACKAGE = bitbucket.org/idomdavis/gofigure
//	 BUILD_TIME = $(shell date +"%Y/%m/%d-%H:%M:%S")
//	 HASH = $(shell git rev-parse HEAD)
//
//	 ifdef BITBUCKET_TAG
//	 TAG = $(BITBUCKET_TAG)
//	 else ifdef BITBUCKET_BRANCH
//	 TAG = $(BITBUCKET_BRANCH)
//	 else
//	 TAG = $(shell git rev-parse --abbrev-ref HEAD)
//	 endif
//
//	 ifeq (, $(TAG))
//	 TAG = dev
//	 endif
//
//	 LDFLAGS = -X $(PACKAGE).BuildTime=$(BUILD_TIME) \
//		-X $(PACKAGE).CommitHash=$(HASH) \
//		-X $(PACKAGE).Identifier=$(TAG)
//
//	 build:
//	 	go build -ldflags "$(LDFLAGS)"
//
// The identifiers are used for reference only and have no impact on the
// operation of gofigure.
//
//nolint:gochecknoglobals // Need to be globals to be set during build.
var (
	Identifier = Dev
	BuildTime  = Unset
	CommitHash = Unset
)

// Unset is used when version data is unset and cannot be inferred.
const Unset = "<unset>"

// Dev identifier, used as a default when the Identifier is unset and cannot be
// inferred.
const Dev = "dev"

// Build returns a build string for the application. This relies on the correct
// ldflags being set.
func Build() string {
	return fmt.Sprintf("%s [%s] (built: %s)", Identifier, CommitHash, BuildTime)
}
