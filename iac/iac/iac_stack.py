from aws_cdk import (
    aws_s3 as s3,
    aws_dynamodb as ddb,
    aws_cognito as cognito,
    aws_certificatemanager as acm,
    core
)

class IacStack(core.Stack):

    def __init__(self, scope: core.Construct, id: str, **kwargs) -> None:
        super().__init__(scope, id, **kwargs)

        # S3: Create a Bucket for Unicorn Pursuit web page, and grant public read:
        bucket = s3.Bucket(self, "www.unicornpursuit.com",
                        bucket_name="www.unicornpursuit.com",
                        access_control=s3.BucketAccessControl.PUBLIC_READ,
                        )

        # Grant public read access to the bucket
        bucket.grant_public_access()

        # DynamoDB: Create Table for Project Info (ID, Owner, Content, Photo and Votes)
        voting = ddb.CfnTable(
            self, "UnicornDynamoDBVoting",
            table_name="UnicornDynamoDBVoting",
            key_schema=[
                ddb.CfnTable.KeySchemaProperty(attribute_name="id",key_type="HASH"),
                ddb.CfnTable.KeySchemaProperty(attribute_name="owner",key_type="RANGE"),
            ],
            
        # In the new DynamoDB, you can't create AttDefProperty for non-key attributes.
            attribute_definitions=[
                ddb.CfnTable.AttributeDefinitionProperty(attribute_name="id",attribute_type="N"),
                ddb.CfnTable.AttributeDefinitionProperty(attribute_name="owner",attribute_type="S"),
#               ddb.CfnTable.AttributeDefinitionProperty(attribute_name="title",attribute_type="S"),
#               ddb.CfnTable.AttributeDefinitionProperty(attribute_name="content",attribute_type="S"),
#               ddb.CfnTable.AttributeDefinitionProperty(attribute_name="photo",attribute_type="S"),
#               ddb.CfnTable.AttributeDefinitionProperty(attribute_name="votes",attribute_type="N"),
            ],
            provisioned_throughput=ddb.CfnTable.ProvisionedThroughputProperty(
                read_capacity_units=5,
                write_capacity_units=5
            )
        )

        # Grant RW writes for Unicorn App in Fargate and relevant Lambdas invoked from API Gateway
        # voting.grant_read_write_data(user)

        # Cognito: Create User Pool
        userpool = cognito.UserPool(
            self, "CognitoUnicornUserPool",
            user_pool_name="CognitoUnicornUserPool",
            self_sign_up_enabled=True,
            
            ## Require username or email for users to sign in
            sign_in_aliases=cognito.SignInAliases(
                username=True,
                email=True,
            ),
            # Require users to give their full name when signing up
            required_attributes=cognito.RequiredAttributes(
                fullname=True,
                email=True,
            ),
            # Verify new sign ups using email
            auto_verify=cognito.AutoVerifiedAttrs(
                email=True,
                phone=False,
            ),
            # Configure a verification email, sent by default Cognito email address (no-reply@verificationemail.com)
            user_verification=cognito.UserVerificationConfig(
                email_subject="Unicorn Pursuit: Verify your email",
                email_body="Hello {username}. Welcome to Unicorn Pursuit! Follow the link {##Verify Email##} to confirm your email address.",
                email_style=cognito.VerificationEmailStyle.LINK,
            ),
            # Set up required password policy
            password_policy=cognito.PasswordPolicy(
                min_length=12,
                require_symbols=True,
                require_lowercase=True,
                require_uppercase=True,
                require_digits=True,
            )
        )

        ## Cognito: Create App Client & create Authentication Flow with User and Password
        client = userpool.add_client(
            "UnicornAppClient",
            user_pool_client_name="UnicornAppClient",
            auth_flows={
                "user_password": True,
                "refresh_token": True,
            },
        )

        client_id = client.user_pool_client_id

        ## Cognito: Create Domain for auth.unicornpursuit.com
        certificate_arn="arn:aws:acm:us-east-1:057097267726:certificate/953cc4d7-f344-403f-87ce-f9b96153a304"
        domain_cert = acm.Certificate.from_certificate_arn(self, "domainCert", certificate_arn)
        userpool.add_domain("UnicornPursuitDomain",
           custom_domain={
                "domain_name": "auth.unicornpursuit.com",
                "certificate": domain_cert
            }
        )