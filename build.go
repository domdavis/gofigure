package gofigure

import "fmt"

// Build identifiers. These are given sensible defaults, but can be overridden
// when an application is compiled using -ldflags. Within a Taskfile this may
// look like:
//
//	build:
//	  cmds:
//	    - task: clean
//	    - |
//	      go build -ldflags "\
//	      -X {{.package}}.BuildTime={{now | date "2006/01/02/15:04:05"}} \
//	      -X {{.package}}.CommitHash={{or .GITHUB_SHA .BITBUCKET_COMMIT "local"}} \
//	      -X {{.package}}.Identifier={{or .GITHUB_REF_NAME .BITBUCKET_TAG .BITBUCKET_BRANCH "local"}}\
//	      " ./...
//	  vars:
//	    package: github.com/domdavis/gofigure
//
// Within a Makefile this may
// look like:
//
//		 PACKAGE = github.com/domdavis/gofigure
//		 BUILD_TIME = $(shell date +"%Y/%m/%d-%H:%M:%S")
//	  ifdef .GITHUB_SHA
//	  HASH = $(GITHUB_SHA)
//	  else ifdef BITBUCKET_COMMIT
//	  HASH = $(BITBUCKET_COMMIT)
//	  else
//		 HASH = $(shell git rev-parse HEAD)
//	  endif
//
//	  ifeq (, $(HASH))
//	  HASH = local
//	  endif
//
//		 ifdef GITHUB_REF_NAME
//	  TAG = $(GITHUB_REF_NAME)
//		 else ifdef BITBUCKET_TAG
//		 TAG = $(BITBUCKET_TAG)
//		 else ifdef BITBUCKET_BRANCH
//		 TAG = $(BITBUCKET_BRANCH)
//		 else
//		 TAG = $(shell git rev-parse --abbrev-ref HEAD)
//		 endif
//
//		 ifeq (, $(TAG))
//		 TAG = local
//		 endif
//
//		 LDFLAGS = -X $(PACKAGE).BuildTime=$(BUILD_TIME) \
//			-X $(PACKAGE).CommitHash=$(HASH) \
//			-X $(PACKAGE).Identifier=$(TAG)
//
//		 build:
//		 	go build -ldflags "$(LDFLAGS)"
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
