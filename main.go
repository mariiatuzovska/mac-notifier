package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sparrc/go-ping"
)

var notifier chan error

type Configuration struct {
	Credentials struct {
		Sender   string
		Passwors string
	}
	Host         string
	SMTPAddress  string
	Recipient    []string
	TestPingHost string
}

func (c *Configuration) valid() error {
	if len(c.Recipient) < 1 {
		return errors.New("Unknown recipient")
	}
	return nil
}

func main() {

	notifier = make(chan error)

	go func(notifier chan error) {

		config := new(Configuration)
		file, err := ioutil.ReadFile("/opt/mac-notifier/config.json")
		notifier <- err
		notifier <- json.Unmarshal(file, &config)
		notifier <- config.valid()

		for {
			_, err := ping.NewPinger(config.TestPingHost)
			if err == nil {
				break
			}
			time.Sleep(time.Duration(15) * time.Second)
		}

		msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: LOGIN INFO: %s\n\n", config.Credentials.Sender, config.Recipient[0], time.Now().UTC().String())

		airportCmd := exec.Command("/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport", "-s")
		airportCmdOutput, err := airportCmd.Output()
		notifier <- err

		lines := strings.Split(string(airportCmdOutput), "\n")
		for _, line := range lines {
			columns := strings.Fields(line)
			if len(columns) > 0 {
				match, _ := regexp.MatchString("^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$", columns[1])
				if match == true {
					power, _ := strconv.Atoi(columns[2])
					msg += fmt.Sprintf("SSID: %s, BSSID: %s, Power: %d. \n", columns[0], columns[1], power)
				}
			}
		}

		notifier <- smtp.SendMail(config.SMTPAddress,
			smtp.PlainAuth("", config.Credentials.Sender, config.Credentials.Passwors, config.Host),
			config.Credentials.Sender, config.Recipient, []byte(msg))

	}(notifier)

	for {
		err := <-notifier
		if err != nil {
			log.Fatal(err)
		}
	}

}
