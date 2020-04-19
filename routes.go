// routes.go

package main

func initializeRoutes() {

	// Use the setUserStatus middleware for every route to set a flag
	// indicating whether the request was from an authenticated user or not
	router.Use(setUserStatus())

	// Handle the index route
	router.GET("/", showIndexPage)

	// Group user related routes together
	userRoutes := router.Group("/u")
	{
		// Handle the GET requests at /u/login
		// Show the login page
		// Ensure that the user is not logged in by using the middleware
		userRoutes.GET("/login", ensureNotLoggedIn(), showLoginPage)

		// Handle POST requests at /u/login
		// Ensure that the user is not logged in by using the middleware
		userRoutes.POST("/login", ensureNotLoggedIn(), performLogin)

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
