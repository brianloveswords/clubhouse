package clubhouse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path"
	"reflect"
	"strconv"
	"time"

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
	ErrServerError      = ErrResponse{500, "Server error"}
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
func (c *Client) UpdateCategory(id int, params *UpdateCategoryParams) (*Category, error) {
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
	// TODO: update this to use a similar algorithm to UpdateEpicParams
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

	body, err := json.Marshal(&bodyp)
	if err != nil {
		return nil, fmt.Errorf("UpdateCategory: could not marshal params, %s", err)
	}
	bytes, err := c.RequestWithBody("PUT", resource, body, nil)
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

// CreateCategory creates a new category.
func (c *Client) CreateCategory(params *CreateCategoryParams) (*Category, error) {
	if params.Type == "" {
		params.Type = CategoryTypeMilestone
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("CreateCategory: could not marshal params, %s", err)
	}
	bytes, err := c.RequestWithBody("POST", "categories", body, nil)
	if err != nil {
		return nil, err
	}
	category := Category{}
	if err := json.Unmarshal(bytes, &category); err != nil {
		return nil, err
	}
	return &category, nil
}

// ListEpics lists all the epics
func (c *Client) ListEpics() ([]Epic, error) {
	c.checkSetup()
	bytes, err := c.Request("GET", "epics")
	if err != nil {
		return nil, err
	}
	epics := []Epic{}
	if err := json.Unmarshal(bytes, &epics); err != nil {
		return nil, err
	}
	return epics, nil
}

// EpicState ...
type EpicState string

// Epic State values
const (
	EpicStateDone       EpicState = "done"
	EpicStateInProgress           = "in progress"
	EpicStateToDo                 = "to do"
)

// Time is a convenience function for getting a pointer to a time.Time
// from an expression
func Time(t time.Time) *time.Time {
	return &t
}

// ResetTime is the sentinel value for indicating that a null value
// should be sent for a time type.
var ResetTime = Time(time.Time{})

// CreateEpicParams ...
type CreateEpicParams struct {
	CompletedAtOverride *time.Time          `json:"completed_at_override,omitempty"`
	CreatedAt           *time.Time          `json:"created_at,omitempty"`
	Deadline            *time.Time          `json:"deadline,omitempty"`
	ExternalID          string              `json:"external_id,omitempty"`
	FollowerIDs         []string            `json:"follower_ids,omitempty"`
	Labels              []CreateLabelParams `json:"labels,omitempty"`
	MilestoneID         int                 `json:"milestone_id,omitempty"`
	Name                string              `json:"name"`
	OwnerIDs            []string            `json:"owner_ids,omitempty"`
	StartedAtOverride   *time.Time          `json:"started_at_override,omitempty"`
	State               EpicState           `json:"state,omitempty"`
	UpdatedAt           *time.Time          `json:"updated_at,omitempty"`
}

// CreateEpic ...
func (c *Client) CreateEpic(params *CreateEpicParams) (*Epic, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("CreateEpic: could not marshal params, %s", err)
	}
	bytes, err := c.RequestWithBody("POST", "epics", body, nil)
	if err != nil {
		return nil, err
	}
	epic := Epic{}
	if err := json.Unmarshal(bytes, &epic); err != nil {
		return nil, err
	}
	return &epic, nil
}

// GetEpic gets an epic by ID
func (c *Client) GetEpic(id int) (*Epic, error) {
	c.checkSetup()
	resource := path.Join("epics", strconv.Itoa(id))
	bytes, err := c.Request("GET", resource)
	if err != nil {
		return nil, err
	}
	epic := Epic{}
	if err := json.Unmarshal(bytes, &epic); err != nil {
		return nil, err
	}
	return &epic, nil
}

// ID ...
func ID(id int) *int {
	return &id
}

// ResetID is the sentinel value for indicating that a null should be
// sent for an ID field
var ResetID = ID(-1)

// String ...
func String(s string) *string {
	return &s
}

// EmptyString is a convenience reference to an empty string.
var EmptyString = String("")

// UpdateEpicParams ...
type UpdateEpicParams struct {
	AfterID             *int
	Archived            *bool
	BeforeID            *int
	CompletedAtOverride *time.Time
	Deadline            *time.Time
	Description         *string
	FollowerIDs         []string
	Labels              []CreateLabelParams
	MilestoneID         *int
	Name                string
	OwnerIDs            []string
	StartedAtOverride   *time.Time
	State               EpicState
}

type updateEpicParamsResolved struct {
	AfterID             *int                `json:"after_id,omitempty"`
	Archived            *bool               `json:"archived,omitempty"`
	BeforeID            *int                `json:"before_id,omitempty"`
	CompletedAtOverride *json.RawMessage    `json:"completed_at_override,omitempty"`
	Deadline            *json.RawMessage    `json:"deadline,omitempty"`
	Description         *string             `json:"description,omitempty"`
	FollowerIDs         []string            `json:"follower_ids,omitempty"`
	Labels              []CreateLabelParams `json:"labels,omitempty"`
	MilestoneID         *json.RawMessage    `json:"milestone_id,omitempty"`
	Name                string              `json:"name,omitempty"`
	OwnerIDs            []string            `json:"owner_ids,omitempty"`
	StartedAtOverride   *json.RawMessage    `json:"started_at_override,omitempty"`
	State               EpicState           `json:"state,omitempty"`
}

type nullable []struct {
	in   interface{}
	out  **json.RawMessage
	null func() bool
}

func (n nullable) Do() {
	for _, f := range n {
		if !reflect.ValueOf(f.in).IsNil() {
			var raw json.RawMessage
			if f.null() {
				raw = json.RawMessage(`null`)
			} else {
				raw, _ = json.Marshal(f.in)
			}
			*f.out = &raw
		}
	}
}

// MarshalJSON ...
func (p UpdateEpicParams) MarshalJSON() ([]byte, error) {
	out := updateEpicParamsResolved{
		Archived:    p.Archived,
		AfterID:     p.AfterID,
		BeforeID:    p.BeforeID,
		Description: p.Description,
		FollowerIDs: p.FollowerIDs,
		Labels:      p.Labels,
		Name:        p.Name,
		OwnerIDs:    p.OwnerIDs,
		State:       p.State,
	}

	nullable{{
		in:   p.CompletedAtOverride,
		out:  &out.CompletedAtOverride,
		null: func() bool { return p.CompletedAtOverride.IsZero() },
	}, {
		in:   p.StartedAtOverride,
		out:  &out.StartedAtOverride,
		null: func() bool { return p.StartedAtOverride.IsZero() },
	}, {
		in:   p.Deadline,
		out:  &out.Deadline,
		null: func() bool { return p.Deadline.IsZero() },
	}, {
		in:   p.MilestoneID,
		out:  &out.MilestoneID,
		null: func() bool { return p.MilestoneID == ResetID },
	}}.Do()

	return json.Marshal(&out)
}

// UpdateEpic ...
func (c *Client) UpdateEpic(id int, params UpdateEpicParams) (*Epic, error) {
	resource := path.Join("epics", strconv.Itoa(id))
	body, err := json.Marshal(&params)
	if err != nil {
		return nil, fmt.Errorf("UpdateEpic: could not marshal params, %s", err)
	}
	bytes, err := c.RequestWithBody("PUT", resource, body, nil)
	if err != nil {
		return nil, err
	}
	epic := Epic{}
	if err := json.Unmarshal(bytes, &epic); err != nil {
		return nil, err
	}
	return &epic, nil
}

// DeleteEpic creates an epic
func (c *Client) DeleteEpic(id int) error {
	resource := path.Join("epics", strconv.Itoa(id))
	_, err := c.Request("DELETE", resource)
	return err
}

// CreateEpicComment ...
func (c *Client) CreateEpicComment(epicID int, params *CreateCommentParams) (*ThreadedComment, error) {
	resource := path.Join("epics", strconv.Itoa(epicID), "comments")
	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("CreateEpicComment: could not marshal params, %s", err)
	}
	bytes, err := c.RequestWithBody("POST", resource, body, nil)
	if err != nil {
		return nil, err
	}
	comment := ThreadedComment{}
	if err := json.Unmarshal(bytes, &comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

// UpdateEpicComment ...
func (c *Client) UpdateEpicComment(
	epicID int,
	commentID int,
	params *UpdateCommentParams,
) (*ThreadedComment, error) {
	resource := path.Join(
		"epics", strconv.Itoa(epicID),
		"comments", strconv.Itoa(commentID),
	)
	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("UpdateEpicComment: could not marshal params, %s", err)
	}
	bytes, err := c.RequestWithBody("PUT", resource, body, nil)
	if err != nil {
		return nil, err
	}
	comment := ThreadedComment{}
	if err := json.Unmarshal(bytes, &comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

// CreateEpicCommentComment ...
func (c *Client) CreateEpicCommentComment(
	epicID int,
	commentID int,
	params *CreateCommentParams,
) (*ThreadedComment, error) {
	resource := path.Join(
		"epics", strconv.Itoa(epicID),
		"comments", strconv.Itoa(commentID),
	)
	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("CreateEpicCommentComment: could not marshal params, %s", err)
	}
	bytes, err := c.RequestWithBody("POST", resource, body, nil)
	if err != nil {
		return nil, err
	}
	comment := ThreadedComment{}
	if err := json.Unmarshal(bytes, &comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

// ListEpicComments ...
func (c *Client) ListEpicComments(epicID int) ([]ThreadedComment, error) {
	resource := path.Join("epics", strconv.Itoa(epicID), "comments")
	bytes, err := c.Request("GET", resource)
	if err != nil {
		return nil, err
	}
	comments := []ThreadedComment{}
	if err := json.Unmarshal(bytes, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

// GetEpicComment ...
func (c *Client) GetEpicComment(epicID, commentID int) (*ThreadedComment, error) {
	resource := path.Join(
		"epics", strconv.Itoa(epicID),
		"comments", strconv.Itoa(commentID),
	)
	bytes, err := c.Request("GET", resource)
	if err != nil {
		return nil, err
	}
	comment := ThreadedComment{}
	if err := json.Unmarshal(bytes, &comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

// DeleteEpicComment creates an epic
func (c *Client) DeleteEpicComment(epicID, commentID int) error {
	resource := path.Join(
		"epics", strconv.Itoa(epicID),
		"comments", strconv.Itoa(commentID),
	)
	_, err := c.Request("DELETE", resource)
	return err
}

// FileUpload ...
type FileUpload struct {
	Name string
	File io.Reader
}

// UploadFiles ...
func (c *Client) UploadFiles(fs []FileUpload) ([]File, error) {
	c.checkSetup()
	resource := "files"
	buf := bytes.NewBuffer([]byte{})
	mp := multipart.NewWriter(buf)
	for i, f := range fs {
		p, err := mp.CreateFormFile(fmt.Sprintf("file%d", i), f.Name)
		if err != nil {
			return nil, fmt.Errorf("UploadFiles: couldn't create form file, %s", err)
		}
		_, err = io.Copy(p, f.File)
		if err != nil {
			return nil, fmt.Errorf("UploadFiles: io.Copy error: %s", err)
		}
	}
	ct := mp.FormDataContentType()
	if err := mp.Close(); err != nil {
		return nil, fmt.Errorf("UploadFiles: multipart writer close error %s", err)
	}
	body := buf.Bytes()
	header := http.Header{}
	header.Add("Content-Type", ct)
	bytes, err := c.RequestWithBody("POST", resource, body, &header)
	if err != nil {
		return nil, fmt.Errorf("UploadFiles: error making request: %s", err)
	}
	files := []File{}
	if err := json.Unmarshal(bytes, &files); err != nil {
		return nil, fmt.Errorf("UploadFiles: error unmarshaling response: %s", err)
	}
	return files, nil
}

// ListFiles ...
func (c *Client) ListFiles() ([]File, error) {
	c.checkSetup()
	resource := "files"
	bytes, err := c.Request("GET", resource)
	if err != nil {
		return nil, fmt.Errorf("ListFiles: error making request: %s", err)
	}
	files := []File{}
	if err := json.Unmarshal(bytes, &files); err != nil {
		return nil, fmt.Errorf("ListFiles: error unmarshaling response: %s", err)
	}
	return files, nil
}

// GetFile ...
func (c *Client) GetFile(id int) (*File, error) {
	resource := path.Join("files", strconv.Itoa(id))
	bytes, err := c.Request("GET", resource)
	if err != nil {
		return nil, fmt.Errorf("GetFile: error making request: %s", err)
	}
	file := File{}
	if err := json.Unmarshal(bytes, &file); err != nil {
		return nil, fmt.Errorf("GetFile: error unmarshaling response: %s", err)
	}
	return &file, nil
}

// UpdateFileParams ...
type UpdateFileParams struct {
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	Description *string    `json:"description,omitempty"`
	ExternalID  *string    `json:"external_id,omitempty"`
	Name        *string    `json:"name,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	UploaderID  *string    `json:"uploader_id,omitempty"`
}

// UpdateFile ...
func (c *Client) UpdateFile(id int, params *UpdateFileParams) (*File, error) {
	resource := path.Join("files", strconv.Itoa(id))
	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("UpdateFile: couldn't marshal params %s", err)
	}
	bytes, err := c.RequestWithBody("PUT", resource, body, nil)
	if err != nil {
		return nil, fmt.Errorf("UpdateFile: error making request: %s", err)
	}
	file := File{}
	if err := json.Unmarshal(bytes, &file); err != nil {
		return nil, fmt.Errorf("UpdateFile: error unmarshaling response: %s", err)
	}
	return &file, nil
}

// DeleteFile ...
func (c *Client) DeleteFile(id int) error {
	c.checkSetup()
	resource := path.Join("files", strconv.Itoa(id))
	_, err := c.Request("DELETE", resource)
	return err
}

// Request makes an HTTP request to the Clubhouse API without a body. See
// RequestWithBody for full documentation.
func (c *Client) Request(method string, endpoint string) ([]byte, error) {
	return c.RequestWithBody(method, endpoint, []byte{}, nil)
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
// If client is missing AuthToken, this method will panic.
func (c *Client) RequestWithBody(
	method string,
	endpoint string,
	content []byte,
	header *http.Header,
) ([]byte, error) {
	var err error

	// finish setup or panic if the client isn't configured correctly
	c.checkSetup()

	url := c.makeURL(endpoint)
	body := bytes.NewBuffer(content)
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, ErrClientRequest{
			Err:    err,
			URL:    url,
			Method: method,
		}
	}

	if header == nil {
		header = &http.Header{}
		header.Add("Content-Type", "application/json")
	}
	req.Header = *header

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

func (c *Client) makeURL(resource string) string {
	uri := path.Join(c.Version, resource)
	uri += "?token=" + c.AuthToken
	return c.RootURL + uri
}
