package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/boromil/timesheet/args"
	timesheetHttp "github.com/boromil/timesheet/http"
	assetfs "github.com/elazarl/go-bindata-assetfs"
)

func main() {
	pa := args.Parse()

	appCtx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	s := &http.Server{Addr: pa.Address, Handler: nil}

	go func() {
		<-signals
		cancelCtx()

		if err := s.Shutdown(appCtx); err != nil {
			log.Printf("error shuting down the server: %v\n", err)
		}
	}()

	timesheetHttp.SetupReverseProxy(pa.SdbDBName, pa.SdbInstanceAddr, pa.SdbAPIKey, pa.SdbAPIValue)
	afs := &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: ""}
	timesheetHttp.SetupBasicHandlers(pa.SdbDBName, afs)
	timesheetHttp.SetupAuthHandlers(pa.SdbDBName, pa.SdbInstanceAddr, pa.SdbAPIKey, pa.SdbAPIValue, pa.Address)

	fmt.Printf("Serving on http://%s/app/\n", pa.Address)
	log.Fatal(s.ListenAndServe())
}
