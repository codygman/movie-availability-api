package api

// TODO: Add the information for generating proper affiliate id's in links
// and post details under here
//TODO: Make more like the itunes.go wrapper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"encoding/json"
	"strings"
	"log"
)

type YoutubeEntry struct {
	Title struct {
		Value string	`json:"$t"`
	}
	Media *struct {
		PlaceHolder string
		MediaContent []struct {
			Url string
		} `json:"media$content"`
		MediaCredit []struct {
			Name string	`json:"$t"`
			Role string
		} `json:"media$credit"`
		// TODO: Make sure the case of this not parsing correctly is handled correctly
		MediaPrice []struct {
			// youtube only allows you to buy the two cheapest purchase/rent options
			Price float64 `json:",string"`
			MediaType string `json:"type"`
		} `json:"media$price"`

	} `json:"media$group"`
	Link []struct {
		Href string
	}
}

type YoutubeResponse struct {
	Feed struct {
		SearchResults struct{ 
			Value int	`json:"$t"`
		} `json:"openSearch$totalResults"`
		Entries []*YoutubeEntry `json:"entry"`
	}
}



func (e YoutubeEntry) GetDirector() string {
	for _, contributor := range e.Media.MediaCredit {
		if contributor.Role == "Director" {
			return contributor.Name
		}
	}
	return ""
}

func (e YoutubeEntry) MediaType() string {
	// only supports Movies atm
	return "Movie"
}

func (e *YoutubeEntry) GetOffers() (map[string]float64, error) {
	// gets the top cheapest offers of all types
	offerMap := make(map[string]float64)
	// TODO: WARNING err might as well be _
	// Make this actually throw errors, it's a bit complex compared to the other code
	var err error

	if e != nil && e.Media != nil {
		for _, offer := range e.Media.MediaPrice {
			// check if key exists, if it does and the price
			// is lower add it to the map
			_, keyExists := offerMap[offer.MediaType]
			if keyExists && offerMap[offer.MediaType] > offer.Price {
				offerMap[offer.MediaType] = offer.Price
			} else if keyExists == false {
				offerMap[offer.MediaType] = offer.Price
			}
		}
	}
	return offerMap, err
}

// This needs to return []YoutubeEntry.
func (r YoutubeResponse) Search(title, director string) ([]YoutubeEntry, error) {
	var (
		matchedEntries []YoutubeEntry
		err error
	)

	if (r.Feed.SearchResults.Value == 0) {
		log.Println("youtube: Found no search results")
		return matchedEntries, err
	}

	for _, entry := range r.Feed.Entries {
		log.Printf("youtube: '%s' == '%s'? %v.\n", entry.Title.Value, title, strings.EqualFold(entry.Title.Value, title))
		if strings.EqualFold(entry.Title.Value, title) {
			matchedEntries = append(matchedEntries, *entry)
		}
	}

	return matchedEntries, err
}

func findMovieEntry(keyword string, director string) ([]YoutubeEntry, []byte, error) {
	var (
		apiUrl *url.URL
		title string
		content []byte
		res YoutubeResponse
		resp *http.Response
		movieEntries []YoutubeEntry
		err error
	)

	keyword = normalizeKeyword(keyword)
	apiUrl, err = url.Parse("http://gdata.youtube.com/feeds/api/videos/?alt=json&v=2&category=Movies&key=AI39si5N8WHKROSipJ5SCza8UFOWQlr7YMFnPh3hNnOYULjgnZgY8c7tX1DWVC11tihIAPRbnOM0_2bUKnXgLb2P0OjlTxCi7A&q=" + url.QueryEscape(keyword))
	if err != nil {
		fmt.Println(err)
	}

	resp, err = http.Get(apiUrl.String())
	if err != nil {
		fmt.Println(err)
		return movieEntries, nil, err
	}
	defer resp.Body.Close()

	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	// see if there is function to split bytes like this so cast isn't needed
	err = json.Unmarshal(content, &res)
	if err != nil {
		log.Println("Json unmarshal error:", err)
	}
	//res.Feed.ListTitleInfo()

	title, err = url.QueryUnescape(keyword)

	movieEntries, err = res.Search(title, director)

	if err != nil {
		return movieEntries, nil, err
	}
	return movieEntries, content, err
}

func getPurchaseType(purchaseType string) (string) {
	if strings.EqualFold(purchaseType, "purchase") {
		return "Buy"
	} else if strings.EqualFold(purchaseType, "rent") {
		return "Rent"
	}
	return ""
}

type Youtube struct {}

func (y Youtube) Query(query string) (SearchResult) {
	var(
		// we are failures until we prove we have succeeded
		movies []Media
		queryResult SearchResult = SearchResult{false, movies, nil, "YouTube"} 
		receivedMovies []YoutubeEntry
		rawResponseData []byte
		err error
	)

	receivedMovies, rawResponseData, err = findMovieEntry(query, "") // put this in the if statement
	if err != nil {
		log.Println(err)
		log.Printf("failed to get movie <query: %s> <service: %s>(# results not yet implemented)", query, "youtube")
		return queryResult
	}

	// iterate through received []YoutubeMovie and
	// turn them into a []Media for SearchResult (queryResult here)
	for _, receivedMovie := range receivedMovies {

		offers, err := receivedMovie.GetOffers()
		if err != nil {
			log.Printf("Error getting offers from youtube: %s", err)
		}

		movie := Media{}
		// build type.Movie for queryResult
		movie.Title = receivedMovie.Title.Value

		// TODO: Elegantly set the Link attribute and do this checking within the struct
		if len(receivedMovie.Link) > 0 {
			movie.Link = receivedMovie.Link[0].Href // clean this up to guarantee a link or failure
			queryResult.Success = true
		}

		// initial movie is made, now iterate through offers and make a movie for each offer with it's respective
		// MediaType rent/buy and price
		// TODO: Instead of coercing this map into the movie make a map called offers an attribute of Media
		for purchaseType, price := range offers {

			movie.PurchaseType = getPurchaseType(purchaseType)
			movie.Price = price 
			movies = append(movies, movie)
		}

	}
	queryResult.Media = movies
	queryResult.RawData = rawResponseData

	return queryResult
}
