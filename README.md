# Unicorn Project

Unicorn Project is an **open source Cloud Learning Plarform**. If you want to improve your skill in Cloud, Serverless, Web development in Cloud (frontend and backend) - **you're in the right place**! Keep reading.

## Unicorn Project Mission: Make you a Cloud Ninja

Cloud Ninja may mean Cloud Developer, Cloud Engineer or a Cloud Architect. Think about what YOU want it to mean. If you are one of the following, you definitely belong here:

- Developer, who wants to go deeper to Cloud and DevOps.
- SysAdmin, who wants to go deeper to Cloud and DevOps.
- Cloud Engineer or Achitect who wants to learn more about end-to-end application design and deployment in the Cloud.

Regardless if you're an absolute beginner, or an experienced professional - we promisse you that you'll learn a bunch.

## Unicorn Project Rules

There is only 1 rule: **our main objective is for you to learn. If you write code, it's not enough to just find the best solution to a problem. Make sure to explain why you did it that way, because like you learn from others - the others will learn from you.**

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
