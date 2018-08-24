package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/boromil/timesheet/args"
	transport "github.com/boromil/timesheet/transport"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"gitlab.com/boromil/goslashdb/slashdb"
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

	transport.SetupReverseProxy(pa.SdbDBName, pa.SdbInstanceAddr, pa.SdbAPIKey, pa.SdbAPIValue)
	afs := &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: ""}
	transport.SetupBasicHandlers(pa.SdbDBName, afs)

	var sdbService slashdb.Service
	sdbService, err := slashdb.NewService(
		pa.SdbInstanceAddr,
		pa.SdbAPIKey,
		pa.SdbAPIValue,
		pa.RefIDPrefix,
		pa.EchoMode,
		&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				ResponseHeaderTimeout: time.Second * 10,
				IdleConnTimeout:       time.Second * 10,
				MaxIdleConns:          30,
				MaxIdleConnsPerHost:   3,
			},
		},
	)
	if err != nil {
		log.Fatalf("error initing SlashDB service: %v\n", err)
	}
	transport.SetupAuthHandlers(sdbService)

	fmt.Printf("Serving on http://%s/app/\n", pa.Address)
	log.Fatal(s.ListenAndServe())
}
