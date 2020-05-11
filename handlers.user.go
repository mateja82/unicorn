// handlers.user.go

package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/aws/aws-sdk-go/aws"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

const flowUsernamePassword = "USER_PASSWORD_AUTH"
const flowRefreshToken = "REFRESH_TOKEN_AUTH"

func showLoginPage(c *gin.Context) {

	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "Login",
	}, "login.html")
}

func performLogin(ce CognitoExample) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Obtain the POSTed username and password values
		username := c.PostForm("username")
		password := c.PostForm("password")

		// Authentication Flow for Login to Cognito
		flow := aws.String(flowUsernamePassword)
		params := map[string]*string{
			"USERNAME": aws.String(username),
			"PASSWORD": aws.String(password),
		}

		authTry := &cognito.InitiateAuthInput{
			AuthFlow:       flow,
			AuthParameters: params,
			ClientId:       aws.String(ce.AppClientID),
		}

		result, err := ce.CognitoClient.InitiateAuth(authTry)

		if err != nil {
			// If the username/password combination is invalid,
			// show the error message on the login page
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"ErrorTitle":   "Login Failed",
				"ErrorMessage": "Invalid credentials provided"})
		} else {
			fmt.Println(result)
			// set Context to OK
			// Check if the username/password combination is valid
			//			if isUserValid(username, password) {
			// If the username/password is valid set the token in a cookie
			token := generateSessionToken()
			c.SetCookie("token", token, 3600, "", "", false, true)
			c.Set("is_logged_in", true)

			render(c, gin.H{
				"title": "Successful Login"}, "login-successful.html")
			return
			//}
		}
	}
	return gin.HandlerFunc(fn)

}

func generateSessionToken() string {
	// We're using a random 16 character string as the session token
	// This is NOT a secure way of generating session tokens
	// DO NOT USE THIS IN PRODUCTION
	return strconv.FormatInt(rand.Int63(), 16)
}

func logout(c *gin.Context) {

	// var sameSiteCookie http.SameSite

	// Clear the cookie
	c.SetCookie("token", "", -1, "", "", false, true)

	// Redirect to the home page
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func showRegistrationPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "Register"}, "register.html")
}

func register(c *gin.Context) {
	// Obtain the POSTed username and password values
	username := c.PostForm("username")
	password := c.PostForm("password")

	// var sameSiteCookie http.SameSite

	if _, err := registerNewUser(username, password); err == nil {
		// If the user is created, set the token in a cookie and log the user in
		token := generateSessionToken()
		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)

		render(c, gin.H{
			"title": "Successful registration & Login"}, "login-successful.html")

	} else {
		// If the username/password combination is invalid,
		// show the error message on the login page
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"ErrorTitle":   "Registration Failed",
			"ErrorMessage": err.Error()})

	}
}
