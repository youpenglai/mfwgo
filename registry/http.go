package registry

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func HttpPost(postUrl string, data string) (result []byte, err error) {
	body := bytes.NewBuffer([]byte(data))

	var resp *http.Response
	resp, err = http.Post(postUrl, "application/json", body)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func HttpPutWithHeader(putUrl string, headers map[string]string, data []byte) (result []byte, err error) {
	client := http.Client{}

	buff := bytes.NewBuffer(data)

	req, e := http.NewRequest("PUT", putUrl, buff)
	if e != nil {
		err = e
		return
	}

	for name, value := range headers {
		req.Header.Set(name, value)
	}

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func HttpPostWithHeader(postUrl string, headers map[string]string, data []byte) (result []byte, err error) {
	client := http.Client{}

	buff := bytes.NewBuffer(data)

	req, e := http.NewRequest("POST", postUrl, buff)
	if e != nil {
		err = e
		return
	}

	for name, value := range headers {
		req.Header.Set(name, value)
	}

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func HttpGetWithHeader(getUrl string, headers map[string]string) (result []byte, err error) {
	client := http.Client{}

	req, e := http.NewRequest("GET", getUrl, nil)
	if e != nil {
		err = e
		return
	}

	for name, value := range headers {
		req.Header.Set(name, value)
	}

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
