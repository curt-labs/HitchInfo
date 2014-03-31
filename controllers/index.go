package controllers

import (
	"bufio"
	"encoding/csv"
	"fmt"
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

func IndexRedirect(rw http.ResponseWriter, req *http.Request) {
	http.Redirect(rw, req, "/", http.StatusMovedPermanently)
	return
}

func MassReplace(rw http.ResponseWriter, req *http.Request) {
	// Open the file for reading
	fi, err := os.Open("static/images.csv")
	if err != nil {
		fmt.Fprintln(rw, "Cound not read import file")
		return
	}

	// close fi on exit
	defer fi.Close()

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(rw, "Cound not read working directory")
		return
	}

	files, err := ioutil.ReadDir(wd + "/static")
	if err != nil {
		fmt.Fprintln(rw, "Cound not read static directory")
		return
	}

	oldFiles, err := ioutil.ReadDir(wd + "/static/old")
	if err != nil {
		fmt.Fprintln(rw, "Cound not read static/old directory")
		return
	}

	// make read buffer
	r := bufio.NewReader(fi)
	reader := csv.NewReader(r)
	for {
		line, err := reader.Read()
		if err != nil || len(line) < 2 {
			break
		}

		for _, file := range files {
			if !file.IsDir() {
				lines, err := readLines(wd+"/static/"+file.Name(), line[0], line[1])
				if err != nil {
					break
				}
				err = writeLines(lines, wd+"/static/"+file.Name())
				if err != nil {
					fmt.Fprintln(rw, "Cound not write file "+wd+"/static/"+file.Name())
					return
				}
			}
		}

		for _, file := range oldFiles {
			if !file.IsDir() {
				lines, err := readLines(wd+"/static/old/"+file.Name(), line[0], line[1])
				if err != nil {
					break
				}
				err = writeLines(lines, wd+"/static/old/"+file.Name())
				if err != nil {
					fmt.Fprintln(rw, "Cound not write file "+wd+"/static/old/"+file.Name())
					return
				}
			}
		}
		fmt.Fprintln(rw, "Replaced %s with %s", line[0], line[1])
	}
}

func readLines(path string, find string, replace string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
		if strings.Contains(strings.ToLower(txt), strings.ToLower(find)) {
			txt = strings.Replace(txt, find, replace, -1)
		}
		lines = append(lines, txt)
	}

	return lines, scanner.Err()
}

func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
