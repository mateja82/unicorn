#!/usr/bin/env python3


from aws_cdk import core
from aws_cdk.core import Tag

from iac.iac_stack import IacStack

app = core.App()
IacStack(app, "iac", env={'region': 'eu-west-1'})
Tag.add(app,"project", "unicorn")
Tag.add(app,"bu", "cloud")

app.synth()
