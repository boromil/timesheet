package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/boromil/timesheet-demo/args"
	timesheetHttp "github.com/boromil/timesheet-demo/http"
	assetfs "github.com/elazarl/go-bindata-assetfs"
)

func main() {
	pa := args.Parse()
	timesheetHttp.SetupReverseProxy(pa.SdbDBName, pa.SdbInstanceAddr, pa.SdbAPIKey, pa.SdbAPIValue)
	afs := &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: ""}
	timesheetHttp.SetupBasicHandlers(pa.SdbDBName, afs)
	timesheetHttp.SetupAuthHandlers(pa.SdbDBName, pa.SdbInstanceAddr, pa.SdbAPIKey, pa.SdbAPIValue, pa.Address)
	fmt.Printf("Serving on http://%s/app/\n", pa.Address)
	log.Fatal(http.ListenAndServe(pa.Address, nil))
}
