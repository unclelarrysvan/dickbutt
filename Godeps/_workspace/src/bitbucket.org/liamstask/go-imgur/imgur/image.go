package imgur

import (
	"fmt"
)

const (
	ThumbSmallSquare = "s" // 90x90
	ThumbBigSquare   = "b" // 160x160
	ThumbSmall       = "t" // 160x160
	ThumbMedium      = "m" // 320x320
	ThumbLarge       = "l" // 640x640
	ThumbHuge        = "h" // 1024x1024
)

// ImageService handles communication with the image related
// methods of the Imgur API.
//
// API docs: https://api.imgur.com/endpoints/image
type ImageService struct {
	client *Client
}

type Image struct {
	Id          string // The ID for the image
	Title       string // The title of the image.
	Description string // Description of the image.
	DateTime    int    // Time inserted into the gallery, epoch time
	MimeType    string // Image MIME type.
	Animated    bool   // is the image animated
	Width       int    // The width of the image in pixels
	Height      int    // The height of the image in pixels
	Size        int    // The size of the image in bytes
	Views       int    // The number of image views
	Bandwidth   int    // Bandwidth consumed by the image in bytes
	DeleteHash  string // OPTIONAL, the deletehash, if you're logged in as the image owner
	Section     string // If the image has been categorized by our backend then this will contain the section the image belongs in. (funny, cats, adviceanimals, wtf, etc)
	Link        string // The direct link to the the image
	Nsfw        bool   // Is the link safe for work
	Gifv        string
	Mp4         string
	Webm        string
	Looping     bool
}

// Info retrieves information about an image.
func (s *ImageService) Info(id string) (*Image, error) {
	url := fmt.Sprintf("image/%s", id)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	type imgResponse struct {
		Data *Image
		Result
	}
	response := &imgResponse{}

	_, err = s.client.Do(req, response)
	if err != nil {
		return response.Data, err
	}

	return response.Data, nil
}
