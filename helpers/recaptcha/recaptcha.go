// Package recaptcha handles reCaptcha (http://www.google.com/recaptcha) form submissions
//
// This package is designed to be called from within an HTTP server or web framework
// which offers reCaptcha form inputs and requires them to be evaluated for correctness
//
// Edit the recaptcha_private_key constant before building and using
package recaptcha

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	RECAPTCHA_SERVER_NAME = "http://www.google.com/recaptcha/api/verify"
	RECAPTCHA_PRIVATE_KEY = "6Levd8cSAAAAAO_tjAPFuXbfzj6l5viTEaz5YjVv"
	RECAPTCHA_PUBLIC_KEY  = "6Levd8cSAAAAAGJDbZzk9hRRqm0ltuS8-TgwadQS"
)

// check uses the client ip address, the challenge code from the reCaptcha form,
// and the client's response input to that challenge to determine whether or not
// the client answered the reCaptcha input question correctly.
// It returns a boolean value indicating whether or not the client answered correctly.
func check(remoteip, challenge, response string) (s string) {
	s = ""
	resp, err := http.PostForm(RECAPTCHA_SERVER_NAME,
		url.Values{"privatekey": {RECAPTCHA_PRIVATE_KEY}, "remoteip": {remoteip}, "challenge": {challenge}, "response": {response}})
	if err != nil {
		log.Println("Post error: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read error: could not read body: %s", err)
	} else {
		s = string(body)
	}
	return
}

// Confirm is the public interface function.
// It calls check, which the client ip address, the challenge code from the reCaptcha form,
// and the client's response input to that challenge to determine whether or not
// the client answered the reCaptcha input question correctly.
// It returns a boolean value indicating whether or not the client answered correctly.
func Confirm(remoteip, challenge, response string) (result bool) {
	result = strings.HasPrefix(check(remoteip, challenge, response), "true")
	return
}

// Init allows the webserver or code evaluating the reCaptcha form input to set the
// reCaptcha private key (string) value, which will be different for every domain.
func Init(theme string) string {
	if theme == "" {
		theme = "clean"
	}
	captcha := fmt.Sprintf("<script type=\"text/javascript\">var RecaptchaOptions = { theme : '%s' };</script>", theme)
	captcha = fmt.Sprintf("%s<script type=\"text/javascript\" src=\"http://www.google.com/recaptcha/api/challenge?k=%s\"></script>", captcha, RECAPTCHA_PUBLIC_KEY)
	captcha = fmt.Sprintf("%s<noscript><iframe src=\"http://www.google.com/recaptcha/api/noscript?k=%s\" height=\"300\" width=\"500\" frameborder=\"0\"></iframe><br><textarea name=\"recaptcha_challenge_field\" rows=\"3\" cols=\"40\"></textarea><input type=\"hidden\" name=\"recaptcha_response_field\" value=\"manual_challenge\"></noscript>", captcha, RECAPTCHA_PUBLIC_KEY)

	return captcha
}
