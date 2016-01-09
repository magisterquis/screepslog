package main

/*
 * websocket.go
 * Handles websocket comms
 * By J. Stuart McMurray
 * Created 20160108
 * Last Modified 20160108
 */

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

/* Log unmarshal buffer */
var logMsgBuf = struct {
	Messages struct {
		Log []string
	}
}{}

/* handleWebsocket connects to the service, auths with username u and
password p, and logs logs.  It only returns on error  */
func handleWebsocket(u, p string) error {
	/* Dial the websocket server */
	ws, err := websocket.Dial(wsURL, "", "http://kd5pbo/screepslog")
	if nil != err {
		return err
	}

	/* Authenticate */
	if err := ws.SetWriteDeadline(time.Now().Add(wto)); nil != err {
		return err
	}
	if _, err := ws.Write([]byte("auth " + token)); nil != err {
		return err
	}

	/* Loop over read messages */
	b := make([]byte, 10240)
	for {
		/* Set read deadline */
		if err := ws.SetReadDeadline(time.Now().Add(
			*rto,
		)); nil != err {
			return err
		}
		/* Read a message */
		b = b[:cap(b)]
		n, err := ws.Read(b)
		b = b[:n]
		if nil != err {
			return err
		}
		/* Print it */
		if err := handleMessage(b, ws); nil != err {
			return err
		}
	}

	return nil
}

/* handleMessage handles a message m from ws */
func handleMessage(m []byte, ws *websocket.Conn) error {
	switch {
	case bytes.HasPrefix(m, []byte("time")):
		/* Ping */
		return nil
	case bytes.HasPrefix(m, []byte("auth ok")):
		return handleAuthOk(ws)
	default:
		if err := handleLogMessage(m); nil != err {
			return err
		}
	}

	return nil
}

/* handleAuthOk requests a subscription to console log messages after an "auth
ok" has been received */
func handleAuthOk(ws *websocket.Conn) error {
	/* Set write timeout */
	if err := ws.SetWriteDeadline(time.Now().Add(wto)); nil != err {
		return err
	}
	/* Ask to subscribe */
	if _, err := ws.Write([]byte(
		"subscribe user:" + uid + "/console",
	)); nil != err {
		return err
	}
	return nil
}

/* handleLogMessage handles a message which (hopefully) contains logs */
func handleLogMessage(m []byte) error {
	/* Make a decoder */
	dec := json.NewDecoder(bytes.NewBuffer(m))
	/* Read to the open [ */
	t, err := dec.Token()
	if nil != err {
		return err
	}
	if d, ok := t.(json.Delim); ok {
		if "[" != d.String() {
			return fmt.Errorf(
				"bad delimiter %q in message %q",
				d,
				m,
			)
		}
	} else {
		return fmt.Errorf("unexpected delimiter (%T): %v", t, t)
	}
	/* Read the first bit, which is the uid */
	var s string
	if err := dec.Decode(&s); nil != err {
		return err
	}
	/* If there's no more, something's wrong */
	if !dec.More() {
		return fmt.Errorf("unexpected end of message")
	}

	/* Try to unmarshal the rest */
	if err := dec.Decode(&logMsgBuf); nil != err {
		return err
	}

	/* Print out the log messages */
	for _, v := range logMsgBuf.Messages.Log {
		log.Printf("%v", v)
	}

	return nil
}
