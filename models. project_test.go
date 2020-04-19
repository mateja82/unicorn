// models.project_test.go

package main

import "testing"

// Test the function that fetches all projects
func TestGetAllProjects(t *testing.T) {
	alist := getAllProjects()

	// Check that the length of the list of projects returned is the
	// same as the length of the global variable holding the list
	if len(alist) != len(projectList) {
		t.Fail()
	}

	// Check that each member is identical
	for i, v := range alist {
		if v.Content != projectList[i].Content ||
			v.ID != projectList[i].ID ||
			v.Title != projectList[i].Title {

			t.Fail()
			break
		}
	}
}

// Test the function that fetche an Project by its ID
func TestGetProjectByID(t *testing.T) {
	a, err := getProjectByID(1)

	if err != nil || a.ID != 1 || a.Title != "Project 1" || a.Content != "Project 1 body" {
		t.Fail()
	}
}

// Test the function that creates a new project
func TestCreateNewProject(t *testing.T) {
	// get the original count of projects
	originalLength := len(getAllProjects())

	// add another project
	a, err := createNewProject("New test title", "New test content", "New S3 URL for Photo")

	// get the new count of projects
	allProjects := getAllProjects()
	newLength := len(allProjects)

	if err != nil || newLength != originalLength+1 ||
		a.Title != "New test title" || a.Content != "New test content" {

		t.Fail()
	}
}
