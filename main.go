// go-angular project main.go
package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"io/ioutil" 
	"net/url"
	"encoding/json"
	"encoding/xml"
)

const (
	FlickrEndPoint = "https://api.flickr.com/services/rest";
	FlickrQuery = "flickr.photos.search";
	FlickrKey = "300d436fa36986e197efe2a62682e05b";
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func FetchPuppies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {	
	tags := ps.ByName("tags")
	page := ps.ByName("page")
	
	baseUrl, err := url.Parse(FlickrEndPoint)
	if err != nil {
		log.Fatal(err)
	}

	params := url.Values{}
	params.Add("method", FlickrQuery)
	params.Add("api_key", FlickrKey)
	params.Add("tags", tags)
	params.Add("per_page", "10")
	params.Add("page", page)
	params.Add("safe_search", "2")
	params.Add("sort", "date-posted-desc")
		
	baseUrl.RawQuery = params.Encode()
	
	resp, err := http.Get(baseUrl.String())
	if err != nil {
		// handle error, send proper error response
		log.Println(err)
	}	
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error, send proper error response
		log.Println(err)
	}

	flickrResponse := struct {
		Stat string `xml:"stat,attr"`
		Err flickrError `xml:"err"`
		Photos SearchResponse `xml:"photos"`
	}{}	
	
	xml.Unmarshal([]byte(body), &flickrResponse)

	//stat := flickrResult.Stat
	//if stat is "ok"	
	if flickrResponse.Stat != "ok" {
		println(flickrResponse.Err.Msg)
		//return error message
	}	
	flickrPhotos := flickrResponse.Photos.Photos
	for i, ph := range flickrPhotos {
		flickrPhotos[i].Thumbnail_T = ph.URL(SizeThumbnail)
		flickrPhotos[i].Large_T = ph.URL(SizeLarge)
	}

	response, err := json.Marshal(flickrPhotos)
	if err != nil {
  		http.Error(w, err.Error(), http.StatusInternalServerError)
	  	  return
	}

	//return json response	
	fmt.Fprint(w, string(response))
}

func main() {
	fmt.Printf("hello, world\n")
	router := httprouter.New()
	
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	router.GET("/fetch/:page/:tags", FetchPuppies)

	log.Fatal(http.ListenAndServe(":8080", router))
}
