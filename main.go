package main

import (
	"os/exec"
	"log"
	"io"
	"bufio"
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
}

func FilterStdout(stdout io.ReadCloser) WinfoList {
	var (
		scanner *bufio.Scanner
		line    string
	)

	scanner = bufio.NewScanner(stdout)
	for scanner.Scan() {
		line = scanner.Text()

	}
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
