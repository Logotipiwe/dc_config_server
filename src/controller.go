package main

import (
	"encoding/json"
	"errors"
	"fmt"
	env "github.com/logotipiwe/dc_go_env_lib"
	. "github.com/logotipiwe/dc_go_utils/src"
	"github.com/logotipiwe/dc_go_utils/src/config"
	"log"
	"net/http"
	"os"
)

func main() {
	err := InitDb()
	if err != nil {
		panic(err)
	}
	logoutUrl := "/api/logout"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		println("/")
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		fmt.Fprintf(w, "Hello, you've requested: %s</br>", r.URL.Path)

		err := authAsAdmin(r)
		if err != nil {
			println(err.Error())
			getLoginForm(w)
			return
		}

		fmt.Fprintf(w, "Welcome admin!</br>")
		fmt.Fprintf(w, "<a href='%s'>Log out</a>", config.GetConfig("SUBPATH")+logoutUrl)
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
		err := authAsAdmin(r)
		if err != nil {
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
		err := authAsAdmin(r)
		if err != nil {
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
		err := authAsAdmin(r)
		if err != nil {
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
		err := authAsAdmin(r)
		if err != nil {
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
		err := authAsAdmin(r)
		if err != nil {
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
		err := authAsAdmin(r)
		if err != nil {
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
		println("/get-props")
		err := authAsMachine(r)
		err2 := authAsAdmin(r)
		if err != nil && err2 != nil {
			fmt.Println("Error when getting props. Cannot auth as machine, neither as admin")
			fmt.Println(err.Error())
			w.WriteHeader(403)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		namespace := r.URL.Query().Get("namespace")
		service := r.URL.Query().Get("service")
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
		println("/export")
		err := authAsMachine(r)
		err2 := authAsAdmin(r)
		if err != nil && err2 != nil {
			w.WriteHeader(403)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		props, err := GetAllProps()
		if err != nil {
			handleErrInController(w, err)
		}
		namespaces, err := GetAllNamespaces()
		if err != nil {
			handleErrInController(w, err)
			return
		}
		services, err := GetAllServices()
		if err != nil {
			handleErrInController(w, err)
			return
		}

		namespaceDtos := Map(namespaces, func(n Namespace) NamespaceDto { return n.toDto() })
		serviceDtos := Map(services, func(s Service) ServiceDto { return s.toDto() })
		propsDtos := Map(props, func(p Property) CSPropertyDto { return p.toDto() })

		err = json.NewEncoder(w).Encode(ImportExportAnswer{namespaceDtos, serviceDtos, propsDtos})
		if err != nil {
			handleErrInController(w, err)
		}
	})

	http.HandleFunc("/api/import", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("/import")
		err := authAsMachine(r)
		err2 := authAsAdmin(r)
		if err != nil && err2 != nil {
			w.WriteHeader(403)
			return
		}
		if r.Method != "POST" {
			handleBadRequest(w, errors.New("only post allowed"))
		}
		var data ImportExportAnswer
		err = json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			handleBadRequest(w, err)
		}
		namespaces := Map(data.Namespaces, func(n NamespaceDto) Namespace { return n.toModel() })
		services := Map(data.Services, func(s ServiceDto) Service { return s.toModel() })
		models := Map(data.Props, csPropToModel)
		err = importConfig(namespaces, services, models)
		if err != nil {
			handleErrInController(w, err)
		}
	})

	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		err = r.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		secret := r.PostFormValue("cs_login")
		if secret != config.GetConfig("AUTH_SECRET") {
			fmt.Println("Wrong secret: ", secret)
		} else {
			setSecretToCookie(w, secret)
		}
		toIndex(w, r)
	})

	http.HandleFunc(logoutUrl, func(w http.ResponseWriter, r *http.Request) {
		err := authAsAdmin(r)
		if err == nil {
			setSecretToCookie(w, "")
		}
		toIndex(w, r)
	})

	println("Ready")
	port := os.Getenv("CONTAINER_PORT")
	fmt.Println("Inner port is " + port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

func toIndex(w http.ResponseWriter, r *http.Request) {
	url := env.GetCurrUrl() + env.GetSubpath()
	fmt.Println("Redirecting to index " + url)
	http.Redirect(w, r, url, 302)
}

func handleBadRequest(w http.ResponseWriter, err error) {
	fmt.Println(err)
	w.WriteHeader(400)
	fmt.Fprintf(w, "{\"ok\": \"false\", \"err\":\"%s\"}", err.Error())
}

func handleErrInController(w http.ResponseWriter, err error) {
	fmt.Println(err)
	w.WriteHeader(500)
	fmt.Fprintf(w, "{\"ok\": \"false\", \"err\":\"%s\"}", err.Error())
}
