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
