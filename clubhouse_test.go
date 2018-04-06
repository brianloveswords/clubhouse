package clubhouse

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"
)

func TestCreateCategoryParams(t *testing.T) {
	fieldtest{{
		Name:   "empty",
		Params: CreateCategoryParams{},
		Expect: "{}",
	}, {
		Name:   "Name",
		Params: CreateCategoryParams{Name: "hey"},
		Expect: `{"name":"hey"}`,
	}, {
		Name:   "Type",
		Params: CreateCategoryParams{Type: "the type"},
		Expect: `{"type":"the type"}`,
	}, {
		Name:   "ExternalID",
		Params: CreateCategoryParams{ExternalID: "an ID"},
		Expect: `{"external_id":"an ID"}`,
	}, {
		Name:   "Color",
		Params: CreateCategoryParams{Color: "red"},
		Expect: `{"color":"red"}`,
	},
	}.Test(t)
}

func TestUpdateCategoryParams(t *testing.T) {
	fieldtest{{
		Name:   "empty",
		Params: UpdateCategoryParams{},
		Expect: "{}",
	}, {
		Name:   "Name",
		Params: UpdateCategoryParams{Name: String("hey")},
		Expect: `{"name":"hey"}`,
	}, {
		Name:   "ExternalID",
		Params: UpdateCategoryParams{Archived: Unarchived},
		Expect: `{"archived":false}`,
	}, {
		Name:   "Color",
		Params: UpdateCategoryParams{Color: String("red")},
		Expect: `{"color":"red"}`,
	}, {
		Name:   "Color: reset",
		Params: UpdateCategoryParams{Color: ResetColor},
		Expect: `{"color":null}`,
	},
	}.Test(t)

}

func TestCRUDCategories(t *testing.T) {
	var (
		c    = makeClient()
		cat  *Category
		cats []Category
		err  error
	)
	t.Run("create", func(_ *testing.T) {
		cat, err = c.CreateCategory(&CreateCategoryParams{
			Name:  fmt.Sprintf("%v", time.Now()),
			Color: "powerful",
		})
		if err != nil {
			fmt.Println(err)
			t.Fatal("did not expect error")
		}
	})
	t.Run("read", func(t *testing.T) {
		getcat, err := c.GetCategory(cat.ID)
		if err != nil {
			t.Fatal("did not expect error", err)
		}
		if getcat.Color != cat.Color {
			t.Errorf("color didn't stick, got %s", getcat.Color)
		}
	})
	t.Run("list", func(t *testing.T) {
		cats, err = c.ListCategories()
		if err != nil {
			t.Fatal("did not expect error", err)
		}
		if len(cats) == 0 {
			t.Fatal("expected to get some categories")
		}
	})
	t.Run("update", func(t *testing.T) {
		upcat, err := c.UpdateCategory(cat.ID, &UpdateCategoryParams{
			Color:    ResetColor,
			Archived: Archived,
		})
		if err != nil {
			t.Fatal("error updating category", err)
		}
		if upcat.Color != "" {
			t.Error("color is wrong")
		}
		if upcat.Archived != true {
			t.Error("should be archived")
		}
	})
	t.Run("delete", func(t *testing.T) {
		for _, category := range cats {
			if err := c.DeleteCategory(category.ID); err != nil {
				t.Fatal("did not expect error deleting category", err)
			}
		}
	})
}

func TestUpdateEpicParams(t *testing.T) {
	testTime, _ := time.Parse(time.RFC3339, "2018-04-20T16:20:00+04:00")
	fieldtest{{
		Name:   "empty",
		Params: UpdateEpicParams{},
		Expect: "{}",
	}, {
		Name:   "AfterID",
		Params: UpdateEpicParams{AfterID: ID(10)},
		Expect: `{"after_id":10}`,
	}, {
		Name:   "BeforeID",
		Params: UpdateEpicParams{BeforeID: ID(113)},
		Expect: `{"before_id":113}`,
	}, {
		Name:   "Archived: unarchived",
		Params: UpdateEpicParams{Archived: Unarchived},
		Expect: `{"archived":false}`,
	}, {
		Name:   "Archived: archived",
		Params: UpdateEpicParams{Archived: Archived},
		Expect: `{"archived":true}`,
	}, {
		Name:   "Deadline: time set",
		Params: UpdateEpicParams{Deadline: &testTime},
		Expect: `{"deadline":"2018-04-20T16:20:00+04:00"}`,
	}, {
		Name:   "Deadline: reset (null)",
		Params: UpdateEpicParams{Deadline: ResetTime},
		Expect: `{"deadline":null}`,
	}, {
		Name:   "CompletedAtOverride: time set",
		Params: UpdateEpicParams{CompletedAtOverride: &testTime},
		Expect: `{"completed_at_override":"2018-04-20T16:20:00+04:00"}`,
	}, {
		Name:   "CompletedAtOverride: reset (null)",
		Params: UpdateEpicParams{CompletedAtOverride: ResetTime},
		Expect: `{"completed_at_override":null}`,
	}, {
		Name:   "Description: value",
		Params: UpdateEpicParams{Description: String("oh hey!")},
		Expect: `{"description":"oh hey!"}`,
	}, {
		Name:   "Description: empty value",
		Params: UpdateEpicParams{Description: EmptyString},
		Expect: `{"description":""}`,
	}, {
		Name:   "FollowerIDs",
		Params: UpdateEpicParams{FollowerIDs: []string{"1", "2"}},
		Expect: `{"follower_ids":["1","2"]}`,
	}, {
		Name:   "Labels",
		Params: UpdateEpicParams{Labels: []CreateLabelParams{{Name: "hi"}}},
		Expect: `{"labels":[{"name":"hi"}]}`,
	}, {
		Name:   "MilestoneID",
		Params: UpdateEpicParams{MilestoneID: ID(124)},
		Expect: `{"milestone_id":124}`,
	}, {
		Name:   "MilestoneID: reset",
		Params: UpdateEpicParams{MilestoneID: ResetID},
		Expect: `{"milestone_id":null}`,
	}, {
		Name:   "Name",
		Params: UpdateEpicParams{Name: "steven"},
		Expect: `{"name":"steven"}`,
	}, {
		Name:   "OwnerIDs",
		Params: UpdateEpicParams{OwnerIDs: []string{"karen", "georgia"}},
		Expect: `{"owner_ids":["karen","georgia"]}`,
	}, {
		Name:   "StartedAtOverride: time set",
		Params: UpdateEpicParams{StartedAtOverride: &testTime},
		Expect: `{"started_at_override":"2018-04-20T16:20:00+04:00"}`,
	}, {
		Name:   "StartedAtOverride: reset (null)",
		Params: UpdateEpicParams{StartedAtOverride: ResetTime},
		Expect: `{"started_at_override":null}`,
	}, {
		Name:   "State",
		Params: UpdateEpicParams{State: "hi"},
		Expect: `{"state":"hi"}`,
	}}.Test(t)
}

func TestCRUDEpics(t *testing.T) {
	var (
		c      = makeClient()
		name   = "new test epic"
		label  = CreateLabelParams{Color: "red", Name: "test epic label"}
		err    error
		epicID int
		epics  []Epic
	)

	t.Run("create", func(t *testing.T) {
		epic, err := c.CreateEpic(&CreateEpicParams{
			Name:      "new test epic",
			CreatedAt: Time(time.Now()),
			State:     EpicStateInProgress,
			Labels:    []CreateLabelParams{label},
		})
		if err != nil {
			t.Fatal("CreateEpic: couldn't create", err)
		}
		if epic == nil {
			t.Fatal("CreateEpic: epic shouldn't be nil")
		}

		epicID = epic.ID

		if epic.Name != name {
			t.Errorf("CreateEpic: name didn't stick, %s != %s", epic.Name, name)
		}
		if epic.State != EpicStateInProgress {
			t.Errorf("CreateEpic: state didn't stick, %s != %s", epic.State, EpicStateInProgress)
		}
	})
	t.Run("list", func(t *testing.T) {
		epics, err = c.ListEpics()
		if err != nil {
			t.Fatal("did not expect error", err)
		}
		if len(epics) == 0 {
			t.Fatal("expected some epics")
		}
		if epics[0].EntityType != "epic" {
			t.Fatal("expected entity type to be epic")
		}
	})
	t.Run("read", func(t *testing.T) {
		epic, err := c.GetEpic(epicID)
		if err != nil {
			t.Fatal("GetEpic: couldn't create", err)
		}
		if epic.Name != name {
			t.Errorf("GetEpic: name didn't stick, %s != %s", epic.Name, name)
		}
		if len(epic.Labels) == 0 {
			t.Fatal("GetEpic: expected labels")
		}
		if epic.Labels[0].Name != label.Name {
			t.Error("GetEpic: label name didn't match")
		}
	})
	t.Run("update", func(t *testing.T) {
		_, err := c.UpdateEpic(epicID, UpdateEpicParams{
			Name: "a different name",
		})
		if err != nil {
			t.Fatal("UpdateEpic: did not expect error updating", err)
		}
	})
	t.Run("delete", func(t *testing.T) {
		if len(epics) == 0 {
			t.Fatal("DeleteEpic: create must have failed")
		}
		for _, e := range epics {
			if err := c.DeleteEpic(e.ID); err != nil {
				t.Error("DeleteEpic: couldn't delete epic", err)
			}
		}
	})
}

func TestCreateCommentParams(t *testing.T) {
	testTime, _ := time.Parse(time.RFC3339, "2018-04-20T16:20:00+04:00")
	fieldtest{{
		Name:   "empty",
		Params: CreateCommentParams{},
		Expect: "{}",
	}, {
		Name:   "AuthorID",
		Params: CreateCommentParams{AuthorID: "some author id"},
		Expect: `{"author_id":"some author id"}`,
	}, {
		Name:   "CreatedAt",
		Params: CreateCommentParams{CreatedAt: &testTime},
		Expect: `{"created_at":"2018-04-20T16:20:00+04:00"}`,
	}, {
		Name:   "ExternalID",
		Params: CreateCommentParams{ExternalID: "extid"},
		Expect: `{"external_id":"extid"}`,
	}, {
		Name:   "Text",
		Params: CreateCommentParams{Text: "angry comment"},
		Expect: `{"text":"angry comment"}`,
	}, {
		Name:   "UpdatedAt",
		Params: CreateCommentParams{UpdatedAt: &testTime},
		Expect: `{"updated_at":"2018-04-20T16:20:00+04:00"}`,
	},
	}.Test(t)
}

func TestCRUDEpicComments(t *testing.T) {
	// make an epic first.
	c := makeClient()
	epic, err := c.CreateEpic(&CreateEpicParams{
		Name: "test epic: comments",
	})
	if err != nil {
		t.Fatal("unexpected error making epic for comments")
	}
	epicID := epic.ID
	defer c.DeleteEpic(epicID)

	var commentID, replyID int
	text := "ur wrong"
	reply := "ur wrongerer"
	t.Run("create", func(t *testing.T) {
		comment, err := c.CreateEpicComment(epicID, &CreateCommentParams{
			Text: text,
		})
		if err != nil {
			t.Fatal("unexpected error making comment", err)
		}
		if comment.Text != text {
			t.Errorf("comment text didn't stick, expected %s got %s", text, comment.Text)
		}
		deepcomment, err := c.CreateEpicCommentComment(epicID, comment.ID, &CreateCommentParams{
			Text: reply,
		})
		if deepcomment.Text != reply {
			t.Errorf("comment text didn't stick, expected %s got %s", reply, deepcomment.Text)
		}

	})
	t.Run("list", func(t *testing.T) {
		comments, err := c.ListEpicComments(epicID)
		if err != nil {
			t.Fatal("unexpected error listing comments", err)
		}
		if len(comments) == 0 {
			t.Fatal("should have gotten at least one comment")
		}
		if comments[0].Text != text {
			t.Errorf("comment text didn't stick, expected %s got %s", text, comments[0].Text)
		}
		commentID = comments[0].ID
	})
	t.Run("read", func(t *testing.T) {
		comment, err := c.GetEpicComment(epicID, commentID)
		if err != nil {
			t.Fatal("unexpected error reading comment", err)
		}
		if comment.Text != text {
			t.Errorf("comment text didn't stick, expected %s got %s", text, comment.Text)
		}
		gotreply := comment.Comments[0].Text
		if gotreply != reply {
			t.Errorf("replytext didn't stick, expected %s got %s", reply, gotreply)
		}
		replyID = comment.Comments[0].ID
	})
	t.Run("update", func(t *testing.T) {
		updated := "n/m sorry"
		reply, err := c.UpdateEpicComment(
			epicID, commentID,
			&UpdateCommentParams{Text: updated},
		)
		if err != nil {
			t.Fatal("unexpected error updating comment", err)
		}
		if reply.Text != updated {
			t.Errorf("comment text didn't stick, expected %s got %s", updated, reply.Text)
		}
	})
	t.Run("delete", func(t *testing.T) {
		if err := c.DeleteEpicComment(epicID, commentID); err != nil {
			t.Fatal("unexpected error deleting comment", err)
		}
	})
}

func TestUpdateFileParams(t *testing.T) {
	testTime, _ := time.Parse(time.RFC3339, "2018-04-20T16:20:00+04:00")
	fieldtest{{
		Name:   "empty",
		Params: UpdateFileParams{},
		Expect: "{}",
	}, {
		Name:   "CreatedAt",
		Params: UpdateFileParams{CreatedAt: &testTime},
		Expect: `{"created_at":"2018-04-20T16:20:00+04:00"}`,
	}, {
		Name:   "Description",
		Params: UpdateFileParams{Description: String("")},
		Expect: `{"description":""}`,
	}, {
		Name:   "ExternalID",
		Params: UpdateFileParams{ExternalID: String("some id")},
		Expect: `{"external_id":"some id"}`,
	}, {
		Name:   "Name",
		Params: UpdateFileParams{Name: String("a name!")},
		Expect: `{"name":"a name!"}`,
	}, {
		Name:   "UpdatedAt",
		Params: UpdateFileParams{UpdatedAt: &testTime},
		Expect: `{"updated_at":"2018-04-20T16:20:00+04:00"}`,
	}, {
		Name:   "UploaderID",
		Params: UpdateFileParams{UploaderID: String("lkajlf")},
		Expect: `{"uploader_id":"lkajlf"}`,
	},
	}.Test(t)
}

func TestCRUDFiles(t *testing.T) {
	c := makeClient()
	f1, err := os.Open("testdata/test-file-1.txt")
	defer f1.Close()
	if (err) != nil {
		t.Fatal("unexpected error opening test-file-1.txt")
	}
	f2, err := os.Open("testdata/test-file-2.txt")
	defer f2.Close()
	if (err) != nil {
		t.Fatal("unexpected error opening test-file-2.txt")
	}

	var files, listed []File
	t.Run("create", func(t *testing.T) {
		files, err = c.UploadFiles([]FileUpload{
			{
				Name: "test-file-1.txt",
				File: f1,
			},
			{
				Name: "test-file-2.txt",
				File: f2,
			},
		})
		if err != nil {
			t.Fatal("unexpected error uploading file", err)
		}

		if len(files) < 2 {
			t.Fatal("expected 2 files, got", len(files))
		}
	})
	t.Run("list", func(t *testing.T) {
		listed, err = c.ListFiles()
		if err != nil {
			t.Fatal("unexpected error listing files", err)
		}

		if len(listed) < 2 {
			t.Fatal("expected 2 files, got", len(files))
		}
	})
	t.Run("read", func(t *testing.T) {
		file, err := c.GetFile(files[0].ID)
		if err != nil {
			t.Fatal("unexpected error getting file by id", err)
		}
		if file.Name != "test-file-1.txt" {
			t.Fatal("expected file name to be test-file-1.txt")
		}
	})
	t.Run("update", func(t *testing.T) {
		file, err := c.UpdateFile(files[0].ID, &UpdateFileParams{
			Name: String("cranberry"),
		})

		if err != nil {
			t.Fatal("error updating name")
		}
		if file.Name != "cranberry" {
			t.Error("expected name to update")
		}
	})
	t.Run("delete", func(t *testing.T) {
		for _, f := range listed {
			if err := c.DeleteFile(f.ID); err != nil {
				t.Fatal("unexpected error deleting file", err)
			}
		}
	})
}

/* helpers */

// func snapshot(t *testing.T, name string, obj interface{}) {
// 	got := fmt.Sprintf("%v", obj)
// 	filename := filepath.Join("testdata", name+".snapshot")
// 	if *update {
// 		fmt.Printf("Updating snapshot %s\n", name)
// 		file, err := os.Create(filename)
// 		defer file.Close()
// 		if err != nil {
// 			panic(fmt.Errorf("could not create file %s: %s", name, err))
// 		}
// 		file.Write([]byte(got))
// 		return
// 	}

// 	if *check {
// 		fmt.Printf("%s: %s\n", name, got)
// 		return
// 	}

// 	expect, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		panic(fmt.Errorf("could not read file %s: %s", name, err))
// 	}

// 	if string(expect) != string(got) {
// 		t.Errorf("mismatch:\n%s \n!= \n%s", got, expect)
// 	}
// }

type credentials struct {
	AuthToken string
}

func mustOpen(p string) io.Reader {
	file, err := os.Open(p)
	if err != nil {
		log.Fatal("could not open file", err)
	}
	return file
}

func loadCredentials() credentials {
	file := mustOpen("secrets.json")
	dec := json.NewDecoder(file)
	creds := credentials{}
	if err := dec.Decode(&creds); err != nil {
		log.Fatal("could not decode secrets.json", err)
	}
	return creds
}

func makeClient() *Client {
	creds := loadCredentials()
	return &Client{
		AuthToken: creds.AuthToken,
	}
}

type fieldtest []struct {
	Name   string
	Params interface{}
	Expect string
}

func (ft fieldtest) Test(t *testing.T) {
	for _, test := range ft {
		t.Run(test.Name, func(t *testing.T) {
			b, err := json.Marshal(&test.Params)
			if err != nil {
				t.Fatal("shouldn't get an error")
			}
			if test.Expect != string(b) {
				t.Errorf("%s != %s", string(b), test.Expect)
			}
		})
	}
}
