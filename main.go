// main.go

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

var router *gin.Engine

// CognitoExample holds internals for auth flow.
type CognitoExample struct {
	CognitoClient *cognito.CognitoIdentityProvider
	RegFlow       *regFlow
	UserPoolID    string
	AppClientID   string
}

type regFlow struct {
	Username string
}

// ProjectExample is a Struct of a Project, including current Votes, , saved in a DynamoDB Votes
type ProjectExample struct {
	ID      int    `json:"id"`
	Title   string `json:"title" validate:"required"`
	Owner   string `json:"owner" validate:"required,email"`
	Content string `json:"content" validate:"required"`
	Photo   string `json:"photo"`
	Votes   int    `json:"votes"`
}

// LoggedInUserEmail will help us identify who is logged in, for when permissions and voting need checks
var LoggedInUserEmail string = "anonymous"

// CurrentVotes represents the votes that the logged in used made, saved in a DynamoDB Users
type CurrentVotes struct {
	Owner  string
	Voted1 int
	Voted2 int
	Voted3 int
	Voted4 int
	Voted5 int
}

func main() {
	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	var StaticContent string

	StaticContent = "/go/src/unicorn/templates/*"
	//StaticContent = "templates/*"

	// Set the router as the default one provided by Gin
	router = gin.Default()

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	router.LoadHTMLGlob(StaticContent)

	// Initialize the routes
	initializeRoutes()

	// Start serving the application
	router.Run()
}

// Render one of HTML, JSON or CSV based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that
// the template name is present
func render(c *gin.Context, data gin.H, templateName string) {
	loggedInInterface, _ := c.Get("is_logged_in")
	data["is_logged_in"] = loggedInInterface.(bool)

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		// Respond with XML
		c.XML(http.StatusOK, data["payload"])
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}
}
