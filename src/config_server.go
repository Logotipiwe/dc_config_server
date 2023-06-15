package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const clientId = "319710408255-ntkf14k8ruk4p98sn2u1ho4j99rpjqja.apps.googleusercontent.com"

type User struct {
	Sub        string `json:"sub"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
	Locale     string `json:"locale"`
}

func main() {
	adminId := os.Getenv("LOGOTIPIWE_GMAIL_ID")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		fmt.Fprintf(w, "Hello, you've requested: %s</br>", r.URL.Path)
		cookie, _ := r.Cookie("access_token")
		var accessToken string
		if cookie != nil {
			accessToken = cookie.Value
		} else {
			accessToken = ""
		}
		if accessToken != "" {
			userData, err := getUserData(accessToken)
			if err != nil {
				getLoginForm(w)
				return
			}
			if userData.Sub != adminId {
				fmt.Fprintf(w, "Sorry, %s, you are not admin here!</br>", userData.Name)
				fmt.Fprintf(w, "<a href='/logout'>Log out</a>")
			} else {
				fmt.Fprintf(w, "Welcome: %s!</br>", userData.Name)
				fmt.Fprintf(w, getAdminPage())
				fmt.Fprintf(w, "<a href='/logout'>Log out</a>")
			}
		} else {
			getLoginForm(w)
		}
	})

	//http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "%s://%s", getScheme(), r.Host)
	//})

	http.HandleFunc("/g_oauth", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		println("Code is: " + code)
		token := exchangeCodeToToken(r, code)
		setATCookie(w, token)
		println(token)
		http.Redirect(w, r, "/", 302)
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		setATCookie(w, "")
		http.Redirect(w, r, "/", 302)
	})

	err := http.ListenAndServe(":"+os.Getenv("CONTAINER_PORT"), nil)
	if err != nil {
		panic("Lol server fell")
	}
}

func getLoginForm(w http.ResponseWriter) {
	loginUrl, _ := url.Parse("https://accounts.google.com/o/oauth2/v2/auth")
	q := loginUrl.Query()

	q.Set("client_id", clientId)
	q.Set("redirect_uri", getCurrHost()+getSubpath()+"/g_oauth")
	q.Set("response_type", "code")
	q.Set("scope", "profile")
	loginUrl.RawQuery = q.Encode()
	fmt.Fprintf(w, "<a href='%s'>%s</a>", loginUrl.String(), loginUrl.String())
}

func setATCookie(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     "access_token",
		Value:    token,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func exchangeCodeToToken(r *http.Request, code string) string {
	postUrl := "https://oauth2.googleapis.com/token"
	data := url.Values{}
	data.Set("client_id", clientId)
	data.Set("client_secret", os.Getenv("G_OAUTH_CLIENT_SECRET"))
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", getCurrHost()+getSubpath()+"/g_oauth")
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, postUrl, strings.NewReader(data.Encode())) // URL-encoded payload
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	var answer map[string]string
	json.NewDecoder(resp.Body).Decode(&answer)
	if resp.StatusCode != 200 {
		fmt.Printf("Got error while exchanging code to token. Status: %d. Body: %v", resp.StatusCode, answer)
	}
	return answer["access_token"]
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
