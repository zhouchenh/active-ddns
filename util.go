package main

import "strings"

// Copied from net/url/url.go
// Copyright 2009 The Go Authors. All rights reserved.
func splitHostPort(hostPort string) (host, port string) {
	host = hostPort
	colon := strings.LastIndexByte(host, ':')
	if colon != -1 && validOptionalPort(host[colon:]) {
		host, port = host[:colon], host[colon+1:]
	}
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = host[1 : len(host)-1]
	}
	return
}

// Copied from net/url/url.go
// Copyright 2009 The Go Authors. All rights reserved.
func validOptionalPort(port string) bool {
	if port == "" {
		return true
	}
	if port[0] != ':' {
		return false
	}
	for _, b := range port[1:] {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}
