/*

Image and text scraper for XCentium test

/api/pageinfo takes url query (full url, http/s required)
and returns list of images, word count, and top 10 words.

*/

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"
)

type response struct {
	Images []image `json:"images"`
	Words  words   `json:"words"`
	Err    string  `json:"error"`
}

type image struct {
	Name   string `json:"name"`
	Desc   string `json:"description"`
	ImgURL string `json:"url"`
}

type words struct {
	TotalCount int      `json:"totalCount"`
	TopWords   wordList `json:"topWords"`
}

type word struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

// Required for sort interface
// Sort by top words, ascending - Top at index 0
type wordList []word

func (w wordList) Less(i, j int) bool {
	if w[i].Count > w[j].Count {
		return true
	}
	if w[i].Count < w[j].Count {
		return false
	}
	return w[i].Count > w[j].Count
}

func (w wordList) Len() int      { return len(w) }
func (w wordList) Swap(i, j int) { w[i], w[j] = w[j], w[i] }

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/pageinfo", getPageInfo).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func getPageInfo(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	q := r.URL.Query()
	urlQuery, ok := q["url"]

	if !ok {
		enc.Encode(response{Err: "URL was not provided."})
		return
	}
	if len(urlQuery) == 0 {
		enc.Encode(response{Err: "URL query is empty."})
		return
	}

	// err is returned on malformed URLs, cleanURL is safe
	cleanURL, err := url.ParseRequestURI(urlQuery[0]) // Client does not allow > 1
	if err != nil {
		enc.Encode(response{Err: err.Error()})
		return
	}

	scrapeResp := getResponse(cleanURL.String())
	enc.Encode(scrapeResp)
}
