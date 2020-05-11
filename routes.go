// routes.go

package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/ssm"
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

	// Get value of Cognito Application Client ID
	ssmsvc := ssm.New(sess, aws.NewConfig().WithRegion("eu-west-1"))

	ClientIDKey := "CognitoAppClientID"
	withDecryption := true
	param, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           &ClientIDKey,
		WithDecryption: &withDecryption,
	})

	// ClientIDValue stores the value of Cognito Application Client ID
	ClientIDValue := *param.Parameter.Value
	fmt.Println(ClientIDValue)

	// Define object with Cognito Stuff
	ce := CognitoExample{
		CognitoClient: cognito.New(sess),
		RegFlow:       &regFlow{},
		UserPoolID:    "CognitoUnicornUserPool",
		AppClientID:   ClientIDValue,
	}

	//Group Global routes (About, Leaderboard)
	globalRoutes := router.Group("/g")
	{
		// Handle the GET requests at /g/about
		// Ensure that the user is logged in by using the middleware
		globalRoutes.GET("/leaderboard", showAboutPage)

		// Handle POST requests at /g/leaderboard
		// Ensure that the user is logged in by using the middleware
		globalRoutes.GET("/about", ensureLoggedIn(), showLeaderboardPage)
	}

	// Group User routes together (Login, Register)
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
		userRoutes.POST("/register", ensureNotLoggedIn(), register(ce))

		// Handle the GET requests at /u/otp
		// Show the login page
		// Ensure that the user is not logged in by using the middleware
		userRoutes.GET("/otp", ensureNotLoggedIn(), showOTPPage)

		// Handle POST requests at /u/otp
		// Ensure that the user is not logged in by using the middleware
		userRoutes.POST("/otp", ensureNotLoggedIn(), OTP(ce))
	}

	// Group Project routes (View, Create)
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
