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

var version = "dev"

var (
	hostname      = kingpin.Flag("hostname", "teamcity hostname").Short('H').Required().String()
	username      = kingpin.Flag("username", "teamcity username").Short('u').Required().String()
	password      = kingpin.Flag("password", "teamcity password").Short('p').String()
	branch        = kingpin.Flag("branch", "Branch for VSC root in teamcity job").Short('b').String()
	jobParams     = kingpin.Flag("job_param", "teamcity job parameters in key=value format").Short('j').Strings()
	sleepDuration = kingpin.Flag("sleep", "sleep duration of pooling teamcity").Default("5s").Duration()
	nowait        = kingpin.Flag("nowait", "Does not wait for queued job to finish").Default("false").Bool()
	configID      = kingpin.Arg("configId", "id of build configuration which you can run").String()
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

	b, err := client.QueueBuild(*configID, *branch, properties)
	if err != nil {
		log.Fatalf("QueueBuild error: %s\n", err)
	}

	log.Println("Build queued (", b.WebURL, ")")

	if *nowait {
		os.Exit(0)
	}

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

	log.Printf("Status text:\n%s\n", b.StatusText)
	if b.Status == "SUCCESS" {
		os.Exit(0)
	}

	os.Exit(1)
}
