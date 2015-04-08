package main

import (
	"bitbucket.org/liamstask/go-imgur/imgur"

	"math/rand"
	"os"
)

var clientID = os.Getenv("IMGUR_CLIENT_ID")
var clientSecret = os.Getenv("IMGUR_SECRET_ID")

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
