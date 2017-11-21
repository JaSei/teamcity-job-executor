package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bgentry/speakeasy"
	"github.com/mkobaly/teamcity"
	"gopkg.in/alecthomas/kingpin.v2"
)

var version = "0.3.0"

var (
	hostname      = kingpin.Flag("hostname", "teamcity hostname").Short('H').Required().String()
	username      = kingpin.Flag("username", "teamcity username").Short('u').Required().String()
	password      = kingpin.Flag("password", "teamcity password").Short('p').String()
	jobParams     = kingpin.Flag("job_param", "teamcity job parameters in key=value format").Short('j').Strings()
	configID      = kingpin.Arg("configId", "id of build configuration which you can run").String()
	sleepDuration = kingpin.Flag("sleep", "sleep duration of pooling teamcity").Default("5s").Duration()
)

func main() {
	kingpin.Version(version)
	kingpin.Parse()

	if len(*password) == 0 {
		pass, err := speakeasy.Ask("Please enter a password: ")
		if err != nil {
			log.Fatal(err)
		}

		*password = pass
	}

	client := teamcity.New(*hostname, *username, *password)

	properties := make(map[string]string)
	for _, pair := range *jobParams {
		keyValue := strings.Split(pair, "=")

		if len(keyValue) != 2 {
			log.Fatalf("Cannot parse job parameter: %s", pair)
		}

		properties[keyValue[0]] = keyValue[1]
	}

	b, err := client.QueueBuild(*configID, "master", properties)
	if err != nil {
		log.Fatalf("QueueBuild error: %s\n", err)
	}

	log.Println("Build queued (", b.WebURL, ")")

	for {
		b, err = client.GetBuild(strconv.FormatInt(b.ID, 10))
		if err != nil {
			log.Fatalf("GetBuild error: %s\n", err)
		}

		log.Println(b)
		time.Sleep(*sleepDuration)

		if b.ComputedState() == teamcity.Finished {
			break
		}
	}

	log.Println(b.StatusText)
	if b.StatusText == "Success" {
		os.Exit(0)
	}

	os.Exit(1)
}
