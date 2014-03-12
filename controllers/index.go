package controllers

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func Index(rw http.ResponseWriter, req *http.Request) {
	dir, err := os.Getwd()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	qry, err := url.QueryUnescape(req.URL.RequestURI())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if qry == "" {
		qry = "index.html"
	}

	contents, err := ioutil.ReadFile(dir + "/static/old" + qry)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	ct := "text/html"
	if strings.Contains(qry, ".gif") {
		ct = "image/gif"
	} else if strings.Contains(qry, ".js") {
		ct = "text/javascript"
	} else if strings.Contains(qry, ".png") {
		ct = "image/png"
	}
	rw.Header().Set("Content-Type", ct)
	rw.Write(contents)
}
