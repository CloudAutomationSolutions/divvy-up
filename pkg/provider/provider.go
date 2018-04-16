package provider

import "strings"

type (
	// Backend is a generic interface for backends
	Backend interface {
		Bootstrap() string
		Distribute() string
	}
)

func cleanPrefix(prefix string) string {
	return strings.Trim(prefix, "/")
}
