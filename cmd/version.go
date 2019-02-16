package cmd

import (
	"fmt"
	"os"

	"github.com/blang/semver"
)

// VersionString to display on --version call
const VersionString = "csvdiff v1.0.0"

const defaultVersion = "1.0-dev"

var version = defaultVersion

// SetVersion will set the version of the cmd package
func SetVersion(_version string) {
	if _version == "" {
		version = defaultVersion
		return
	}

	v, err := semver.Make(_version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "built with wrong version tag\n")
		version = defaultVersion
		return
	}

	if err = v.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "built with wrong version tag\n")
		version = defaultVersion
		return
	}

	version = _version
}

// Version will return the set version of cmd package
func Version() string {
	if version == "" {
		return defaultVersion
	}

	return version
}
