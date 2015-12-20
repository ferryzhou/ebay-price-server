package ebayprice

import (
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
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
	http.HandleFunc("/api/keywords/", handler)
	http.HandleFunc("/test", handlerTest)
}

func handlerTest(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	resp, err := client.Get("https://www.google.com/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "HTTP GET returned status %v<br/>%v", resp.Status, string(bs))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	keywords := r.URL.Path[len("/api/keywords/"):]
	fmt.Printf("keywords: %v", keywords)
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	e := NewEBay(ebay_appid, client)
	response, err := e.FindItemsByKeywords(GLOBAL_ID_EBAY_US, keywords, 10)
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
