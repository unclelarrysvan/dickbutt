package imgur

import (
	"fmt"
	// "log"
)

// GalleryService handles communication with the gallery related
// methods of the Imgur API.
//
// API docs: https://api.imgur.com/endpoints/gallery
type GalleryService struct {
	client *Client
}

// Many of the Gallery endpoints return mixtures of
// GalleryImage and GalleryAlbum, and these two structs
// share many elements, so they are combined into a single struct.
// Test the IsAlbum field to determine which type an instance is.
type GalleryImageAlbum struct {
	// Common fields
	ID           string `json:"id,omitempty"`
	Title        string `json:"title,omitempty"`
	Description  string `json:"description,omitempty"`
	DateTime     int    `json:"datetime,omitempty"`
	Views        int    `json:"views,omitempty"`
	Vote         string `json:"vote,omitempty"`
	Section      string `json:"section,omitempty"`
	AccountUrl   string `json:"account_url,omitempty"`
	Ups          int    `json:"ups,omitempty"`
	Downs        int    `json:"downs,omitempty"`
	Score        int    `json:"score,omitempty"`
	Link         string `json:"link,omitempty"`
	IsAlbum      bool   `json:"is_album,omitempty"`
	Nsfw         bool   `json:"nsfw,omitempty"`
	CommentCount int    `json:"comment_count,omitempty"`

	// Image only fields
	Bandwidth  int    `json:"bandwidth,omitempty"`
	DeleteHash string `json:"deletehash,omitempty"`
	Animated   bool   `json:"animated,omitempty"`
	MimeType   string `json:"type,omitempty"`
	Width      int    `json:"width,omitempty"`
	Height     int    `json:"height,omitempty"`
	Size       int    `json:"size,omitempty"`
	Gifv       string `json:"gifv,omitempty"`
	Webm       string `json:"webm,omitempty"`

	// Album only fields
	Cover       string  `json:"cover,omitempty"`
	CoverWidth  int     `json:"cover_width,omitempty"`
	CoverHeight int     `json:"cover_height,omitempty"`
	Privacy     string  `json:"privacy,omitempty"`
	Layout      string  `json:"layout,omitempty"`
	ImagesCount int     `json:"images_count,omitempty"`
	Images      []Image `json:"images,omitempty"`
}

type galleryImageAlbumResult struct {
	Data    []GalleryImageAlbum
	Status  int
	Success bool
}

// Return a gallery of the specified section, sort, window page, etc.
func (s *GalleryService) gallery(route, sort, window, paramStr string, page int) ([]GalleryImageAlbum, error) {
	if page < 0 {
		page = 0
	}

	response := &galleryImageAlbumResult{}

	if route == "" {
		return response.Data, fmt.Errorf("route must be provided")
	}

	url := "gallery/" + route

	if sort != "" {
		url = url + "/" + sort
		if window != "" {
			url = url + "/" + window
		}
		// Avoiding pulling in all of strconv
		url = url + fmt.Sprintf("/%d", page)
	} else {
		return response.Data, fmt.Errorf("sort must be provided to gallery() method")
		if window != "" {
		}
	}

	if paramStr != "" {
		url = url + paramStr
	}

	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return response.Data, err
	}

	_, err = s.client.Do(req, response)
	if err != nil {
		return response.Data, err
	}

	return response.Data, nil
}

// Returns the main gallery, as if the user had simply navigated to imgur.com
func (s *GalleryService) Main(section, sort, window string, page int) ([]GalleryImageAlbum, error) {

	if section == "" {
		section = "hot"
	}

	if sort == "" {
		sort = "viral"
	}

	if window == "" {
		window = "day"
	}

	return s.gallery(section, sort, window, "", page)
}

// Returns a subreddit gallery (requires clientID & clientSecret)
func (s *GalleryService) Subreddit(subreddit, sort, window string, page int) ([]GalleryImageAlbum, error) {

	// no default for subreddit. Currently let the user fail on their own if it isn't provided or is invalid

	if sort == "" {
		sort = "time"
	}

	if window == "" {
		window = "week"
	}

	route := fmt.Sprintf("r/%s", subreddit)

	return s.gallery(route, sort, window, "", page)
}

// Returns the memes gallery (requires clientID & clientSecret)
func (s *GalleryService) Memes(sort, window string, page int) ([]GalleryImageAlbum, error) {

	if sort == "" {
		sort = "viral"
	}

	if window == "" {
		window = "week"
	}

	return s.gallery("g/memes", sort, window, "", page)
}

// Search searches the gallery with a given query string.
func (s *GalleryService) Search(q string, sort string, page int) ([]GalleryImageAlbum, error) {
	// optional    time | viral - defaults to time
	if sort == "" {
		sort = "time"
	}

	searchQuery := fmt.Sprintf("?q=%s", q)

	return s.gallery("search", sort, "", searchQuery, page)
}

// Random returns a random set of gallery images.
func (s *GalleryService) Random(page int) ([]GalleryImageAlbum, error) {
	// optional    integer - the data paging number
	if page < 0 {
		page = 0
	}

	response := &galleryImageAlbumResult{}

	url := fmt.Sprintf("gallery/random/random/%d", page)

	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return response.Data, err
	}

	_, err = s.client.Do(req, response)
	if err != nil {
		return response.Data, err
	}

	return response.Data, nil
}
