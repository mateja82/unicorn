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
				"ErrorMessage": "Invalid credentials"})
		} else {
			fmt.Println(result)
			// Set Token and ïs_logged_in Boolean, and jump to a different Menu
			// NEXT STEP: Import TOKEN variable from Login

			token := generateSessionToken()

			c.SetCookie("token", token, 3600, "", "", false, true)
			c.Set("is_logged_in", true)

			render(c, gin.H{
				"title": "Successful Login"}, "login-successful.html")
			return
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

func register(ce CognitoExample) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		// Obtain the POSTed username and password values
		username := c.PostForm("username")
		password := c.PostForm("password")
		fullname := c.PostForm("name")
		emailAddress := c.PostForm("email")
		phoneNumber := c.PostForm("phone_number")

		user := &cognito.SignUpInput{
			Username: aws.String(username),
			Password: aws.String(password),
			ClientId: aws.String(ce.AppClientID),
			UserAttributes: []*cognito.AttributeType{
				{
					Name:  aws.String("email"),
					Value: aws.String(emailAddress),
				},
				{
					Name:  aws.String("name"),
					Value: aws.String(fullname),
				},
				{
					Name:  aws.String("phone_number"),
					Value: aws.String(phoneNumber),
				},
			},
		}

		_, err := ce.CognitoClient.SignUp(user)
		if err != nil {
			fmt.Println(err)
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"ErrorTitle":   "Registration Failed",
				"ErrorMessage": err.Error()})
		} else {
			ce.RegFlow.Username = username
			render(c, gin.H{
				"title": "Redirected to One Time Password"}, "otp.html")

		}

	}
	return gin.HandlerFunc(fn)
}

func showOTPPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "OTP"}, "otp.html")
}

// OTP is Handling One Time Password
func OTP(ce CognitoExample) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		otp := c.PostForm("otp")

		user := &cognito.ConfirmSignUpInput{
			ConfirmationCode: aws.String(otp),
			Username:         aws.String(ce.RegFlow.Username),
			ClientId:         aws.String(ce.AppClientID),
		}

		result, err := ce.CognitoClient.ConfirmSignUp(user)

		if err != nil {
			fmt.Println(err)
			c.HTML(http.StatusBadRequest, "otp.html", gin.H{
				"ErrorTitle":   "OTP Failed",
				"ErrorMessage": err.Error()})
		} else {
			fmt.Println(result)
			token := generateSessionToken()

			c.SetCookie("token", token, 3600, "", "", false, true)
			c.Set("is_logged_in", true)

			render(c, gin.H{
				"title": "Successful registration & Login"}, "login-successful.html")

		}
	}

	return gin.HandlerFunc(fn)
}

// Global Pages are temporally here, cause I wasn't sure where to put them
func showAboutPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "About"}, "about.html")
}

func showLeaderboardPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "Leaderboard"}, "leaderboard.html")
}
