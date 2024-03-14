package todoist_test

import (
	"context"
	"flag"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/haleyrc/todoist"
)

var debug = flag.Bool("debug", false, "Enable verbose logging for the client")

func TestClient(t *testing.T) {
	ctx := context.Background()
	c := getTestClient(t)

	// CreateProject
	t.Log("creating project")
	project, err := c.CreateProject(ctx, todoist.ProjectParams{
		Name:       "SDK Test",
		Color:      todoist.ColorTaupe,
		IsFavorite: true,
		ViewStyle:  todoist.ViewStyleBoard,
	})
	if err != nil {
		t.Fatal(err)
	}

	// GetProject
	t.Log("getting project")
	time.Sleep(time.Second)
	projectResult, err := c.GetProject(ctx, project.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(project, projectResult) {
		t.Errorf(
			"Expected projects to be equal, but they weren't.\n\tWant: %#v\n\tGot:  %#v",
			project,
			projectResult,
		)
	}

	// GetProjects
	t.Log("getting projects")
	time.Sleep(time.Second)
	projectsResult, err := c.GetProjects(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if found := findProject(projectsResult, project.ID); found == nil {
		t.Errorf("Expected projects to contain %s, but they didn't.", project.Name)
	} else {
		if !reflect.DeepEqual(project, found) {
			t.Errorf(
				"Expected projects to be equal, but they weren't.\n\tWant: %#v\n\tGot:  %#v",
				project,
				found,
			)
		}
	}

	// CreateTask
	t.Log("creating task")
	time.Sleep(time.Second)
	task, err := c.CreateTask(ctx, todoist.TaskParams{
		Content:      "This is a test task",
		Description:  "This is the description.",
		Labels:       []string{"Buy"},
		Priority:     4,
		DueString:    "in 3 days at 9:30AM",
		Duration:     30,
		DurationUnit: "minute",
		ProjectID:    project.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	// CloseTask
	t.Log("closing task")
	time.Sleep(time.Second)
	if err := c.CloseTask(ctx, task.ID); err != nil {
		t.Fatal(err)
	}

	// ReopenTask
	t.Log("reopening task")
	time.Sleep(time.Second)
	if err := c.ReopenTask(ctx, task.ID); err != nil {
		t.Fatal(err)
	}

	// DeleteTask
	t.Log("deleting task")
	time.Sleep(time.Second)
	if err := c.DeleteTask(ctx, task.ID); err != nil {
		t.Fatal(err)
	}

	// DeleteProject
	t.Log("deleting project")
	time.Sleep(time.Second)
	if err := c.DeleteProject(ctx, project.ID); err != nil {
		t.Fatal(err)
	}
}

func getTestClient(t *testing.T) *todoist.Client {
	apiToken := os.Getenv("API_TOKEN")
	if apiToken == "" {
		t.Skip("set API_TOKEN to run this test")
	}

	c, err := todoist.NewClient(apiToken)
	if err != nil {
		t.Fatal(err)
	}

	c.Verbose = *debug

	return c
}

func findProject(projects []todoist.Project, id string) *todoist.Project {
	for _, project := range projects {
		if project.ID == id {
			return &project
		}
	}
	return nil
}
