package main

import "bitbucket.org/taruti/pbkdf2"

func authenticate(username, password string) bool {
	if username == "" || password == "" {
		return false
	}
	ph, err := getPassHash(username)
	if err != nil {
		return false
	}
	return pbkdf2.MatchPassword(password, ph)
}
