package main

import (
	"encoding/json"
	"errors"
	"fmt"
	env "github.com/logotipiwe/dc_go_env_lib"
	"net/http"
	"net/url"
	"os"
)

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
				fmt.Fprintf(w, "<a href='/oauth2/logout?redirect=%v'>Log out</a>", url.QueryEscape(env.GetPathToApp()))
				fmt.Fprint(w, getAdminPage())
			}
		} else {
			getLoginForm(w)
		}
	})

	http.HandleFunc("/api/create-service", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		name := r.PostFormValue("name")
		service := NewService(name)
		service.save()
		toIndex(w, r)
	})

	http.HandleFunc("/api/create-prop", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		serviceId := r.PostFormValue("service")
		namespaceId := r.PostFormValue("namespace")
		name := r.PostFormValue("name")
		value := r.PostFormValue("value")

		prop := NewProperty(name, value, namespaceId, serviceId)
		err := prop.save()
		if err != nil {
			println(fmt.Printf("Error saving prop: %v", err))
		}

		toIndex(w, r)
	})

	http.HandleFunc("/api/delete-prop", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id := r.PostFormValue("id")
		err := DeleteProperty(id)
		if err != nil {
			println(err.Error())
		}
		toIndex(w, r)
	})

	http.HandleFunc("/api/save-prop", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id := r.PostFormValue("id")
		name := r.PostFormValue("name")
		value := r.PostFormValue("value")
		println(name + " --- " + value)
		prop := GetProp(id)
		prop.Name = name
		prop.Value = value
		err := prop.save()
		if err != nil {
			println(err.Error())
		}
		toIndex(w, r)
	})

	http.HandleFunc("/api/activate-prop", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id := r.PostFormValue("id")
		prop := GetProp(id)
		prop.Active = true
		err := prop.save()
		if err != nil {
			println(err.Error())
		}
		toIndex(w, r)
	})

	http.HandleFunc("/api/deactivate-prop", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id := r.PostFormValue("id")
		prop := GetProp(id)
		prop.Active = false
		err := prop.save()
		if err != nil {
			println(err.Error())
		}
		toIndex(w, r)
	})

	err := http.ListenAndServe(":"+os.Getenv("CONTAINER_PORT"), nil)
	if err != nil {
		panic("Lol server fell")
	}
}

func toIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, os.Getenv("SUBPATH"), 302)
}

func getLoginForm(w http.ResponseWriter) {
	loginUrl, _ := url.Parse(env.GetCurrUrl() + "/oauth2/auth")
	q := loginUrl.Query()
	q.Set("redirect", env.GetCurrUrl()+env.GetSubpath())
	loginUrl.RawQuery = q.Encode()

	fmt.Fprintf(w, "<a href='%s'>%s</a>", loginUrl.String(), loginUrl.String())
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
