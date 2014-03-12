package rest

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

func Get(path string) (buf []byte, err error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	buf, err = ioutil.ReadAll(resp.Body)
	return
}

func Post(path string, vals url.Values) (buf []byte, err error) {
	resp, err := http.PostForm(path, vals)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	buf, err = ioutil.ReadAll(resp.Body)

	return buf, err
}
