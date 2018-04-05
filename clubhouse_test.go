package clubhouse

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	check  = flag.Bool("check", false, "check the new snapshots")
	update = flag.Bool("update", false, "update test snapshot")
)

func TestListCategories(t *testing.T) {
	c := makeClient()
	categories, err := c.ListCategories()
	if err != nil {
		t.Error("did not expect error", err)
	}
	output := fmt.Sprintf("len: %d; entity type: %s",
		len(categories), categories[0].EntityType)
	snapshot(t, "ListCategories", output)
}

func TestGetCategory(t *testing.T) {
	c := makeClient()
	knownID := 17
	category, err := c.GetCategory(knownID)
	if err != nil {
		t.Error("did not expect error", err)
	}
	out := func(c *Category) string {
		return fmt.Sprintf("name:%#v archived:%#v color:%#v",
			c.Name, c.Archived, c.Color)
	}
	snapshot(t, "GetCategory", out(category))
}

func TestUpdateCategory(t *testing.T) {
	c := makeClient()
	knownID := 17

	newColor, err := c.UpdateCategory(knownID, UpdateCategoryParams{
		Color:    Color("#00ff00"),
		Archived: Archived,
	})
	if err != nil {
		t.Error("did not expect error", err)
	}
	if newColor.Color != "#00ff00" {
		t.Error("color didn't take", newColor.Color)
	}
	if !newColor.Archived {
		t.Error("should have archived newColor")
	}

	newArchive, err := c.UpdateCategory(knownID, UpdateCategoryParams{
		Archived: Unarchived,
	})
	if err != nil {
		t.Error("did not expect error", err)
	}
	if newArchive.Archived != false {
		t.Error("archive didn't take, should be false:", newArchive.Archived)
	}
	if newArchive.Color != "#00ff00" {
		t.Error("color didn't stick through archive", newArchive.Color)
	}

	resetColor, err := c.UpdateCategory(knownID, UpdateCategoryParams{
		Color: ResetColor,
	})
	if err != nil {
		t.Error("did not expect error", err)
	}
	if resetColor.Color != "" {
		t.Error("resetting color didn't take", resetColor.Color)
	}

	out := func(c *Category) string {
		return fmt.Sprintf("name:%#v archived:%#v color:%#v",
			c.Name, c.Archived, c.Color)
	}

	snapshot(t, "UpdateCategory_newColor", out(newColor))
	snapshot(t, "UpdateCategory_newArchive", out(newArchive))
	snapshot(t, "UpdateCategory_resetColor", out(resetColor))
}

func TestCreateAndDeleteCategory(t *testing.T) {
	c := makeClient()
	newcat, err := c.CreateCategory(CreateCategoryParams{
		Name:  "Hammes Fistkicker",
		Color: "powerful",
	})
	if err != nil {
		t.Error("did not expect error", err)
	}

	if newcat.Color != "powerful" {
		t.Error("color didn't take", newcat.Color)
	}
	if newcat.Name != "Hammes Fistkicker" {
		t.Error("name didn't take", newcat.Name)
	}

	if err := c.DeleteCategory(newcat.ID); err != nil {
		t.Error("did not expect error when deleting", err)
	}

	category, err := c.GetCategory(newcat.ID)
	if err == nil {
		t.Error("*expected* error trying to find category", category)
	}

	interr := err.(ErrClientRequest)
	if interr.Err != ErrResourceNotFound {
		t.Error("expected a 404 error", category)
	}
}

func TestListEpics(t *testing.T) {
	c := makeClient()
	epics, err := c.ListEpics()
	if err != nil {
		t.Fatal("did not expect error", err)
	}
	if len(epics) == 0 {
		t.Fatal("expected some epics")
	}
	if epics[0].EntityType != "epic" {
		t.Fatal("expected entity type to be epic")
	}
	output := fmt.Sprintf("len: %d; entity type: %s",
		len(epics), epics[0].EntityType)
	snapshot(t, "ListEpics", output)
}

func TestUpdateEpicParams(t *testing.T) {
	testTime, _ := time.Parse(time.RFC3339, "2018-04-20T16:20:00+04:00")

	type fieldtest struct {
		Name   string
		Params UpdateEpicParams
		Expect string
	}
	fieldtests := []fieldtest{
		{
			Name:   "empty",
			Params: UpdateEpicParams{},
			Expect: "{}",
		},
		{
			Name:   "AfterID",
			Params: UpdateEpicParams{AfterID: ID(10)},
			Expect: `{"after_id":10}`,
		},
		{
			Name:   "BeforeID",
			Params: UpdateEpicParams{BeforeID: ID(113)},
			Expect: `{"before_id":113}`,
		},
		{
			Name:   "Archived: unarchived",
			Params: UpdateEpicParams{Archived: Unarchived},
			Expect: `{"archived":false}`,
		},
		{
			Name:   "Archived: archived",
			Params: UpdateEpicParams{Archived: Archived},
			Expect: `{"archived":true}`,
		},
		{
			Name:   "Deadline: time set",
			Params: UpdateEpicParams{Deadline: &testTime},
			Expect: `{"deadline":"2018-04-20T16:20:00+04:00"}`,
		},
		{
			Name:   "Deadline: reset (null)",
			Params: UpdateEpicParams{Deadline: ResetTime},
			Expect: `{"deadline":null}`,
		},
		{
			Name:   "CompletedAtOverride: time set",
			Params: UpdateEpicParams{CompletedAtOverride: &testTime},
			Expect: `{"completed_at_override":"2018-04-20T16:20:00+04:00"}`,
		},
		{
			Name:   "CompletedAtOverride: reset (null)",
			Params: UpdateEpicParams{CompletedAtOverride: ResetTime},
			Expect: `{"completed_at_override":null}`,
		},
		{
			Name:   "Description: value",
			Params: UpdateEpicParams{Description: String("oh hey!")},
			Expect: `{"description":"oh hey!"}`,
		},
		{
			Name:   "Description: empty value",
			Params: UpdateEpicParams{Description: EmptyString},
			Expect: `{"description":""}`,
		},
		{
			Name:   "FollowerIDs",
			Params: UpdateEpicParams{FollowerIDs: []string{"1", "2"}},
			Expect: `{"follower_ids":["1","2"]}`,
		},
		{
			Name:   "Labels",
			Params: UpdateEpicParams{Labels: []CreateLabelParams{{Name: "hi"}}},
			Expect: `{"labels":[{"name":"hi"}]}`,
		},
		{
			Name:   "MilestoneID",
			Params: UpdateEpicParams{MilestoneID: ID(124)},
			Expect: `{"milestone_id":124}`,
		},
		{
			Name:   "MilestoneID: reset",
			Params: UpdateEpicParams{MilestoneID: ResetID},
			Expect: `{"milestone_id":null}`,
		},
		{
			Name:   "Name",
			Params: UpdateEpicParams{Name: "steven"},
			Expect: `{"name":"steven"}`,
		},
		{
			Name:   "OwnerIDs",
			Params: UpdateEpicParams{OwnerIDs: []string{"karen", "georgia"}},
			Expect: `{"owner_ids":["karen","georgia"]}`,
		},
		{
			Name:   "StartedAtOverride: time set",
			Params: UpdateEpicParams{StartedAtOverride: &testTime},
			Expect: `{"started_at_override":"2018-04-20T16:20:00+04:00"}`,
		},
		{
			Name:   "StartedAtOverride: reset (null)",
			Params: UpdateEpicParams{StartedAtOverride: ResetTime},
			Expect: `{"started_at_override":null}`,
		},
		{
			Name:   "State",
			Params: UpdateEpicParams{State: "hi"},
			Expect: `{"state":"hi"}`,
		},
	}

	for _, test := range fieldtests {
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

func TestCRUDEpics(t *testing.T) {
	c := makeClient()
	name := "new test epic"
	label := CreateLabelParams{Color: "red", Name: "test epic label"}
	var epicID int
	t.Run("create", func(t *testing.T) {
		epic, err := c.CreateEpic(CreateEpicParams{
			Name:      "new test epic",
			CreatedAt: Time(time.Now()),
			State:     EpicStateInProgress,
			Labels: []CreateLabelParams{
				label,
			},
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
		_, err := c.UpdateEpic(epicID, UpdateEpicParams{})
		if err != nil {
			t.Fatal("UpdateEpic: did not expect error updating", err)
		}
	})
	t.Run("delete", func(t *testing.T) {
		if epicID == 0 {
			t.Fatal("DeleteEpic: create must have failed")
		}

		if err := c.DeleteEpic(epicID); err != nil {
			t.Error("DeleteEpic: couldn't delete epic", err)
		}
	})
}

/* helpers */

func snapshot(t *testing.T, name string, obj interface{}) {
	got := fmt.Sprintf("%v", obj)
	filename := filepath.Join("testdata", name+".snapshot")
	if *update {
		fmt.Printf("Updating snapshot %s\n", name)
		file, err := os.Create(filename)
		defer file.Close()
		if err != nil {
			panic(fmt.Errorf("could not create file %s: %s", name, err))
		}
		file.Write([]byte(got))
		return
	}

	if *check {
		fmt.Printf("%s: %s\n", name, got)
		return
	}

	expect, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Errorf("could not read file %s: %s", name, err))
	}

	if string(expect) != string(got) {
		t.Errorf("mismatch:\n%s \n!= \n%s", got, expect)
	}
}

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
