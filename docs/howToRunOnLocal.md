## How to run on KubeS3 Local

1. Run a local cluster
2. Export your aws keys to this dir
3. use `make install` to install the crd in the cluster
4. Build the operator binary using `make build`
5. You will see a binary in bin folder
6. Run that binary using bin/<binary name>
7. Now install the cr S3Bucket in the cluster and you should see the logs of creating a S3 bucket
8. TO verify check the AWS Console and see if the S3 bucket is present or not
9. Create a S3Data object with the path and BucketName.
10. Create a pod with data present in the path provided in the S3Data object path and you should see the same data present on S3 bucket in AWS Console.