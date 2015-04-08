package imgur

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the GitHub client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

// setup sets up a test HTTP server along with an imgur.Client that is
// configured to talk to that test server.  Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func imgurTestSetup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// imgur client configured to use test server
	client = NewClient(nil, "clientID", "clientSecret")
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
	// client.UploadURL = url
}

// teardown closes the test HTTP server.
func imgurTestTeardown() {
	server.Close()
}

// test that the http method is what we expect
func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}
