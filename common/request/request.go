package request

import (
	"io/ioutil"
	"net/http"
	"encoding/json"
	"bytes"
	"github.com/zhifeichen/bluesky-protocol/common/xlogger"
	"errors"
	"io"
	"strings"
	"net/url"
	"fmt"
)

type HttpMethod string

const (
	HTTP_GET    HttpMethod = "GET"
	HTTP_POST   HttpMethod = "POST"
	HTTP_PUT    HttpMethod = "PUT"
	HTTP_DELETE HttpMethod = "DELETE"
)

func requestGet(uri string, params map[string]string, headers map[string]string) (string, error) {
	client := &http.Client{}
	ps := make([]string, 0, len(params))
	for k, v := range params {
		ps = append(ps, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
	}
	uri = uri + "?" + strings.Join(ps, "&")
	xlogger.Debugf("GET %s no body", uri)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		// handle error
		xlogger.Error(err)
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		xlogger.Error(err)
		return "", err
	}

	return string(body), nil
}

func requestPost(uri, method string, params map[string]string, headers map[string]string) (string, error) {
	client := &http.Client{}
	var r io.Reader
	if params != nil {
		reqBody, err := json.Marshal(params)

		if err != nil {
			xlogger.Error(err)
			return "", err
		}
		xlogger.Debugf("%s %s %s", uri, method, string(reqBody))
		r = bytes.NewReader(reqBody)
	} else {
		xlogger.Debugf("%s %s no body", uri, method)
	}

	req, err := http.NewRequest(method, uri, r)
	if err != nil {
		// handle error
		xlogger.Error(err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		xlogger.Errorf("%s %s 错误 StatusCode:%d", method, uri, resp.StatusCode)
		return "", errors.New("http response status error")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		xlogger.Error(err)
		return "", err
	}

	return string(body), nil
}

//Request
func Request(url string, method HttpMethod, params map[string]string, headers map[string]string) (string, error) {

	switch(method) {
	case HTTP_GET:
		return requestGet(url, params, headers)
	case HTTP_POST:
		fallthrough
	case HTTP_PUT:
		fallthrough
	case HTTP_DELETE:
		fallthrough
	default:
		return requestPost(url, string(method), params, headers)
	}
}

func Fetch(uri string, method HttpMethod, body io.Reader, headers map[string]string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(string(method), uri, body)
	if err != nil {
		// handle error
		xlogger.Error(err)
		return "", err
	}
	if method != HTTP_GET {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		xlogger.Errorf("%s %s 错误 StatusCode:%d", method, uri, resp.StatusCode)
		return "", errors.New("http response status error")
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		xlogger.Error(err)
		return "", err
	}

	return string(res), nil
}
