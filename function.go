package function

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/vmihailenco/msgpack" // Import MessagePack package

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"

	"github.com/disintegration/imaging"
)

func init() {
	functions.HTTP("Resize", resize)
}

// Define a map of allowed MIME types and their corresponding formats
var mimesAllowed = map[string]string{
	"image/gif":    "gif",
	"image/x-icon": "ico",
	"image/jpeg":   "jpeg",
	"image/webp":   "webp",
	"image/png":    "png",
}

// ImageData represents the structure of the incoming JSON payload
type ImageData struct {
	Filename string `msgpack:"filename"`
	Bucket   string `msgpack:"bucket"`
	Width    int    `msgpack:"width"`
	Height   int    `msgpack:"height"`
}

func sendResponse(w http.ResponseWriter, message string, statusCode int, width, height int) {
	responseData := map[string]interface{}{
		"message":     message,
		"status_code": statusCode,
		"width":       width,
		"height":      height,
	}

	// Serialize the response to MessagePack
	responseBytes, err := msgpack.Marshal(responseData)
	if err != nil {
		http.Error(w, "Failed to create response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/msgpack")
	w.WriteHeader(statusCode)
	w.Write(responseBytes)
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	sendResponse(w, message, statusCode, 0, 0)
}

func resize(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Parse the request data using MessagePack
	imageData := &ImageData{}
	decoder := msgpack.NewDecoder(r.Body)
	if err := decoder.Decode(imageData); err != nil {
		sendError(w, "Failed to parse request data", http.StatusBadRequest)
		return
	}

	if imageData.Filename == "" || imageData.Bucket == "" ||
		(imageData.Width == 0 && imageData.Height == 0) {
		sendError(w, "Incorrect filename, bucket, width, or height", http.StatusBadRequest)
		return
	}

	if strings.HasSuffix(imageData.Filename, "-thumbnail") {
		sendResponse(w, "Can't resize a thumbnail", http.StatusOK, 0, 0)
		return
	}

	// Initialize the Google Cloud Storage client
	client, err := storage.NewClient(ctx)
	if err != nil {
		sendError(w, "Failed to create client", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Get a handle to the source bucket
	bucket := client.Bucket(imageData.Bucket)
	obj := bucket.Object(imageData.Filename)

	rc, err := obj.NewReader(ctx)
	if err != nil {
		sendError(w, "Failed to read source file", http.StatusInternalServerError)
		return
	}
	defer rc.Close()

	// Determine the MIME type of the file
	buf := make([]byte, 512)
	_, err = io.ReadFull(rc, buf)
	if err != nil {
		sendError(w, "Failed to determine MIME type", http.StatusInternalServerError)
		return
	}

	mime := http.DetectContentType(buf)
	extension, allowed := mimesAllowed[mime]
	if !allowed {
		sendError(w, "File type not allowed", http.StatusBadRequest)
		return
	}

	// Read the full file content
	rc, err = obj.NewReader(ctx)
	if err != nil {
		sendError(w, "Failed to read source file", http.StatusInternalServerError)
		return
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		sendError(w, "Failed to read source file", http.StatusInternalServerError)
		return
	}

	// Open the image and determine its current width and height
	srcImage, err := imaging.Decode(bytes.NewReader(content))
	if err != nil {
		sendError(w, "Failed to decode source image", http.StatusInternalServerError)
		return
	}

	currentWidth := srcImage.Bounds().Dx()
	currentHeight := srcImage.Bounds().Dy()

	// Calculate new dimensions while maintaining the aspect ratio
	newWidth, newHeight := calculateNewDimensions(currentWidth, currentHeight, imageData.Width, imageData.Height)

	// Resize the image
	resizedImage := imaging.Resize(srcImage, newWidth, newHeight, imaging.Lanczos)

	// Create a .meta file with the MIME type
	metaData := fmt.Sprintf(`{"mime":"%s"}`, mime)
	metaObj := bucket.Object(fmt.Sprintf("%s.meta", imageData.Filename))
	metaWc := metaObj.NewWriter(ctx)
	if _, err := metaWc.Write([]byte(metaData)); err != nil {
		sendError(w, "Failed to write .meta content", http.StatusInternalServerError)
		return
	}
	metaWc.Close()

	// Upload the resized image to the destination bucket
	destObj := bucket.Object(fmt.Sprintf("%s-%dx%d.%s", imageData.Filename, newWidth, newHeight, extension))
	destWc := destObj.NewWriter(ctx)
	if err := imaging.Encode(destWc, resizedImage, imaging.Format(extension)); err != nil {
		sendError(w, "Failed to encode image", http.StatusInternalServerError)
		return
	}
	destWc.Close()

	sendResponse(w, "Ok", http.StatusOK, newWidth, newHeight)
}

func calculateNewDimensions(currentWidth, currentHeight, desiredWidth, desiredHeight int) (int, int) {
	if desiredWidth == 0 {
		// Calculate new width while maintaining the aspect ratio
		newWidth := int(float64(desiredHeight) / float64(currentHeight) * float64(currentWidth))
		return newWidth, desiredHeight
	}

	if desiredHeight == 0 {
		// Calculate new height while maintaining the aspect ratio
		newHeight := int(float64(desiredWidth) / float64(currentWidth) * float64(currentHeight))
		return desiredWidth, newHeight
	}

	// Use the desired width and height without maintaining the aspect ratio
	return desiredWidth, desiredHeight
}
