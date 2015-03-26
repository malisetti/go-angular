// go-angular project main.go
package main

/*
"database/sql"
"fmt"
*/
import (
	"encoding/json"
	"encoding/xml"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	FlickrEndPoint = "https://api.flickr.com/services/rest"
	FlickrQuery    = "flickr.photos.search"
	FlickrKey      = "300d436fa36986e197efe2a62682e05b"
	DatabaseName   = "puppies.db"
	PathPrefix     = "/pups"
)

// badRequest is handled by setting the status code in the reply to StatusBadRequest.
type badRequest struct{ error }

// notFound is handled by setting the status code in the reply to StatusNotFound.
type notFound struct{ error }

// errorHandler wraps a function returning an error by handling the error and returning a http.Handler.
// If the error is of the one of the types defined above, it is handled as described for every type.
// If the error is of another type, it is considered as an internal error and its message is logged.
func errorHandler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			return
		}
		switch err.(type) {
		case badRequest:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case notFound:
			http.Error(w, "task not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(w, "oops", http.StatusInternalServerError)
		}
	}
}

func UpdatePuppy(w http.ResponseWriter, r *http.Request) {
	var v Vote
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		//return badRequest{err}

		println("some error");
	}

	imageManager := NewImageManager()
	println(v.ID)
	image, exists := imageManager.Find(v.ID)

	if exists == true{
		println(image.ID)
	}else{
		println("does not exist")
	}

}

func ListPuppies(w http.ResponseWriter, r *http.Request) {
	page := mux.Vars(r)["page"]
	if page == "" {
		page = "1"
	}
	//tags := mux.Vars(r)["tags"]
	tags := "puppies,dogs,cute"

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
		Stat   string         `xml:"stat,attr"`
		Err    flickrError    `xml:"err"`
		Photos SearchResponse `xml:"photos"`
	}{}

	xml.Unmarshal([]byte(body), &flickrResponse)

	//stat := flickrResult.Stat
	//if stat is "ok"
	if flickrResponse.Stat != "ok" {
		println(flickrResponse.Err.Msg)
		//return error message
	}

	searchResponse := flickrResponse.Photos
	flickrPhotos := searchResponse.Photos

	imageManager := NewImageManager()
	for _, ph := range flickrPhotos {
		imageManager.Save(imageManager.NewImage(ph))
	}

	puppiesResponse := imageManager.GetPuppiesResponse(&searchResponse)
	response, err := json.Marshal(puppiesResponse)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func main() {
	r := mux.NewRouter().StrictSlash(false)

	pups := r.Path(PathPrefix).Subrouter()
	pups.Methods("GET").HandlerFunc(ListPuppies)

	pupsPerPage := r.Path(PathPrefix + "/{page}").Subrouter()
	pupsPerPage.Methods("GET").HandlerFunc(ListPuppies)

	pupsUpdate := r.Path(PathPrefix).Subrouter()
	pupsUpdate.Methods("PUT").HandlerFunc(UpdatePuppy)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
