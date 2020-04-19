from aws_cdk import(
        core,
        aws_s3 as s3,
        aws_dynamodb as ddb,
)

class UnicornProdStack(core.Stack):

    def __init__(self, scope: core.Construct, id: str, **kwargs) -> None:
        super().__init__(scope, id, **kwargs)

        # The code that defines your stack goes here
        bucket = s3.Bucket(self, "s3-prod-unicorn",
                        bucket_name="s3-prod-unicorn",
                        access_control=s3.BucketAccessControl.PUBLIC_READ,
                        )

        bucket.grant_public_access()

        # Create DynamoDB for Voting and Leaderboard

        voting = ddb.CfnTable(
                        self, "dynamodb-prod-unicorn",
                        table_name="dynamodb-prod-unicorn",
                        key_schema=[
                            ddb.CfnTable.KeySchemaProperty(attribute_name="project",key_type="HASH"),
                            ddb.CfnTable.KeySchemaProperty(attribute_name="owner",key_type="RANGE")
                            ],
                        attribute_definitions=[
                            ddb.CfnTable.AttributeDefinitionProperty(attribute_name="project",attribute_type="S"),
                            ddb.CfnTable.AttributeDefinitionProperty(attribute_name="owner",attribute_type="S"),
                            ddb.CfnTable.AttributeDefinitionProperty(attribute_name="votes",attribute_type="N"),
                            ddb.CfnTable.AttributeDefinitionProperty(attribute_name="members",attribute_type="S"),
                        ],
                        provisioned_throughput=ddb.CfnTable.ProvisionedThroughputProperty(
                            read_capacity_units=5,
                            write_capacity_units=5
                        )

        )