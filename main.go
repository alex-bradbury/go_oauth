package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	githubOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:3000/callback",
		ClientID:     goDotEnvVariable("GITHUB_CLIENT_ID"),
		ClientSecret: goDotEnvVariable("GITHUB_CLIENT_SECRET"),
		Endpoint:     github.Endpoint,
	}
	state = "random"
)

func main() {
	// Hello world, the web server

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/auth", handleAuth)
	http.HandleFunc("/callback", handleCallback)

	log.Println("Listing for requests at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
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
	http.ServeFile(w, r, "./static/callback.html")

	// //Encode the data
	// postBody, _ := json.Marshal(map[string]string{
	// 	"client_id":  "Toby",
	// 	"client_secret": "Toby@example.com",
	// 	"code": "Toby@example.com",
	// 	"redirect_uri": "Toby@example.com",
	// })
	// responseBody := bytes.NewBuffer(postBody)

	// resp, err := http.Get("https://github.com/login/oauth/authorize?token=" + token.AccessToken)
	// if err != nil {
	// 	fmt.Println("Could not create request: "+err.Error())
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	// defer resp.Body.Close()
	// content, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println("Could not read content: "+err.Error())
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	// fmt.Fprintf(w, "Response: %s", content)
}

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
