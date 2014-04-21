package main

import (
	"fmt"
	"log"
	"net/url"
	"sort"
	"errors"
	"net/http"
	"encoding/gob"
	"os"
)

var cmdLogin = &Command {
	Run: runLogin,
	Name: "login",
}

func runLogin(cmd *Command, args []string) {
	var username string
	var password string
	fmt.Printf("Enter username: ")
	_, err := fmt.Scanln(&username)
	checkInputError("Username is required", err)
	fmt.Printf("Enter password: ")
	_, err = fmt.Scanln(&password)
	checkInputError("Password is required", err)
	cookies, err := login(username, password)
	checkInputError("", err)
	err = saveCookies(cookies, configFile)
	if err != nil {
		log.Fatalf("Internal error, could not save details")
	}
	log.Println("Logged in")
}

func checkInputError(errorMessage string, err error) {
	switch {
	case err == nil:
	case err.Error() == "unexpected newline":
		log.Fatalf(errorMessage)
	default:
		log.Fatalf(err.Error())
	}
}

func login(username string, password string) (token []*http.Cookie, err error) {
	apiEndpoint := fmt.Sprintf("%s/accounts/login/", apiURL)
	client.PostForm(apiEndpoint, url.Values{
		"username": {username},
		"password": {password},
	})
	// Peek into the cookie jar, see if we have a session id there
	baseUrl, _ := url.Parse(apiURL)
	cookies := client.Jar.Cookies(baseUrl)
	sessionIdx := sort.Search(len(cookies),
		func(i int) bool {
			return cookies[i].Name == "sessionid"
		})
	if sessionIdx < len(cookies) {
		return cookies, nil
	} else {
		return nil, errors.New("Login unsuccessful")
	}
}

func saveCookies(cookies []*http.Cookie, configFile string) (err error) {
	file, err := os.Create(configFile)
	enc := gob.NewEncoder(file)
	err = enc.Encode(cookies)
	defer func() {
		file.Close()
	}()
	return nil
}

func tryPrimeCookieJar() (err error) {
	file, err := os.Open(configFile)
	dec := gob.NewDecoder(file)
	var cookies []*http.Cookie
	err = dec.Decode(&cookies)
	baseUrl, err := url.Parse(apiURL)
	client.Jar.SetCookies(baseUrl, cookies)
	return nil
}
