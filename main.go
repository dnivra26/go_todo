package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/newrelic/go-agent"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	DB_USER     = "dms"
	DB_PASSWORD = "postgres"
	DB_NAME     = "dms"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	if txn, ok := w.(newrelic.Transaction); ok {
		segment := newrelic.StartSegment(txn, "Go Time call")
		response, err := http.Get("http://localhost:9091")
		segment.End()
		segment1 := newrelic.StartSegment(txn, "DB call")
		dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			DB_USER, DB_PASSWORD, DB_NAME)
		db, err := sql.Open("postgres", dbinfo)
		rows, err := db.Query("SELECT * FROM document")
		var count int
		rows.Scan(&count)
		fmt.Println(count)
		defer db.Close()
		segment1.End()

		if err != nil {
			fmt.Printf("%s", err)
		} else {
			defer response.Body.Close()
			contents, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Printf("%s", err)
			}
			fmt.Fprintf(w, "Hello world!"+string(contents)) // send data to client side
		}
	}

}

func main() {
	config := newrelic.NewConfig("Go ToDo", "909dd7e6d21da9499f9de1c8d5b3aa2995342975")
	app, err := newrelic.NewApplication(config)
	txn := app.StartTransaction("todo", nil, nil)
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/", sayhelloName))
	http.Handle("/metrics", promhttp.Handler())
	error := http.ListenAndServe(":9090", nil) // set listen port
	if error != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	defer txn.End()
}
