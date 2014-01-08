package api

import (
	"testing"
	"net/url"
	"strings"
	//"fmt"
)

/*
func queryMovies(movieQuery string) ([]SearchResult) {
	// high level test to make sure each
	// wrapper is queryable

	ch := make(chan SearchResult)
	results := []SearchResult{}

	yt := Youtube{}
	itunes := Itunes{}
	netflix := Netflix{}
	hulu := Hulu{}

	go func() { ch <- yt.Query(movieQuery) } ()
	go func() { ch <- itunes.Query(movieQuery) } ()
	go func() { ch <- netflix.Query(movieQuery) } ()
	go func() { ch <- hulu.Query(movieQuery) } ()

	for i:=0; i<5; i++ {
		result := <-ch
		results = append(results, result)
	}

	return results
}
*/

func TestHuluApi(t *testing.T) {
	service := "hulu"
	movieQuery := url.QueryEscape("Arrow")
	hulu := Hulu{}
	searchResult := hulu.Query(movieQuery)


	// test that results are valid
	for _, movie := range searchResult.Media {
		if movie.Title != "Arrow" {
			t.Errorf("invalid result for %s: '%s' != 'Arrow'", service, movie.Title)
		}
	}
	t.Logf("All %s results validated", service)

	if len(searchResult.Media) > 0 {
		t.Logf("Got results back for %s: %v", service, searchResult)
	} else {
		t.Errorf("Failure: %v", searchResult)
	}

}

func TestYoutubeApi(t *testing.T) {
	service := "youtube"
	movieQuery := url.QueryEscape("Fight Club")
	youtube := Youtube{}
	searchResult := youtube.Query(movieQuery)

	// test that results are valid
	for _, movie := range searchResult.Media {
		if movie.Title != "Fight Club" {
			t.Errorf("invalid result for %s: '%s' != 'Fight Club'", service, movie.Title)
		}
	}
	t.Logf("All %s results validated", service)

	// test that offers are valid
	// TODO: Requires each Media type to have an offers map

	if len(searchResult.Media) > 0 {
		t.Logf("Got results back for %s: %v", service, searchResult)
	} else {
		t.Errorf("Failure: %v", searchResult)
	}

}

func TestItunesApi(t *testing.T) {
	service := "itunes"

	movieQuery := url.QueryEscape("Mulan")
	itunes := Itunes{}
	searchResult := itunes.Query(movieQuery)

	// test that results are valid
	for _, movie := range searchResult.Media {
		if movie.Title != "Mulan" {
			t.Errorf("invalid result for %s: '%s' != 'Mulan'", service, movie.Title)
		}
	}
	t.Logf("All %s results validated", service)

	if len(searchResult.Media) > 0 {
		t.Logf("Got results back for %s: %v", service, searchResult)
	} else {
		t.Errorf("Failure: %v", searchResult)
	}

}

func TestNetflixApi(t *testing.T) {
	service := "netflix"

	movieQuery := url.QueryEscape("Thor")
	netflix := Netflix{}
	searchResult := netflix.Query(movieQuery)

	// test that results are valid
	for _, movie := range searchResult.Media {
		if movie.Title != "Thor" {
			t.Errorf("invalid result for %s: '%s' != 'Thor'", service, movie.Title)
		}
	}
	t.Logf("All %s results validated", service)

	if len(searchResult.Media) > 0 {
		t.Logf("Got results back for %s: %v", service, searchResult)
		t.Logf("WARNING: THIS COULD BE A FALSE POSITIVE IF THIS MOVIE IS NO LONGER AVAILABLE FOR 'INSTANT PLAY'")
	} else {
		t.Logf("WARNING: THIS COULD BE A FALSE ERROR. NETFLIX COULD HAVE ADDED THIS PREVIOUSLY 'DVD-ONLY' MOVIE INTO THEIR STREAMING COLLECTION")
		t.Errorf("Failure: %v", searchResult)
	}

}

func TestNetflixApiStreamingOnly(t *testing.T) {
	/*
	Use known (as of 07/11/13) dvd-only movie searches
	and verify that we don't show them as available
	*/
	movieQueryString := "argo"
	movieQuery := url.QueryEscape(movieQueryString)
	netflix := Netflix{}
	searchResult := netflix.Query(movieQuery)
	if !searchResult.Success {
		t.Error("Error getting query from netflix")
	}

	// if we find argo it's an error
	for _, movie := range searchResult.Media {
		t.Log("TESTING: ", movie.Title)
		if strings.EqualFold(movieQueryString, movie.Title) {
			t.Logf("WARNING: THIS COULD BE A FALSE ERROR. NETFLIX COULD HAVE ADDED THIS PREVIOUSLY 'DVD-ONLY' MOVIE INTO THEIR STREAMING COLLECTION")
			//netflixSearchUrl := 
			t.Logf("Check manually at this url to be sure: %v", "http://movies.netflix.com/WiSearch?raw_query=" + movieQueryString + "&ac_category_type=none&ac_rel_posn=-1&ac_abs_posn=-1&v1=argo&search_submit=")
			t.Errorf("Movie query for dvd-only movie '%v' found", movieQueryString)
		} else {
			t.Logf("Dvd-only movie '%v' not found", movieQueryString)
		}
	}

	movieQueryString = "stolen"
	movieQuery = url.QueryEscape(movieQueryString)
	netflix = Netflix{}
	searchResult = netflix.Query(movieQuery)
	if !searchResult.Success {
		t.Error("Error getting query from netflix")
	}
	if len(searchResult.Media) < 1 {
		t.Error("No search results for known instant movie")
	}

	// if we don't find inception it's an error
	for _, movie := range searchResult.Media {
		if !strings.EqualFold(movieQueryString, movie.Title) {
			t.Logf("WARNING: THIS COULD BE A FALSE ERROR. NETFLIX COULD HAVE ADDED THIS PREVIOUSLY 'DVD-ONLY' MOVIE INTO THEIR STREAMING COLLECTION")
			t.Logf("Check manually at this url to be sure: %v", "http://movies.netflix.com/WiSearch?raw_query=" + movieQueryString + "&ac_category_type=none&ac_rel_posn=-1&ac_abs_posn=-1&v1=argo&search_submit=")
			t.Errorf("Movie query for instant movie '%v' not found", movieQueryString)
		} else {
			t.Logf("instant movie '%v' found", movieQueryString)
		}
	}
}

func TestCrackleApi(t *testing.T) {
	service := "crackle"

	movieQuery := "Assassination Games"
	crackle := Crackle{}
	searchResult := crackle.Query(movieQuery)

	if len(searchResult.Media) > 0 {
		t.Logf("Got results back for %s: %v", service, searchResult)
	} else {
		t.Errorf("Crackle failed to get results: %v", searchResult)
	}
}

func TestAmazonApi(t *testing.T) {
	service := "amazon"

	//movieQuery := url.QueryEscape("Fight Club")
	amazon := Amazon{}
	searchResult := amazon.Query("Fight Club")

	// test that results are valid
	for _, movie := range searchResult.Media {
		if movie.Title != "Fight Club" {
			t.Errorf("invalid result for %s: '%s' != 'Fight Club'", service, movie.Title)
		}
	}
	t.Logf("All %s results validated", service)

	if len(searchResult.Media) > 0 {
		t.Logf("Got results back for %s: %v", service, searchResult)
	} else {
		t.Errorf("Failure: %v", searchResult)
	}

}

/*
// TODO: pivotal: https://www.pivotaltracker.com/story/show/50713393
func TestAmazonMediaTypeTv(t *testing.T) {
	amazon := Amazon{}
	tvQuery := "Dr Who"

	tvSearchResult := amazon.Query(tvQuery)
	if len(tvSearchResult.Media) < 1 {
		t.Fatalf("No results found for '%s'", tvQuery)
	}
	for _, show := range tvSearchResult.Media {
		t.Log("Media type is:", show.MediaType)
		if show.MediaType != "TvShow" {
			t.Fatalf("Failure: %v", show)
		} else {
			t.Logf("breaking bad productgroup: %s", show.MediaType)
		}
	}

}
*/
func TestAmazonPriceInfo(t *testing.T) {
	amazon := Amazon{}
	movieQuery := "inception"
	movieResult := amazon.Query(movieQuery)
	movie := movieResult.Media[0]
	priceInfo := getPricingInformation(movie.Link)
	if priceInfo.PurchaseInfo["rent"] != 1.99 || priceInfo.PurchaseInfo["buy"] != 9.99 {
		t.Fatalf("%s prices incorrect: %v", movieQuery, priceInfo.PurchaseInfo)
	} else {
		t.Logf("%s prices correct: %v", movieQuery, priceInfo.PurchaseInfo)
	}



}


/*
tests to re-implement:
func TestNonExistentMovie(t *testing.T) {
func TestYoutubeSearchResultFailure(t *testing.T) {
func TestHuluResourceLink(t *testing.T) {
*/
