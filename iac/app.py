#!/usr/bin/env python3
# https://docs.aws.amazon.com/cdk/latest/guide/work-with-cdk-python.html

from aws_cdk import core

from iac.iac_stack import IacStack


app = core.App()
IacStack(app, "iac")

app.synth()
