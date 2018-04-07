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
	ErrUnauthorized     = ErrResponse{401, "Unauthorized"}
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

// String ...
func String(s string) *string {
	return &s
}

// ID ...
func ID(id int) *int {
	return &id
}

// Int ...
func Int(i int) *int {
	return &i
}

// Time is a convenience function for getting a pointer to a time.Time
// from an expression
func Time(t time.Time) *time.Time {
	return &t
}

// TODO: fill out docs
var (
	Archived        = &ptrue
	Unarchived      = &pfalse
	ShowThermometer = &ptrue
	HideThermometer = &pfalse
	ResetID         = ID(-1)
	ResetEstimate   = ID(-1)
	ResetTime       = Time(time.Time{})
	ResetColor      = String("")
	EmptyString     = String("")

	ptrue  = true
	pfalse = false
)

// Client represents a Clubhouse API client
type Client struct {
	AuthToken  string
	RootURL    string
	Version    string
	HTTPClient *http.Client
	Limiter    ratelimit.Limiter
}

// ListCategories returns a list of all categories and their attributes
func (c *Client) ListCategories() (Categories, error) {
	categories := Categories{}
	if err := c.getResource(&categories); err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategory returns information about the selected category
func (c *Client) GetCategory(id int) (*Category, error) {
	category := Category{ID: id}
	if err := c.getResource(&category); err != nil {
		return nil, err
	}
	return &category, nil
}

// UpdateCategory allows you to replace a Category name with another
// name. If you try to name a Category something that already exists,
// you will get an ErrUnprocessable error.
func (c *Client) UpdateCategory(id int, params *UpdateCategoryParams) (*Category, error) {
	category := Category{ID: id}
	if err := c.updateResource(&category, params); err != nil {
		return nil, err
	}
	return &category, nil
}

// DeleteCategory deletes a category
func (c *Client) DeleteCategory(id int) error {
	category := Category{ID: id}
	return c.deleteResource(&category)
}

// CreateCategory creates a new category. If Category is given a name
// that already exists, you will get an ErrUnprocessable error.
func (c *Client) CreateCategory(params *CreateCategoryParams) (*Category, error) {
	if params.Type == "" {
		params.Type = CategoryTypeMilestone
	}
	category := Category{}
	if err := c.createResource(&category, params); err != nil {
		return nil, err
	}
	return &category, nil
}

// ListEpics lists all the epics
func (c *Client) ListEpics() (Epics, error) {
	epics := Epics{}
	if err := c.getResource(&epics); err != nil {
		return nil, err
	}
	return epics, nil
}

// CreateEpic ...
func (c *Client) CreateEpic(params *CreateEpicParams) (*Epic, error) {
	epic := Epic{}
	if err := c.createResource(&epic, params); err != nil {
		return nil, err
	}
	return &epic, nil
}

// GetEpic gets an epic by ID
func (c *Client) GetEpic(id int) (*Epic, error) {
	epic := Epic{ID: id}
	if err := c.getResource(&epic); err != nil {
		return nil, err
	}
	return &epic, nil
}

// UpdateEpic ...
func (c *Client) UpdateEpic(id int, params UpdateEpicParams) (*Epic, error) {
	epic := Epic{ID: id}
	if err := c.updateResource(&epic, params); err != nil {
		return nil, err
	}
	return &epic, nil
}

// DeleteEpic ...
func (c *Client) DeleteEpic(id int) error {
	epic := Epic{ID: id}
	return c.deleteResource(&epic)
}

// CreateEpicComment ...
func (c *Client) CreateEpicComment(epicID int, params *CreateCommentParams) (*ThreadedComment, error) {
	epic := Epic{ID: epicID}
	comment := ThreadedComment{parent: &epic}
	if err := c.createResource(&comment, params); err != nil {
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
	epic := Epic{ID: epicID}
	comment := ThreadedComment{ID: commentID, parent: &epic}
	if err := c.updateResource(&comment, params); err != nil {
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
	epic := Epic{ID: epicID}
	comment := ThreadedComment{ID: commentID, parent: &epic}
	if err := c.createResource(&comment, params); err != nil {
		return nil, err
	}
	return &comment, nil
}

// ListEpicComments ...
func (c *Client) ListEpicComments(epicID int) (ThreadedComments, error) {
	resource := path.Join("epics", strconv.Itoa(epicID), "comments")
	bytes, err := c.Request("GET", resource)
	if err != nil {
		return nil, err
	}
	comments := ThreadedComments{}
	if err := json.Unmarshal(bytes, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

// GetEpicComment ...
func (c *Client) GetEpicComment(epicID, commentID int) (*ThreadedComment, error) {
	epic := Epic{ID: epicID}
	comment := ThreadedComment{ID: commentID, parent: &epic}
	if err := c.getResource(&comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

// DeleteEpicComment ...
func (c *Client) DeleteEpicComment(epicID, commentID int) error {
	epic := Epic{ID: epicID}
	comment := ThreadedComment{ID: commentID, parent: &epic}
	return c.deleteResource(&comment)
}

// FileUpload ...
type FileUpload struct {
	Name string
	File io.Reader
}

// UploadFiles ...
func (c *Client) UploadFiles(fs []FileUpload) ([]File, error) {
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
func (c *Client) ListFiles() (Files, error) {
	files := Files{}
	if err := c.getResource(&files); err != nil {
		return nil, err
	}
	return files, nil
}

// GetFile ...
func (c *Client) GetFile(id int) (*File, error) {
	file := File{ID: id}
	if err := c.getResource(&file); err != nil {
		return nil, err
	}
	return &file, nil
}

// UpdateFile ...
func (c *Client) UpdateFile(id int, params *UpdateFileParams) (*File, error) {
	file := File{ID: id}
	if err := c.updateResource(&file, params); err != nil {
		return nil, err
	}
	return &file, nil
}

// DeleteFile ...
func (c *Client) DeleteFile(id int) error {
	file := File{ID: id}
	return c.deleteResource(&file)
}

// ListLabels ...
func (c *Client) ListLabels() (Labels, error) {
	labels := Labels{}
	err := c.getResource(&labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

// CreateLabel ...
func (c *Client) CreateLabel(params *CreateLabelParams) (*Label, error) {
	label := Label{}
	err := c.createResource(&label, params)
	if err != nil {
		return nil, err
	}
	return &label, nil
}

// GetLabel ...
func (c *Client) GetLabel(id int) (*Label, error) {
	label := Label{ID: id}
	err := c.getResource(&label)
	if err != nil {
		return nil, err
	}
	return &label, nil
}

// DeleteLabel ...
func (c *Client) DeleteLabel(id int) error {
	label := Label{ID: id}
	return c.deleteResource(&label)
}

// UpdateLabel ...
func (c *Client) UpdateLabel(id int, params *UpdateLabelParams) (*Label, error) {
	label := Label{ID: id}
	err := c.updateResource(&label, params)
	if err != nil {
		return nil, err
	}
	return &label, nil
}

// ListMembers ...
func (c *Client) ListMembers() (Members, error) {
	members := Members{}
	err := c.getResource(&members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

// GetMember ...
func (c *Client) GetMember(id string) (*Member, error) {
	member := Member{ID: id}
	err := c.getResource(&member)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// ListMilestones ...
func (c *Client) ListMilestones() (Milestones, error) {
	milestones := Milestones{}
	err := c.getResource(&milestones)
	if err != nil {
		return nil, err
	}
	return milestones, nil
}

// CreateMilestone ...
func (c *Client) CreateMilestone(params *CreateMilestoneParams) (*Milestone, error) {
	milestone := Milestone{}
	err := c.createResource(&milestone, params)
	if err != nil {
		return nil, err
	}
	return &milestone, nil
}

// GetMilestone ...
func (c *Client) GetMilestone(id int) (*Milestone, error) {
	milestone := Milestone{ID: id}
	err := c.getResource(&milestone)
	if err != nil {
		return nil, err
	}
	return &milestone, nil
}

// DeleteMilestone ...
func (c *Client) DeleteMilestone(id int) error {
	milestone := Milestone{ID: id}
	return c.deleteResource(&milestone)
}

// UpdateMilestone ...
func (c *Client) UpdateMilestone(id int, params *UpdateMilestoneParams) (*Milestone, error) {
	milestone := Milestone{ID: id}
	err := c.updateResource(&milestone, params)
	if err != nil {
		return nil, err
	}
	return &milestone, nil
}

// ListProjects ...
func (c *Client) ListProjects() (Projects, error) {
	projects := Projects{}
	err := c.getResource(&projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

// CreateProject ...
func (c *Client) CreateProject(params *CreateProjectParams) (*Project, error) {
	project := Project{}
	err := c.createResource(&project, params)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// GetProject ...
func (c *Client) GetProject(id int) (*Project, error) {
	project := Project{ID: id}
	err := c.getResource(&project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// DeleteProject ...
func (c *Client) DeleteProject(id int) error {
	project := Project{ID: id}
	return c.deleteResource(&project)
}

// UpdateProject ...
func (c *Client) UpdateProject(id int, params *UpdateProjectParams) (*Project, error) {
	project := Project{ID: id}
	err := c.updateResource(&project, params)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// ListRepositories ...
func (c *Client) ListRepositories() (Repositories, error) {
	repos := Repositories{}
	err := c.getResource(&repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

// GetRepository ...
func (c *Client) GetRepository(id int) (*Repository, error) {
	repo := Repository{ID: id}
	err := c.getResource(&repo)
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

// CreateStory ...
func (c *Client) CreateStory(params *CreateStoryParams) (*Story, error) {
	story := Story{}
	err := c.createResource(&story, params)
	if err != nil {
		return nil, err
	}
	return &story, nil
}

// GetStory ...
func (c *Client) GetStory(id int) (*Story, error) {
	story := Story{ID: id}
	err := c.getResource(&story)
	if err != nil {
		return nil, err
	}
	return &story, nil
}

// DeleteStory ...
func (c *Client) DeleteStory(id int) error {
	story := Story{ID: id}
	return c.deleteResource(&story)
}

// UpdateStory ...
func (c *Client) UpdateStory(id int, params *UpdateStoryParams) (*Story, error) {
	story := Story{ID: id}
	err := c.updateResource(&story, params)
	if err != nil {
		return nil, err
	}
	return &story, nil
}

// Request makes an HTTP request to the Clubhouse API without a body. See
// RequestWithBody for full documentation.
func (c *Client) Request(method string, endpoint string) ([]byte, error) {
	return c.RequestWithBody(method, endpoint, []byte{}, nil)
}

// ErrClientRequest is returned when the client runs into
// problems making a request.
type ErrClientRequest struct {
	Err          error
	Method       string
	URL          string
	Request      *http.Request
	Response     *http.Response
	RequestBody  []byte
	ResponseBody []byte
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
	// finish setup or panic if the client isn't configured correctly
	c.checkSetup()

	url := c.makeURL(endpoint)
	body := bytes.NewBuffer(content)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, ErrClientRequest{
			Err:         err,
			URL:         url,
			Method:      method,
			Request:     req,
			RequestBody: content,
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
			Err:         err,
			URL:         url,
			Method:      method,
			Request:     req,
			RequestBody: content,
		}
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrClientRequest{
			Err:          err,
			URL:          url,
			Method:       method,
			Request:      req,
			RequestBody:  content,
			Response:     resp,
			ResponseBody: bytes,
		}
	}

	switch resp.StatusCode {
	case 400:
		err = ErrSchemaMismatch
	case 401:
		err = ErrUnauthorized
	case 404:
		err = ErrResourceNotFound
	case 422:
		err = ErrUnprocessable
	case 500:
		err = ErrServerError
	}

	if err != nil {
		return nil, ErrClientRequest{
			Err:          err,
			URL:          url,
			Method:       method,
			Request:      req,
			RequestBody:  content,
			Response:     resp,
			ResponseBody: bytes,
		}
	}
	return bytes, nil
}

func (c *Client) getResource(r Resource) error {
	bytes, err := c.Request("GET", r.MakeURL())
	if err != nil {
		return err
	}
	if len(bytes) > 0 {
		return json.Unmarshal(bytes, r)
	}
	return nil
}
func (c *Client) deleteResource(r Resource) error {
	_, err := c.Request("DELETE", r.MakeURL())
	return err
}
func (c *Client) createResource(r Resource, params interface{}) error {
	return c.createOrUpdateResource("POST", r, params)
}
func (c *Client) updateResource(r Resource, params interface{}) error {
	return c.createOrUpdateResource("PUT", r, params)
}
func (c *Client) createOrUpdateResource(m string, r Resource, p interface{}) error {
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("could not marshal params, %s", err)
	}
	bytes, err := c.RequestWithBody(m, r.MakeURL(), body, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, &r)
}

func (c *Client) checkSetup() {
	if c.AuthToken == "" {
		panic("clubhouse: Client missing AuthToken")
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
