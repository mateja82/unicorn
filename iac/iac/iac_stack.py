from aws_cdk import(
        core,
)

class IacStack(core.Stack):

    def __init__(self, scope: core.Construct, id: str, **kwargs) -> None:
        super().__init__(scope, id, **kwargs)

        # Create S3 Bucket, a new object of class Bucket
        #bucket = s3.Bucket(self, "S3BucketUnicorn",
        #                bucket_name="S3BucketUnicorn",
        #                access_control=s3.BucketAccessControl.PUBLIC_READ,
        #                )
        # Since we'll be using this bucket to store publically available images, we can give Public Access.
        #bucket.grant_public_access()


