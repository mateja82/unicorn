# unicorn

## Unicorn Web for Startup Pitch Experience

Web written in Go.

Architecture:
- Front and App - Golang in Docker, AWS Fargate + ALB with SSL
- API Gateway
- DynamoDB for Leaderboards and Session Info
- Aurora for User Info
- S3 for Static Content
