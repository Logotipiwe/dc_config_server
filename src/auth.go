package main

import (
	"encoding/json"
	"errors"
	"fmt"
	env "github.com/logotipiwe/dc_go_env_lib"
	"net/http"
	"net/url"
)

type User struct {
	Sub        string `json:"sub"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
	Locale     string `json:"locale"`
}

func getUserData(accessToken string) (User, error) {
	bearer := "Bearer " + accessToken
	getUrl := "https://www.googleapis.com/oauth2/v3/userinfo"
	request, _ := http.NewRequest("GET", getUrl, nil)
	request.Header.Add("Authorization", bearer)

	client := &http.Client{}
	res, _ := client.Do(request)
	defer res.Body.Close()
	var answer User
	json.NewDecoder(res.Body).Decode(&answer)
	if answer.Sub != "" {
		return answer, nil
	} else {
		return answer, errors.New("WTF HUH")
	}
}

func getLoginForm(w http.ResponseWriter) {
	loginUrl, _ := url.Parse(env.GetCurrUrl() + "/oauth2/auth")
	q := loginUrl.Query()
	q.Set("redirect", env.GetCurrUrl()+env.GetSubpath())
	loginUrl.RawQuery = q.Encode()

	fmt.Fprintf(w, "<a href='%s'>%s</a>", loginUrl.String(), loginUrl.String())
}
