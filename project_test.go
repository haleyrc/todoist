package todoist_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/haleyrc/todoist"
)

func TestClient_GetProjects(t *testing.T) {
	ctx := context.Background()

	apiToken := os.Getenv("API_TOKEN")
	if apiToken == "" {
		t.Skip("set API_TOKEN to run this test")
	}

	c, err := todoist.NewClient(apiToken)
	if err != nil {
		t.Fatal(err)
	}

	projects, err := c.GetProjects(ctx)
	if err != nil {
		t.Fatal(err)
	}

	bytes, _ := json.MarshalIndent(projects, "", "  ")
	fmt.Println(string(bytes))
}
