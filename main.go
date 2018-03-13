package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/newrelic/go-agent"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	if txn, ok := w.(newrelic.Transaction); ok {
		segment := newrelic.StartSegment(txn, "Go Time call")
		response, err := http.Get("http://localhost:9091")
		segment.End()

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
