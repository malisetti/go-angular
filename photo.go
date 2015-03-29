package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strconv"
	"strings"
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
	DatabaseName    = "puppies.sqlite"
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
	UpVotes   int    `json:"upvotes"`
	DownVotes int    `json:"downvotes"`
}

type PuppiesResponse struct {
	Page    int      `json:"page"`
	Pages   int      `json:"pages"`
	PerPage int      `json:"perpage"`
	Total   int      `json:"total"`
	Images  []*Image `json:"images"`
}

type ImageManager struct {
	images []*Image
	db     *sql.DB
}

type Vote struct {
	ID string `json:"id"`
	VT bool   `json:"vt"`
}

func NewImageManager() *ImageManager {
	return &ImageManager{}
}

func (m *ImageManager) GetPuppiesResponse(searchResponse *SearchResponse) *PuppiesResponse {
	page, err := strconv.Atoi(searchResponse.Page)
	pages, err := strconv.Atoi(searchResponse.Pages)
	perPage, err := strconv.Atoi(searchResponse.PerPage)
	total, err := strconv.Atoi(searchResponse.Total)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &PuppiesResponse{page, pages, perPage, total, m.images}
}

func (m *ImageManager) NewImage(photo Photo) *Image {
	return &Image{photo.ID, photo.Title, photo.URL(SizeThumbnail), photo.URL(SizeLarge), 0, 0}
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

func (m *ImageManager) Find(ID string) (*Image, bool) {
	for _, im := range m.images {
		if im.ID == ID {
			return im, true
		}
	}

	return nil, false
}

func (m *ImageManager) Update(image *Image, upOrDown bool) (int, int) {
	if upOrDown == true {
		image.UpVotes++
	} else {
		image.DownVotes--
	}

	for _, im := range m.images {
		if im.ID == image.ID {
			im.UpVotes = image.UpVotes
			im.DownVotes = image.DownVotes
		}
	}

	return image.UpVotes, image.DownVotes
}

func (m *ImageManager) UpdateVotes(puppy_id int, up_vote bool) {
	sqlStmt := "update votes set "
	if up_vote == true {
		sqlStmt += " up_votes = up_votes + 1"
	} else {
		sqlStmt += " down_votes = down_votes + 1"
	}

	sqlStmt += " where puppy_id = ?"

	stmt, err := m.db.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	res, _ := stmt.Exec(puppy_id)

	affect, _ := res.RowsAffected()

	fmt.Println(affect)

	return
}

func (m *ImageManager) GetPuppiesCount() int {
	query := "select count(id) from votes"

	rows, err := m.db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		rows.Scan(&count)
	}

	return count
}

func (m *ImageManager) GetPuppiesByMostVotes(pageId int) []*Image {
	perPage := 10
	if pageId != 0{
		pageId-- 
	}
	start := perPage * pageId
	query := "select * from votes order by up_votes desc limit ?,?"

	stmt, err := m.db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()
	rows, err := stmt.Query(start, perPage)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var rs []*Image

	for rows.Next() {
		var dbImage Image
		var id int
		rows.Scan(&id, &dbImage.ID, &dbImage.Title, &dbImage.Thumbnail, &dbImage.Large, &dbImage.UpVotes, &dbImage.DownVotes)
		rs = append(rs, &dbImage)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return rs

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

func (m *ImageManager) InitDB(removeDb bool) error {
	if removeDb == true {
		os.Remove("./" + DatabaseName)
	}

	db, err := sql.Open("sqlite3", "./"+DatabaseName)
	if err != nil {
		log.Fatal(err)
		return err
	}

	m.db = db
	return nil
}

func (m *ImageManager) GetDB() *sql.DB {
	return m.db
}

func (m *ImageManager) CreateTables() {
	createSqlStmt := `
	create table if not exists votes (id integer not null primary key, puppy_id integer unique, title string, thumbnail string, large string, up_votes integer, down_votes integer);
	delete from votes;
	`
	_, err := m.db.Exec(createSqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, createSqlStmt)
	}
}

func (m *ImageManager) InsertPuppies(images []*Image) {
	tx, err := m.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into votes(puppy_id, title, thumbnail, large, up_votes, down_votes) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	for _, im := range images {
		_, err = stmt.Exec(im.ID, im.Title, im.Thumbnail, im.Large, im.UpVotes, im.DownVotes)
		if err != nil {
			log.Fatal(err)
		}
	}

	tx.Commit()
}

func (m *ImageManager) FindOldPuppies(ids []string) []*Image {

	//sqlStmt = "select * from votes where puppy_id in (?" + strings.Repeat(",?", len(ids)-1) + ")"

	query := fmt.Sprintf("select * from votes where puppy_id in (%s)",
		strings.Join(strings.Split(strings.Repeat("?", len(ids)), ""), ","))

	stmt, err := m.db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}

	var params []interface{}
	for _, id := range ids {
		params = append(params, id)
	}
	defer stmt.Close()
	rows, err := stmt.Query(params...)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var rs []*Image

	for rows.Next() {
		var dbImage Image
		var id int
		rows.Scan(&id, &dbImage.ID, &dbImage.Title, &dbImage.Thumbnail, &dbImage.Large, &dbImage.UpVotes, &dbImage.DownVotes)
		rs = append(rs, &dbImage)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return rs
}
