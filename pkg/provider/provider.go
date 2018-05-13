package provider

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"

	"gopkg.in/yaml.v2"
)

type (
	// Backend is a generic interface for cloud provider backends
	Backend interface {
		// This method provisions the divvy-up backend onto the cloud provider.
		// It takes in a set of files and parameters which represent to manifests to be applied on the cloud provider
		Bootstrap(userSpecifiedParametersLocation string)

		// This method has to be able to get a path to a file as an input.
		// It should return an url which can be shared by the user.
		// At the returned url the easily accessible but secure data should be present.
		Distribute(filePath string) string
	}

	BootstrapParameterElement struct {
		Key   string `yaml:"key"`
		Value string `yaml:"value"`
	}

	BootstrapConfig struct {
		TemplateFile       string                      `yaml:"file"`
		BoottrapParameters []BootstrapParameterElement `yaml:"parameters"`
	}
)

// to be used for local files only
func readFile(filePath string) []byte {
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Cannot read %s: %s", filePath, err.Error())
	}
	return contents
}

// to be used for paths that include protocol. This supports file:// and https://
func readFileWithSchema(filePath string) []byte {
	t := &http.Transport{}
	// TODO: Make sure we support windows
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	c := &http.Client{Transport: t}

	resp, err := c.Get(filePath)
	if err != nil {
		log.Fatalf("Cannot fetch file %s: %s", filePath, err.Error())
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Cannot read %s: %s", filePath, err.Error())
	}
	return contents

}

func getUserSpecifiedBootstrapConfig(filePath string) []BootstrapConfig {
	// TODO: User pointers for bootstrapParameters
	contents := readFileWithSchema(filePath)

	var userParameters []BootstrapConfig
	yaml.Unmarshal(contents, &userParameters)

	return userParameters
}

func crashIfMissingFlags(v *viper.Viper, flags []string) {
	missing := []string{}
	for _, flag := range flags {
		if !v.IsSet(flag) {
			missing = append(missing, fmt.Sprintf("%s", flag))
		}
	}
	if len(missing) > 0 {
		log.Fatal("Missing mandatory config field(s): ", strings.Join(missing, ", "))
	}
}

func crashIfMissingSection(section string) {
	if !viper.IsSet(section) {
		log.Fatal("Missing configuration section for backend: ", section)
	}

}
