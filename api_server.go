package main
import (
	"fmt"
	"net/http"
	"api"
	"log"
	"encoding/json"
	"runtime"

	// non go-source packages
	"github.com/simonz05/godis/redis"
)


func queryAll(w http.ResponseWriter, r *http.Request) {

	var (
		inRedis bool
		receivedMovie []byte
		err error
		// maps are special and must be initialized with make
		resultsMap map[string][]api.Media = make(map[string][]api.Media)
		//results []api.Media
		redisMovie redis.Elem
	)

	//allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	w.Header().Set("Content-Type", "application/json")

	//movieQuery := keyword
	movieQuery := r.FormValue("keyword")

	c := redis.New("", 0, "")

	// check redis for keyword, if exists... return!!
	inRedis, err = c.Exists(movieQuery)
	if err != nil {
		// TODO: Change to logging!
		//fmt.Println(err)
	}

	if inRedis {
		redisMovie, err = c.Get(movieQuery)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Got %s from cache", movieQuery)
		receivedMovie = redisMovie.Bytes()
	} else if inRedis == false {

		wrappers := []api.ApiWrapper{
			new(api.Youtube), 
			new(api.Itunes), 
			new(api.Hulu),
			new(api.Crackle),
		}
		ch := make(chan api.SearchResult, len(wrappers))

		for _, wrapper := range wrappers {
			go func(wrapper api.ApiWrapper) { ch <- wrapper.Query(movieQuery) }(wrapper)
		}

		for _ = range wrappers {
			result := <-ch
			if result.Success {
				resultsMap[result.Service] = result.Media
			}
		}

		// the code below writes the json blob to the writer responsible for
		// writing code directly to the browser
		// look at using an io.Multiwriter here to avoid casting
		//if err := json.NewEncoder(w).Encode(&resultsMap); err != nil {

		if receivedMovie, err = json.Marshal(resultsMap); err != nil {
			log.Println(err)
		}

		// set in redis
		// TODO: Update to only set key in redis if at least one result is found
		if err = c.Set(movieQuery, receivedMovie); err != nil {
			log.Println("failed to set redis key for <query: %s>", movieQuery)
		}

	}
	fmt.Fprint(w, string(receivedMovie))
}


func showRawData(w http.ResponseWriter, r *http.Request) {
	// SECURITY WARNING:
	// Raw request data sometimes includes sensitive data
	// such as api/password/secret keys. It may be a good idea to
	// either scrub sensitive data OR only allow this view to be
	// shown to authorized users

	keyword := r.FormValue("keyword")
	if keyword == "" {
		fmt.Fprintf(w, string("Request malformed"))
		return
	}
	provider := r.FormValue("provider")

	var wrapper api.ApiWrapper

	switch provider {
	case "amazon":
		w.Header().Set("Content-Type", "application/xml")
		// wrapper = new(api.Amazon)
	case "hulu":
		w.Header().Set("Content-Type", "application/xml")
		wrapper = new(api.Hulu)
	case "itunes":
		w.Header().Set("Content-Type", "application/json")
		wrapper = new(api.Itunes)
	case "netflix":
		w.Header().Set("Content-Type", "application/json")
		// wrapper = new(api.Netflix)
	case "youtube":
		w.Header().Set("Content-Type", "application/json")
		wrapper = new(api.Youtube)
	}

	if wrapper == nil {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "That provider does not exist")
		return
	}

	searchResult := wrapper.Query(keyword)
	if searchResult.RawData != nil {
		fmt.Fprint(w, string(searchResult.RawData))
	} else {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "No results or not implemented in provider wrapper '%s'", provider)
	}
	fmt.Println("done")

}

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	http.HandleFunc("/query/all", queryAll)
	http.HandleFunc("/query/provider_response", showRawData)

	http.ListenAndServe("localhost:9999", nil)
}

