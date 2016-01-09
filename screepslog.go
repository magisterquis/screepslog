package main

/*
 * screepslog.go
 * Program to watch screeps console output
 * By J. Stuart McMurray
 * Created 20160108
 * Last Modified 20160108
 */

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

/* This program is a pretty direct ripoff of
 * https://github.com/TooAngel/screeps-cli */

const (
	/* URLs */
	wsURL   = "wss://screeps.com/socket/websocket"  /* Websocket */
	authURL = "https://screeps.com/api/auth/signin" /* Authentication */
	uidURL  = "https://screeps.com/api/auth/me"     /* User ID */
	/* Read/Write timeout */
	wto = 8 * time.Second /* Write timeout */
)

var (
	/* Turn off or on verbose logging */
	enableVerbose = flag.Bool(
		"v",
		false,
		"Enable erbose logging",
	)
	/* Network read timeouts */
	rto = flag.Duration(
		"w",
		time.Hour,
		"Exit if no log entries have been recevied after this amount "+
			"of time",
	)
)

/* Globals, because that's how web devs roll >:( */
var (
	token string
	uid   string
)

func main() {
	var (
		uname = flag.String(
			"u",
			"kd5pbo@gmail.com",
			"Screeps `username`, often an email address",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %v [-v] [-u username]

After reading a newline-terminated password on stdin, prints screeps console
output to stdout.

Options:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* Get password */
	var pass []byte
	var err error
	if terminal.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Printf("Password: ")
		pass, err = terminal.ReadPassword(int(os.Stdin.Fd()))
	} else {
		pass, err = bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
	if nil != err {
		log.Fatalf("Unable to read password: %v", err)
	}

	/* Remove trailing newline characters from the password */
	pass = bytes.TrimRight(pass, "\r\n")

	/* Get Token and UID. */
	token, err = getToken(*uname, string(pass))
	if nil != err {
		log.Fatalf("Unable to get token: %v", err)
	}
	verbose("Got token: %v", token)
	uid, err = getUID(token)
	if nil != err {
		log.Fatalf("Unable to get user ID: %v", err)
	}
	verbose("Got user ID: %v", uid)

	/* Handle websocket */
	if err := handleWebsocket(*uname, string(pass)); nil != err {
		log.Fatalf("Error: %v", err)
	}

}

/* verbose prints, log.Printf-style, if -v was given */
func verbose(f string, a ...interface{}) {
	if *enableVerbose {
		log.Printf(f, a...)
	}
}
