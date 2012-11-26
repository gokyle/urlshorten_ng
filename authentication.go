package main

import "bitbucket.org/taruti/pbkdf2"

const saltSize = 16

var admin_user string

func authenticate(username, password string) bool {
	if !check_auth {
		return true
	}
	if username == "" || password == "" {
		return false
	}
	ph, err := getPassHash(username)
	if err != nil {
		return false
	}
	return pbkdf2.MatchPassword(password, ph)
}
