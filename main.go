package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/brunetto/goutils/debug"
)

var Debug = false

func main() {
	var (
		airportCmd *exec.Cmd
		err        error
		stdout     io.ReadCloser
	)

	fmt.Println("Scanning...")

	debug.LogDebug(Debug, "Creating command")
	airportCmd = exec.Command("airport", "-s")

	debug.LogDebug(Debug, "Setting STDOUT")
	stdout, err = airportCmd.StdoutPipe()
	if err != nil {
		log.Fatal("Can't connect to mysqldump stdout.")
	}

	debug.LogDebug(Debug, "Setting STDERR")
	airportCmd.Stderr = os.Stderr

	debug.LogDebug(Debug, "Start command")
	err = airportCmd.Start()
	if err != nil {
		log.Fatal("Can't start airport: ", err)
	}

	debug.LogDebug(Debug, "Start intercepting output and wait")
	wil := FilterStdout(stdout)
	err = airportCmd.Wait()
	if err != nil {
		log.Fatal("Wait dumping to finish: ", err)
	}
	debug.LogDebug(Debug, "Finished")

	debug.LogDebug(Debug, "Print result")
	fmt.Println()
	for k, _ := range wil {
		fmt.Println(k)
	}
}

func FilterStdout(stdout io.ReadCloser) WinfoList {
	if Debug {
		defer debug.TimeMe(time.Now())
	}

	var (
		scanner *bufio.Scanner
		line    string
		wil     = WinfoList{}
		reg     = regexp.MustCompile(`(?P<Name>\S+)\s+` +
			`(?P<APMAC>([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2}))\s+` +
			`(?P<Strength>-+\d+)\s+` +
			`(?P<Channel>\d+)\s+` +
			`(?P<HT>\S+)\s+` +
			`(?P<CC>\S+)\s+` +
			`(?P<Security>[A-Za-z0-9]+)` +
			`(?P<SecurityGroup>\((?P<Auth>\S+)\/(?P<Unicast>\S+)\/(?P<Group>\S+)\)){0,1}`)
		match []string
		names map[string]int
	)

	scanner = bufio.NewScanner(stdout)
	for scanner.Scan() {
		line = scanner.Text()
		match = reg.FindStringSubmatch(line)
		if len(match) == 0 {
			continue
		}

		names = map[string]int{}
		for i, n := range reg.SubexpNames() {
			if i != 0 { // first name is empty because index 0 is the whole match
				names[n] = i
			}
		}

		_, exists := wil[match[names["Name"]]]
		if !exists {
			wil[match[names["Name"]]] = []Winfo{}
		}
		wil[match[names["Name"]]] = append(wil[match[names["Name"]]], Winfo{
			Name:     match[names["Name"]],
			APMAC:    match[names["APMAC"]],
			Strength: match[names["Strength"]],
			Channel:  match[names["Channel"]],
			HT:       match[names["HT"]],
			CC:       match[names["CC"]],
			Auth:     match[names["Auth"]],
			Unicast:  match[names["Unicast"]],
			Group:    match[names["Group"]],
		})

	}

	return wil
}

type Winfo struct {
	Name     string // (E)SSID
	APMAC    string // BSSID
	Strength string // RSSI
	Channel  string
	HT       string
	CC       string
	Security string
	Auth     string
	Unicast  string
	Group    string
}

type WinfoList map[string][]Winfo
