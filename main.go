package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/yohcop/openid-go"
)

const domain = "localhost"
const port = ":8080"

// Load the templates once
var templateDir = "./"
var indexTemplate = template.Must(template.ParseFiles(templateDir + "index.html"))

// NoOpDiscoveryCache implements the DiscoveryCache interface and doesn't cache anything.
// For a simple website, I'm not sure you need a cache.
type NoOpDiscoveryCache struct{}

// Put is a no op.
func (n *NoOpDiscoveryCache) Put(id string, info openid.DiscoveredInfo) {}

// Get always returns nil.
func (n *NoOpDiscoveryCache) Get(id string) openid.DiscoveredInfo {
	return nil
}

var nonceStore = openid.NewSimpleNonceStore()
var discoveryCache = &NoOpDiscoveryCache{}

// indexHandler serves up the index template with the "Sign in through STEAM" button.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplate.Execute(w, nil)
}

// discoverHandler calls the Steam openid API and redirects to steam for login.
func discoverHandler(w http.ResponseWriter, r *http.Request) {
	url, err := openid.RedirectURL(
		"http://steamcommunity.com/openid",
		"http://"+domain+port+"/openidcallback",
		"http://"+domain+port+"/")

	if err != nil {
		log.Printf("Error creating redirect URL: %q\n", err)
	} else {
		http.Redirect(w, r, url, 303)
	}
}

// callbackHandler handles the response back from Steam. It verifies the callback and then renders
// the index template with the logged in user's id.
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	fullURL := "http://" + domain + port + r.URL.String()

	id, err := openid.Verify(fullURL, discoveryCache, nonceStore)
	if err != nil {
		log.Printf("Error verifying: %q\n", err)
	} else {
		log.Printf("NonceStore: %+v\n", nonceStore)
		data := make(map[string]string)
		data["user"] = id
		indexTemplate.Execute(w, data)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/discover", discoverHandler)
	http.HandleFunc("/openidcallback", callbackHandler)
	http.ListenAndServe(port, nil)
}
