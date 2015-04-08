package main

import (
	"bitbucket.org/liamstask/go-imgur/imgur"

	"math/rand"
)

const (
	//clientID     = os.Getenv("IMGUR_CLIENT_ID")
	//clientSecret = os.Getenv("IMGUR_SECRET_ID")
	clientID     = "e29200a15c6a770"
	clientSecret = "b7a486c0a691fd8a34d026045d67a5fa6e4e18f2"
)

var client = imgur.NewClient(nil, clientID, clientSecret)

func ImgurSearcher(image string) (url string) {
	results, err := client.Gallery.Search(image, "top", 0)

	if err != nil {
		url = "http://s.imgur.com/images/OverCapacity_700.png"
		return
	}

	if len(results) <= 0 {
		url = "http://s.imgur.com/images/OverCapacity_700.png"
		return
	}

	image_index := rand.Intn(len(results))
	images := results[image_index]
	if images.IsAlbum {
		if len(images.Images) > 0 {
			url = images.Images[0].Link
		} else {
			url = getFirstImage(results)
		}
	} else {
		url = images.Link
	}
	return
}

func getFirstImage(images []imgur.GalleryImageAlbum) (url string) {
	for _, image := range images {
		if !image.IsAlbum {
			url = image.Link
			return
		}
	}
	url = "http://s.imgur.com/images/OverCapacity_700.png"
	return
}
