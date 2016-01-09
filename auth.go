package main

/*
 * auth.go
 * Handle authentication
 * By J. Stuart McMurray
 * Created 20160108
 * Last Modified 20160108
 */

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

/* getToken gets the authentication token for the given username and
password */
func getToken(u, p string) (string, error) {
	/* Send the request */
	res, err := http.PostForm(authURL, url.Values(map[string][]string{
		"email":    {u},
		"password": {p},
	}))
	if nil != err {
		return "", err
	}
	defer res.Body.Close()
	/* Read body into buffer */
	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, res.Body)
	if nil != err {
		return "", err
	}
	/* If we're not authorized, :( */
	if "Unauthorized" == buf.String() {
		return "", fmt.Errorf("unauthorized")
	}
	/* Struct into which to decode token */
	rs := struct {
		Ok    int
		Token string
	}{}
	/* Unmarshal the JSON */
	if err := json.Unmarshal(buf.Bytes(), &rs); nil != err {
		return "", err
	}

	return rs.Token, nil
}

/* getUID gets the user's UID for the websocket connection given the token t */
func getUID(t string) (string, error) {
	/* Roll the request */
	req, err := http.NewRequest("GET", uidURL, nil)
	if nil != err {
		return "", err
	}
	req.Header.Add("X-Token", t)
	req.Header.Add("X-Username", t)
	/* Send the request */
	res, err := http.DefaultClient.Do(req)
	if nil != err {
		return "", err
	}
	defer res.Body.Close()
	/* Buffer the output */
	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, res.Body)
	/* Struct into which to unmarshal the uid */
	rs := struct {
		Error string
		Id    string `json:"_id"`
	}{}
	/* Unmarshal the json */
	if err := json.Unmarshal(buf.Bytes(), &rs); nil != err {
		return "", err
	}
	/* If there's an error, note it */
	if "" != rs.Error {
		return "", fmt.Errorf("%v", rs.Error)
	}
	return rs.Id, err
}
