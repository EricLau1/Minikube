package api

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
)

var (
	port = flag.String("p", ":8080", "set port")
)

func init() {
	flag.Parse()
}

func Run() {
	log.Printf("Api listening on http://localhost%s\n", *port)
	http.HandleFunc("/", hello)
	http.HandleFunc("/fail", fail)
	log.Fatal(http.ListenAndServe(*port, nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF8")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "Hello World!"})
}

func fail(w http.ResponseWriter, r *http.Request) {
	log.Fatal("server down")
	os.Exit(1)
}
