package ebayprice

import (
	"encoding/json"
	"fmt"
	"github.com/heatxsink/go-ebay"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	ebay_appid = ""
)

func init() {
	bs, err := ioutil.ReadFile("ebay.pem")
	if err != nil {
		log.Fatalf("failed to read ebay appid: %v", err)
	}
	ebay_appid = strings.TrimSpace(string(bs))
	log.Printf("got appid: %q", ebay_appid)
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	keywords := r.URL.Path[len("/api/keywords/"):]
	fmt.Printf("keywords: %v", keywords)
	e := ebay.New(ebay_appid)
	response, err := e.FindItemsByKeywords(ebay.GLOBAL_ID_EBAY_US, keywords, 10)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "ERROR: ", err)
		return
	}
	outgoingJSON, err := json.Marshal(response)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(outgoingJSON))
}
