package api

// TODO: Must build itunes links according to this:
// https://www.apple.com/itunes/affiliates/resources/documentation/linking-to-the-itunes-music-store.html#AffiliateEncodingLinkShare
// example link:
//

import (
	"fmt"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"strings"
	"log"
)


type itunesMovie struct {
	DirectorName string	`json:"artistName"`
	Title string		`json:"trackName"`
	Price float64		`json:"trackPrice"`
	Url string		`json:"trackViewUrl"`
}

type ItunesResponse struct {
	// TODO: get actual number of results from itunes
	Movies []*itunesMovie		`json:"results"`
	RawData []byte
}

func (s ItunesResponse) listTitles() {
	for _, m := range s.Movies {
		fmt.Println(m.Title)
	}
}

func (s ItunesResponse) Search(title, director string) ([]itunesMovie) {
	// Searching on a search seems redundant, but is needed
	// rename to filter?
	// so we can verify we are giving correct data
	var (
		movies []itunesMovie
	)

	for _, movie := range s.Movies {
		log.Printf("itunes: '%s' == '%s'? %v.\n", movie.Title, title, strings.EqualFold(movie.Title, title))
		if strings.EqualFold(movie.Title, title) {
			movies = append(movies, *movie)
		}
	}

	return movies
}

func findMovie(keyword string, director string) ([]itunesMovie, []byte, error) {
	keyword = normalizeKeyword(keyword)
	apiUrl, err := url.Parse("https://itunes.apple.com/search?media=movie&term=" + url.QueryEscape(keyword))
	if err != nil {
		log.Println(err)
	}

	resp, err := http.Get(apiUrl.String())
	if err != nil {
		log.Println(err)
		//return Result{}
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	// see if there is function to split bytes like this so cast isn't needed
	var results ItunesResponse
	// TODO: Fix entire structure so that this duplication isn't necessary
	// and I can set the SearchResult directly

	err = json.Unmarshal(content, &results)

	if err != nil {
		log.Println("itunes Json unmarshal:", err)
	}

	title, _ := url.QueryUnescape(keyword)

	// this is only necessary for some api's which don't have a great search
	movieEntry := results.Search(title, director)
	return movieEntry, content, err
}

type Itunes struct {}

func (i Itunes) Query(query string) (SearchResult) {
	var(
		// update this to use a reference
		movies []Media
		queryResult SearchResult = SearchResult{false, movies, nil, "iTunes"} 
	)

	// change to okay or not check!
	receivedMovies, rawResponseData, err := findMovie(query, "")
	if err != nil {
		// if there was an actual error in getting the query, then we failed to get
		// the query or there was an issue on their side
		return queryResult
	}

	for _, receivedMovie := range receivedMovies {
		if receivedMovie.Title == "" {
			// TODO: get actual number of results from itunes
			log.Printf("Failed to match <query: %s> <service: %s> (# results not implemented)", query, "itunes")
		}

		// build type.Media for queryResult
		//{receivedMovie.Title, receivedMovie.Price, receivedMovie.Url, "itunes"}
		movie := Media{}
		movie.Title = receivedMovie.Title
		movie.Price = receivedMovie.Price
		movie.Link = receivedMovie.Url
		movie.PurchaseType = "Buy"

		movies = append(movies, movie)
	}

	if (len(movies) > 0) {
		queryResult.Success = true
	}

	queryResult.Media = movies
	queryResult.RawData = rawResponseData

	return queryResult
}
