#!/usr/bin/env python3

from aws_cdk import core
from aws_cdk.core import Tag

from unicorn_prod.unicorn_prod_stack import UnicornProdStack

app = core.App()
UnicornProdStack(app, "unicorn-prod", env={'region': 'eu-west-1'})
Tag.add(app,"project", "unicorn")
Tag.add(app,"bu", "cloud")
Tag.add(app,"environment", "prod")

app.synth()
