package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"go_starter/src/pkg/dummy"
	"go_starter/src/pkg/env"
	"go_starter/version"
	"os"
	"runtime"
)

type Configurations struct {
	Dummy  dummy.Config
	Client dummy.Client
}

var config *Configurations

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
}

func initConfig() error {
	// init default config
	config = &Configurations{
		Dummy: dummy.Config{
			Username:  env.GetStringOrDie("GO_STARTER_DUMMY_USERNAME"),
			Password:  env.GetStringOrDie("GO_STARTER_DUMMY_PASSWORD"),
			Debug:     env.GetBoolOrDie("GO_STARTER_DUMMY_DEBUG"),
			UserAgent: env.GetStringOrDie("GO_STARTER_DUMMY_USER_AGENT"),
			BaseUrl:   env.GetStringOrDie("GO_STARTER_DUMMY_BASE_URL"),
			UrlPath:   env.GetStringOrDie("GO_STARTER_DUMMY_URL_PATH"),
		},
	}
	return nil
}

func run(dryRun bool) error {
	if dryRun {
		logrus.Infof("Dry run initiated.")
	}
	client := dummy.NewClient(config.Dummy)
	logrus.Infof("Client BaseUrl/Path is: %s/%s", config.Dummy.BaseUrl, config.Dummy.UrlPath)

	status, err := client.PostExample()
	if err != nil {
		return err
	}
	logrus.Infof("Status is %s", status)

	return nil
}

func main() {
	vrs := flag.Bool("version", false, "Logs current version")
	dryRun := flag.Bool("dry-run", false, "Fakes delivery to file store")
	flag.Parse()
	if *vrs {
		fmt.Printf("go_starter Version: %s\nGo Version: %s\nGo OS/ARCH: %s %s\n",
			version.Version,
			runtime.Version(),
			runtime.GOOS,
			runtime.GOARCH)
		os.Exit(0)
	}

	err := initConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	err = run(*dryRun)
	if err != nil {
		logrus.Fatal(err)
	}
}
