package args

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

// ParsedArgs - container for parsed CLI args.
type ParsedArgs struct {
	Interf          string
	Port            uint
	Address         string
	SdbInstanceAddr string
	SdbDBName       string
	SdbAPIKey       string
	SdbAPIValue     string
}

// ParseArgs - parses the command line arguments and populates the ParsedArgs object with the outcome
func Parse() *ParsedArgs {
	pa := &ParsedArgs{}

	flag.StringVar(&pa.Interf, "net-interface", "localhost", "network interface to serve on")
	flag.UintVar(&pa.Port, "port", 8000, "local port to serve on")
	flag.StringVar(&pa.SdbInstanceAddr, "sdb-address", "https://demo.slashdb.com", "SlashDB instance address")
	flag.StringVar(&pa.SdbDBName, "sdb-dbname", "timesheet", "SlashDB DB name i.e. https://demo.slashdb.com/db/>>timesheet<<")
	var sdbAPIKey string
	flag.StringVar(
		&sdbAPIKey,
		"sdb-apikey", "apikey:timesheet-api-key", "SlashDB user API key, key and value separated by single ':'",
	)
	flag.Parse()

	pa.Address = fmt.Sprintf("%s:%d", pa.Interf, pa.Port)
	// extract SlashDB API key
	if tmp := strings.Split(sdbAPIKey, ":"); len(tmp) != 2 {
		log.Fatalln(fmt.Errorf("expected key, value pair, got: %s", sdbAPIKey))
	} else {
		pa.SdbAPIKey, pa.SdbAPIValue = tmp[0], tmp[1]
	}

	return pa
}
