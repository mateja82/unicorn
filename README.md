# Unicorn Project

Unicorn Project is an **open source Cloud Learning Plarform**. If you want to improve your skill in Cloud, Serverless, Web development in Cloud (frontend and backend) - **you're in the right place**! Keep reading.

## Unicorn Project Mission: Make you a Cloud Ninja

Cloud Ninja may mean Cloud Developer, Cloud Engineer or a Cloud Architect. Think about what YOU want it to mean. If you are one of the following, you definitely belong here:

- Developer, who wants to go deeper to Cloud and DevOps.
- SysAdmin, who wants to go deeper to Cloud and DevOps.
- Cloud Engineer or Achitect who wants to learn more about end-to-end application design and deployment in the Cloud.

Regardless if you're an absolute beginner, or an experienced professional - we promisse you that you'll learn a bunch.

## Unicorn Project Rules

### Rule 1: Focus on WHY and HOW

Our main objective is for you to learn. If you write code, it's not enough to just find the best solution to a problem. **Make sure to explain why you did it that way, because like you learn from others - the others will learn from you.**

### Rule 2: Less is more

**The two principles of Unicorn are focus and simplicity**:

- **FOCUS**: We go deep into, and really learn the technologies that interest us. This means our priority is expertise in one thing, not doing "hello world" to learn 20 different technologies that do the same thing.
- **SIMPLICITY**: We go step by step. To start building, you need essential info, not 30 hours of deep dive videos. You'll get to the deep dark stuff, but the idea is to start small, and build your level up, without getting scared of how complex the technology with all it's details is. **Don't copy and paste huge chunks of text. No one likes to read dry text. Rather make an effort and try to make it as short as possible.** We all appreciate short text that simply "explains it all".

## Unicorn Project = Unicorn Pursuit + Unicorn Workshop

Unicorn Project has two components:

- **Unicorn Pursuit is an open source Voting Web Platform**, we are building on the up-and-coming technology stack. These are the skills that will be in the highest demand within the next 5 to 10 years. More details [here](https://www.matscloud.com/docs/unicorn-project/).
- **Unicorn Workshop is a step by step Cloud Ninja Training**, which explains each step of building Unicorn Pursuit. More details about the workshop can be found [here](https://www.matscloud.com/docs/unicorn-project/workshop/).

## Technology Stack

The beauty of Unicorn Project is that it allows you to choose the part of technology stack you want to learn, and maybe even contribute on. Basic architecture stack consists of:

- **AWS Serverless** (Lambda, Cognito, DynamoDB, Gargate etc.).
- All AWS resources provisioned using **AWS CDK (Cloud Development Kit) - Python**.
- **Go (golang)**.
- **Gin** web development framework for Go.

For more details about [Unicorn Pursuit Architecture, click here](https://www.matscloud.com/docs/unicorn-project/architecture/).

Yes, there will probably be a Unicorn Pursuit version on Google Cloud. If you are interested in this, please let us know.

## Repository Structure

Tree shown below:

```
unicorn
├── iac
│   ├── cdk.out
│   └── iac
│       └── __pycache__
└── templates
```

The `main.go` and the rest of the Go packages are in the `root` folder.

Go HTML Templates are in the `template` folder.

AWS CDK with all it's components is in `iac` folder.

## Documentation

The entire documentation is currently stored on Mats Cloud (a personal blog). The idea is to migrate the content to Unicorn Pursuit, once the web is made.

You can find all the information about Unicorn Pursuit and Unicorn Workshop [on Mats Cloud](https://www.matscloud.com/docs/).

**Please feel free to start your Cloud Ninja journey, contribute if you want, and give feedback so we can improve**.

## Requirements & Dev Environment

You can find all the requirements, and how to build your Dev environment [here](https://www.matscloud.com/docs/dev-environment/). Based on the technology that interests you most, you may not need to install the entire set of tools.

## I want to contribute, what do I do?

Glad you asked! First, fork a project, and make sure to create your branch. Try naming it as the feature you're working on.

Now, about what feature to work on...

There are 2 CANBAN boards associated with the Unicorn Git repository. You can visit the **Backlog** column, and either choose a ticket, move it to "In Progress" and start working on it, or you can propose an improvement by adding a ticket to the "Backlog".

If it's an **issue** or an improvement of an existing functionality, create an Issue, be sure to associate it with the correct Project, and it will automatically be added to **To Do** column.

The two projects are:

- **AWS Infrastructure**: [where we use AWS CDK for managing AWS resources via Python](https://github.com/mateja82/unicorn/projects/2)
- **Unicorn Pursuit Web App**: [which is our web app, done in Golang, React and AWS SDK](https://github.com/mateja82/unicorn/projects/1)

## How do I run Unicorn Pursuit in my environment?

You need to set up:

- AWS CLI v2 (make sure you run `aws configure` and configure the authentication, and the desired AWS Region)
- AWS CDK
- Go, with the correct GOPATH and GOROOT

You can follow [the guide on how to set up a dev environment here](https://www.matscloud.com/docs/dev-environment/), or set the above yourself.

Clone/Fork this repo to `$HOME/go/src/`, and first deploy the AWS infrastructure. Before you deploy everything with AWS CDK, be sure you've deployed the initial CDK CloudFormation Stack, and go to `iac` folder to check out if CDK "works":

```
cd $HOME/go/src/unicorn/iac
cdk ls
UnicornIaC
```

If your CDK isn't working, [make sure all dependencies found here are correctly deployed](https://docs.aws.amazon.com/cdk/latest/guide/getting_started.html), and come back.

Once CDK works, go for deployment:
```
cdk deploy
```

This will deploy everything: S3 Bucket, DynamoDB Table, Cognito User Pool. The only thing you need to do "manually", is configure the Cognito Application Cliend in SSM Parameter Store, because it's confidential, and you shouldn't dispose it in your code. You can [find how to do that here](https://www.matscloud.com/docs/cloud-sdk/ssm-for-credentials/), it's pretty simple.

Once you're done with that, Unicorn Pursuit should run smoothly:

```
cd $HOME/go/src/unicorn
go run unicorn
```
You should be able to access your App from your browser, port 8080: `http://localhost:8080`.

Don't forget that you first need to create a user. Your phone number will be verified, and you can go on and create a few Projects.

## Workshop Agenda

True beauty of this workshop is that it's designed so that **you choose the technologies you want to focus on, there is no need to follow the proposed oreder**. The Unicorn App is ready and deployable, so just choose your tech. Hopefully you'll end up contributing.

||||
|-------------------|-----------------|------|
||**SECTION 1**| **Let's set up the Tools and Platforms** |
|1.1| [Build your Dev Environment](https://www.matscloud.com/docs/dev-environment/) | Before we start writting any code, you need to have the right tools. This includes AWS Account, GitHub, AWS CLI, AWS CDK, AWS SDK, Go etc.|
|1.2|[Get familiar with GitHub Unicorn Repository](https://github.com/mateja82/unicorn)|Unicorn Project is an open source Cloud Learning Plarform. All the code we're handling here is public. Read the contend of `README.md` to get acquainted with how code is organized.|
|1.3|[Let's read the Requirements, and you understand the Web App we're building](https://www.matscloud.com/docs/unicorn-project/requirements/)|Remember this: before starting the design of any application, make sure you fully understand the requirements. How else will you know if your design meets the requirements if you don't even know the requirements? This includes SLAs, Backup, how users need to authenticate, regulation and compliance etc.|
|1.4|[Let's get familiar with the Unicorn Pursuit Web App Architecture](https://www.matscloud.com/docs/unicorn-project/architecture/)|Understand the technology stack, and what each component does in the final architecture.|
|1.5| [Intro to AWS CloudFormation](https://www.matscloud.com/docs/cloud-infrastructure/cloudformation/)  | How to write CloudFormation YAML Templates and create Stacks. Before going deep into individual AWS Services, we need to get familiar with configuration tools we'll be using. |
|1.6| [Start with AWS CDK (Cloud Development Kit)](https://www.matscloud.com/docs/cloud-infrastructure/aws-cdk/)| Learn how to seize the best of both, indicative and declarative programming. How to create CloudFormation Stacks using Python with special library called AWS CDK (Cloud Development Kit). AWS CDK with Python is our official deployment tool, which makes this exercise very important. |
||||
||**SECTION 2**|**AWS: Cloud Infrastructure** |
|2.1|[Understand Regions, TAGs, Naming](https://www.matscloud.com/docs/cloud-infrastructure/)| Understand basic AWS concepts, including Availability Zones, Regions, Tagging, Naming, etc. as well as the rules we'll be using them in this Workshop|
|2.2|[S3 as Object Storage](https://www.matscloud.com/docs/cloud-infrastructure/s3/)|S3 will be used to store our static content. Learn how to configure S3 bucket.|
|2.3|[DynamoDB as NoSQL Database](https://www.matscloud.com/docs/cloud-infrastructure/dynamodb/)|DynamoDB is a NoSQL high performant database that fits perfectly in our example of Voting and forming Leaderboards for Unicorn Pursuit|
|2.4|[Cognito for Authentication](https://www.matscloud.com/docs/cloud-infrastructure/cognito/)|Even though we could just locally handle users and passwords, we'd prefer to do something enterprise production ready, with managed user pool and MFA, which we will achieve.|
|2.5|[EC2](https://www.matscloud.com/docs/cloud-infrastructure/ec2/)|Learn everythink about AWS IaaS and Virtual Machines|
|2.6|[ECS, managed Container Service - Fargate](https://www.matscloud.com/docs/cloud-infrastructure/ecs-fargate)|To make our app run the same on any platform, we will dockerize it, and deploy on container management platform with Fargate|
|2.7|[CDN with CloudFront](https://www.matscloud.com/docs/cloud-infrastructure/cloudfront/)|To enable caching on the edge, and improve user experience wherever our users are, we will use CloudFront, a Content Delivery Network by AWS.|
|2.8|[IAM - Identity Access Management](https://www.matscloud.com/docs/cloud-infrastructure/iam/)|Learn how to manage users, and their access to the AWS resources. |
|2.9|[KMS - Key Management Service](https://www.matscloud.com/docs/cloud-infrastructure/kms/)|Learn how to use  KMS to encrypt data to help protect against improper access.|
|2.10|[API Gateway](https://www.matscloud.com/docs/cloud-infrastructure/api-gateway/)|API Gateway allows you to publish, maintain, monitor and secure your APIs. API Gateway basically exposes REST API HTTPS endpoints, and manages what each API endpoint “maps to” in the backend, what it triggers.|
|2.11|[Messaging Services: SQS, SNS](https://www.matscloud.com/docs/cloud-infrastructure/sqs-sns/)|Messaging protocols can be very useful when you want to notify either users or other services of an event.|
||||
||**SECTION 3**|**Let's dive into Web Development in GO**|
|3.1|[Start with Go basics](https://www.matscloud.com/docs/web-development/golang/)|Go is an open source programming language that makes it easy to build simple, reliable, and efficient software. We need to learn the basic constructs, such as packages, routes, templates etc.|
|3.2|[Gin Web development framework for Go](https://www.matscloud.com/docs/web-development/gin-framework/)|Learn how to make Go Web development a bit more user-friendly and more performant using Gin framework.|
|3.3|[Containerize Unicorn with Dockerfile](https://www.matscloud.com/docs/web-development/docker/)|How to containerize Golang Web App. Learn how to create a Dockerfile|
|3.4|[Create Frontend with HTML, CSS, JS and React](https://www.matscloud.com/docs/web-development/html-css/)||
|3.5|[Testing in Go](https://www.matscloud.com/docs/web-development/testing/)||
||||
||**SECTION 4**|[**AWS SDK (Software Development Kit) for Go and Python**](https://www.matscloud.com/docs/cloud-sdk/)|
|4.1|[Create AWS Session using AWS SDK](https://www.matscloud.com/docs/cloud-sdk/aws-session/)||
|4.2|[Use Systems Manager to store credentials](https://www.matscloud.com/docs/cloud-sdk/ssm-for-credentials/)|How to safely store your passwords and other confidential info using AWS Systems Manager Parameter Store, and retreive the value to a variable in your code. Learn to NEVER store your passwords or access keys in plane text, anywhere.|
|4.3|[Configure storing static files in S3](https://www.matscloud.com/docs/cloud-sdk/go-and-s3/)|How to store and consume S3 objects from Golang. In Unicorn Pursuit, we'll store Project Images and Diagrams in an S3 Bucket.|
|4.4|[Handle Unicorn Users with Cognito](https://www.matscloud.com/docs/cloud-sdk/go-and-cognito/)|We created the Cognito User Pool and App Client using AWS CDK. Now it’s time to create Login, Logout and Register page in Unicorn Pursuit, so our users can authenticate.|
|4.5|[Unicorn Projects and Voting with DynamoDB](https://www.matscloud.com/docs/cloud-sdk/go-and-dynamodb/)||
|4.6|[Build Unicorn container in Fargate](https://www.matscloud.com/docs/cloud-sdk/docker-in-fargate/)|  How to containerize the Go application, and deploy it in ECS managed container orchestrator - Fargate.|
|4.7|[Create required Lambdas](https://www.matscloud.com/docs/cloud-sdk/lambda/)|Lambda is always handy to implement the functions we don't have "out of the box". Let's use AWS SDK for Python, and fill the gaps.|
||||
||**SECTION 5**|**[Site Reliability Engineering (SRE)](https://www.matscloud.com/docs/sre/)**|
|5.1|[Why SRE could be the perfect implementation of DevOps](https://www.matscloud.com/docs/sre/class-sre-implements-devops/)||
|5.2|[Day 2 Operations: SLOs, SLAs, Monitoring, Alarms](https://www.matscloud.com/docs/sre/slos-and-day-2-ops/)||
||||
||**SECTION 6**|**[Understand Shared Responsability model in Cloud](https://www.matscloud.com/docs/shared-responsability/)**|
||||
||**SECTION 7**|**[Get Certified](https://www.matscloud.com/docs/get-certified/)**|

## What now?

**What do I do next?**: [Continue to Mat's Cloud and the read Unicorn Pursuit App Requirements, to understand what we're building](https://www.matscloud.com/docs/unicorn-project/requirements/)

You will find further instructions there.