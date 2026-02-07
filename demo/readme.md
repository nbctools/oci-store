# Running locally

# Testing 

## Pushing an image to AWS S3 bucket
``` shell
docker pull redis:latest
env AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test \
go run ../.  s3 push -r us-east-1 -e s3.localhost.localstack.cloud:4566 oci-images-bucket/demo/redis:v1 -i redis:latest
```

## Testing GCS Storage

Push image:

``` shell
# run `gcloud auth application-default login` to get default app credentials
go run ../.  gcs push oci-store-bucket/demo/redis:v1 -i redis:latest

```

## Pushign an image to Azure blob storage

``` shell
# export AZURE_STORAGE_ACCOUNT=... \
#     AZURE_STORAGE_KEY=... \
#     AZURE_SECRET=... \
#     AZURE_CLIENT_ID=... \
#     AZURE_TENANT_ID=...
source .env
go run ../. azure push mycontainer/demo/redis:stable -i redis:latest
```
