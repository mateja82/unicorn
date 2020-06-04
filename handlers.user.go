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

	"github.com/go-playground/validator"
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
		username := c.PostForm("email")

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
			//fmt.Println(result)
			fmt.Println(result.AuthenticationResult.AccessToken)
			// Set Token and ïs_logged_in Boolean, and jump to a different Menu
			// NEXT STEP: Import TOKEN variable from Login

			token := generateSessionToken()
			fmt.Println(token)

			c.SetCookie("token", token, 3600, "", "", false, true)
			c.Set("is_logged_in", true)

			// set global user variable to email used to log in
			LoggedInUserEmail = username
			fmt.Println("User just logged in: " + LoggedInUserEmail)

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
		username := c.PostForm("email")
		password := c.PostForm("password")
		fullname := c.PostForm("name")
		emailAddress := c.PostForm("email")
		phoneNumber := c.PostForm("phone_number")

		// Verify username (email) and Password using Validator package
		validate := validator.New()
		userErr := validate.Var(username, "required,email")

		if userErr != nil {
			fmt.Println(userErr.Error())

			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"ErrorTitle":   "e-mail format is not valid. DevDetails: ",
				"ErrorMessage": userErr.Error()})
		} else {

			passErr := validate.Var(password, "required,min=8")

			if passErr != nil {
				fmt.Println(passErr.Error())

				c.HTML(http.StatusBadRequest, "register.html", gin.H{
					"ErrorTitle":   "Password doesn't meet security requirements. DevDetails: ",
					"ErrorMessage": passErr.Error()})
			} else {

				// create a Cognito user object
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
				// attempt SignUp with the created user object
				_, err := ce.CognitoClient.SignUp(user)
				if err != nil {
					fmt.Println(err)
					c.HTML(http.StatusBadRequest, "register.html", gin.H{
						"ErrorTitle":   "Registration Failed",
						"ErrorMessage": err.Error()})
				} else {
					// is Success, store username in the regFlow and redirect to OTP
					ce.RegFlow.Username = username
					render(c, gin.H{
						"title": "Redirected to One Time Password"}, "otp.html")

				}
			}
		}

	}
	return gin.HandlerFunc(fn)
}

func showOTPPage(c *gin.Context) {
	render(c, gin.H{
		"title": "OTP"}, "otp.html")
}

// OTP is Handling One Time Password
func OTP(ce CognitoExample) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		// We need to get the code that user typed in the html form
		otp := c.PostForm("otp")

		// Creating a user object to ConfirmSignUp
		user := &cognito.ConfirmSignUpInput{
			ConfirmationCode: aws.String(otp),
			Username:         aws.String(ce.RegFlow.Username),
			ClientId:         aws.String(ce.AppClientID),
		}

		// Confirm SignUp with Cognito Application Client
		result, err := ce.CognitoClient.ConfirmSignUp(user)
		if err != nil {
			fmt.Println(err)
			c.HTML(http.StatusBadRequest, "otp.html", gin.H{
				"ErrorTitle":   "OTP Failed",
				"ErrorMessage": err.Error()})
		} else {

			// If OK, generate Token, set "is_logged_in" to True, and redirect to Login Success.
			fmt.Println(result)
			token := generateSessionToken()

			c.SetCookie("token", token, 3600, "", "", false, true)
			c.Set("is_logged_in", true)

			// set global user variable to email used to log in
			LoggedInUserEmail = ce.RegFlow.Username
			fmt.Println("User just Registered and logged in: " + LoggedInUserEmail)

			render(c, gin.H{
				"title": "Successful registration & Login"}, "login-successful.html")

		}
	}

	return gin.HandlerFunc(fn)
}

// "About" page is temporarily within Users section, should be moved to Global
func showAboutPage(c *gin.Context) {
	render(c, gin.H{
		"title": "About"}, "about.html")
}
