package clubhouse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"

	"go.uber.org/ratelimit"
)

// ErrResponse ...
type ErrResponse struct {
	Code    int
	Message string
}

func (e ErrResponse) Error() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}

// Errors
var (
	ErrSchemaMismatch   = ErrResponse{400, "Schema mismatch"}
	ErrResourceNotFound = ErrResponse{404, "Resource does not exist"}
	ErrUnprocessable    = ErrResponse{422, "Unprocessable"}
)

// Defaults
var (
	// Root URL for the API
	DefaultRootURL = "https://api.clubhouse.io/api/"

	// Current version as of 04-2018 is v2
	DefaultVersion = "v2"

	// Clubhouse API is 200/minute, so 3.333 every second, which we
	// round down to 3 since we need to use an int
	DefaultLimiter = RateLimiter(3)

	DefaultHTTPClient = http.DefaultClient
)

// RateLimiter makes a new rate limiter using n as the number of
// requests per second that is allowed. If 0 is passed, the limiter will
// be unlimited.
func RateLimiter(n int) ratelimit.Limiter {
	if n == 0 {
		return ratelimit.NewUnlimited()
	}
	return ratelimit.New(n)
}

// Client represents a Clubhouse API client
type Client struct {
	AuthToken  string
	RootURL    string
	Version    string
	HTTPClient *http.Client
	Limiter    ratelimit.Limiter
}

// ListCategories returns a list of all categories and their attributes
func (c *Client) ListCategories() ([]Category, error) {
	c.checkSetup()
	bytes, err := c.Request("GET", "categories")
	if err != nil {
		return nil, err
	}
	categories := []Category{}
	if err := json.Unmarshal(bytes, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategory returns information about the selected category
func (c *Client) GetCategory(id int) (*Category, error) {
	c.checkSetup()
	resource := path.Join("categories", strconv.Itoa(id))
	bytes, err := c.Request("GET", resource)
	if err != nil {
		return nil, err
	}
	category := Category{}
	if err := json.Unmarshal(bytes, &category); err != nil {
		return nil, err
	}
	return &category, nil
}

// Color returns a pointer to a color string for use in params
func Color(c string) *string {
	return &c
}

// ResetColor will set a blank color
var ResetColor = Color("")

// Archive status
var ()
var (
	archived   = true
	unarchived = false
	Archived   = &archived
	Unarchived = &unarchived
)

// UpdateCategoryParams contains the parameters for UpdateCategory
// requests.
type UpdateCategoryParams struct {
	Archived *bool   `json:"archived,omitempty"`
	Color    *string `json:"color"`
	Name     string  `json:"name,omitempty"`
}

type updateCategoryParamsWithoutColor struct {
	Archived *bool   `json:"archived,omitempty"`
	Color    *string `json:"color,omitempty"`
	Name     string  `json:"name,omitempty"`
}

// UpdateCategory allows you to replace a Category name with another
// name. If you try to name a Category something that already exists,
// you will get an ErrUnprocessable error.
func (c *Client) UpdateCategory(id int, params UpdateCategoryParams) (*Category, error) {
	c.checkSetup()
	resource := path.Join("categories", strconv.Itoa(id))

	// We want this function to be ergonomic: if the user passes a
	// value, we assume they mean to update it. If they omit that value,
	// we assume they don't want to update it and leave it at the
	// current value. We have to handle Color specially because in order
	// to reset the Color, we have to send {"color": null}, *but* we
	// only want to send that if the user has explicitly indicated that
	// they want to reset the color, *not* in cases where the Color has
	// been left out of the parameter list.
	var bodyp interface{}
	if params.Color == nil {
		bodyp = updateCategoryParamsWithoutColor{
			Archived: params.Archived,
			Name:     params.Name,
		}
	} else {
		if *params.Color == "" {
			params.Color = nil
		}
		bodyp = params
	}

	bodybytes, err := json.Marshal(bodyp)
	if err != nil {
		return nil, fmt.Errorf("UpdateCategory: could not marshal params, %s", err)
	}

	body := bytes.NewBuffer(bodybytes)
	bytes, err := c.RequestWithBody("PUT", resource, body)
	if err != nil {
		return nil, err
	}
	category := Category{}
	if err := json.Unmarshal(bytes, &category); err != nil {
		return nil, err
	}
	return &category, nil
}

// DeleteCategory deletes a category
func (c *Client) DeleteCategory(id int) error {
	resource := path.Join("categories", strconv.Itoa(id))
	_, err := c.Request("DELETE", resource)
	return err
}

// CreateCategory creates a new category
func (c *Client) CreateCategory(params CreateCategoryParams) (*Category, error) {
	if params.Type == "" {
		params.Type = CategoryTypeMilestone
	}

	bodybytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("CreateCategory: could not marshal params, %s", err)
	}

	body := bytes.NewBuffer(bodybytes)
	bytes, err := c.RequestWithBody("POST", "categories", body)
	if err != nil {
		return nil, err
	}
	category := Category{}
	if err := json.Unmarshal(bytes, &category); err != nil {
		return nil, err
	}
	return &category, nil
}

// Request makes an HTTP request to the Clubhouse API without a body. See
// RequestWithBody for full documentation.
func (c *Client) Request(
	method string,
	endpoint string,
) ([]byte, error) {
	return c.RequestWithBody(method, endpoint, http.NoBody)
}

// ErrClientRequest is returned when the client runs into
// problems making a request.
type ErrClientRequest struct {
	Err    error
	Method string
	URL    string
}

func (e ErrClientRequest) Error() string {
	return fmt.Sprintf("clubhouse client request error: %s %s: %s", e.Method, e.URL, e.Err)
}

// RequestWithBody makes an HTTP request to the Clubhouse API.
//
// Ideally you should be able to use the table type methods
// (List/Get/Update/Delete) and shouldn't have to use this too much.
//
// endpoint will be combined with the client's RootlURL, Version and
// BaseID, to create the complete URL. endpoint is expected to already
// be encoded; if necessary, use url.PathEscape before passing
// RequestWithBody.
//
// options takes a value that satisfies the QueryEncoder interface,
// which is a type that has an `Encode() string` method. See the Options
// type in this package, but also know that you can always use an
// instance of url.Values.
//
// If client is missing APIKey or BaseID, this method will panic.
func (c *Client) RequestWithBody(
	method string,
	endpoint string,
	body io.Reader,
) ([]byte, error) {
	var err error

	// finish setup or panic if the client isn't configured correctly
	c.checkSetup()

	url := c.makeURL(endpoint)
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, ErrClientRequest{
			Err:    err,
			URL:    url,
			Method: method,
		}
	}

	c.makeHeader(req)

	// Take() will block until we can safely make the next request
	// without going over the rate limit
	c.Limiter.Take()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, ErrClientRequest{
			Err:    err,
			URL:    url,
			Method: method,
		}
	}

	switch resp.StatusCode {
	case 400:
		err = ErrSchemaMismatch
	case 404:
		err = ErrResourceNotFound
	case 422:
		err = ErrUnprocessable
	}

	if err != nil {
		return nil, ErrClientRequest{
			Err:    err,
			URL:    url,
			Method: method,
		}
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrClientRequest{
			Err:    err,
			URL:    url,
			Method: method,
		}
	}

	return bytes, nil
}

func (c *Client) checkSetup() {
	if c.AuthToken == "" {
		panic("clubhouse: Client missing APIKey")
	}
	if c.HTTPClient == nil {
		c.HTTPClient = DefaultHTTPClient
	}
	if c.Version == "" {
		c.Version = DefaultVersion
	}
	if c.RootURL == "" {
		c.RootURL = DefaultRootURL
	}
	if c.Limiter == nil {
		c.Limiter = DefaultLimiter
	}
}

func (c *Client) makeHeader(r *http.Request) {
	r.Header = http.Header{}
	r.Header.Add("Content-Type", "application/json")
}

func (c *Client) makeURL(resource string) string {
	uri := path.Join(c.Version, resource)
	uri += "?token=" + c.AuthToken
	return c.RootURL + uri
}
