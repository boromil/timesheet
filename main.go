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

	transport "github.com/boromil/timesheet/transport"
	"gitlab.com/boromil/goslashdb/slashdb"
)

func main() {
	parsedArgs := Parse()

	appCtx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	s := &http.Server{Addr: parsedArgs.Address, Handler: nil}

	go func() {
		<-signals
		cancelCtx()

		if err := s.Shutdown(appCtx); err != nil {
			log.Printf("error shuting down the server: %v\n", err)
		}
	}()

	fmt.Println(AssetNames())

	err := transport.SetupReverseProxy(
		parsedArgs.SdbDBName,
		parsedArgs.SdbInstanceAddr,
		parsedArgs.SdbAPIKey,
		parsedArgs.SdbAPIValue,
	)
	if err != nil {
		log.Fatalf("transport.SetupReverseProxy: %v", err)
	}

	err = transport.SetupBasicHandlers(parsedArgs.SdbDBName, AssetFile())
	if err != nil {
		log.Fatalf("transport.SetupBasicHandlers: %v", err)
	}

	externalHTTPClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			ResponseHeaderTimeout: time.Second * 10,
			IdleConnTimeout:       time.Second * 10,
			MaxIdleConns:          30,
			MaxIdleConnsPerHost:   3,
		},
	}

	sdbService, err := slashdb.NewService(
		parsedArgs.SdbInstanceAddr,
		parsedArgs.SdbAPIKey,
		parsedArgs.SdbAPIValue,
		parsedArgs.RefIDPrefix,
		parsedArgs.EchoMode,
		externalHTTPClient,
	)
	if err != nil {
		log.Fatalf("error initing SlashDB service: %v\n", err)
	}
	transport.Init(sdbService)

	fmt.Printf("Serving on http://%s/app/\n", parsedArgs.Address)
	log.Fatal(s.ListenAndServe())
}
