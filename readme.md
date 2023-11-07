# Google Cloud Function for Image Resizing

This Google Cloud Function is designed to resize images from one Google Cloud Storage bucket and save the resized image to the same bucket with the desired dimensions. It allows you to specify the source filename, bucket, and the dimensions (width and height) for the resized image.

## Execute locally

```bash
FUNCTION_TARGET=TransferFile LOCAL_ONLY=true go run cmd/main.go 
```

## Prerequisites

Before using this function, make sure you have the following prerequisites in place:

1. **Google Cloud Storage Bucket**: You need a Google Cloud Storage bucket where the source image is stored and where the resized image will be saved. Ensure you have the appropriate access permissions for the bucket.

2. **Google Cloud Functions**: You should have a Google Cloud Functions project set up.

3. **MessagePack Package**: This function uses the `github.com/vmihailenco/msgpack` package for working with MessagePack data.

## Function Overview

The `TransferImage` function allows you to resize an image stored in a Google Cloud Storage bucket and save the resized image back to the same bucket. Here's how it works:

1. The function parses the incoming request data using MessagePack. The request should include the source filename, the source bucket, and the desired dimensions for the resized image.

2. It checks for the correctness of the provided data and ensures that both the filename and bucket are specified, and at least one of the dimensions (width or height) is provided.

3. If the source filename ends with "-thumbnail," the function returns a response indicating that it can't resize a thumbnail.

4. The function initializes the Google Cloud Storage client and gets a handle to the source bucket and the source image object.

5. It reads the image content and determines the MIME type of the image.

6. If the MIME type is allowed (based on the predefined list of allowed types), the function extracts the file extension based on the MIME type.

7. It decodes the image, calculates the new dimensions while maintaining the aspect ratio, and resizes the image.

8. The function creates a .meta file in the same bucket to store the MIME type information.

9. The resized image is saved to the same bucket with a filename indicating the new dimensions.

10. The function responds with a success message and the dimensions of the resized image.

## Usage

You can trigger the `TransferImage` function by making an HTTP request with a MessagePack-encoded JSON payload containing the source filename, source bucket, and the desired dimensions. For example:

```json
{
  "filename": "source-image.jpg",
  "bucket": "your-bucket-name",
  "width": 300,
  "height": 200
}
