# Google Cloud Function for File Transfer

This Google Cloud Function is designed to transfer a file from one Google Cloud Storage bucket to another. It allows you to specify the source and destination buckets as query parameters when making an HTTP request.

## Execute locally

```bash
FUNCTION_TARGET=TransferFile LOCAL_ONLY=true go run cmd/main.go 
```

## Prerequisites

Before using this function, make sure you have the following prerequisites in place:

1. **Google Cloud Storage Buckets**: You need two Google Cloud Storage buckets - one for the source and another for the destination. Ensure you have the appropriate access permissions for both buckets.

2. **Google Cloud Functions**: You should have a Google Cloud Functions project set up.

3. **MessagePack Package**: This function uses the `github.com/vmihailenco/msgpack` package for working with MessagePack data.

4. **Environment Variables**: You should set environment variables for specifying the source and destination buckets based on your environment (e.g., testing or production).

## Function Overview

The `TransferFile` function allows you to transfer a file between two Google Cloud Storage buckets. Here's how it works:

1. The function extracts the source and destination bucket names from the query parameters in the HTTP request.

2. It checks the environment to determine which source bucket to use. The source bucket is determined based on the value of the `ENVIRONMENT` environment variable, which should be set to either "test" or "prod."

3. If the source bucket is not one of the predefined values (e.g., "test" or "prod"), an error response is sent.

4. The function initializes the Google Cloud Storage client and gets handles to the source and destination buckets.

5. It specifies the file to transfer (you can change the filename as needed).

6. The file is copied from the source bucket to the destination bucket.

7. The function sends a response indicating the status of the operation.

## Usage

You can trigger the `TransferFile` function by making an HTTP request with the following query parameters:

- `sourceBucket`: The source Google Cloud Storage bucket name.
- `destinationBucket`: The destination Google Cloud Storage bucket name.

Additionally, you need to set environment variables for specifying the source bucket based on your environment. For example:

- `TEST_SOURCE_BUCKET`: The source bucket for the "test" environment.
- `PROD_SOURCE_BUCKET`: The source bucket for the "prod" environment.

Ensure that the environment variables are correctly configured based on your needs.

## Error Handling

The function handles errors and returns appropriate error responses when the source bucket is invalid or when the required query parameters are missing.

## Disclaimer

This code serves as a basic example for educational purposes. You should adapt it to your specific requirements and ensure proper security, error handling, and access control in a production environment.
