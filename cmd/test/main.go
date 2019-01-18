package main

import (
	"fmt"
	"log"
	"os"

	"github.com/haleyrc/todoist"
)

func main() {
	c := todoist.NewClient(os.Getenv("TODOIST_KEY"))
	projects, err := c.AllProjects()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%d\n", len(projects))

	p, err := projects.FindName("Work")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%#v\n", p)

	work, err := c.GetProject(p.ID)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%#v\n", work)

	tasks, err := c.ActiveTasks()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%d\n", len(tasks))

	workTasks, err := c.ActiveTasks(todoist.WithProjectID(2191185494))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%d\n", len(workTasks))

	p1Tasks, err := c.ActiveTasks(todoist.WithRawFilter("p1"))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%#v\n", p1Tasks)

	task, err := c.GetTask(3006570584)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%#v\n", task)

	labelTasks, err := c.ActiveTasks(todoist.WithLabelID(2151777880))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%#v\n", labelTasks)

	labels, err := c.AllLabels()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%d\n", len(labels))

	label, err := c.GetLabel(2151777880)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%#v\n", label)
}
