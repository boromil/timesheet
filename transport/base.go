package transport

import (
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	assetfs "github.com/elazarl/go-bindata-assetfs"
)

func SetupBasicHandlers(sdbDBName string, afs *assetfs.AssetFS) {
	tmpData := struct{ SdbDBName string }{SdbDBName: sdbDBName}
	http.HandleFunc("/app/", func(w http.ResponseWriter, r *http.Request) {
		indexTmpl := template.New("index.html")
		data, err := afs.Asset("index.html")
		if err != nil {
			log.Printf("afs.Asset: %v\n", err)
		}
		if _, err = indexTmpl.Parse(string(data)); err != nil {
			log.Printf("indexTmpl.Parse: %v\n", err)
		}
		if err = indexTmpl.Execute(w, tmpData); err != nil {
			log.Printf("indexTmpl.Execute: %v\n", err)
		}
	})
	http.Handle("/app/static/", http.StripPrefix("/app/static/", http.FileServer(afs)))
}

func SetupReverseProxy(sdbDBName, sdbInstanceAddr, sdbAPIKey, sdbAPIValue string) {
	// get address for the SlashDB instance and parse the URL
	url, err := url.Parse(sdbInstanceAddr)
	if err != nil {
		log.Fatalf("failed to parse sdbInstanceAddr: %v\n", err)
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
}
