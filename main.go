// go-angular project main.go
package main

import (
	"encoding/json"
	"encoding/xml"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const (
	FlickrEndPoint = "https://api.flickr.com/services/rest"
	FlickrQuery    = "flickr.photos.search"
	FlickrKey      = "300d436fa36986e197efe2a62682e05b"
	PathPrefix     = "/pups"
	TopPupsPrefix  = "/top"
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

func ListTopPuppies(w http.ResponseWriter, r *http.Request) {
	page := mux.Vars(r)["page"]
	if page == "" {
		page = "0"
	}

	imageManager := NewImageManager()
	dbError := imageManager.InitDB(false)
	if dbError != nil {
		log.Printf("%q\n", dbError)
		return
	}

	defer imageManager.GetDB().Close()

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}

	puppies := imageManager.GetPuppiesByMostVotes(pageInt)
	count := imageManager.GetPuppiesCount()

	perPage := 10
	pages := count / perPage

	searchResponse := PuppiesResponse{pageInt, pages, perPage, count, puppies}

	response, err := json.Marshal(searchResponse)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

func UpdatePuppy(w http.ResponseWriter, r *http.Request) {
	var v Vote
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		//return badRequest{err}

		println("some error")
	}

	imageManager := NewImageManager()

	dbError := imageManager.InitDB(false)
	if dbError != nil {
		log.Printf("%q\n", dbError)
		return
	}

	defer imageManager.GetDB().Close()
	id, err := strconv.Atoi(v.ID)
	imageManager.UpdateVotes(id, v.VT)

	response, err := json.Marshal(v)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
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

	var tempIDs []string
	imageManager := NewImageManager()
	for _, ph := range flickrPhotos {
		img := imageManager.NewImage(ph)
		imageManager.Save(img)
		tempIDs = append(tempIDs, img.ID)
	}

	dbError := imageManager.InitDB(false)
	if dbError != nil {
		log.Printf("%q\n", dbError)
		return
	}
	defer imageManager.GetDB().Close()

	all := imageManager.All()

	dbPuppies := imageManager.FindOldPuppies(tempIDs)

	var newPuppies []*Image

	if len(dbPuppies) == 0 {
		imageManager.InsertPuppies(all)
	} else {
		for _, puppy := range dbPuppies {
			id := puppy.ID
			for _, allP := range all {
				//allPID, _ := strconv.Atoi(allP.ID)
				if allP.ID == id {
					allP.DownVotes = puppy.DownVotes
					allP.UpVotes = puppy.UpVotes
				} else {
					exists := true
					var existingPuppy *Image
					for _, np := range newPuppies {
						//nPID, _ := strconv.Atoi(np.ID)
						if np.ID == allP.ID {
							exists = false
							existingPuppy = allP
							break
						}
					}
					if exists == false {
						newPuppies = append(newPuppies, existingPuppy)
					}
				}
			}
		}

		imageManager.InsertPuppies(newPuppies)
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
	imageManager := NewImageManager()
	dbError := imageManager.InitDB(false)

	defer imageManager.GetDB().Close()

	if dbError != nil {
		log.Printf("%q\n", dbError)
		return
	} else {
		imageManager.CreateTables()
	}

	r := mux.NewRouter().StrictSlash(false)

	pups := r.Path(PathPrefix).Subrouter()
	pups.Methods("GET").HandlerFunc(ListPuppies)

	pupsPerPage := r.Path(PathPrefix + "/{page}").Subrouter()
	pupsPerPage.Methods("GET").HandlerFunc(ListPuppies)

	topPups := r.Path(TopPupsPrefix).Subrouter()
	topPups.Methods("GET").HandlerFunc(ListTopPuppies)

	topPupsPerPage := r.Path(TopPupsPrefix + "/{page}").Subrouter()
	topPupsPerPage.Methods("GET").HandlerFunc(ListTopPuppies)

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
