package transport

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// SetupBasicHandlers setups page index and static assets filestystem
func SetupBasicHandlers(sdbDBName string, afs http.FileSystem) error {
	tmplData := struct{ SdbDBName string }{SdbDBName: sdbDBName}
	indexTmpl := template.New("index.html")
	indexFile, err := afs.Open("templates/index.html")
	if err != nil {
		return fmt.Errorf("afs.Open: %w", err)
	}
	defer indexFile.Close()

	templateData, err := ioutil.ReadAll(indexFile)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %w", err)
	}

	if _, err = indexTmpl.Parse(string(templateData)); err != nil {
		return fmt.Errorf("indexTmpl.Parse: %w", err)
	}

	http.HandleFunc("/app/", func(w http.ResponseWriter, r *http.Request) {
		if err = indexTmpl.Execute(w, tmplData); err != nil {
			log.Printf("indexTmpl.Execute: %v\n", err)
			return
		}
	})
	http.Handle("/app/static/", http.StripPrefix("/app/static/", http.FileServer(afs)))

	return nil
}

// SetupReverseProxy setups the reverse proxy to SlashDB instance
func SetupReverseProxy(
	sdbDBName,
	sdbInstanceAddr,
	sdbAPIKey,
	sdbAPIValue string,
) error {
	// get address for the SlashDB instance and parse the URL
	url, err := url.Parse(sdbInstanceAddr)
	if err != nil {
		return fmt.Errorf("failed to parse sdbInstanceAddr: %w", err)
	}

	// create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)
	// make it play nice with https endpoints, also add some timeouts
	proxy.Transport = defaultTransport

	proxyHandler := func(w http.ResponseWriter, r *http.Request) {
		// set API key header
		r.Header.Set(sdbAPIKey, sdbAPIValue)
		// set CORS headers for easy proxy to SDB communication
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Accept, Origin, Content-Type, Content-Length, X-Requested-With, Accept-Encoding, X-CSRF-Token, Authorization",
		)
		log.Printf("passing on a %q request to: %q\n", r.Method, r.URL.String())
		proxy.ServeHTTP(w, r)
	}
	// bind the proxy handler to "/"
	http.HandleFunc("/", authorizationMiddleware(sdbDBName, proxyHandler, nil))

	return nil
}
