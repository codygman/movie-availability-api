package api

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"encoding/json"
	"log"
	"bytes"
	"strings"
	"fmt"
)

type CrackleItem struct {
	Title string `json:"Title"`
	RootChannel string `json:"RootChannel"`
}

type Crackle struct {
	Count int `json:"Count"`
	Items []CrackleItem `json:"Items"`
}

// generates link since XItemId is unreliable
// This is where things get tricky, they use a custom url prettifyer it seems
// (silver lining: the resulting generated url isn't picky and seems to do an
// "I'm feeling lucky type search)
// here are some rules I've observed
// replace " " with "_"
// replace "'" with ""

func (item CrackleItem) Link() string {
	// url
	urlString := fmt.Sprintf("http://www.crackle.com/c/%s", item.Title)

	replacer := strings.NewReplacer( " ", "_", "'", "")
	urlString = replacer.Replace(urlString)
	
	url, err := url.Parse(urlString)
	if err != nil {
		log.Fatalf("Failed to parse url from title '%s'", item.Title)
	}
	return url.String()
}

// Filters all shows in c.Items that match given title
func (c Crackle) Filter(title string) []CrackleItem {
	items := []CrackleItem{}

	for _, item := range c.Items {
		log.Println("crackle: ", item)
		if item.Title == "" {
			log.Printf("crackle: Received blank result filtering by title '%s'", title)
			break
		}

		titlesEqual := strings.EqualFold(item.Title, title)
		log.Printf("'%s' == '%s'? %v.", item.Title, title, titlesEqual)
		if titlesEqual {
			items = append(items, item)
		}
	}
	log.Printf("crackle: Returning %d items from filtering %s", len(items), title)
	return items
}

// Queries crackle's api when given a keyword:
// http://api.crackle.com/Service.svc/search/all/:UrlEncodedQueryHere/US?format=json
// coerce them into a SearchResult full of []Media
func  (c Crackle) Query(keyword string) SearchResult {
	var (
		receivedItems []CrackleItem
		shows []Media
		searchResult = SearchResult{false, shows, nil, "Crackle"}
	)
	
	// the escaped keyword doesn't work, lol
	// queryUrl := fmt.Sprintf("http://api.crackle.com/Service.svc/search/all/%s/US?format=json", url.QueryEscape(keyword))
	log.Printf("keyword for crackle: '%s'", keyword)
	queryUrl := fmt.Sprintf("http://api.crackle.com/Service.svc/search/all/%s/US?format=json", keyword)
	resp, err := http.Get(queryUrl)
	if err != nil {
		log.Println(err)
		return searchResult

	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return searchResult
	}
	searchResult.RawData = content

	// decode json into Crackle structure
	crackle := new(Crackle)
	if err := json.NewDecoder(bytes.NewReader(content)).Decode(&crackle); err != nil {
		log.Printf("Error decoding json from url: %s; error: %s", queryUrl, err)
		return searchResult
	}

	log.Printf("There are %d results found at %s", len(crackle.Items), queryUrl)

	// filter
	// TODO URGENT: temporary until a parent structure for Movie and Show is settled upon
	receivedItems = crackle.Filter(keyword)

	for _, receivedShow := range receivedItems {
		if receivedShow.Title == "" {
			// no show name found, continue on to next iteration
			continue
		} else {
			// build a show, and add it to shows []Media

			show := Media{}
			show.Title = receivedShow.Title
			// as of this writing 04/17/2013 crackle has these attributes but isn't using them
			//show.EpisodeTitle = receivedShow.Episode
			show.PurchaseType = "Free"
			show.Link = receivedShow.Link()

			shows = append(shows, show)
		}
	}
	searchResult.Media = shows
	searchResult.Success = len(shows) > 0

	return searchResult
}
