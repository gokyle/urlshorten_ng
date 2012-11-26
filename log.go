package main

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
)

const timestamp_fmt = "2006-01-02T15:04:05Z"

var ipScrub = regexp.MustCompile("^([^:]):.*$")

func LogRequest(page *Page, r *http.Request) {
	client_ip := ipScrub.ReplaceAllString(r.RemoteAddr, "$1")
	timestamp := time.Now().Format(timestamp_fmt)
	log_line := fmt.Sprintf("%s %s %s %s", client_ip, timestamp,
		r.Method, r.URL.Path)
	fmt.Println(log_line)
}
