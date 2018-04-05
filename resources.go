package clubhouse

import "time"

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
	Name       string       `json:"name"`
	Type       CategoryType `json:"type"`
}

// CreateCommentParams represents request parameters for creating a
// Comment on a Clubhouse Story.
type CreateCommentParams struct {
	AuthorID   string    `json:"author_id"`
	CreatedAt  time.Time `json:"created_at"`
	ExternalID string    `json:"external_id"`
	Text       string    `json:"text"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateLabelParams represents request parameters for creating a Label
// on a Clubhouse story.
type CreateLabelParams struct {
	Color      string `json:"color,omitempty"`
	ExternalID string `json:"external_id,omitempty"`
	Name       string `json:"name"`
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
	ObjectID  int       `json:"object_id"`
	SubjectID int       `json:"subject_id"`
	Verb      StoryVerb `json:"verb"`
}

// StoryType represents the type of story
type StoryType string

// Valid states for StoryType
const (
	TypeBug     StoryType = "bug"
	TypeChore             = "chore"
	TypeFeature           = "feature"
)

// CreateStoryParams is used to create multiple stories in a single
// request.
type CreateStoryParams struct {
	Comments            []CreateCommentParams   `json:"comments"`
	CompletedAtOverride time.Time               `json:"completed_at_override"`
	CreatedAt           time.Time               `json:"created_at"`
	Deadline            time.Time               `json:"deadline"`
	Description         string                  `json:"description"`
	EpicID              int                     `json:"epic_id"`
	Estimate            int                     `json:"estimate"`
	ExternalID          string                  `json:"external_id"`
	FileIDs             []string                `json:"follower_ids"`
	Labels              []CreateLabelParams     `json:"labels"`
	LinkedFileIDs       []int                   `json:"linked_file_ids"`
	Name                []int                   `json:"name"`
	OwnerIDs            []string                `json:"owner_ids"`
	ProjectID           int                     `json:"project_id"`
	RequestedByID       string                  `json:"requested_by_id"`
	StartedAtOverride   time.Time               `json:"started_at_override"`
	StoryLinks          []CreateStoryLinkParams `json:"story_links"`
	StoryType           StoryType               `json:"story_type"`
	Tasks               []CreateTaskParams      `json:"tasks"`
	UpdatedAt           time.Time               `json:"updated_at"`
	WorkflowStateID     int                     `json:"workflow_state_id"`
}

// CreateTaskParams request parameters for creating a Task on a Story.
type CreateTaskParams struct {
	Complete    bool      `json:"complete"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	ExternalID  string    `json:"external_id"`
	OwnerIDs    []string  `json:"owner_ids"`
	UpdatedAt   time.Time `json:"updated_at"`
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
	State               string            `json:"state"`
	Stats               EpicStats         `json:"stats"`
	UpdatedAt           time.Time         `json:"updated_at"`
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
	State               string     `json:"state"`
	UpdatedAt           time.Time  `json:"updated_at"`
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
	Completed           time.Time        `json:"completed"`
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
	ID          string    `json:"id"`
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
