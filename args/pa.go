package args

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

// ParsedArgs - container for parsed CLI args.
type ParsedArgs struct {
	Port uint
	Interface,
	Address,
	SdbInstanceAddr,
	SdbDBName,
	SdbAPIKey,
	SdbAPIValue,
	RefIDPrefix string
	EchoMode bool
}

// Parse - parses the command line arguments and populates the ParsedArgs object with the outcome
func Parse() *ParsedArgs {
	pa := &ParsedArgs{}

	flag.StringVar(&pa.Interface, "net-interface", "localhost", "network interface to serve on")
	flag.UintVar(&pa.Port, "port", 8000, "local port to serve on")
	flag.StringVar(&pa.SdbInstanceAddr, "sdb-address", "https://demo.slashdb.com", "SlashDB instance address")
	flag.StringVar(&pa.SdbDBName, "sdb-dbname", "timesheet", "SlashDB DB name i.e. https://demo.slashdb.com/db/>>timesheet<<")
	flag.StringVar(&pa.RefIDPrefix, "sdb-ref-id-prefix", "__href", "SlashDB's object ref URL prefix")
	flag.BoolVar(&pa.EchoMode, "echo-mode", true, "printout SlashDB's connection info - usefull for debugging")

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
