// models.project.go

package main

import "errors"

type project struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Photo   string `json:"photo"`
	Votes   int    `json:"votes"`
}

// For this demo, we're storing the project list in memory
// In a real application, this list will most likely be fetched
// from a database or from static files
var projectList = []project{
	project{ID: 1, Title: "Flying Gas Stations", Content: "We can build them in the air... sky is the limit!", Photo: "none", Votes: 0},
	project{ID: 2, Title: "Corporate Gym Spinning Class creating energy", Content: "Just imagine the potential...", Photo: "none", Votes: 0},
}

// Return a list of all the projects
func getAllProjects() []project {
	return projectList
}

// Fetch an project based on the ID supplied
func getProjectByID(id int) (*project, error) {
	for _, a := range projectList {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, errors.New("Project not found")
}

// Create a new project with the title and content provided
func createNewProject(title, content, photo string) (*project, error) {
	// Set the ID of a new project to one more than the number of projects
	a := project{ID: len(projectList) + 1, Title: title, Content: content, Photo: photo, Votes: 0}

	// Add the project to the list of projects
	projectList = append(projectList, a)

	return &a, nil
}
