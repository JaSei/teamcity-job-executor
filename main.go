package main

import (
	"github.com/bgentry/speakeasy"
	"github.com/mkobaly/teamcity"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
	"strconv"
	"time"
)

const VERSION = "0.2.0"

var (
	hostname      = kingpin.Flag("hostname", "teamcity hostname").Short('H').Required().String()
	username      = kingpin.Flag("username", "teamcity username").Short('u').Required().String()
	password      = kingpin.Flag("password", "teamcity password").Short('p').String()
	configId      = kingpin.Arg("configId", "id of build configuration which you can run").String()
	sleepDuration = kingpin.Flag("sleep", "sleep duration of pooling teamcity").Default("5s").Duration()
)

func main() {
	kingpin.Version(VERSION)
	kingpin.Parse()

	if len(*password) == 0 {
		pass, err := speakeasy.Ask("Please enter a password: ")
		if err != nil {
			log.Fatal(err)
		}

		*password = pass
	}

	client := teamcity.New(*hostname, *username, *password)

	b, err := client.QueueBuild(*configId, "master", nil)
	if err != nil {
		log.Fatal("QueueBuild error: %s\n", err)
	}

	log.Println("Build queued (", b.WebURL, ")")

	for {
		b, err = client.GetBuild(strconv.FormatInt(b.ID, 10))
		if err != nil {
			log.Fatal("GetBuild error: %s\n", err)
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
