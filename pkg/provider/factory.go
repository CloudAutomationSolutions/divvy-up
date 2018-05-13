package provider

import (
	"log"

	"github.com/spf13/viper"
)

func BackendFromFlag(backendFlag string) Backend {
	// use one of these objects and add NewBackend to the interface https://github.com/spf13/viper#extract-sub-tree
	supportedBackends := []string{"amazon"}
	var backend Backend
	var v *viper.Viper

	switch backendFlag {
	case "amazon":
		crashIfMissingSection("amazon")
		v = viper.Sub("amazon")
		backend = newAmazonBackend(v)
	default:
		log.Printf("Please choose one of the supported backends: %v", supportedBackends)
		log.Fatal("Unsupported provider backend: ", backendFlag)
	}

	return backend
}

func newAmazonBackend(v *viper.Viper) Backend {
	v.SetDefault("prefix", "")
	v.SetDefault("endpoint", "")

	crashIfMissingFlags(v, []string{"bucket", "region"})

	return Backend(NewAmazonBackend(
		v.GetString("bucket"),
		v.GetString("region"),
		v.GetString("prefix"),
		v.GetString("endpoint"),
	))
}
