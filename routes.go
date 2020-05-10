// routes.go

package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func initializeRoutes() {

	// Use the setUserStatus middleware for every route to set a flag
	// indicating whether the request was from an authenticated user or not
	router.Use(setUserStatus())

	// Handle the index route
	router.GET("/", showIndexPage)

	// Create AWS Session
	conf := &aws.Config{Region: aws.String("eu-west-1")}
	sess, err := session.NewSession(conf)
	if err != nil {
		panic(err)
	}

	// Define object with Cognito Stuff
	ce := CognitoExample{
		CognitoClient: cognito.New(sess),
		RegFlow:       &regFlow{},
		UserPoolID:    "CognitoUnicornUserPool",
		AppClientID:   "4un2qodp09fojc5bm7ibb6a8u6",
	}

	// Group user related routes together
	userRoutes := router.Group("/u")
	{
		// Handle the GET requests at /u/login
		// Show the login page
		// Ensure that the user is not logged in by using the middleware
		userRoutes.GET("/login", ensureNotLoggedIn(), showLoginPage)

		// Handle POST requests at /u/login
		// Ensure that the user is not logged in by using the middleware
		userRoutes.POST("/login", ensureNotLoggedIn(), performLogin(ce))

		// Handle GET requests at /u/logout
		// Ensure that the user is logged in by using the middleware
		userRoutes.GET("/logout", ensureLoggedIn(), logout)

		// Handle the GET requests at /u/register
		// Show the registration page
		// Ensure that the user is not logged in by using the middleware
		userRoutes.GET("/register", ensureNotLoggedIn(), showRegistrationPage)

		// Handle POST requests at /u/register
		// Ensure that the user is not logged in by using the middleware
		userRoutes.POST("/register", ensureNotLoggedIn(), register)
	}

	// Group project related routes together
	projectRoutes := router.Group("/project")
	{
		// Handle GET requests at /project/view/some_project_id
		projectRoutes.GET("/view/:project_id", getProject)

		// Handle the GET requests at /project/create
		// Show the project creation page
		// Ensure that the user is logged in by using the middleware
		projectRoutes.GET("/create", ensureLoggedIn(), showProjectCreationPage)

		// Handle POST requests at /project/create
		// Ensure that the user is logged in by using the middleware
		projectRoutes.POST("/create", ensureLoggedIn(), createProject)
	}
}
