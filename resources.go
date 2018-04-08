package clubhouse

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"
)

// Resource ...
type Resource interface {
	MakeURL() string
}

// State ...
type State string

// State values
const (
	StateDone       State = "done"
	StateInProgress       = "in progress"
	StateToDo             = "to do"
)

// See https://clubhouse.io/api/rest/v2/#Resources for complete
// documentation

// Branch refers to a GitHub branch. Branches are feature branches
// associated with Clubhouse Stories.
type Branch struct {
	CreatedAt       time.Time     `json:"created_at"`
	Deleted         bool          `json:"deleted"`
	EntityType      string        `json:"entity_type"`
	ID              int           `json:"id"`
	MergedBranchIDs []int         `json:"merged_branch_ids"`
	Name            string        `json:"name"`
	Persistent      bool          `json:"persistent"`
	PullRequests    []PullRequest `json:"pull_requests"`
	RepositoryID    int           `json:"repository_id"`
	UpdatedAt       time.Time     `json:"updated_at"`
	URL             string        `json:"url"`
}

// Category can be used to associate Milestones.
type Category struct {
	Archived   bool      `json:"archived"`
	Color      string    `json:"color"`
	CreatedAt  time.Time `json:"created_at"`
	EntityType string    `json:"entity_type"`
	ExternalID string    `json:"external_id"`
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// MakeURL ...
func (c *Category) MakeURL() string {
	if c.ID == 0 && c.Name == "" {
		return "categories"
	}
	return path.Join("categories", strconv.Itoa(c.ID))
}

// Categories is a Category slice
type Categories []Category

// MakeURL ...
func (c *Categories) MakeURL() string {
	return "categories"
}

// Comment is any note added within the Comment field of a Story.
type Comment struct {
	AuthorID   string    `json:"author_id"`
	CreatedAt  time.Time `json:"created_at"`
	EntityType string    `json:"entity_type"`
	ExternalID string    `json:"external_id"`
	ID         int       `json:"id"`
	MentionIDs []string  `json:"mention_ids"`
	Position   int       `json:"position"`
	StoryID    int       `json:"story_id"`
	Text       string    `json:"text"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Commit refers to a GitHub commit and all associated details.
type Commit struct {
	AuthorEmail     string    `json:"author_email"`
	AuthorID        string    `json:"author_id"`
	AuthorIdentity  Identity  `json:"author_identity"`
	CreatedAt       time.Time `json:"created_at"`
	EntityType      string    `json:"entity_type"`
	Hash            string    `json:"hash"`
	ID              int       `json:"id"`
	MergedBranchIDs []int     `json:"merged_branch_ids"`
	Message         string    `json:"message"`
	RepositoryID    int       `json:"repository_id"`
	Timestamp       time.Time `json:"timestamp"`
	URL             string    `json:"url"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CategoryType enum for CreateCategoryParams
type CategoryType string

// Currently the only valid value
var (
	CategoryTypeMilestone CategoryType = "milestone"
)

// CreateCategoryParams represents request parameters for creating a
// Category with a Milestone.
type CreateCategoryParams struct {
	Color      string       `json:"color,omitempty"`
	ExternalID string       `json:"external_id,omitempty"`
	Name       string       `json:"name,omitempty"`
	Type       CategoryType `json:"type,omitempty"`
}

// UpdateCategoryParams contains the parameters for UpdateCategory
// requests.
type UpdateCategoryParams struct {
	Archived *bool
	Color    *string
	Name     *string
}
type updateCategoryParamsResolved struct {
	Archived *bool            `json:"archived,omitempty"`
	Color    *json.RawMessage `json:"color,omitempty"`
	Name     *string          `json:"name,omitempty"`
}

// MarshalJSON ...
func (p UpdateCategoryParams) MarshalJSON() ([]byte, error) {
	out := updateCategoryParamsResolved{
		Archived: p.Archived,
		Name:     p.Name,
	}
	nullable{{
		in:   p.Color,
		out:  &out.Color,
		null: func() bool { return p.Color == ResetColor },
	}}.Do()
	return json.Marshal(&out)
}

// CreateCommentParams represents request parameters for creating a
// Comment on a Clubhouse Story.
type CreateCommentParams struct {
	AuthorID   string     `json:"author_id,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	ExternalID string     `json:"external_id,omitempty"`
	Text       string     `json:"text,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

// UpdateCommentParams ...
type UpdateCommentParams struct {
	Text string `json:"text"`
}

// StoryVerb represents the verb connecting two stories together
type StoryVerb string

// Valid values for StoryVerb
const (
	VerbBlocks     StoryVerb = "blocks"
	VerbDuplicates           = "duplicates"
	VerbRelatesTo            = "relates to"
)

// CreateStoryLinkParams represents request parameters for creating a
// Story Link within a Story.
type CreateStoryLinkParams struct {
	ObjectID  int       `json:"object_id,omitempty"`
	SubjectID int       `json:"subject_id,omitempty"`
	Verb      StoryVerb `json:"verb,omitempty"`
}

// StoryType represents the type of story
type StoryType string

// Valid states for StoryType
const (
	StoryTypeBug     StoryType = "bug"
	StoryTypeChore             = "chore"
	StoryTypeFeature           = "feature"
)

// CreateStoryParams is used to create multiple stories in a single
// request.
type CreateStoryParams struct {
	Comments            []CreateCommentParams   `json:"comments,omitempty"`
	CompletedAtOverride *time.Time              `json:"completed_at_override,omitempty"`
	CreatedAt           *time.Time              `json:"created_at,omitempty"`
	Deadline            *time.Time              `json:"deadline,omitempty"`
	Description         string                  `json:"description,omitempty"`
	EpicID              int                     `json:"epic_id,omitempty"`
	Estimate            int                     `json:"estimate,omitempty"`
	ExternalID          string                  `json:"external_id,omitempty"`
	FileIDs             []int                   `json:"file_ids,omitempty"`
	FollowerIDs         []string                `json:"follower_ids,omitempty"`
	Labels              []CreateLabelParams     `json:"labels,omitempty"`
	LinkedFileIDs       []int                   `json:"linked_file_ids,omitempty"`
	Name                string                  `json:"name,omitempty"`
	OwnerIDs            []string                `json:"owner_ids,omitempty"`
	ProjectID           int                     `json:"project_id,omitempty"`
	RequestedByID       string                  `json:"requested_by_id,omitempty"`
	StartedAtOverride   *time.Time              `json:"started_at_override,omitempty"`
	StoryLinks          []CreateStoryLinkParams `json:"story_links,omitempty"`
	StoryType           StoryType               `json:"story_type,omitempty"`
	Tasks               []CreateTaskParams      `json:"tasks,omitempty"`
	UpdatedAt           *time.Time              `json:"updated_at,omitempty"`
	WorkflowStateID     int                     `json:"workflow_state_id,omitempty"`
}

// CreateTaskParams request parameters for creating a Task on a Story.
type CreateTaskParams struct {
	Complete    bool       `json:"complete,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	Description string     `json:"description,omitempty"`
	ExternalID  string     `json:"external_id,omitempty"`
	OwnerIDs    []string   `json:"owner_ids,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// UpdateStoriesParams ...
type UpdateStoriesParams struct {
	AfterID           *int
	Archived          *bool
	Deadline          *time.Time
	EpicID            *int
	Estimate          *int
	FollowerIDsAdd    []string
	FollowerIDsRemove []string
	LabelsAdd         []CreateLabelParams
	LabelsRemove      []CreateLabelParams
	LinkedFileIDs     []int
	OwnerIDsAdd       []string
	OwnerIDsRemove    []string
	ProjectID         *int
	RequestedByID     *string
	StoryIDs          []int
	StoryType         StoryType
	WorkflowStateID   *int
}
type updateStoriesParamsResolved struct {
	AfterID           *int                `json:"after_id,omitempty"`
	Archived          *bool               `json:"archived,omitempty"`
	Deadline          *json.RawMessage    `json:"deadline,omitempty"`
	EpicID            *json.RawMessage    `json:"epic_id,omitempty"`
	Estimate          *json.RawMessage    `json:"estimate,omitempty"`
	FollowerIDsAdd    []string            `json:"follower_ids_add,omitempty"`
	FollowerIDsRemove []string            `json:"follower_ids_remove,omitempty"`
	LabelsAdd         []CreateLabelParams `json:"labels_add,omitempty"`
	LabelsRemove      []CreateLabelParams `json:"labels_remove,omitempty"`
	LinkedFileIDs     []int               `json:"linked_file_ids,omitempty"`
	OwnerIDsAdd       []string            `json:"owner_ids_add,omitempty"`
	OwnerIDsRemove    []string            `json:"owner_ids_remove,omitempty"`
	ProjectID         *int                `json:"project_id,omitempty"`
	RequestedByID     *string             `json:"requested_by_id,omitempty"`
	StoryIDs          []int               `json:"story_ids,omitempty"`
	StoryType         StoryType           `json:"story_type,omitempty"`
	WorkflowStateID   *int                `json:"workflow_state_id,omitempty"`
}

// MarshalJSON ...
func (p UpdateStoriesParams) MarshalJSON() ([]byte, error) {
	out := updateStoriesParamsResolved{
		AfterID:           p.AfterID,
		Archived:          p.Archived,
		FollowerIDsAdd:    p.FollowerIDsAdd,
		FollowerIDsRemove: p.FollowerIDsRemove,
		LabelsAdd:         p.LabelsAdd,
		LabelsRemove:      p.LabelsRemove,
		LinkedFileIDs:     p.LinkedFileIDs,
		OwnerIDsAdd:       p.OwnerIDsAdd,
		OwnerIDsRemove:    p.OwnerIDsRemove,
		ProjectID:         p.ProjectID,
		RequestedByID:     p.RequestedByID,
		StoryIDs:          p.StoryIDs,
		StoryType:         p.StoryType,
		WorkflowStateID:   p.WorkflowStateID,
	}
	nullable{{
		in:   p.Deadline,
		out:  &out.Deadline,
		null: func() bool { return p.Deadline == ResetTime },
	}, {
		in:   p.EpicID,
		out:  &out.EpicID,
		null: func() bool { return p.EpicID == ResetID },
	}, {
		in:   p.Estimate,
		out:  &out.Estimate,
		null: func() bool { return p.Estimate == ResetEstimate },
	}}.Do()
	return json.Marshal(&out)
}

// UpdateStoryParams ...
type UpdateStoryParams struct {
	AfterID             *int
	Archived            *bool
	BeforeID            *int
	BranchIDs           []int
	CommitIDs           []int
	CompletedAtOverride *time.Time
	Deadline            *time.Time
	Description         *string
	EpicID              *int
	Estimate            *int
	FileIDs             []int
	FollowerIDs         []string
	Labels              []CreateLabelParams
	LinkedFileIDs       []int
	Name                *string
	OwnerIDs            []string
	ProjectID           *int
	RequestedByID       *string
	StartedAtOverride   *time.Time
	StoryType           StoryType
	WorkflowStateID     *int
}
type updateStoryParamsResolved struct {
	AfterID             *int                `json:"after_id,omitempty"`
	Archived            *bool               `json:"archived,omitempty"`
	BeforeID            *int                `json:"before_id,omitempty"`
	BranchIDs           []int               `json:"branch_ids,omitempty"`
	CommitIDs           []int               `json:"commit_ids,omitempty"`
	CompletedAtOverride *json.RawMessage    `json:"completed_at_override,omitempty"`
	Deadline            *json.RawMessage    `json:"deadline,omitempty"`
	Description         *string             `json:"description,omitempty"`
	EpicID              *json.RawMessage    `json:"epic_id,omitempty"`
	Estimate            *json.RawMessage    `json:"estimate,omitempty"`
	FileIDs             []int               `json:"file_ids,omitempty"`
	FollowerIDs         []string            `json:"follower_ids,omitempty"`
	Labels              []CreateLabelParams `json:"labels,omitempty"`
	LinkedFileIDs       []int               `json:"linked_file_ids,omitempty"`
	Name                *string             `json:"name,omitempty"`
	OwnerIDs            []string            `json:"owner_ids,omitempty"`
	ProjectID           *int                `json:"project_id,omitempty"`
	RequestedByID       *string             `json:"requested_by_id,omitempty"`
	StartedAtOverride   *json.RawMessage    `json:"started_at_override,omitempty"`
	StoryType           StoryType           `json:"story_type,omitempty"`
	WorkflowStateID     *int                `json:"workflow_state_id,omitempty"`
}

// MarshalJSON ...
func (p UpdateStoryParams) MarshalJSON() ([]byte, error) {
	out := updateStoryParamsResolved{
		AfterID:         p.AfterID,
		Archived:        p.Archived,
		BeforeID:        p.BeforeID,
		BranchIDs:       p.BranchIDs,
		CommitIDs:       p.CommitIDs,
		Description:     p.Description,
		FileIDs:         p.FileIDs,
		FollowerIDs:     p.FollowerIDs,
		Labels:          p.Labels,
		LinkedFileIDs:   p.LinkedFileIDs,
		Name:            p.Name,
		OwnerIDs:        p.OwnerIDs,
		ProjectID:       p.ProjectID,
		RequestedByID:   p.RequestedByID,
		StoryType:       p.StoryType,
		WorkflowStateID: p.WorkflowStateID,
	}
	nullable{{
		in:   p.Deadline,
		out:  &out.Deadline,
		null: func() bool { return p.Deadline == ResetTime },
	}, {
		in:   p.CompletedAtOverride,
		out:  &out.CompletedAtOverride,
		null: func() bool { return p.CompletedAtOverride == ResetTime },
	}, {
		in:   p.EpicID,
		out:  &out.EpicID,
		null: func() bool { return p.EpicID == ResetID },
	}, {
		in:   p.Estimate,
		out:  &out.Estimate,
		null: func() bool { return p.Estimate == ResetEstimate },
	}, {
		in:   p.StartedAtOverride,
		out:  &out.StartedAtOverride,
		null: func() bool { return p.StartedAtOverride == ResetTime },
	}}.Do()
	return json.Marshal(&out)
}

// Epic is a collection of stories that together might make up a
// release, a milestone, or some other large initiative that your
// organization is working on.
type Epic struct {
	Archived            bool              `json:"archived"`
	Comments            []ThreadedComment `json:"comments"`
	Completed           bool              `json:"completed"`
	CompletedAt         time.Time         `json:"completed_at"`
	CompletedAtOverride time.Time         `json:"completed_at_override"`
	CreatedAt           time.Time         `json:"created_at"`
	Deadline            time.Time         `json:"deadline"`
	Description         string            `json:"description"`
	EntityType          string            `json:"entity_type"`
	ExternalID          string            `json:"external_id"`
	FollowerIDs         []string          `json:"follower_ids"`
	ID                  int               `json:"id"`
	Labels              []Label           `json:"labels"`
	MilestoneID         int               `json:"milestone_id"`
	Name                string            `json:"name"`
	OwnerIDs            []string          `json:"owner_ids"`
	Position            int               `json:"position"`
	ProjectIDs          []int             `json:"project_ids"`
	Started             bool              `json:"started"`
	StartedAt           time.Time         `json:"started_at"`
	StartedAtOverride   time.Time         `json:"started_at_override"`
	State               State             `json:"state"`
	Stats               EpicStats         `json:"stats"`
	UpdatedAt           time.Time         `json:"updated_at"`
}

// MakeURL ...
func (e *Epic) MakeURL() string {
	if e.ID == 0 && e.Name == "" {
		return "epics"
	}
	return path.Join("epics", strconv.Itoa(e.ID))
}

// Epics ...
type Epics []Epic

// MakeURL ...
func (e Epics) MakeURL() string {
	return "epics"
}

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
	State               State               `json:"state,omitempty"`
	UpdatedAt           *time.Time          `json:"updated_at,omitempty"`
}

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
	State               State
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
	State               State               `json:"state,omitempty"`
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

// EpicStats represents a group of calculated values for an Epic.
type EpicStats struct {
	LastStoryUpdate       time.Time `json:"last_story_update"`
	NumPoints             int       `json:"num_points"`
	NumPointsDone         int       `json:"num_points_done"`
	NumPointsStarted      int       `json:"num_points_started"`
	NumPointsUnstarted    int       `json:"num_points_unstarted"`
	NumStoriesDone        int       `json:"num_stories_done"`
	NumStoriesStarted     int       `json:"num_stories_started"`
	NumStoriesUnestimated int       `json:"num_stories_unestimated"`
	NumStoriesUnstarted   int       `json:"num_stories_unstarted"`
}

// File is any document uploaded to your Clubhouse. Files attached from a third-party service can be accessed using the Linked Files endpoint.
type File struct {
	ContentType  string    `json:"content_type"`
	CreatedAt    time.Time `json:"created_at"`
	Description  string    `json:"description"`
	EntityType   string    `json:"entity_type"`
	ExternalID   string    `json:"external_id"`
	Filename     string    `json:"filename"`
	ID           int       `json:"id"`
	MentionIDs   []string  `json:"mention_ids"`
	Name         string    `json:"name"`
	Size         int       `json:"size"`
	StoryIDs     []int     `json:"story_ids"`
	ThumbnailURL string    `json:"thumbnail_url"`
	UpdatedAt    time.Time `json:"updated_at"`
	UploaderID   string    `json:"uploader_id"`
	URL          string    `json:"url"`
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

// MakeURL ...
func (f File) MakeURL() string {
	if f.ID == 0 {
		return "files"
	}
	return path.Join("files", strconv.Itoa(f.ID))
}

// Files ...
type Files []File

// MakeURL ...
func (f Files) MakeURL() string {
	return "files"
}

// Icon is used to attach images to Organizations, Members, and Loading
// screens in the Clubhouse web application.
type Icon struct {
	CreatedAt  time.Time `json:"created_at"`
	EntityType string    `json:"entity_type"`
	ID         string    `json:"id"`
	UpdatedAt  time.Time `json:"updated_at"`
	URL        string    `json:"url"`
}

// Identity is your GitHub login. Clubhouse uses Identity to attempt to
// connect your GitHub actions with your Clubhouse Stories.
type Identity struct {
	EntityType string `json:"entity_type"`
	Name       string `json:"name"`
	Type       string `json:"type"`
}

// Label can be used to associate and filter Stories and Epics, and also create new Workspaces.
type Label struct {
	Archived   bool       `json:"archived"`
	Color      string     `json:"color"`
	CreatedAt  time.Time  `json:"created_at"`
	EntityType string     `json:"entity_type"`
	ExternalID string     `json:"external_id"`
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Stats      LabelStats `json:"stats"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// MakeURL ...
func (l Label) MakeURL() string {
	if l.ID == 0 {
		return "labels"
	}
	return path.Join("labels", strconv.Itoa(l.ID))
}

// Labels ...
type Labels []Label

// MakeURL ...
func (l Labels) MakeURL() string {
	return "labels"
}

// CreateLabelParams represents request parameters for creating a Label
// on a Clubhouse story.
type CreateLabelParams struct {
	Color      string `json:"color,omitempty"`
	ExternalID string `json:"external_id,omitempty"`
	Name       string `json:"name,omitempty"`
}

// UpdateLabelParams ...
type UpdateLabelParams struct {
	Archived *bool
	Color    *string
	Name     *string
}
type updateLabelParamsResolved struct {
	Archived *bool            `json:"archived,omitempty"`
	Color    *json.RawMessage `json:"color,omitempty"`
	Name     *string          `json:"name,omitempty"`
}

// MarshalJSON ...
func (p UpdateLabelParams) MarshalJSON() ([]byte, error) {
	out := updateLabelParamsResolved{
		Archived: p.Archived,
		Name:     p.Name,
	}
	nullable{{
		in:   p.Color,
		out:  &out.Color,
		null: func() bool { return p.Color == ResetColor },
	}}.Do()
	return json.Marshal(&out)
}

// LabelStats represents a group of calculated values for a Label.
type LabelStats struct {
	NumEpics              int `json:"num_epics"`
	NumPointsCompleted    int `json:"num_points_completed"`
	NumPointsInProgress   int `json:"num_points_in_progress"`
	NumPointsTotal        int `json:"num_points_total"`
	NumStoriesCompleted   int `json:"num_stories_completed"`
	NumStoriesInProgress  int `json:"num_stories_in_progress"`
	NumStoriesTotal       int `json:"num_stories_total"`
	NumStoriesUnestimated int `json:"num_stories_unestimated"`
}

// LinkedFile represents files that are stored on a third-party website
// and linked to one or more Stories. Clubhouse currently supports
// linking files from Google Drive, Dropbox, Box, and by URL.
type LinkedFile struct {
	ContentType  string    `json:"content_type"`
	CreatedAt    time.Time `json:"created_at"`
	Description  string    `json:"description"`
	EntityType   string    `json:"entity_type"`
	ID           int       `json:"id"`
	MentionIDs   []string  `json:"mention_ids"`
	Name         string    `json:"name"`
	Size         int       `json:"size"`
	StoryIDs     []int     `json:"story_ids"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Type         string    `json:"type"`
	UpdatedAt    time.Time `json:"updated_at"`
	UploaderID   string    `json:"uploader_id"`
	URL          string    `json:"url"`
}

// Member represents details about individual Clubhouse user within the
// Clubhouse organization that has issued the token.
type Member struct {
	CreatedAt  time.Time `json:"created_at"`
	Disabled   bool      `json:"disabled"`
	EntityType string    `json:"entity_type"`
	ID         string    `json:"id"`
	Profile    Profile   `json:"profile"`
	Role       string    `json:"role"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// MakeURL ...
func (m Member) MakeURL() string {
	if m.ID == "" {
		return "members"
	}
	return path.Join("members", m.ID)
}

// Members ...
type Members []Member

// MakeURL ...
func (m Members) MakeURL() string {
	return "members"
}

// Milestone is a collection of Epics that represent a release or some
// other large initiative that your organization is working on.
type Milestone struct {
	Categories          []Category `json:"categories"`
	Completed           bool       `json:"completed"`
	CompletedAt         time.Time  `json:"completed_at"`
	CompletedAtOverride time.Time  `json:"completed_at_override"`
	Description         string     `json:"description"`
	EntityType          string     `json:"entity_type"`
	ID                  int        `json:"id"`
	Name                string     `json:"name"`
	Position            int        `json:"position"`
	Started             bool       `json:"started"`
	StartedAt           time.Time  `json:"started_at"`
	StartedAtOverride   time.Time  `json:"started_at_override"`
	State               State      `json:"state"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// MakeURL ...
func (m Milestone) MakeURL() string {
	if m.ID == 0 {
		return "milestones"
	}
	return path.Join("milestones", strconv.Itoa(m.ID))
}

// Milestones ...
type Milestones []Milestone

// MakeURL ...
func (m Milestones) MakeURL() string {
	return "milestones"
}

// CreateMilestoneParams ...
type CreateMilestoneParams struct {
	Categories          []CreateCategoryParams `json:"categories,omitempty"`
	CompletedAtOverride *time.Time             `json:"completed_at_override,omitempty"`
	Description         string                 `json:"description,omitempty"`
	Name                string                 `json:"name,omitempty"`
	StartedAtOverride   *time.Time             `json:"started_at_override,omitempty"`
	State               State                  `json:"state,omitempty"`
}

// UpdateMilestoneParams ...
type UpdateMilestoneParams struct {
	AfterID             *int
	BeforeID            *int
	Categories          []CreateCategoryParams
	CompletedAtOverride *time.Time
	Description         *string
	Name                *string
	StartedAtOverride   *time.Time
	State               State
}
type updateMilestoneParamsResolved struct {
	AfterID             *int                   `json:"after_id,omitempty"`
	BeforeID            *int                   `json:"before_id,omitempty"`
	Categories          []CreateCategoryParams `json:"categories,omitempty"`
	CompletedAtOverride *json.RawMessage       `json:"completed_at_override,omitempty"`
	Description         *string                `json:"description,omitempty"`
	Name                *string                `json:"name,omitempty"`
	StartedAtOverride   *json.RawMessage       `json:"started_at_override,omitempty"`
	State               State                  `json:"state,omitempty"`
}

// MarshalJSON ...
func (p UpdateMilestoneParams) MarshalJSON() ([]byte, error) {
	out := updateMilestoneParamsResolved{
		AfterID:     p.AfterID,
		BeforeID:    p.BeforeID,
		Categories:  p.Categories,
		Description: p.Description,
		Name:        p.Name,
		State:       p.State,
	}
	nullable{{
		in:   p.CompletedAtOverride,
		out:  &out.CompletedAtOverride,
		null: func() bool { return p.CompletedAtOverride == ResetTime },
	}, {
		in:   p.StartedAtOverride,
		out:  &out.StartedAtOverride,
		null: func() bool { return p.StartedAtOverride == ResetTime },
	}}.Do()
	return json.Marshal(&out)
}

// Profile represents details about individual Clubhouse user’s profile
// within the Clubhouse organization that has issued the token.
type Profile struct {
	Deactivated            bool   `json:"deactivated"`
	DisplayIcon            Icon   `json:"display_icon"`
	EmailAddress           string `json:"email_address"`
	EntityType             string `json:"entity_type"`
	GravatarHash           string `json:"gravatar_hash"`
	ID                     string `json:"id"`
	MentionName            string `json:"mention_name"`
	Name                   string `json:"name"`
	TwoFactorAuthActivated bool   `json:"two_factor_auth_activated"`
}

// Project typically map to teams (such as Frontend, Backend, Mobile,
// Devops, etc) but can represent any open-ended product, component, or
// initiative.
type Project struct {
	Abbreviation      string       `json:"abbreviation"`
	Archived          bool         `json:"archived"`
	Color             string       `json:"color"`
	CreatedAt         time.Time    `json:"created_at"`
	DaysToThermometer int          `json:"days_to_thermometer"`
	Description       string       `json:"description"`
	EntityType        string       `json:"entity_type"`
	ExternalID        string       `json:"external_id"`
	FollowerIDs       []string     `json:"follower_ids"`
	ID                int          `json:"id"`
	IterationLength   int          `json:"iteration_length"`
	Name              string       `json:"name"`
	ShowThermometer   bool         `json:"show_thermometer"`
	StartTime         time.Time    `json:"start_time"`
	Stats             ProjectStats `json:"stats"`
	TeamID            int          `json:"team_id"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

// MakeURL ...
func (p Project) MakeURL() string {
	if p.ID == 0 {
		return "projects"
	}
	return path.Join("projects", strconv.Itoa(p.ID))
}

// Projects ...
type Projects []Project

// MakeURL ...
func (p Projects) MakeURL() string {
	return "projects"
}

// CreateProjectParams ...
type CreateProjectParams struct {
	Abbreviation    string     `json:"abbreviation,omitempty"`
	Color           string     `json:"color,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	Description     string     `json:"description,omitempty"`
	ExternalID      string     `json:"external_id,omitempty"`
	FollowerIDs     []string   `json:"follower_ids,omitempty"`
	IterationLength int        `json:"iteration_length,omitempty"`
	Name            string     `json:"name,omitempty"`
	StartTime       *time.Time `json:"start_time,omitempty"`
	TeamID          int        `json:"team_id,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}

// UpdateProjectParams ...
type UpdateProjectParams struct {
	Abbreviation      *string  `json:"abbreviation,omitempty"`
	Archived          *bool    `json:"archived,omitempty"`
	Color             *string  `json:"color,omitempty"`
	DaysToThermometer *int     `json:"days_to_thermometer,omitempty"`
	Description       *string  `json:"description,omitempty"`
	FollowerIDs       []string `json:"follower_ids,omitempty"`
	Name              *string  `json:"name,omitempty"`
	ShowThermometer   *bool    `json:"show_thermometer,omitempty"`
	TeamID            *int     `json:"team_id,omitempty"`
}

// ProjectStats represents a group of calculated values for an Project.
type ProjectStats struct {
	NumPoints  int `json:"num_points"`
	NumStories int `json:"num_stories"`
}

// PullRequest corresponds to a GitHub Pull Request attached to a Clubhouse story.
type PullRequest struct {
	BranchID       int       `json:"branch_id"`
	Closed         bool      `json:"closed"`
	CreatedAt      time.Time `json:"created_at"`
	EntityType     string    `json:"entity_type"`
	ID             int       `json:"id"`
	NumAdded       int       `json:"num_added"`
	NumCommits     int       `json:"num_commits"`
	NumRemoved     int       `json:"num_removed"`
	Number         int       `json:"number"`
	TargetBranchID int       `json:"target_branch_id"`
	Title          string    `json:"title"`
	UpdatedAt      time.Time `json:"updated_at"`
	URL            string    `json:"url"`
}

// Repository refers to a GitHub repository.
type Repository struct {
	CreatedAt  time.Time `json:"created_at"`
	EntityType string    `json:"entity_type"`
	ExternalID string    `json:"external_id"`
	FullName   string    `json:"full_name"`
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	UpdatedAt  time.Time `json:"updated_at"`
	URL        string    `json:"url"`
}

// MakeURL ...
func (r Repository) MakeURL() string {
	if r.ID == 0 {
		return "repositories"
	}
	return path.Join("repositories", strconv.Itoa(r.ID))
}

// Repositories ...
type Repositories []Repository

// MakeURL ...
func (m Repositories) MakeURL() string {
	return "repositories"
}

// SearchQuery ...
type SearchQuery struct {
	Raw           string
	Epic          string
	Estimate      int
	HasAttachment bool
	HasComment    bool
	HasDeadline   bool
	HasEpic       bool
	HasTask       bool
	ID            int
	IsArchived    bool
	IsBlocked     bool
	IsBlocker     bool
	IsDone        bool
	IsOverdue     bool
	IsStarted     bool
	IsUnestimated bool
	IsUnstarted   bool
	Label         []string
	Owner         []string
	Project       string
	Requester     string
	State         string
	Text          string
	Type          StoryType
	Inversions    SearchQueryInversions
}

// SearchQueryInversions ...
type SearchQueryInversions struct {
	Epic          []string
	Estimate      []int
	HasAttachment bool
	HasComment    bool
	HasDeadline   bool
	HasEpic       bool
	HasTask       bool
	ID            []int
	IsArchived    bool
	IsBlocked     bool
	IsBlocker     bool
	IsDone        bool
	IsOverdue     bool
	IsStarted     bool
	IsUnestimated bool
	IsUnstarted   bool
	Label         []string
	Owner         []string
	Project       []string
	Requester     []string
	State         []string
	Text          []string
	Type          []StoryType
}

// MarshalJSON ...
func (q SearchQuery) MarshalJSON() ([]byte, error) {
	if q.Raw != "" {
		return json.Marshal(q.Raw)
	}

	parts := []string{}
	if q.Epic != "" {
		parts = append(parts, fmt.Sprintf(`epic:"%s"`, q.Epic))
	}
	if q.Estimate != 0 {
		parts = append(parts, fmt.Sprintf(`estimate:%d`, q.Estimate))
	}
	if q.HasAttachment {
		parts = append(parts, "has:attachment")
	}
	if q.HasComment {
		parts = append(parts, "has:comment")
	}
	if q.HasDeadline {
		parts = append(parts, "has:deadline")
	}
	if q.HasEpic {
		parts = append(parts, "has:epic")
	}
	if q.HasTask {
		parts = append(parts, "has:task")
	}
	if q.ID != 0 {
		parts = append(parts, fmt.Sprintf("id:%d", q.ID))
	}
	if q.IsArchived {
		parts = append(parts, "is:archived")
	}
	if q.IsBlocked {
		parts = append(parts, "is:blocked")
	}
	if q.IsBlocker {
		parts = append(parts, "is:blocker")
	}
	if q.IsDone {
		parts = append(parts, "is:done")
	}
	if q.IsOverdue {
		parts = append(parts, "is:overdue")
	}
	if q.IsStarted {
		parts = append(parts, "is:started")
	}
	if q.IsUnestimated {
		parts = append(parts, "is:unestimated")
	}
	if q.IsUnstarted {
		parts = append(parts, "is:unstarted")
	}
	if q.Label != nil {
		for _, e := range q.Label {
			parts = append(parts, fmt.Sprintf(`label:"%s"`, e))
		}
	}
	if q.Owner != nil {
		for _, e := range q.Owner {
			parts = append(parts, fmt.Sprintf(`owner:"%s"`, e))
		}
	}
	if q.Project != "" {
		parts = append(parts, fmt.Sprintf(`project:"%s"`, q.Project))
	}
	if q.Requester != "" {
		parts = append(parts, fmt.Sprintf(`requester:"%s"`, q.Requester))
	}
	if q.State != "" {
		parts = append(parts, fmt.Sprintf(`state:"%s"`, q.State))
	}
	if q.Text != "" {
		parts = append(parts, fmt.Sprintf(`"%s"`, q.Text))
	}
	if q.Type != "" {
		parts = append(parts, fmt.Sprintf(`type:%s`, q.Type))
	}

	if q.Inversions.Epic != nil {
		for _, e := range q.Inversions.Epic {
			parts = append(parts, fmt.Sprintf(`-epic:"%s"`, e))
		}
	}
	if q.Inversions.Estimate != nil {
		for _, e := range q.Inversions.Estimate {
			parts = append(parts, fmt.Sprintf(`-estimate:%d`, e))
		}
	}
	if q.Inversions.HasAttachment {
		parts = append(parts, "-has:attachment")
	}
	if q.Inversions.HasComment {
		parts = append(parts, "-has:comment")
	}
	if q.Inversions.HasDeadline {
		parts = append(parts, "-has:deadline")
	}
	if q.Inversions.HasEpic {
		parts = append(parts, "-has:epic")
	}
	if q.Inversions.HasTask {
		parts = append(parts, "-has:task")
	}
	if q.Inversions.ID != nil {
		for _, e := range q.Inversions.ID {
			parts = append(parts, fmt.Sprintf(`-id:%d`, e))
		}
	}
	if q.Inversions.IsArchived {
		parts = append(parts, "-is:archived")
	}
	if q.Inversions.IsBlocked {
		parts = append(parts, "-is:blocked")
	}
	if q.Inversions.IsBlocker {
		parts = append(parts, "-is:blocker")
	}
	if q.Inversions.IsDone {
		parts = append(parts, "-is:done")
	}
	if q.Inversions.IsOverdue {
		parts = append(parts, "-is:overdue")
	}
	if q.Inversions.IsStarted {
		parts = append(parts, "-is:started")
	}
	if q.Inversions.IsUnestimated {
		parts = append(parts, "-is:unestimated")
	}
	if q.Inversions.IsUnstarted {
		parts = append(parts, "-is:unstarted")
	}
	if q.Inversions.Label != nil {
		for _, e := range q.Inversions.Label {
			parts = append(parts, fmt.Sprintf(`-label:"%s"`, e))
		}
	}
	if q.Inversions.Owner != nil {
		for _, e := range q.Inversions.Owner {
			parts = append(parts, fmt.Sprintf(`-owner:"%s"`, e))
		}
	}
	if q.Inversions.Project != nil {
		for _, e := range q.Inversions.Project {
			parts = append(parts, fmt.Sprintf(`-project:"%s"`, e))
		}
	}
	if q.Inversions.Requester != nil {
		for _, e := range q.Inversions.Requester {
			parts = append(parts, fmt.Sprintf(`-requester:"%s"`, e))
		}
	}
	if q.Inversions.State != nil {
		for _, e := range q.Inversions.State {
			parts = append(parts, fmt.Sprintf(`-state:"%s"`, e))
		}
	}
	if q.Inversions.Text != nil {
		for _, e := range q.Inversions.Text {
			parts = append(parts, fmt.Sprintf(`-"%s"`, e))
		}
	}
	if q.Inversions.Type != nil {
		for _, e := range q.Inversions.Type {
			parts = append(parts, fmt.Sprintf(`-type:%s`, e))
		}
	}
	return json.Marshal(strings.Join(parts, " "))
}

// SearchParams ...
type SearchParams struct {
	Next     string       `json:"next,omitempty"`
	PageSize int          `json:"page_size,omitempty"`
	Query    *SearchQuery `json:"query,omitempty"`
}

// SearchResults represents the results of the search query.
type SearchResults struct {
	Data  []StorySearch `json:"data"`
	Next  string        `json:"next"`
	Total int           `json:"total"`
}

// Story the standard unit of work in Clubhouse and represent individual
// features, bugs, and chores.
type Story struct {
	AppURL              string           `json:"app_url"`
	Archived            bool             `json:"archived"`
	Blocked             bool             `json:"blocked"`
	Blocker             bool             `json:"blocker"`
	Branches            []Branch         `json:"branches"`
	Comments            []Comment        `json:"comments"`
	Commits             []Commit         `json:"commits"`
	Completed           bool             `json:"completed"`
	CompletedAt         time.Time        `json:"completed_at"`
	CompletedAtOverride time.Time        `json:"completed_at_override"`
	CreatedAt           time.Time        `json:"created_at"`
	Deadline            time.Time        `json:"deadline"`
	Description         string           `json:"description"`
	EntityType          string           `json:"entity_type"`
	EpicID              int              `json:"epic_id"`
	Estimate            int              `json:"estimate"`
	ExternalID          string           `json:"external_id"`
	Files               []File           `json:"files"`
	FollowerIDs         []string         `json:"follower_ids"`
	ID                  int              `json:"id"`
	Labels              []Label          `json:"labels"`
	LinkedFiles         []LinkedFile     `json:"linked_files"`
	MovedAt             time.Time        `json:"moved_at"`
	Name                string           `json:"name"`
	OwnerIDs            []string         `json:"owner_ids"`
	Position            int              `json:"position"`
	ProjectID           int              `json:"project_id"`
	RequestedByID       string           `json:"requested_by_id"`
	Started             bool             `json:"started"`
	StartedAt           time.Time        `json:"started_at"`
	StartedAtOverride   time.Time        `json:"started_at_override"`
	StoryLinks          []TypedStoryLink `json:"story_links"`
	StoryType           StoryType        `json:"story_type"`
	Tasks               []Task           `json:"tasks"`
	UpdatedAt           time.Time        `json:"updated_at"`
	WorflowStateID      int              `json:"worflow_state_id"`
}

// MakeURL ...
func (s *Story) MakeURL() string {
	if s.ID == 0 && s.Name == "" {
		return "stories"
	}
	return path.Join("stories", strconv.Itoa(s.ID))
}

// Stories is a Story slice
type Stories []Story

// MakeURL ...
func (s *Stories) MakeURL() string {
	return "stories"
}

// StoryLink represents a semantic relationships between two
// stories. Relationship types are relates to, blocks / blocked by, and
// duplicates / is duplicated by. The format is subject -> link ->
// object, or for example “story 5 blocks story 6”.
type StoryLink struct {
	CreatedAt  time.Time `json:"created_at"`
	EntityType string    `json:"entity_type"`
	ID         int       `json:"id"`
	ObjectID   int       `json:"object_id"`
	SubjectID  int       `json:"subject_id"`
	UpdatedAt  time.Time `json:"updated_at"`
	Verb       StoryVerb `json:"verb"`
}

// StorySearch ...
type StorySearch struct {
	AppURL              string           `json:"app_url"`
	Archived            bool             `json:"archived"`
	Blocked             bool             `json:"blocked"`
	Blocker             bool             `json:"blocker"`
	Completed           bool             `json:"completed"`
	CompletedAt         time.Time        `json:"completed_at"`
	CompletedAtOverride time.Time        `json:"completed_at_override"`
	CreatedAt           time.Time        `json:"created_at"`
	Deadline            time.Time        `json:"deadline"`
	Description         string           `json:"description"`
	EntityType          string           `json:"entity_type"`
	EpicID              int              `json:"epic_id"`
	Estimate            int              `json:"estimate"`
	ExternalID          string           `json:"external_id"`
	FollowerIDs         []string         `json:"follower_ids"`
	ID                  int              `json:"id"`
	Labels              []Label          `json:"labels"`
	MovedAt             time.Time        `json:"moved_at"`
	Name                string           `json:"name"`
	OwnerIDs            []string         `json:"owner_ids"`
	Position            int              `json:"position"`
	ProjectID           int              `json:"project_id"`
	RequestedByID       string           `json:"requested_by_id"`
	Started             bool             `json:"started"`
	StartedAt           time.Time        `json:"started_at"`
	StartedAtOverride   time.Time        `json:"started_at_override"`
	StoryLinks          []TypedStoryLink `json:"story_links"`
	StoryType           StoryType        `json:"story_type"`
	UpdatedAt           time.Time        `json:"updated_at"`
	WorkflowStateID     int              `json:"workflow_state_id"`
}

// StorySlim is a pared down version of the Story resource.
type StorySlim struct {
	AppURL              string           `json:"app_url"`
	Archived            bool             `json:"archived"`
	Blocked             bool             `json:"blocked"`
	Blocker             bool             `json:"blocker"`
	CommentIDs          []int            `json:"comment_ids"`
	Completed           bool             `json:"completed"`
	CompletedAtOverride time.Time        `json:"completed_at_override"`
	CreatedAt           time.Time        `json:"created_at"`
	Deadline            time.Time        `json:"deadline"`
	EntityType          string           `json:"entity_type"`
	EpicID              int              `json:"epic_id"`
	Estimate            int              `json:"estimate"`
	ExternalID          string           `json:"external_id"`
	FileIDs             []int            `json:"file_ids"`
	FollowerIDs         []string         `json:"follower_ids"`
	ID                  int              `json:"id"`
	Labels              []Label          `json:"labels"`
	LinkedFileIDs       []int            `json:"linked_file_ids"`
	MovedAt             time.Time        `json:"moved_at"`
	Name                string           `json:"name"`
	OwnerIDs            []string         `json:"owner_ids"`
	Position            int              `json:"position"`
	ProjectID           int              `json:"project_id"`
	RequestedByID       string           `json:"requested_by_id"`
	Started             bool             `json:"started"`
	StartedAt           time.Time        `json:"started_at"`
	StartedAtOverride   time.Time        `json:"started_at_override"`
	StoryLinks          []TypedStoryLink `json:"story_links"`
	StoryType           StoryType        `json:"story_type"`
	TaskIDs             []int            `json:"task_ids"`
	UpdatedAt           time.Time        `json:"updated_at"`
	WorkflowStateID     int              `json:"workflow_state_id"`
}

// Task ...
type Task struct {
	Complete    bool      `json:"complete"`
	CompletedAt time.Time `json:"completed_at"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	EntityType  string    `json:"entity_type"`
	ExternalID  string    `json:"external_id"`
	ID          int       `json:"id"`
	MentionIDs  []string  `json:"mention_ids"`
	OwnerIDs    []string  `json:"owner_ids"`
	Position    int       `json:"position"`
	StoryID     int       `json:"story_id"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Team contains a group of projects within the same Workspace
type Team struct {
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	EntityType  string    `json:"entity_type"`
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Position    int       `json:"position"`
	ProjectIDs  []int     `json:"project_ids"`
	UpdatedAt   time.Time `json:"updated_at"`
	Workflow    Workflow  `json:"workflow"`
}

// ThreadedComment represents Comments associated with Epic Discussions.
type ThreadedComment struct {
	parent     Resource
	AuthorID   string            `json:"author_id"`
	Comments   []ThreadedComment `json:"comments"`
	CreatedAt  time.Time         `json:"created_at"`
	Deleted    bool              `json:"deleted"`
	EntityType string            `json:"entity_type"`
	ExternalID string            `json:"external_ids"`
	ID         int               `json:"id"`
	MentionIDs []string          `json:"mention_ids"`
	Text       string            `json:"text"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

// MakeURL ...
func (c ThreadedComment) MakeURL() string {
	base := c.parent.MakeURL()
	if c.ID == 0 {
		return path.Join(base, "comments")
	}
	return path.Join(base, "comments", strconv.Itoa(c.ID))
}

// ThreadedComments ...
type ThreadedComments []ThreadedComment

// TypedStoryLink represents the type of Story Link. The string can be
// subject or object.
type TypedStoryLink struct {
	CreatedAt  time.Time `json:"created_at"`
	EntityType string    `json:"entity_type"`
	ID         int       `json:"id"`
	ObjectID   int       `json:"object_id"`
	SubjectID  int       `json:"subject_id"`
	Type       string    `json:"type"`
	UpdatedAt  time.Time `json:"updated_at"`
	Verb       string    `json:"verb"`
}

// Workflow is the array of defined Workflow States. Workflow can be
// queried using the API but must be updated in the Clubhouse UI.
type Workflow struct {
	CreatedAt      time.Time       `json:"created_at"`
	DefaultStateID int             `json:"default_state_id"`
	Description    string          `json:"description"`
	EntityType     string          `json:"entity_type"`
	ID             int             `json:"id"`
	Name           string          `json:"name"`
	States         []WorkflowState `json:"states"`
	TeamID         int             `json:"team_id"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// WorkflowState is any of the at least 3 columns. Workflow States
// correspond to one of 3 types: Unstarted, Started, or Done.
type WorkflowState struct {
	Color       string    `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	EntityType  string    `json:"entity_type"`
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	NumStories  int       `json:"num_stories"`
	Position    int       `json:"position"`
	Type        string    `json:"type"`
	UpdatedAt   time.Time `json:"updated_at"`
	Verb        string    `json:"verb"`
}
