package main

import (
	"fmt"
	"github.com/Logotipiwe/dc_go_auth_lib/auth"
	env "github.com/logotipiwe/dc_go_env_lib"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	adminId := os.Getenv("LOGOTIPIWE_GMAIL_ID")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		println("/")
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		_, err := fmt.Fprintf(w, "Hello, you've requested: %s</br>", r.URL.Path)
		if err != nil {
			log.Fatalln(err)
		}
		cookie, _ := r.Cookie("access_token")
		var accessToken string
		if cookie != nil {
			accessToken = cookie.Value
		} else {
			accessToken = ""
		}
		if accessToken != "" {
			userData, err := auth.FetchUserData(r)
			if err != nil {
				getLoginForm(w)
				return
			}
			if userData.Id != adminId {
				_, err := fmt.Fprintf(w, "Sorry, %s, you are not admin here!</br> <a href='/logout'>Log out</a>", userData.Name)
				if err != nil {
					log.Fatalln(err)
				}
			} else {
				_, err := fmt.Fprintf(w, "Welcome: %s!</br>", userData.Name)
				if err != nil {
					log.Fatalln(err)
				}
				_, err = fmt.Fprintf(w, "<a href='/oauth2/logout?redirect=%v'>Log out</a>", url.QueryEscape(env.GetPathToApp()))
				if err != nil {
					log.Fatalln(err)
				}
				_, err = fmt.Fprint(w, getAdminPage())
				if err != nil {
					log.Fatalln(err)
				}
			}
		} else {
			getLoginForm(w)
		}
	})

	http.HandleFunc("/api/create-service", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		name := r.PostFormValue("name")
		service := CreateService(name)
		err = service.save()
		if err != nil {
			println(fmt.Sprintf("Error saving service, %v", err.Error()))
		} else {
			println(fmt.Sprintf("Service with name %v created!", name))
		}
		toIndex(w, r)
	})

	http.HandleFunc("/api/create-prop", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		serviceId := r.PostFormValue("service")
		namespaceId := r.PostFormValue("namespace")
		name := r.PostFormValue("name")
		value := r.PostFormValue("value")

		prop := CreateProperty(name, value, namespaceId, serviceId)
		err = prop.save()
		if err != nil {
			println(fmt.Sprintf("Error saving prop: %v", err))
		}
		println(fmt.Sprintf("Prop created!"))

		toIndex(w, r)
	})

	http.HandleFunc("/api/delete-prop", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		id := r.PostFormValue("id")
		err = DeleteProperty(id)
		if err != nil {
			println(err.Error())
		}
		toIndex(w, r)
	})

	http.HandleFunc("/api/save-prop", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		id := r.PostFormValue("id")
		name := r.PostFormValue("name")
		value := r.PostFormValue("value")
		println(name + " --- " + value)
		prop := GetProp(id)
		prop.Name = name
		prop.Value = value
		err = prop.save()
		if err != nil {
			println(err.Error())
		}
		toIndex(w, r)
	})

	http.HandleFunc("/api/activate-prop", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		id := r.PostFormValue("id")
		prop := GetProp(id)
		prop.Active = true
		err = prop.save()
		if err != nil {
			println(err.Error())
		}
		toIndex(w, r)
	})

	http.HandleFunc("/api/deactivate-prop", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		id := r.PostFormValue("id")
		prop := GetProp(id)
		prop.Active = false
		err = prop.save()
		if err != nil {
			println(err.Error())
		}
		toIndex(w, r)
	})

	println("Ready")
	err := http.ListenAndServe(":"+os.Getenv("CONTAINER_PORT"), nil)
	if err != nil {
		panic("Lol server fell")
	}
}

func toIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, env.GetSubpath(), 302)
}
