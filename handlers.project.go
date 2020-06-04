// handlers.project.go

package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/go-playground/validator"
)

//from AWS CDK, set variables for DynamoDB Table and S3 bucket all the Project handling functions will use
var tableName = "UnicornDynamoDBVoting"
var tableUsersName = "UnicornDynamoDBUsers"
var bucket = "www.unicornpursuit.com"

// CurrentProjectNumber is used globally as the current number of projects
var CurrentProjectNumber int

func loadProjectsDynamoDB(ddbsvc *dynamodb.DynamoDB) {

	proj := expression.NamesList(expression.Name("id"), expression.Name("title"), expression.Name("owner"), expression.Name("content"), expression.Name("photo"), expression.Name("votes"))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()

	if err != nil {
		fmt.Println("Error building expression:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	// Make the DynamoDB Query API call
	result, err := ddbsvc.Scan(params)
	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		os.Exit(1)
	} else {
		// Unmarshall and sort the results
		numItems := 1
		for _, i := range result.Items {
			item := ProjectExample{}

			err = dynamodbattribute.UnmarshalMap(i, &item)
			if err != nil {
				fmt.Println("Got error unmarshalling:")
				fmt.Println(err.Error())
				os.Exit(1)
			}
			if item.ID > 0 {
				numItems++
			}
			// Load Projects to memory, to avoid consulting DynamoDB for everything
			loadNewProject(item.ID, item.Title, item.Owner, item.Content, item.Photo, item.Votes)
		}
		CurrentProjectNumber = numItems
		fmt.Println(CurrentProjectNumber)
	}

}

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

func showUserVotes(usrsvc *dynamodb.DynamoDB) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		CurrentVotes := getCurrentVotes(usrsvc)

		Voted1 := "Your 1-point vote is still available"
		Voted2 := "Your 2-point vote is still available"
		Voted3 := "Your 3-point vote is still available"
		Voted4 := "Your 4-point vote is still available"
		Voted5 := "Your 5-point vote is still available"

		if CurrentVotes.Voted1 != 0 {
			if project, err := getProjectByID(CurrentVotes.Voted1); err == nil {
				Voted1 = "You've given 1 point vote to a project: " + project.Title
			} else {
				// If the project is not found, abort with an error
				c.AbortWithError(http.StatusNotFound, err)
			}
		}
		if CurrentVotes.Voted2 != 0 {
			if project, err := getProjectByID(CurrentVotes.Voted2); err == nil {
				Voted2 = "You've given 2 point vote to a project: " + project.Title
			} else {
				// If the project is not found, abort with an error
				c.AbortWithError(http.StatusNotFound, err)
			}
		}
		if CurrentVotes.Voted3 != 0 {
			if project, err := getProjectByID(CurrentVotes.Voted3); err == nil {
				Voted3 = "You've given 3 point vote to a project: " + project.Title

			} else {
				// If the project is not found, abort with an error
				c.AbortWithError(http.StatusNotFound, err)
			}
		}
		if CurrentVotes.Voted4 != 0 {
			if project, err := getProjectByID(CurrentVotes.Voted4); err == nil {
				Voted4 = "You've given 4 point vote to a project: " + project.Title
			} else {
				// If the project is not found, abort with an error
				c.AbortWithError(http.StatusNotFound, err)
			}
		}
		if CurrentVotes.Voted5 != 0 {
			if project, err := getProjectByID(CurrentVotes.Voted5); err == nil {
				Voted5 = "You've given 5 point vote to a project: " + project.Title
			} else {
				// If the project is not found, abort with an error
				c.AbortWithError(http.StatusNotFound, err)
			}
		}
		// Call the render function with the name of the template to render
		fmt.Println(Voted1 + Voted2)
		render(c, gin.H{
			"title":  "My Votes",
			"voted1": Voted1, "voted2": Voted2, "voted3": Voted3, "voted4": Voted4, "voted5": Voted5}, "votes.html")

	}
	return gin.HandlerFunc(fn)

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
			// If the project is not found:
			c.HTML(http.StatusBadRequest, "index.html", gin.H{
				"ErrorTitle":   "Project Error, are you logged in?",
				"ErrorMessage": err.Error()})
		}
	} else {
		// If an invalid project ID is specified in the URL, abort with an error
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"ErrorTitle":   "Project Error, are you logged in?",
			"ErrorMessage": err.Error()})
	}
}

func voteForProject(ddbsvc *dynamodb.DynamoDB, usrsvc *dynamodb.DynamoDB) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		// Retrieve information of users current votes from usrsvc DynamoDB
		CurrentVotes := getCurrentVotes(usrsvc)

		votes := c.PostForm("votes")

		if projectID, err := strconv.Atoi(c.Param("project_id")); err == nil {
			// Check if the project exists
			if project, err := getProjectByID(projectID); err == nil {
				// Make sure loggedInUser can still vote.
				// votesInt is a INT version of the string "votes".
				if votesInt, err := strconv.Atoi(votes); err == nil {

					// First we need to check if the user has already Voted for this particular project (projectID)
					if projectID == CurrentVotes.Voted1 || projectID == CurrentVotes.Voted2 || projectID == CurrentVotes.Voted3 || projectID == CurrentVotes.Voted4 || projectID == CurrentVotes.Voted5 {
						var projectIDError error = errors.New("You're not tryint to cheat, are you? Check out your previous votes below")
						fmt.Println("Project Voted!")
						render(c, gin.H{
							"title":   "My Votes",
							"payload": CurrentVotes, "ErrorTitle": "You've already Voted for this project",
							"ErrorMessage": projectIDError.Error()}, "votes.html")

					} else {
						AlreadyVoted := c.Param("project_id")
						VotedBoolean := false
						// Switch number of votes to the AlreadyVoted struct, retrieved from Users DynamoDB
						switch votesInt {
						case 1:
							if CurrentVotes.Voted1 != 0 {
								VotedBoolean = true
								fmt.Println("User already assigned 1 point to:" + AlreadyVoted)
							}

						case 2:
							if CurrentVotes.Voted2 != 0 {
								VotedBoolean = true
								fmt.Println("User already assigned 2 points to:" + AlreadyVoted)
							}
						case 3:
							if CurrentVotes.Voted3 != 0 {
								VotedBoolean = true
								fmt.Println("User already assigned 3 points to:" + AlreadyVoted)
							}
						case 4:
							if CurrentVotes.Voted4 != 0 {
								VotedBoolean = true
								fmt.Println("User already assigned 4 points to:" + AlreadyVoted)
							}
						case 5:
							if CurrentVotes.Voted5 != 0 {
								VotedBoolean = true
								fmt.Println("User already assigned 5 points to:" + AlreadyVoted)
							}
						default:
							VotedBoolean = false
						}

						if VotedBoolean == true {

							// Send message to the user that the vote cannot be done:
							var projectError error = errors.New("You're not tryint to cheat, are you? Check out your previous votes below")
							fmt.Println("Voting Error in number of points!")
							render(c, gin.H{
								"title":   "My Votes",
								"payload": CurrentVotes, "ErrorTitle": "Number of points assigned",
								"ErrorMessage": projectError.Error()}, "votes.html")

						} else {

							// Convert ID to String, required to pass it using UpdateItem function
							// We need ID and Owner as Primary Key to identify an Item we want to update
							ID := strconv.Itoa(project.ID)
							Owner := project.Owner

							// "r" is the Votes value user wants to be added to a project
							input := &dynamodb.UpdateItemInput{
								ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
									":r": {
										N: aws.String(votes),
									},
								},
								TableName: aws.String(tableName),
								Key: map[string]*dynamodb.AttributeValue{
									"id": {
										N: aws.String(ID),
									},
									"owner": {
										S: aws.String(Owner),
									},
								},
								ReturnValues: aws.String("UPDATED_NEW"),
								// Using UpdateExpression, we will add the "votes" to the current value of the field (current number of votes)
								UpdateExpression: aws.String("set votes = votes + :r"),
							}

							_, err := ddbsvc.UpdateItem(input)
							if err != nil {
								c.HTML(http.StatusBadRequest, "project.html", gin.H{
									"ErrorTitle":   "Error updating Votes",
									"ErrorMessage": err.Error()})
							}

							// Don't forget to Scan DynamoDB into a local memory again

							projectList = nil
							loadProjectsDynamoDB(ddbsvc)

							// get project again, to be sure Vote value is updated
							project, err := getProjectByID(projectID)
							if err != nil {
								c.HTML(http.StatusBadRequest, "project.html", gin.H{
									"ErrorTitle":   "Error updating Votes",
									"ErrorMessage": err.Error()})
							} else {
								// Update vote in Users Database, since all Checks for vote validation have shown OK
								updateUsersDatabase(usrsvc, projectID, votesInt)

								// Redirect to Voting Successful
								render(c, gin.H{
									"title":   "You've Voted",
									"payload": project}, "voting-successful.html")
							}

						}
					}

				} else {
					// If the project is not found, abort with an error
					c.AbortWithError(http.StatusNotFound, err)
				}
			} else {
				// If the project is not found, abort with an error
				c.AbortWithError(http.StatusNotFound, err)
			}

		} else {
			// If an invalid project ID is specified in the URL, abort with an error
			c.AbortWithStatus(http.StatusNotFound)
		}
	}
	return gin.HandlerFunc(fn)
}

func showLeaderboardPage(c *gin.Context) {

	// Get the sorted list of projects, starting with currently highest Voted
	projects := getSortedProjects()

	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title":   "Home Page",
		"payload": projects}, "leaderboard.html")
}

func createProject(ddbsvc *dynamodb.DynamoDB, sess *session.Session) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		// Set the next project ID
		id := CurrentProjectNumber
		fmt.Println("Creating project with ID:")
		fmt.Println(id)

		// Obtain the POSTed project values
		title := c.PostForm("title")

		// Get owner as an email of the logged in user
		owner := c.PostForm("owner")
		content := c.PostForm("content")

		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.HTML(http.StatusBadRequest, "create-project.html", gin.H{
				"ErrorTitle":   "File Formation error, you must submit a photo to represent your project. DevDetails:",
				"ErrorMessage": err.Error()})
		} else {
			fmt.Println(fileHeader.Filename)
		}

		f, err := fileHeader.Open()
		if err != nil {
			c.HTML(http.StatusBadRequest, "create-project.html", gin.H{
				"ErrorTitle":   "File Opening error, you must submit a valid photo to represent your project. DevDetails:",
				"ErrorMessage": err.Error()})
			return
		}

		// create a DynamoDB Item. Most information is retrieved from HTML "create-project.html"
		// Owner is the email address of the logged in user [to be implemented].
		// result.Location is a return URL value of S3 Photo Upload function.
		// Votes = 0, because we want all new projects to start with 0 votes.
		project := ProjectExample{
			ID:      id,
			Title:   title,
			Owner:   owner,
			Content: content,
			Photo:   "",
			// Set Votes to 0, as it's a new project
			Votes: 0,
		}

		//Validate the parameters
		validate := validator.New()
		valErr := validate.Struct(project)
		if valErr != nil {
			c.HTML(http.StatusBadRequest, "create-project.html", gin.H{
				"ErrorTitle":   "Error in Parameter Format, all fields are required. DevDetails: ",
				"ErrorMessage": valErr.Error()})
			fmt.Println(valErr.Error())
		} else {

			// Create an S3 Uploader
			uploader := s3manager.NewUploader(sess)

			// Upload
			result, err := uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(fileHeader.Filename),
				Body:   f,
			})
			if err != nil {
				c.HTML(http.StatusBadRequest, "create-project.html", gin.H{
					"ErrorTitle":   "S3 Upload Failed",
					"ErrorMessage": err.Error()})
			} else {
				// Success, print URL to Console
				fmt.Println("Successfully uploaded to", result.Location)
				project.Photo = result.Location

				// Marshall new project into a map of AttributeValue objects.
				av, err := dynamodbattribute.MarshalMap(project)
				if err != nil {
					c.HTML(http.StatusBadRequest, "create-project.html", gin.H{
						"ErrorTitle":   "Error Marshalling a new projects",
						"ErrorMessage": err.Error()})
				} else {

					input := &dynamodb.PutItemInput{
						Item:      av,
						TableName: aws.String(tableName),
					}

					// attempt PutItem with the created project object
					_, err = ddbsvc.PutItem(input)
					if err != nil {
						fmt.Println(err)
						c.HTML(http.StatusBadRequest, "create-project.html", gin.H{
							"ErrorTitle":   "Project Creation Failed",
							"ErrorMessage": err.Error()})
					} else {
						// Success, store project in DynamoDB and redirect to OK
						CurrentProjectNumber++
						loadNewProject(project.ID, project.Title, project.Owner, project.Content, project.Photo, project.Votes)
						render(c, gin.H{
							"title": "Project Created with Success"},
							"submission-successful.html")
					}
				}
			}
		}

	}
	return gin.HandlerFunc(fn)
}

func getCurrentVotes(svc *dynamodb.DynamoDB) CurrentVotes {

	owner := LoggedInUserEmail
	currentUsersVotes := CurrentVotes{}

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableUsersName),
		Key: map[string]*dynamodb.AttributeValue{
			"Owner": {
				S: aws.String(owner),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result)
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &currentUsersVotes)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	} else {
		if currentUsersVotes.Owner != LoggedInUserEmail {
			// Create an Item with Owner = LoggedInUserEmail in Users DynamoDB, so that user can vote, and we can register those votes. Set all Voted values to 0.
			createUserInUserDB(svc, LoggedInUserEmail)
		}
		fmt.Println(currentUsersVotes)
	}

	return currentUsersVotes
}

func createUserInUserDB(svc *dynamodb.DynamoDB, LoggedInUserEmail string) {

	var item CurrentVotes
	fmt.Println("Creating new user in UserDB: " + LoggedInUserEmail)

	item.Owner = LoggedInUserEmail
	item.Voted1 = 0
	item.Voted2 = 0
	item.Voted3 = 0
	item.Voted4 = 0
	item.Voted5 = 0

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println("Got error marshalling new movie item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableUsersName),
	}
	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

}

func updateUsersDatabase(svc *dynamodb.DynamoDB, projectID int, votes int) {

	voted := ""
	switch votes {
	case 1:
		voted = "Voted1"
	case 2:
		voted = "Voted2"
	case 3:
		voted = "Voted3"
	case 4:
		voted = "Voted4"
	case 5:
		voted = "Voted5"
	}

	project := strconv.Itoa(projectID)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				N: aws.String(project),
			},
		},
		TableName: aws.String(tableUsersName),
		Key: map[string]*dynamodb.AttributeValue{
			"Owner": {
				S: aws.String(LoggedInUserEmail),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set " + voted + " = :r"),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
