package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Logotipiwe/dc_go_auth_lib/auth"
	env "github.com/logotipiwe/dc_go_env_lib"
	. "github.com/logotipiwe/dc_go_utils/src"
	"github.com/logotipiwe/dc_go_utils/src/config"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	err := InitDb()
	if err != nil {
		panic(err)
	}
	idpUrl := config.GetConfig("IDP_HOST") + config.GetConfig("IDP_SUBPATH")
	logoutUrl := idpUrl + "/logout?redirect=" + url.QueryEscape(env.GetPathToApp())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		println("/")
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		fmt.Fprintf(w, "Hello, you've requested: %s</br>", r.URL.Path)

		admin, err := authAsAdmin(r)
		if err != nil && err.Error() == "not admin" {
			fmt.Fprintf(w, "Sorry, %s, you are not admin here!</br> <a href='%s'>Log out</a>", admin.Name, logoutUrl)
			return
		}
		if err != nil {
			println(err.Error())
			getLoginForm(w)
			return
		}

		fmt.Fprintf(w, "Welcome: %s!</br>", admin.Name)
		fmt.Fprintf(w, "<a href='%s'>Log out</a>", logoutUrl)
		adminPage, err := getAdminPage()
		if err != nil {
			println(err.Error())
			w.WriteHeader(500)
			return
		}
		fmt.Fprint(w, adminPage)
	})

	http.HandleFunc("/api/create-service", func(w http.ResponseWriter, r *http.Request) {
		println("/create-service")
		_, err := authAsAdmin(r)
		if err != nil && err.Error() == "not admin" {
			toIndex(w, r)
			return
		}
		err = r.ParseForm()
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
		println("/create-prop")
		_, err := authAsAdmin(r)
		if err != nil && err.Error() == "not admin" {
			toIndex(w, r)
			return
		}
		err = r.ParseForm()
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
		println("/delete-prop")
		_, err := authAsAdmin(r)
		if err != nil && err.Error() == "not admin" {
			toIndex(w, r)
			return
		}
		err = r.ParseForm()
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
		println("/save-prop")
		_, err := authAsAdmin(r)
		if err != nil && err.Error() == "not admin" {
			toIndex(w, r)
			return
		}
		err = r.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		id := r.PostFormValue("id")
		name := r.PostFormValue("name")
		value := r.PostFormValue("value")
		println(name + " --- " + value)
		prop, err := GetProp(id)
		prop.Name = name
		prop.Value = value
		err = prop.save()
		if err != nil {
			println(err.Error())
		}
		toIndex(w, r)
	})

	http.HandleFunc("/api/activate-prop", func(w http.ResponseWriter, r *http.Request) {
		println("/activate-prop")
		_, err := authAsAdmin(r)
		if err != nil && err.Error() == "not admin" {
			toIndex(w, r)
			return
		}
		err = r.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		id := r.PostFormValue("id")
		prop, err := GetProp(id)
		prop.Active = true
		err = prop.save()
		if err != nil {
			println(err.Error())
		}
		toIndex(w, r)
	})

	http.HandleFunc("/api/deactivate-prop", func(w http.ResponseWriter, r *http.Request) {
		println("/deactivate-prop")
		_, err := authAsAdmin(r)
		if err != nil && err.Error() == "not admin" {
			toIndex(w, r)
			return
		}
		err = r.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		id := r.PostFormValue("id")
		prop, err := GetProp(id)
		prop.Active = false
		err = prop.save()
		if err != nil {
			println(err.Error())
		}
		toIndex(w, r)
	})

	http.HandleFunc("/api/get-config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		println("/get-props")
		namespace := r.URL.Query().Get("namespace")
		if namespace == "" {
			w.WriteHeader(400)
			fmt.Fprint(w, "{\"error\": \"namespace is empty\"}")
			return
		}
		service := r.URL.Query().Get("service")
		if service == "" {
			w.WriteHeader(400)
			fmt.Fprint(w, "{\"error\": \"service is empty\"}")
			return
		}
		props, err := GetPropsByNamespaceAndService(namespace, service)
		if err != nil {
			println(fmt.Sprintf("Err getting props: %s", err.Error()))
			w.WriteHeader(500)
			return
		}
		propsDtos := Map(props, func(p Property) CSPropertyDto {
			return p.toDto()
		})
		err = json.NewEncoder(w).Encode(propsDtos)
		if err != nil {
			println("Error encoding props ;" + err.Error())
		}
	})

	//TODO auth all requests
	http.HandleFunc("/api/export", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		println("/export")
		props, err := GetAllProps()
		if err != nil {
			handleErrInController(w, err)
		}

		dtos := Map(props, func(p Property) CSPropertyDto {
			return p.toDto()
		})
		err = json.NewEncoder(w).Encode(PropsAnswer{dtos})
		if err != nil {
			handleErrInController(w, err)
		}
	})

	//TODO import services and namespaces
	http.HandleFunc("/api/import", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			handleBadRequest(w, errors.New("only post allowed"))
		}
		var data PropsAnswer
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			handleBadRequest(w, err)
		}
		models := Map(data.Props, csPropToModel)
		err = importProps(models)
		if err != nil {
			handleErrInController(w, err)
		}
	})

	println("Ready")
	err = http.ListenAndServe(":"+os.Getenv("CONTAINER_PORT"), nil)
	if err != nil {
		panic("Lol server fell")
	}
}

func toIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, env.GetSubpath(), 302)
}

func handleBadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(400)
	fmt.Fprintf(w, "{\"ok\": \"false\", \"err\":\"%s\"}", err.Error())
}

func handleErrInController(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	fmt.Fprintf(w, "{\"ok\": \"false\", \"err\":\"%s\"}", err.Error())
}

func authAsAdmin(r *http.Request) (auth.DcUser, error) {
	adminId := config.GetConfig("LOGOTIPIWE_GMAIL_ID")
	userData, err := auth.FetchUserData(r)
	if err != nil {
		return userData, err
	}
	if userData.Id != adminId {
		return userData, errors.New("not admin")
	}
	return userData, nil
}
