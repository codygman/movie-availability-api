package api

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"encoding/xml"
	"log"
	"bytes"
	"strings"
	"fmt"
)


// link structure
// "http://www.jdoqocy.com/click-5907747-10947797"
type huluShow struct {
	XMLName 	xml.Name `xml:"video"`
	Name 	 	string `xml:"show>name"`
	EpisodeTitle 	string `xml:"title"`
	UrlSlug 	string `xml:"canonical-name"`
	PlusOnly 	bool `xml:"is-subscriber-only"`
	HasHD 		bool `xml:"has-hd""`
}

type Hulu struct {
	XMLName 	xml.Name `xml:"results"`
	NumResults 	int64 `xml:"count"`
	Results 	[]huluShow `xml:"videos>video"`
}

func genAffiliateLink(slug string) string {
	showLink := fmt.Sprintf("http://www.hulu.com/%s", slug)

	params := url.Values{}
	params.Add("url", showLink)
	params.Add("cjsku", slug)

	baseUrl := "" // http://youraffiliatelink.com/link?
	link := baseUrl + params.Encode()

	return link
}

func genResourceLink(keyword string) string {
	baseUrl := "http://m.hulu.com/search?"
	params := url.Values{}

	// must always be hulu
	params.Add("dp_identifier", "hulu")
	params.Add("query", url.QueryEscape(keyword))
	params.Add("items_per_page", "10")

	queryUrl := baseUrl + params.Encode()

	return queryUrl
}
// TODO: Return standard format of []huluShow, err
func (h Hulu) Filter(name string) []huluShow {
	shows := []huluShow{}

	for _, show := range h.Results {

		if show.Name == "" {
			log.Println("hulu: Received blank result")
			break
		}
		alreadyExists := false
		for _, s := range shows {
			if s == show {
				alreadyExists = true
			}
		}

		namesEqual := strings.EqualFold(show.Name, name)
		log.Printf("hulu: '%s' == '%s'? %v.\n", show.Name, name, namesEqual)
		if namesEqual && !alreadyExists {
			shows = append(shows, show)
		}
	}
	return shows
}

// TODO: Need to make a "Show" type. What's a good name for a parent type
// for shows and movies? Media? Title, probably.
func queryHulu(keyword string) *SearchResult {
	var (
		receivedShows []huluShow
		shows []Media
		searchResult = SearchResult{false, shows, nil, "Hulu"}
	)
	
	queryUrl := genResourceLink(keyword)
	resp, err := http.Get(queryUrl)
	if err != nil {
		log.Println(err)
		return &searchResult

	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return &searchResult
	}
	searchResult.RawData = content

	// decode xml into Hulu structure
	hulu := new(Hulu)
	if err := xml.NewDecoder(bytes.NewReader(content)).Decode(&hulu); err != nil {
		log.Println(err)
		return &searchResult
	}

	// filter
	// TODO URGENT: temporary until a parent structure for Movie and Show is settled upon
	receivedShows = hulu.Filter(keyword)

	for _, receivedShow := range receivedShows {
		if receivedShow.Name == "" {
			// no show name found, continue on to next iteration
			continue
		} else {
			// build a show, and add it to shows []Media

			show := Media{}
			show.Title = receivedShow.Name
			show.Link = receivedShow.UrlSlug // genAffiliateLink(receivedShow.UrlSlug)
			show.HasHD = receivedShow.HasHD

			shows = append(shows, show)
		}
	}
	searchResult.Media = shows
	searchResult.Success = len(shows) > 0

	return &searchResult
}

func  (h Hulu) Query(query string) SearchResult {
	// query hulu
	sr := queryHulu(query)
	// return search result
	return *sr
}
