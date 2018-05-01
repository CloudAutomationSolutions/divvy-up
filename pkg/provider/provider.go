package provider

import "strings"

type (
	// Backend is a generic interface for cloud provider backends
	Backend interface {
		Bootstrap() string

		// This method should be able to get a path to a file as an input.
		// It should return an url which can be shared by the user.
		// At the returned url the easily accessible but secure data should be present.
		Distribute(filename string) string
	}
)

func cleanPrefix(prefix string) string {
	return strings.Trim(prefix, "/")
}
