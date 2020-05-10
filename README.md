# unicorn

## Unicorn Web for Startup Pitch Experience

Web written in Go.

Architecture:

- Front and App - Golang in Docker, AWS Fargate + ALB with SSL
- Cognito for Authentication and Authorization
- API Gateway
- DynamoDB for Leader boards and Session Info
- Aurora for User Info
- S3 for Static Content (images etc)


