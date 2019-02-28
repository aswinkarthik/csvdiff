package cmd

const defaultVersion = "1.0-dev"

var version = defaultVersion

// SetVersion will set the version of the cmd package
func SetVersion(_version string) {
	if _version == "" {
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
