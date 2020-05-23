// models.project.go

package main

import (
	"errors"
)

// To avoid scanning DynamoDB each time we need to show project info, we'll load all to memory into an array of projects called projectList
var projectList = []ProjectExample{
}

// Return a list of all the projects
func getAllProjects() []ProjectExample {
	return projectList
}

// Fetch an project based on the ID supplied
func getProjectByID(id int) (*ProjectExample, error) {
	for _, a := range projectList {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, errors.New("Project not found")
}

// Load a new project with the title and content provided
func loadNewProject(id int, title string, owner string, content string, photo string, votes int) (*ProjectExample, error) {
	// Set the ID of a new project to one more than the number of projects
	a := ProjectExample{ID: id, Owner: owner, Title: title, Content: content, Photo: photo, Votes: votes}

	// Add the project to the list of projects, making sure each ID appears only once
	counter := 0
	//projectList = nil

	for _, a := range projectList {
		if a.ID == id {
			counter++
		}
	}
	if counter < 1 {
			projectList = append(projectList, a)
		} else {
			//update item

		}

	return &a, nil
}
