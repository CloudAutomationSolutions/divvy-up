package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cloudautomationsolutions/divvy-up/pkg/provider"
	"github.com/urfave/cli"
)

var (
	echo = fmt.Printf
	exit = os.Exit

	// Version is the semantic version (added at compile time)
	Version string

	// Revision is the git commit id (added at compile time)
	Revision string
)

func main() {
	app := cli.NewApp()
	app.Name = "divvy-up"
	app.Version = fmt.Sprintf("%s (build %s)", Version, Revision)
	app.Usage = "Secure file sharing system based on your cloud infrastructure."
	app.Flags = globalFlags
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Robert Tingirică",
			Email: "robert.tingirica@gmail.com",
		},
	}
	app.Commands = cliCommands
	app.Run(os.Args)
}

func handler(action string, c *cli.Context) {
	backend := backendFromContext(c)

	switch action {
	case "distribute":
		crashIfContextMissingFlags(c, []string{"file", "expiration"})
		url := backend.Distribute(c.String("file"))
		echo("Access your file at: %s", url)
	case "bootstrap":
		crashIfContextMissingFlags(c, []string{"method"})
		backend.Bootstrap()
	default:
		log.Fatal("Unsupported action: ", action)
	}
}

func backendFromContext(c *cli.Context) provider.Backend {
	crashIfContextMissingFlags(c, []string{"backend"})
	var backend provider.Backend

	backendFlag := strings.ToLower(c.GlobalString("backend"))
	switch backendFlag {
	case "amazon":
		backend = amazonBackendFromContext(c)
	default:
		log.Fatal("Unsupported provider backend: ", backendFlag)
	}

	return backend
}

func amazonBackendFromContext(c *cli.Context) provider.Backend {
	crashIfContextMissingFlags(c, []string{"amazon-bucket", "amazon-region"})
	return provider.Backend(provider.NewAmazonS3Backend(
		c.GlobalString("amazon-bucket"),
		c.GlobalString("amazon-region"),
		c.GlobalString("amazon-prefix"),
		c.GlobalString("amazon-endpoint"),
	))
}

func crashIfContextMissingFlags(c *cli.Context, flags []string) {
	missing := []string{}
	for _, flag := range flags {
		if c.String(flag) == "" && c.GlobalString(flag) == "" {
			missing = append(missing, fmt.Sprintf("--%s", flag))
		}
	}
	if len(missing) > 0 {
		log.Fatal("Missing mandatory flags(s): ", strings.Join(missing, ", "))
	}
}

var cliCommands = []cli.Command{
	cli.Command{
		Name:        "bootstrap",
		Aliases:     []string{"b"},
		Category:    "Initial setup",
		Description: "Bootstrap the provider account with the required resources",
		Flags:       bootstrapFlags,
		Action: func(c *cli.Context) error {
			handler("bootstrap", c)
			return nil
		},
	},
	cli.Command{
		Name:        "distribute",
		Aliases:     []string{"d"},
		Category:    "Sharing",
		Description: "Share a file in a secure way using your cloud provider",
		Flags:       distributeFlags,
		Action: func(c *cli.Context) error {
			handler("distribute", c)
			return nil
		},
	},
}

var globalFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "backend, b",
		Usage:  "the backend to be used for hosting and delivering the secret file",
		Value:  "amazon",
		EnvVar: "DIVVY_BACKEND",
	},
	cli.StringFlag{
		Name:   "amazon-bucket",
		Usage:  "the S3 bucket to be used for hosting and delivering the secret file",
		EnvVar: "DIVVY_AWS_BUCKET",
	},
	cli.StringFlag{
		Name:   "amazon-region",
		Usage:  "the region where the bucket will be present",
		EnvVar: "DIVVY_AWS_REGION",
	},
	cli.StringFlag{
		Name:   "amazon-prefix",
		Usage:  "the prefix the files will be added with",
		EnvVar: "DIVVY_AWS_BUCKET",
	},
	cli.StringFlag{
		Name:   "amazon-endpoint",
		Usage:  "the AWS api endpoint to be used",
		EnvVar: "DIVVY_AWS_REGION",
	},
}

var bootstrapFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "method",
		Value: "cloudformation",
		Usage: "The file which holds your secrets",
	},
}

var distributeFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "file, f",
		Usage: "The file which holds your secrets",
	},
	cli.IntFlag{
		Name:  "expiration, e",
		Value: 1800,
		Usage: "The time it takes for the file to expire",
	},
}
