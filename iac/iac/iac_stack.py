from aws_cdk import (
    aws_s3 as s3,
    aws_dynamodb as ddb,
    aws_cognito as cognito,
    aws_certificatemanager as acm,
    aws_ecs as ecs,
    aws_ecs_patterns as ecs_patterns,
    aws_ec2 as ec2,
    aws_ecr as ecr,
    aws_iam as iam,

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
            ],
            provisioned_throughput=ddb.CfnTable.ProvisionedThroughputProperty(
                read_capacity_units=5,
                write_capacity_units=5
            )
        )

        # Grant RW writes for Unicorn App in Fargate and relevant Lambdas invoked from API Gateway
        # voting.grant_read_write_data(user)

        # Second DynamoDB table called "users" for storing who voted for whom
        # Example: user1@cepsa.com gave 5 points to project 323, 4 points to 111 etc.
        users = ddb.Table(
            self, "UnicornDynamoDBUsers",
            table_name="UnicornDynamoDBUsers",
            partition_key=ddb.Attribute(
                name="owner",
                type=ddb.AttributeType.STRING
            )
        )

        # Cognito: Create User Pool
        userpool = cognito.UserPool(
            self, "CognitoUnicornUserPool",
            user_pool_name="CognitoUnicornUserPool",
            self_sign_up_enabled=True,
            
            ## Require username or email for users to sign in
            sign_in_aliases=cognito.SignInAliases(
                username=False,
                email=True,
            ),
            # Require users to give their full name when signing up
            required_attributes=cognito.RequiredAttributes(
                fullname=True,
                email=True,
                phone_number=True
            ),
            # Verify new sign ups using email
            auto_verify=cognito.AutoVerifiedAttrs(
                email=False,
                phone=True,
            ),
            # Configure OTP Settings ()
            user_verification=cognito.UserVerificationConfig(
                sms_message="Hey Unicorn Hunter, welcome to Unicorn Pursuit! Your OTP is {####}",
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
            generate_secret=False,
            
            ## We'll allow both Flows, Implicit and Authorization Code, and decide in the app which to use.
            auth_flows=cognito.AuthFlow(
                admin_user_password=False,
                custom=False,
                refresh_token=True,
                user_password=True,
                user_srp=False
                ),
        )

        ## Fargate: Create ECS:Fargate with ECR uploaded image

        vpc = ec2.Vpc(self, "UnicornVPC", max_azs=2)

        cluster = ecs.Cluster(self, "UnicornCluster", vpc=vpc)

        repo = ecr.Repository(self, "unicorn", repository_name="unicorn")

        fargate_role = iam.Role.from_role_arn(self, "ecsTaskExecutionRoleAdminTest", "arn:aws:iam::057097267726:role/ecsTaskExecutionRoleAdminTest",
           mutable=True
        )

        fargate_service = ecs_patterns.ApplicationLoadBalancedFargateService(self, "UnicornFargateService",
            cluster=cluster,
            cpu=512,
            desired_count=1,
            task_image_options=ecs_patterns.ApplicationLoadBalancedTaskImageOptions(
                image=ecs.ContainerImage.from_registry("057097267726.dkr.ecr.eu-west-1.amazonaws.com/unicorn"),
                # image=ecs.ContainerImage.from_registry(repo.repository_uri_for_tag()),
                container_port=8080,
                container_name="unicorn",
                execution_role=fargate_role,
                ),
                
            memory_limit_mib=1024,
            public_load_balancer=True   
        )

        fargate_service.service.connections.security_groups[0].add_ingress_rule(
            peer = ec2.Peer.ipv4(vpc.vpc_cidr_block),
            connection = ec2.Port.tcp(8080),
            description="Allow http inbound from VPC"
        )





