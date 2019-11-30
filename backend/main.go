package main

import (
	"log"
	"net/http"
	"os"
	"sg-api/db"
	myhttp "sg-api/http"

	"github.com/jessevdk/go-flags"
	"rsc.io/quote"
)

var (
	BuildTime  string
	CommitHash string
)

func main() {

	log.Println(quote.Go())
	var opts db.Opts
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	log.Printf("Options: %#v\n", opts)
	opts.BuildTime = BuildTime
	opts.CommitHash = CommitHash

	// Enable line numbers in logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// initialize server
	srv := myhttp.NewServer()
	mydb, err := db.NewDB(&opts)
	if err != nil {
		log.Fatal(err)
	}
	srv.Db = mydb
	log.Println("Hooray. API runs locally at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", srv))
}
