package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/freshman-tech/news-demo-starter-files/news"
	"github.com/joho/godotenv"
)

var tpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}

// pagination functions
func (s *Search) IsLastPage() bool {
	return s.NextPage >= s.TotalPages
}

// previous page exist ?? button
func (s *Search) CurrentPage() int {
	if s.NextPage == 1 {
		return s.NextPage
	}

	return s.NextPage - 1
}
func (s *Search) PreviousPage() int {
	return s.CurrentPage() - 1
}

type Search struct {
	Query      string
	NextPage   int
	TotalPages int
	Results    *news.Results
}

func searchHandler(newsapi *news.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()
		searchQuery := params.Get("q")
		page := params.Get("page")
		if page == "" {
			page = "1"
		}

		results, err := newsapi.FetchEverything(searchQuery, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		nextPage, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		search := &Search{
			Query:      searchQuery,
			NextPage:   nextPage,
			TotalPages: int(math.Ceil(float64(results.TotalResults) / float64(newsapi.PageSize))),
			Results:    results,
		}
		// pagination control - -  - - - -
		if ok := !search.IsLastPage(); ok {
			search.NextPage++
		}
		// -------------
		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf.WriteTo(w)

		fmt.Printf("%+v", results)
	}
}

// IPAPIResponse represents the structure of the response from ip-api.com
type IPAPIResponse struct {
	Country string `json:"country"`
}

// determineLocation determines the location based on the user's IP address
func determineLocation(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	// Call the IP geolocation service (ip-api.com in this example)
	response, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		log.Println("Error getting location:", err)
		return "India" // Default location
	}
	defer response.Body.Close()

	// Decode the JSON response
	var ipAPIResponse IPAPIResponse
	err = json.NewDecoder(response.Body).Decode(&ipAPIResponse)
	if err != nil {
		log.Println("Error decoding location response:", err)
		return "India" // Default location
	}

	return ipAPIResponse.Country
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Fatal("Env: apiKey must be set")
	}

	myClient := &http.Client{Timeout: 10 * time.Second}
	newsapi := news.NewClient(myClient, apiKey, 20)

	fs := http.FileServer(http.Dir("assets"))

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/search", searchHandler(newsapi))
	// mux.HandleFunc("/", indexHandler)
	// Handler for root path, redirects based on location
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		location := determineLocation(r)
		redirectURL := "/search?q=" + location // You can modify this based on your routing logic
		http.Redirect(w, r, redirectURL, http.StatusPermanentRedirect)
	})
	http.ListenAndServe(":"+port, mux)
}
