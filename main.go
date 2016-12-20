package main

import (
	"os/exec"
	"log"
	"io"
	"bufio"
	"regexp"
)

func main () {
	var (
		airportCmd *exec.Cmd
		err error
		stdout        io.ReadCloser
	)



	airportCmd = exec.Command("airport", "-s")
	stdout, err = airportCmd.StdoutPipe();
	if err != nil {
		log.Fatal("Can't connect to mysqldump stdout.")
	}
	FilterStdout(stdout)
	err = airportCmd.Wait()
	if err != nil {
		log.Fatal("Wait dumping to finish: ", err)
	}
}

func FilterStdout(stdout io.ReadCloser) WinfoList {
	var (
		scanner *bufio.Scanner
		line    string
		wil = WinfoList{}
		reg = regexp.MustCompile(	`(?P<Name>\S+)\s+` +
						`(?P<APMAC>([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2}))\s+` +
						`(?P<Strength>-+\d+)\s+` +
						`(?P<Channel>\d+)\s+` +
						`(?P<HT>\S+)\s+` +
						`(?P<CC>\S+)\s+` +
						`(?P<Security>\S+)` +
						`(?P<SecurityGroup>\((?P<Auth>\S+)\/(?P<Unicast>\S+)\/(?P<Group>\S+)\)){0,1}`)
		match []string
		names map[string]int
	)

	scanner = bufio.NewScanner(stdout)
	for scanner.Scan() {
		line = scanner.Text()
		match = reg.FindStringSubmatch(line)
		names = map[string]int{}
		for i, n := range reg.SubexpNames() {
			if i != 0 { // first name is empty because index 0 is the whole match
				names[n] = i
			}
		}
		tmp := Winfo{
			Name: match[names["Name"]],
			APMAC: match[names["APMAC"]],
			Strength: match[names["Strength"]],
			Channel: match[names["Channel"]],
			HT: match[names["HT"]],
			CC: match[names["CC"]],
		}

	}

	return wil
}

type Winfo struct {
	Name string  // (E)SSID
	APMAC string // BSSID
	Strength string // RSSI
	Channel string
	HT string
	CC string
	Security string
	Auth string
	Unicast string
	Group string
}

type WinfoList map[string][]Winfo
