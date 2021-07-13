package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	githubOauthConfig = &oauth2.Config{
		RedirectURL:  goDotEnvVariable("CALLBACK_URL"),
		ClientID:     goDotEnvVariable("GITHUB_CLIENT_ID"),
		ClientSecret: goDotEnvVariable("GITHUB_CLIENT_SECRET"),
		Endpoint:     github.Endpoint,
	}
	state = "random"
)

func main() {
	// Hello world, the web server

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/auth", handleAuth)
	http.HandleFunc("/callback", handleCallback)

	log.Println("Listing for requests at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	url := githubOauthConfig.AuthCodeURL(state)
	//url := "https://github.com/login/oauth/authorize?client_id=" + goDotEnvVariable("GITHUB_CLIENT_ID")
	fmt.Println(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != state {
		fmt.Println("State is not valid")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Println(r.FormValue("code"))
	token, err := githubOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		fmt.Println("Could not get token: " + err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Println("Token: " + token.AccessToken)
	client := httpClient{*http.DefaultClient, token.AccessToken}
	//http.ServeFile(w, r, "./static/callback.html")

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		fmt.Println("Could not create request: " + err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Could not read content: " + err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var noScopeResp noScope
	json.Unmarshal([]byte(content), &noScopeResp)

	tmpl := template.Must(template.ParseFiles("./static/callback.html"))
	tmpl.Execute(w, noScopeResp)
}
