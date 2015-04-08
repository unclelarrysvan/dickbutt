package imgur

import (
	"bytes"
	"encoding/json"
	"fmt"
	// "log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	libraryVersion = "0.1"
	userAgent      = "go-imgur/" + libraryVersion

	defaultBaseURL = "https://api.imgur.com/3/"

	hdrUserRateLimit       = "X-RateLimit-UserLimit"       // per user limit
	hdrUserRateRemaining   = "X-RateLimit-UserRemaining"   // per user remaining
	hdrUserRateReset       = "X-RateLimit-UserReset"       // timestamp (unix epoch) for when the user credits will be reset.
	hdrClientRateLimit     = "X-RateLimit-ClientLimit"     // total application limit
	hdrClientRateRemaining = "X-RateLimit-ClientRemaining" // total application remaining
)

type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	BaseURL *url.URL

	// User agent used when communicating with the imgur API.
	UserAgent string

	// Rate specifies the current rate limit for the client as determined by the
	// most recent API call.  If the client is used in a multi-user application,
	// this rate may not always be up-to-date.  Call RateLimit() to check the
	// current rate.
	Rate Rate

	// Services used for talking to different parts of the API.
	Gallery *GalleryService
	Image   *ImageService

	clientID     string
	clientSecret string
}

type Result struct {
	Status  int
	Success bool
}

// Response is a Imgur API response.  This wraps the standard http.Response
// returned from Imgur and provides convenient access to things like
// pagination links.
type Response struct {
	*http.Response

	// These fields provide the page values for paginating through a set of
	// results.  Any or all of these may be set to the zero value for
	// responses that are not part of a paginated set, or for which there
	// are no additional pages.

	NextPage  int
	PrevPage  int
	FirstPage int
	LastPage  int

	Rate
}

// NewClient returns a new Imgur API client.  If a nil httpClient is
// provided, http.DefaultClient will be used.  To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the goauth2 library).
func NewClient(httpClient *http.Client, apiId, apiSecret string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient,
		BaseURL:      baseURL,
		UserAgent:    userAgent,
		clientID:     fmt.Sprintf("Client-ID %s", apiId),
		clientSecret: apiSecret,
	}
	c.Gallery = &GalleryService{client: c}
	c.Image = &ImageService{client: c}

	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.  If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.UserAgent)
	req.Header.Add("Authorization", c.clientID)
	return req, nil
}

// Do sends an API request and returns the API response.  The API response is
// decoded and stored in the value pointed to by v, or returned as an error if
// an API error has occurred.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)

	c.Rate = response.Rate

	err = CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return response, err
}

// newResponse creats a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	resp := &Response{Response: r}
	// resp.populatePageValues()
	resp.populateRate()
	return resp
}

type ErrorResponse struct {
	Data struct {
		Error   string `json:"error"`
		Request string `json:"request"`
		Method  string `json:"method"`
	}
	Status  int
	Success bool
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("err code %v for request '%v': %v",
		e.Status, e.Data.Request, e.Data.Error)
}

// CheckResponse checks the API response for errors, and returns them if
// present.  A response is considered an error if it has a status code outside
// the 200 range.  API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.  Any other
// response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	e := &ErrorResponse{}
	if err := json.NewDecoder(r.Body).Decode(e); err != nil {
		return err
	}
	return e
}

// populateRate parses the rate related headers and populates the response Rate.
func (r *Response) populateRate() {
	if ulimit := r.Header.Get(hdrUserRateLimit); ulimit != "" {
		r.Rate.UserLimit, _ = strconv.Atoi(ulimit)
	}

	if uremaining := r.Header.Get(hdrUserRateRemaining); uremaining != "" {
		r.Rate.UserRemaining, _ = strconv.Atoi(uremaining)
	}

	if reset := r.Header.Get(hdrUserRateReset); reset != "" {
		if v, _ := strconv.ParseInt(reset, 10, 64); v != 0 {
			r.Rate.UserReset = time.Unix(v, 0)
		}
	}

	if climit := r.Header.Get(hdrClientRateLimit); climit != "" {
		r.Rate.ClientRemaining, _ = strconv.Atoi(climit)
	}

	if cremaining := r.Header.Get(hdrClientRateRemaining); cremaining != "" {
		r.Rate.ClientRemaining, _ = strconv.Atoi(cremaining)
	}
}

// Rate represents the rate limit for the current client.
// Each application can allow approximately 1,250 uploads per day
// or approximately 12,500 requests per day.
type Rate struct {
	// The number of requests per hour the client is currently limited to.
	UserLimit int

	// The number of remaining requests the client can make this hour.
	UserRemaining int

	// The time at which the current rate limit will reset.
	UserReset time.Time

	// The number of requests allowed for your application this month.
	ClientLimit int

	// The number of remaining requests your application can make this month.
	ClientRemaining int
}

// RateLimit returns the rate limit for the current client.
func (c *Client) RateLimit() (*Rate, *Response, error) {
	req, err := c.NewRequest("GET", "credits", nil)
	if err != nil {
		return nil, nil, err
	}

	// API response wrapper to a rate limit request.
	type rateResponse struct {
		Data struct {
			UserLimit       int
			UserRemaining   int
			UserReset       int64
			ClientLimit     int
			ClientRemaining int
		}
		Status  int
		Success bool
	}

	rr := &rateResponse{}
	resp, err := c.Do(req, rr)
	if err != nil {
		return nil, nil, err
	}

	rate := &Rate{
		UserLimit:       rr.Data.UserLimit,
		UserRemaining:   rr.Data.UserRemaining,
		UserReset:       time.Unix(rr.Data.UserReset, 0),
		ClientLimit:     rr.Data.ClientLimit,
		ClientRemaining: rr.Data.ClientRemaining,
	}
	return rate, resp, err
}
