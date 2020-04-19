// handlers.project.go

package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func showIndexPage(c *gin.Context) {
	projects := getAllProjects()

	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title":   "Home Page",
		"payload": projects}, "index.html")
}

func showProjectCreationPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "Create New Project"}, "create-project.html")
}

func getProject(c *gin.Context) {
	// Check if the project ID is valid
	if projectID, err := strconv.Atoi(c.Param("project_id")); err == nil {
		// Check if the project exists
		if project, err := getProjectByID(projectID); err == nil {
			// Call the render function with the title, project and the name of the
			// template
			render(c, gin.H{
				"title":   project.Title,
				"payload": project}, "project.html")

		} else {
			// If the project is not found, abort with an error
			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {
		// If an invalid project ID is specified in the URL, abort with an error
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func createProject(c *gin.Context) {
	// Obtain the POSTed title and content values
	title := c.PostForm("title")
	content := c.PostForm("content")
	photo := c.PostForm("photo")

	if a, err := createNewProject(title, content, photo); err == nil {
		// If the project is created successfully, show success message
		render(c, gin.H{
			"title":   "Submission Successful",
			"payload": a}, "submission-successful.html")
	} else {
		// if there was an error while creating the project, abort with an error
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
