package main

import (
	"fmt"
	env "github.com/logotipiwe/dc_go_env_lib"
	"net/http"
	"net/url"
)

func getLoginForm(w http.ResponseWriter) {
	loginUrl, _ := url.Parse(env.GetCurrUrl() + "/oauth2/auth")
	q := loginUrl.Query()
	q.Set("redirect", env.GetCurrUrl()+env.GetSubpath())
	loginUrl.RawQuery = q.Encode()

	fmt.Fprintf(w, "<a href='%s'>%s</a>", loginUrl.String(), loginUrl.String())
}
