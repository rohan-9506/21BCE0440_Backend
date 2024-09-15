package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/api/upload", UploadHandler)
	return r
}

func TestUploadHandler(t *testing.T) {
	router := setupRouter()

	// Create a buffer to simulate a file upload
	fileContent := bytes.NewBufferString("This is a test file content")
	req := httptest.NewRequest(http.MethodPost, "/api/upload", fileContent)
	req.Header.Set("Content-Type", "multipart/form-data")

	// Add a file to the request
	formFile, err := req.MultipartForm()
	if err != nil {
		t.Fatal(err)
	}
	file := formFile.File["file"][0]
	req.AddFile(file)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response body
	expectedResponse := `{"message":"File uploaded successfully"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}
