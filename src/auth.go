package main

import (
	"errors"
	"fmt"
	"github.com/logotipiwe/dc_go_utils/src/config"
	"net/http"
)

func getLoginForm(w http.ResponseWriter) {
	page, err := getLoginPage()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, page)
}

func authAsAdmin(r *http.Request) error {
	secret, err := getSecretFromCookie(r)
	if err != nil {
		return err
	}
	if secret != config.GetConfig("AUTH_SECRET") {
		return errors.New("wrong auth secret")
	}
	return nil
}
func authAsMachine(r *http.Request) error {
	mToken := config.GetConfig("M_TOKEN")
	providedToken := r.URL.Query().Get("mToken")
	if mToken != providedToken {
		return errors.New("not a machine")
	}
	return nil
}

func setSecretToCookie(w http.ResponseWriter, secret string) {
	configPath := config.GetConfig("SUBPATH")
	var path string
	if configPath == "" {
		path = "/"
	} else {
		path = configPath
	}
	cookie := http.Cookie{
		Name:     "cs_secret",
		Value:    secret,
		HttpOnly: true,
		Path:     path,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
}

func getSecretFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("cs_secret")
	if err != nil {
		return "", err
	}
	var secret string
	if cookie != nil {
		secret = cookie.Value
	} else {
		secret = ""
	}
	return secret, nil
}
