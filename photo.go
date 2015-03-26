package main

import (
	"fmt"
)

// Image sizes supported by Flickr.  See
// http://www.flickr.com/services/api/misc.urls.html for more information.
const (
	SizeSmallSquare = "s"
	SizeThumbnail   = "t"
	SizeSmall       = "m"
	SizeMedium500   = "-"
	SizeMedium640   = "z"
	SizeLarge       = "b"
	SizeOriginal    = "o"
)

// Response for photo search requests.
type SearchResponse struct {
	Page    string  `xml:"page,attr"`
	Pages   string  `xml:"pages,attr"`
	PerPage string  `xml:"perpage,attr"`
	Total   string  `xml:"total,attr"`
	Photos  []Photo `xml:"photo"`
}

// Represents a Flickr photo.
type Photo struct {
	ID          string `xml:"id,attr"`
	Owner       string `xml:"owner,attr"`
	Secret      string `xml:"secret,attr"`
	Server      string `xml:"server,attr"`
	Farm        string `xml:"farm,attr"`
	Title       string `xml:"title,attr"`
	IsPublic    string `xml:"ispublic,attr"`
	IsFriend    string `xml:"isfriend,attr"`
	IsFamily    string `xml:"isfamily,attr"`
	Thumbnail_T string `xml:"thumbnail_t,attr"`
	Large_T     string `xml:"large_t,attr"`
}

type flickrError struct {
	Code string `xml:"code,attr"`
	Msg  string `xml:"msg,attr"`
}

type Image struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	Large     string `json:"large"`
	Votes     int    `json:"votes"`
}

type PuppiesResponse struct {
	Page    string   `json:"page"`
	Pages   string   `json:"pages"`
	PerPage string   `json:"perpage"`
	Total   string   `json:"total"`
	Images  []*Image `json:"images"`
}

type ImageManager struct {
	images []*Image
}

func NewImageManager() *ImageManager {
	return &ImageManager{}
}

func (m *ImageManager) GetPuppiesResponse(searchResponse *SearchResponse) *PuppiesResponse {
	return &PuppiesResponse{searchResponse.Page, searchResponse.Pages, searchResponse.PerPage, searchResponse.Total, m.images}
}

func (m *ImageManager) NewImage(photo Photo) *Image {
	return &Image{photo.ID, photo.Title, photo.URL(SizeThumbnail), photo.URL(SizeLarge), 0}
}

func (m *ImageManager) Save(image *Image) error {
	for _, im := range m.images {
		if im.ID == image.ID {
			return nil
		}
	}

	m.images = append(m.images, cloneImage(image))
	return nil
}

func (m *ImageManager) Update(image *Image, upOrDown bool) int {
	if upOrDown == true {
		image.Votes++
	} else {
		image.Votes--
	}

	for _, im := range m.images {
		if im.ID == image.ID {
			im.Votes = image.Votes
		}
	}

	return image.Votes
}

// All returns the list of all the Tasks in the TaskManager.
func (m *ImageManager) All() []*Image {
	return m.images
}

func cloneImage(i *Image) *Image {
	c := *i
	return &c
}

// Returns the URL to this photo in the specified size.
func (p *Photo) URL(size string) string {
	if size == "-" {
		return fmt.Sprintf("http://farm%s.static.flickr.com/%s/%s_%s.jpg",
			p.Farm, p.Server, p.ID, p.Secret)
	}
	return fmt.Sprintf("http://farm%s.static.flickr.com/%s/%s_%s_%s.jpg",
		p.Farm, p.Server, p.ID, p.Secret, size)
}
