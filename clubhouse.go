package clubhouse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"strconv"
	"time"

	"go.uber.org/ratelimit"
)

// We use this a lot so let's alias it.
var itoa = strconv.Itoa

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

// Defaults. You can override any of these to change the default for all
// clients created.
var (
	// Root URL for the API
	DefaultRootURL = "https://api.clubhouse.io/api/"

	// Current version as of 04-2018 is v2
	DefaultVersion = "v2"

	// Clubhouse API is 200/minute, so 3.333 every second, which we
	// round down to 3 since we need to use an int
	DefaultLimiter = RateLimiter(3)

	// DefaultHTTP client is, perhaps unsurprisingly, the default http
	// client.
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

// CreateCategory creates a new category. If Category is given a name
// that already exists, you will get an ErrUnprocessable error.
func (c *Client) CreateCategory(params *CreateCategoryParams) (*Category, error) {
	resource := Category{}
	uri := "categories"

	if params.Type == "" {
		params.Type = CategoryTypeMilestone
	}

	err := c.RequestResource("POST", &resource, uri, params)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// ListCategories returns a list of all categories and their attributes
func (c *Client) ListCategories() ([]Category, error) {
	resource := []Category{}
	uri := "categories"
	if err := c.RequestResource("GET", &resource, uri, nil); err != nil {
		return nil, err
	}
	return resource, nil
}

// GetCategory returns information about the selected category
func (c *Client) GetCategory(id int) (*Category, error) {
	resource := Category{}
	uri := path.Join("categories", itoa(id))
	if err := c.RequestResource("GET", &resource, uri, nil); err != nil {
		return nil, err
	}
	return &resource, nil
}

// UpdateCategory allows you to replace a Category name with another
// name. If you try to name a Category something that already exists,
// you will get an ErrUnprocessable error.
func (c *Client) UpdateCategory(id int, params *UpdateCategoryParams) (*Category, error) {
	resource := Category{}
	uri := path.Join("categories", itoa(id))
	if err := c.RequestResource("PUT", &resource, uri, params); err != nil {
		return nil, err
	}
	return &resource, nil
}

// DeleteCategory deletes a category
func (c *Client) DeleteCategory(id int) error {
	uri := path.Join("categories", itoa(id))
	return c.RequestResource("DELETE", nil, uri, nil)
}

// ListEpics lists all the epics
func (c *Client) ListEpics() ([]Epic, error) {
	resource := []Epic{}
	uri := "epics"
	err := c.RequestResource("GET", &resource, uri, nil)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

// CreateEpic ...
func (c *Client) CreateEpic(params *CreateEpicParams) (*Epic, error) {
	resource := Epic{}
	uri := "epics"
	err := c.RequestResource("POST", &resource, uri, params)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetEpic gets an epic by ID
func (c *Client) GetEpic(id int) (*Epic, error) {
	resource := Epic{}
	uri := path.Join("epics", itoa(id))
	err := c.RequestResource("GET", &resource, uri, nil)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// UpdateEpic ...
func (c *Client) UpdateEpic(id int, params UpdateEpicParams) (*Epic, error) {
	resource := Epic{}
	uri := path.Join("epics", itoa(id))
	err := c.RequestResource("PUT", &resource, uri, params)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// DeleteEpic ...
func (c *Client) DeleteEpic(id int) error {
	uri := path.Join("epics", itoa(id))
	return c.RequestResource("DELETE", nil, uri, nil)
}

// CreateEpicComment ...
func (c *Client) CreateEpicComment(epicID int, params *CreateCommentParams) (*ThreadedComment, error) {
	resource := ThreadedComment{}
	uri := path.Join("epics", itoa(epicID), "comments")
	err := c.RequestResource("POST", &resource, uri, params)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// UpdateEpicComment ...
func (c *Client) UpdateEpicComment(
	epicID int,
	commentID int,
	params *UpdateCommentParams,
) (*ThreadedComment, error) {
	resource := ThreadedComment{}
	uri := path.Join("epics", itoa(epicID), "comments", itoa(commentID))
	err := c.RequestResource("PUT", &resource, uri, params)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// CreateEpicCommentComment ...
func (c *Client) CreateEpicCommentComment(
	epicID int,
	commentID int,
	params *CreateCommentParams,
) (*ThreadedComment, error) {
	resource := ThreadedComment{}
	uri := path.Join("epics", itoa(epicID), "comments", itoa(commentID))
	err := c.RequestResource("POST", &resource, uri, params)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// ListEpicComments ...
func (c *Client) ListEpicComments(epicID int) ([]ThreadedComment, error) {
	resource := []ThreadedComment{}
	uri := path.Join("epics", itoa(epicID), "comments")
	err := c.RequestResource("GET", &resource, uri, nil)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

// GetEpicComment ...
func (c *Client) GetEpicComment(epicID, commentID int) (*ThreadedComment, error) {
	resource := ThreadedComment{}
	uri := path.Join("epics", itoa(epicID), "comments", itoa(commentID))
	err := c.RequestResource("GET", &resource, uri, nil)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// DeleteEpicComment ...
func (c *Client) DeleteEpicComment(epicID, commentID int) error {
	uri := path.Join("epics", itoa(epicID), "comments", itoa(commentID))
	return c.RequestResource("DELETE", nil, uri, nil)
}

// FileUpload ...
type FileUpload struct {
	Name string
	File io.Reader
}

// UploadFiles ...
func (c *Client) UploadFiles(fs []FileUpload) ([]File, error) {
	// FIXME: break this method up. the first half of it can be broken
	// into a function that reqturns (body, header, error)

	resource := "files"
	buf := bytes.NewBuffer([]byte{})
	mp := multipart.NewWriter(buf)

	// use a non-random boundary when we're in test mode so the outgoing
	// payload is predictable and repreated exactly.
	if os.Getenv("CLUBHOUSE_TEST_MODE") == "true" {
		if err := mp.SetBoundary("predictableclubhousetestingboundarywowow"); err != nil {
			log.Fatal("UploadFiles: error setting boundary", err)
		}
	}

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
	bytes, err := c.HTTPRequest("POST", resource, body, &header)
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
	resource := []File{}
	uri := path.Join("files")
	err := c.RequestResource("GET", &resource, uri, nil)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

// GetFile ...
func (c *Client) GetFile(id int) (*File, error) {
	resource := File{}
	uri := path.Join("files", itoa(id))
	err := c.RequestResource("GET", &resource, uri, nil)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// UpdateFile ...
func (c *Client) UpdateFile(id int, params *UpdateFileParams) (*File, error) {
	resource := File{}
	uri := path.Join("files", itoa(id))
	err := c.RequestResource("PUT", &resource, uri, params)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// DeleteFile ...
func (c *Client) DeleteFile(id int) error {
	uri := path.Join("files", itoa(id))
	return c.RequestResource("DELETE", nil, uri, nil)
}

// CreateLabel ...
func (c *Client) CreateLabel(params *CreateLabelParams) (*Label, error) {
	resource := Label{}
	uri := path.Join("labels")
	err := c.RequestResource("POST", &resource, uri, params)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// ListLabels ...
func (c *Client) ListLabels() ([]Label, error) {
	resource := []Label{}
	uri := path.Join("labels")
	err := c.RequestResource("GET", &resource, uri, nil)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

// GetLabel ...
func (c *Client) GetLabel(id int) (*Label, error) {
	resource := Label{}
	uri := path.Join("labels", itoa(id))
	err := c.RequestResource("GET", &resource, uri, nil)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// UpdateLabel ...
func (c *Client) UpdateLabel(id int, params *UpdateLabelParams) (*Label, error) {
	resource := Label{}
	uri := path.Join("labels", itoa(id))
	err := c.RequestResource("PUT", &resource, uri, params)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// DeleteLabel ...
func (c *Client) DeleteLabel(id int) error {
	uri := path.Join("labels", itoa(id))
	return c.RequestResource("DELETE", nil, uri, nil)
}

// ListMembers ...
func (c *Client) ListMembers() ([]Member, error) {
	resource := []Member{}
	uri := path.Join("members")
	err := c.RequestResource("GET", &resource, uri, nil)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

// GetMember ...
func (c *Client) GetMember(uuid string) (*Member, error) {
	resource := Member{}
	uri := path.Join("members", uuid)
	err := c.RequestResource("GET", &resource, uri, nil)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// CreateMilestone ...
func (c *Client) CreateMilestone(params *CreateMilestoneParams) (*Milestone, error) {
	resource := Milestone{}
	uri := path.Join("milestones")
	err := c.RequestResource("POST", &resource, uri, params)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// ListMilestones ...
func (c *Client) ListMilestones() ([]Milestone, error) {
	resource := []Milestone{}
	uri := path.Join("milestones")
	err := c.RequestResource("GET", &resource, uri, nil)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

// GetMilestone ...
func (c *Client) GetMilestone(id int) (*Milestone, error) {
	resource := Milestone{}
	uri := path.Join("milestones", itoa(id))
	err := c.RequestResource("GET", &resource, uri, nil)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// UpdateMilestone ...
func (c *Client) UpdateMilestone(id int, params *UpdateMilestoneParams) (*Milestone, error) {
	resource := Milestone{}
	uri := path.Join("milestones", itoa(id))
	err := c.RequestResource("PUT", &resource, uri, params)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// DeleteMilestone ...
func (c *Client) DeleteMilestone(id int) error {
	uri := path.Join("milestones", itoa(id))
	return c.RequestResource("DELETE", nil, uri, nil)
}

// CreateProject ...
func (c *Client) CreateProject(params *CreateProjectParams) (*Project, error) {
	resource := Project{}
	uri := path.Join("projects")
	err := c.RequestResource("POST", &resource, uri, params)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// ListProjects ...
func (c *Client) ListProjects() ([]Project, error) {
	resource := []Project{}
	uri := path.Join("projects")
	err := c.RequestResource("GET", &resource, uri, nil)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

// GetProject ...
func (c *Client) GetProject(id int) (*Project, error) {
	resource := Project{}
	uri := path.Join("projects", itoa(id))
	err := c.RequestResource("GET", &resource, uri, nil)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// UpdateProject ...
func (c *Client) UpdateProject(id int, params *UpdateProjectParams) (*Project, error) {
	resource := Project{}
	uri := path.Join("projects", itoa(id))
	err := c.RequestResource("PUT", &resource, uri, params)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// DeleteProject ...
func (c *Client) DeleteProject(id int) error {
	uri := path.Join("projects", itoa(id))
	return c.RequestResource("DELETE", nil, uri, nil)
}

// ListRepositories ...
func (c *Client) ListRepositories() ([]Repository, error) {
	resource := []Repository{}
	uri := path.Join("repositories")
	err := c.RequestResource("GET", &resource, uri, nil)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

// GetRepository ...
func (c *Client) GetRepository(id int) (*Repository, error) {
	resource := Repository{}
	uri := path.Join("repositories", itoa(id))
	err := c.RequestResource("GET", &resource, uri, nil)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// CreateStory ...
func (c *Client) CreateStory(params *CreateStoryParams) (*Story, error) {
	resource := Story{}
	uri := path.Join("stories")
	err := c.RequestResource("POST", &resource, uri, params)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

type createStoriesParam struct {
	Stories []CreateStoryParams `json:"stories"`
}

// CreateStories ...
func (c *Client) CreateStories(plist []CreateStoryParams) ([]StorySlim, error) {
	resource := []StorySlim{}
	uri := path.Join("stories", "bulk")
	params := createStoriesParam{Stories: plist}
	err := c.RequestResource("POST", &resource, uri, params)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

// GetStory ...
func (c *Client) GetStory(id int) (*Story, error) {
	resource := Story{}
	uri := path.Join("stories", itoa(id))
	err := c.RequestResource("GET", &resource, uri, nil)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// DeleteStory ...
func (c *Client) DeleteStory(id int) error {
	uri := path.Join("stories", itoa(id))
	return c.RequestResource("DELETE", nil, uri, nil)
}

type deleteStoriesParam struct {
	StoryIDs []int `json:"story_ids"`
}

// DeleteStories ...
func (c *Client) DeleteStories(ids []int) error {
	uri := path.Join("stories", "bulk")
	params := deleteStoriesParam{StoryIDs: ids}
	return c.RequestResource("DELETE", nil, uri, params)
}

// UpdateStory ...
func (c *Client) UpdateStory(id int, params *UpdateStoryParams) (*Story, error) {
	resource := Story{}
	uri := path.Join("stories", itoa(id))
	err := c.RequestResource("PUT", &resource, uri, params)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// UpdateStories ...
func (c *Client) UpdateStories(params *UpdateStoriesParams) ([]StorySlim, error) {
	resource := []StorySlim{}
	uri := path.Join("stories", "bulk")
	err := c.RequestResource("PUT", &resource, uri, params)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

// SearchStories ...
func (c *Client) SearchStories(params *SearchParams) (*SearchResults, error) {
	resource := SearchResults{}
	uri := path.Join("search", "stories")
	err := c.RequestResource("GET", &resource, uri, params)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// SearchStoriesAll ...
func (c *Client) SearchStoriesAll(params *SearchParams) ([]StorySearch, error) {
	collected := []StorySearch{}

	for {
		page, err := c.SearchStories(params)
		if err != nil {
			return nil, err
		}
		collected = append(collected, page.Data...)
		if page.Next == "" {
			break
		}

		// the clubhouse API returns the whole URL to use as the "next"
		// token. unfortunately, that doesn't really work for us, so we
		// parse the URL and extract just the "next" query var from it
		urlparts, err := url.Parse(page.Next)
		if err != nil {
			return nil, fmt.Errorf("error parsing next page url %s", err)
		}
		next := urlparts.Query().Get("next")
		params.Next = next
	}
	return collected, nil
}

// CreateStoryLink ...
func (c *Client) CreateStoryLink(params *CreateStoryLinkParams) (*StoryLink, error) {
	resource := StoryLink{}
	uri := "story-links"
	err := c.RequestResource("POST", &resource, uri, params)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetStoryLink ...
func (c *Client) GetStoryLink(id int) (*StoryLink, error) {
	resource := StoryLink{}
	uri := path.Join("story-links", itoa(id))
	err := c.RequestResource("GET", &resource, uri, nil)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// DeleteStoryLink ...
func (c *Client) DeleteStoryLink(id int) error {
	uri := path.Join("story-links", itoa(id))
	return c.RequestResource("DELETE", nil, uri, nil)
}

// ListTeams ...
func (c *Client) ListTeams() ([]Team, error) {
	resource := []Team{}
	uri := "teams"
	err := c.RequestResource("GET", &resource, uri, nil)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

// GetTeam ...
func (c *Client) GetTeam(id int) (*Team, error) {
	resource := Team{}
	uri := path.Join("teams", itoa(id))
	err := c.RequestResource("GET", &resource, uri, nil)
	if err != nil {
		return nil, err
	}
	return &resource, nil
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

type errMessage struct {
	Message string
}

// HTTPRequest makes an HTTP request to the Clubhouse API.
//
// Ideally you should be able to use the table type methods
// (List/Get/Update/Delete) and shouldn't have to use this too much.
//
// endpoint will be combined with the client's RootlURL, Version and
// BaseID, to create the complete URL. endpoint is expected to already
// be encoded; if necessary, use url.PathEscape before passing
// HTTPRequest.
//
// options takes a value that satisfies the QueryEncoder interface,
// which is a type that has an `Encode() string` method. See the Options
// type in this package, but also know that you can always use an
// instance of url.Values.
//
// If client is missing AuthToken, this method will panic.
func (c *Client) HTTPRequest(
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

	// TODO: I bet 400 errors also have a message?
	if err != nil {
		if err == ErrUnprocessable {
			message := errMessage{}
			jsonerr := json.Unmarshal(bytes, &message)
			if jsonerr == nil {
				err = fmt.Errorf("%s: %s", err, message.Message)
			}
		}

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

func (c *Client) RequestResource(
	method string,
	resource interface{},
	uri string,
	params interface{},
) error {
	var (
		body = []byte{}
		err  error
	)
	if params != nil {
		body, err = json.Marshal(params)
		if err != nil {
			return fmt.Errorf("could not marshal params, %s", err)
		}

		if os.Getenv("CLUBHOUSE_DEBUG") == "true" {
			log.Print("body", string(body))
		}
	}
	bytes, err := c.HTTPRequest(method, uri, body, nil)
	if err != nil {
		return err
	}
	if resource != nil {
		return json.Unmarshal(bytes, &resource)
	}
	return nil
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
	urlparts, err := url.Parse(c.RootURL)
	if err != nil {
		panic(fmt.Errorf("could not parse RootURL %s", err))
	}
	urlparts.Path = path.Join(urlparts.Path, c.Version, resource)
	urlparts.RawQuery = "token=" + c.AuthToken
	return urlparts.String()
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
