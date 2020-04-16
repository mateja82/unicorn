from aws_cdk import(
        core,
        aws_s3 as s3,
        aws_lambda as _lambda
)

class IacStack(core.Stack):

    def __init__(self, scope: core.Construct, id: str, **kwargs) -> None:
        super().__init__(scope, id, **kwargs)

        # The code that defines your stack goes here

        bucket = s3.Bucket(self, "s3-unicorn-bucket",
                        bucket_name="s3-unicorn-bucket",
                        access_control=s3.BucketAccessControl.PUBLIC_READ,
                        )

        bucket.grant_public_access()
