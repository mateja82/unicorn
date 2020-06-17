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

        # Let's start with creating an IAM Service Role, later to be assumed by our ECS Fargate Container
        # After creating any resource, we'll be attaching IAM policies to this role using the `fargate_role`.
        fargate_role = iam.Role(
            self,
            "ecsTaskExecutionRole",
            role_name="ecsTaskExecutionRole",
            assumed_by=iam.ServicePrincipal("ecs-tasks.amazonaws.com"),
            description="Custom Role assumed by ECS Fargate (container)"
        )

        # S3: Create a Bucket for Unicorn Pursuit web page, and grant public read:
        bucket = s3.Bucket(self, "www.unicornpursuit.com",
                        bucket_name="www.unicornpursuit.com",
                        access_control=s3.BucketAccessControl.PUBLIC_READ,
                        )

        # Grant public read access to the bucket
        bucket.grant_public_access()


        # Grant S3 Read/Write access to our Fargate Container
        fargate_role.add_to_policy(statement=iam.PolicyStatement(
            resources=[bucket.bucket_arn],
            actions=["s3:*"]
        ))

        # DynamoDB: Create Table for Project Info (ID, Owner, Content, Photo and Votes)
        voting_ddb = ddb.CfnTable(
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
        

        # Second DynamoDB table called "users" for storing who voted for whom
        # Example: user1@cepsa.com gave 5 points to project 323, 4 points to 111 etc.
        users_ddb = ddb.Table(
            self, "UnicornDynamoDBUsers",
            table_name="UnicornDynamoDBUsers",
            partition_key=ddb.Attribute(
                name="Owner",
                type=ddb.AttributeType.STRING
            )
        )

        # Grant RW writes for Unicorn App in Fargate
        fargate_role.add_to_policy(statement=iam.PolicyStatement(
            resources=[voting_ddb.attr_arn,
            users_ddb.table_arn],
            actions=["dynamodb:*"]
        ))

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
        userpool.add_client(
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

        # Grant Cognito Access to Fargate. Include SSM, so Client App ID can be retrived.
        fargate_role.add_to_policy(statement=iam.PolicyStatement(
            resources=["*"],
            actions=["ssm:*"]
        ))

        fargate_role.add_to_policy(statement=iam.PolicyStatement(
            resources=[userpool.user_pool_arn],
            actions=["cognito-identity:*","cognito-idp:*","cognito-sync:*"]
        ))

        ## Fargate: Create ECS:Fargate with ECR uploaded image
        vpc = ec2.Vpc(self, "UnicornVPC", max_azs=2, nat_gateways=None)

        """   VPC with optimal NAT GW usage
        vpc_lowcost = ec2.Vpc(self, "LowCostVPC",
            max_azs=2,
            cidr="10.7.0.0/16",
            nat_gateways=None,
            subnet_configuration=[ec2.SubnetConfiguration(
                               subnet_type=ec2.SubnetType.PUBLIC,
                               name="Public",
                               cidr_mask=24,
                           ), ec2.SubnetConfiguration(
                               subnet_type=ec2.SubnetType.ISOLATED,
                               name="Private",
                               cidr_mask=24,
                           )
                           ],     
        ) """
        
        linux_ami = ec2.AmazonLinuxImage(generation=ec2.AmazonLinuxGeneration.AMAZON_LINUX,
                                 edition=ec2.AmazonLinuxEdition.STANDARD,
                                 virtualization=ec2.AmazonLinuxVirt.HVM,
                                 storage=ec2.AmazonLinuxStorage.GENERAL_PURPOSE
                                 )

        nat_ec2 = ec2.Instance(self, "NAT", 
            instance_name="NAT",
            vpc=vpc,
            vpc_subnets=ec2.SubnetSelection(subnet_type=ec2.SubnetType.PUBLIC),
            instance_type=ec2.InstanceType(instance_type_identifier="t3.nano"),
            machine_image=linux_ami,
            user_data=ec2.UserData.for_linux(),
            source_dest_check=False,
        )

        # Configure Linux Instance to act as NAT Instance
        nat_ec2.user_data.add_commands("sysctl -w net.ipv4.ip_forward=1")
        nat_ec2.user_data.add_commands("/sbin/iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE")

        # Add a static route to the ISOLATED subnet, pointing 0.0.0.0/0 to a NAT EC2
        selection = vpc.select_subnets(
            subnet_type=ec2.SubnetType.PRIVATE
        )
    
        for subnet in selection.subnets:
            subnet.add_route("DefaultNAT", router_id=nat_ec2.instance_id, router_type=ec2.RouterType.INSTANCE, destination_cidr_block="0.0.0.0/0") 

        # Create ECS Cluster
        cluster = ecs.Cluster(self, "UnicornCluster", vpc=vpc)
        ecr.Repository(self, "unicorn", repository_name="unicorn")

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

        # Update NAT EC2 Security Group, to allow only HTTPS from Fargate Service Security Group.
        nat_ec2.connections.security_groups[0].add_ingress_rule(
            peer = fargate_service.service.connections.security_groups[0],
            connection = ec2.Port.tcp(443),
            description="Allow https from Fargate Service"
        )

        # Grant ECR Access to Fargate by attaching an existing ReadOnly policy. so that Unicorn Docker Image can be pulled.
        #fargate_role.add_managed_policy(iam.ManagedPolicy("AmazonEC2ContainerRegistryReadOnly"))
        fargate_role.add_to_policy(statement=iam.PolicyStatement(
            resources=["*"],
            actions=["ecr:*"]
        ))
