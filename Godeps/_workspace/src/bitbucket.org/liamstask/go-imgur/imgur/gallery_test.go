package imgur

import (
	"fmt"
	"net/http"
	"testing"
)

const (
	gallerySearchResponse = `{"status":200, "success":true, "data": [{"id":"Wu0zw","title":"Sad Keanu is Now an Action Figure!","datetime":1390524925,"views":770,"account_url":"literallyannperkins","ups":27,"downs":2,"score":25,"link":"http://imgur.com/a/Wu0zw","is_album":true,"cover":"LAobnLK","cover_width":625,"cover_height":944,"privacy":"public","layout":"blog"}]}`

	galleryRandomResponse = `{"status":200, "success":true, "data": [{"id":"1YUpmrH","title":"My cousin had a costume wedding, my grandma went as Princess Leia.","datetime":1390086751,"views":615140,"section":"funny","account_url":"bpony","ups":2919,"downs":48,"score":3532,"link":"http://i.imgur.com/1YUpmrH.jpg","bandwidth":120289396720,"type":"image/jpeg","width":639,"height":852,"size":195548}]}`

	galleryMainResponse = `{"data":[{"id":"Lh6MuPp","title":"God dammit Trebek!","description":null,"datetime":1404776619,"type":"image\/jpeg","animated":false,"width":941,"height":2791,"size":281110,"views":144459,"bandwidth":40608869490,"vote":null,"favorite":false,"nsfw":false,"section":"funny","account_url":null,"link":"http:\/\/i.imgur.com\/Lh6MuPp.jpg","account_id":null,"ups":2117,"downs":27,"score":2163,"is_album":false}],"success":true,"status":200}`

	gallerySubredditResponse = `{"data":[{"id":"4C06G","title":"Dad teaches son how to ride bike","description":null,"datetime":1404779006,"type":"image\/gif","animated":true,"width":300,"height":229,"size":2047447,"views":14031,"bandwidth":28727728857,"favorite":false,"nsfw":false,"section":"Unexpected","link":"http:\/\/i.imgur.com\/4C06G.gif","reddit_comments":"\/r\/Unexpected\/comments\/2a3km9\/dad_teaches_son_how_to_ride_bike\/","account_id":null,"ups":2,"downs":0,"score":13967}],"success":true,"status":200}`

	galleryMemesResponse = `{"data":[{"id":"zHQ2rzI","title":"Of course, but maybe...","description":null,"datetime":1404778199,"type":"image\/png","animated":false,"width":610,"height":679,"size":395612,"views":22,"bandwidth":8703464,"vote":null,"favorite":false,"nsfw":false,"section":null,"account_url":"kJerAFK","link":"http:\/\/i.imgur.com\/zHQ2rzI.png","subtype":"Of course Louie","account_id":3765705,"ups":3,"downs":0,"score":3,"is_album":false}],"success":true,"status":200}`
)

func TestGallerySearch(t *testing.T) {
	imgurTestSetup()
	defer imgurTestTeardown()

	mux.HandleFunc("/gallery/search/time/0", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, gallerySearchResponse)
	})

	albumImgs, err := client.Gallery.Search("searchterm", "time", 0)
	if err != nil {
		t.Errorf("Gallery.Search returned error: %v", err)
	}

	if len(albumImgs) == 0 {
		t.Errorf("Gallery.Search returned %+v, want %+v", len(albumImgs), 0)
	}

	want := "Wu0zw"
	if albumImgs[0].ID != want {
		t.Errorf("Gallery.Search returned %+v, want %+v", albumImgs[0].ID, want)
	}
}

func TestGalleryRandom(t *testing.T) {
	imgurTestSetup()
	defer imgurTestTeardown()

	mux.HandleFunc("/gallery/random/random/0", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, galleryRandomResponse)
	})

	albumImgs, err := client.Gallery.Random(0)
	if err != nil {
		t.Errorf("Gallery.Random returned error: %v", err)
	}

	if len(albumImgs) == 0 {
		t.Errorf("Gallery.Random returned %+v, want %+v", len(albumImgs), 0)
	}

	want := "1YUpmrH"
	if albumImgs[0].ID != want {
		t.Errorf("Gallery.Random returned %+v, want %+v", albumImgs[0].ID, want)
	}
}

func TestGalleryMain(t *testing.T) {
	imgurTestSetup()
	defer imgurTestTeardown()

	mux.HandleFunc("/gallery/hot/viral/day/0", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, galleryMainResponse)
	})

	albumImgs, err := client.Gallery.Main("", "", "", 0)
	if err != nil {
		t.Errorf("Gallery.Main returned error: %v", err)
	}

	if len(albumImgs) == 0 {
		t.Errorf("Gallery.Main returned %+v, want %+v", len(albumImgs), 0)
	}

	want := "Lh6MuPp"
	if albumImgs[0].ID != want {
		t.Errorf("Gallery.Main returned %+v, want %+v", albumImgs[0].ID, want)
	}
}

func TestGallerySubreddit(t *testing.T) {
	imgurTestSetup()
	defer imgurTestTeardown()

	mux.HandleFunc("/gallery/r/unexpected/time/week/0", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, gallerySubredditResponse)
	})

	albumImgs, err := client.Gallery.Subreddit("unexpected", "", "", 0)
	if err != nil {
		t.Errorf("Gallery.Subreddit returned error: %v", err)
	}

	if len(albumImgs) == 0 {
		t.Errorf("Gallery.Subreddit returned %+v, want %+v", len(albumImgs), 0)
	}

	want := "4C06G"
	if albumImgs[0].ID != want {
		t.Errorf("Gallery.Subreddit returned %+v, want %+v", albumImgs[0].ID, want)
	}
}

func TestGalleryMemes(t *testing.T) {
	imgurTestSetup()
	defer imgurTestTeardown()

	mux.HandleFunc("/gallery/g/memes/viral/week/0", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, galleryMemesResponse)
	})

	albumImgs, err := client.Gallery.Memes("", "", 0)
	if err != nil {
		t.Errorf("Gallery.Memes returned error: %v", err)
	}

	if len(albumImgs) == 0 {
		t.Errorf("Gallery.Memes returned %+v, want %+v", len(albumImgs), 0)
	}

	want := "zHQ2rzI"
	if albumImgs[0].ID != want {
		t.Errorf("Gallery.Memes returned %+v, want %+v", albumImgs[0].ID, want)
	}
}
