#!/usr/bin/env python3

from aws_cdk import core

# Import Tag library, so that you can have all resources within the Stack Tagged
from aws_cdk.core import Tag
 
from iac.iac_stack import IacStack

app = core.App()

# Next line is where you define CloudFormation Stack Name and the Region
IacStack(app, "UnicornIaC", env={'region': 'eu-west-1'})

# Define the TAGs all resources will have. Example, key: project, value: unicorn.
Tag.add(app,"project", "unicorn")
Tag.add(app,"bu", "cloud")
Tag.add(app,"environment", "prod")

app.synth()
