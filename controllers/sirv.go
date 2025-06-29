package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// SirvBaseURL is the base URL for accessing uploaded files on Sirv.
// You should replace "your-sirv-account-name" with your actual Sirv account name.
const SirvBaseURL = "https://old-souqs.sirv.com/" // Updated to your Sirv account base URL

// AuthRequest represents the JSON payload for the Sirv authentication request.
type AuthRequest struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

// AuthResponse represents the expected JSON response from the Sirv authentication endpoint.
type AuthResponse struct {
	Token      string `json:"token"`
	ExpiresIn  int    `json:"expiresIn"`  // Optional: to see how long the token is valid
	StatusCode int    `json:"statusCode"` // Optional: for error responses
	Error      string `json:"error"`      // Optional: for error responses
	Message    string `json:"message"`    // Corrected tag to "message"
}

// getSirvToken retrieves an authentication token from the Sirv API.
func GetSirvToken() (string, error) {
	// Initialize an HTTP client with a timeout for robustness.
	client := &http.Client{
		Timeout: 10 * time.Second, // Set a reasonable timeout
	}

	// Define the authentication request payload.
	authPayload := AuthRequest{
		ClientID:     "YJSw6mQ8yagO4n37YEXPhKto3kE",
		ClientSecret: "i0G1wKuzM+qa7VLV3PCaZJjwyONW+J4bdZNoCM+WUgpSdFktUZNR3SqDDLFUxtvrm0/HVLOxlPRwORLl9L70xg==",
	}

	// Marshal the struct into a JSON byte array.
	jsonPayload, err := json.Marshal(authPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON payload: %v", err)
	}

	// Create a new POST request with the JSON payload.
	// Corrected the URL string.
	req, err := http.NewRequest("POST", "https://api.sirv.com/v2/token", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set the Content-Type header to application/json, which is crucial for the Sirv API.
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request.
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request failed: %v", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed.

	// Read the response body.
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Print the raw response for debugging purposes.
	fmt.Println("Sirv raw response:", string(bodyBytes))

	// Check the HTTP status code directly.
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Sirv API HTTP error: Status %d, Response: %s", resp.StatusCode, string(bodyBytes))
	}

	// Unmarshal the JSON response into our AuthResponse struct.
	var result AuthResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse token response JSON: %v", err)
	}

	// Check for API-specific errors in the response body (e.g., if statusCode is present and non-2xx).
	if result.StatusCode != 0 && (result.StatusCode < 200 || result.StatusCode >= 300) {
		return "", fmt.Errorf("Sirv API error in response body: Status %d, Error: %s, Message: %s", result.StatusCode, result.Error, result.Message)
	}

	// Check if the token was successfully received.
	if result.Token == "" {
		return "", fmt.Errorf("Sirv token is empty in the response")
	}

	return result.Token, nil
}

// uploadToSirv uploads a file to Sirv.
func UploadToSirv(filePath, fileName, token string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	// Construct the upload URL. Ensure the path starts with a slash for Sirv.
	// Sirv expects the filename parameter to be the full path on Sirv.
	uploadURL := "https://api.sirv.com/v2/files/upload?filename=/Products/" + fileName
	req, err := http.NewRequest("POST", uploadURL, file)
	if err != nil {
		return fmt.Errorf("failed to create upload request: %w", err)
	}

	// Set necessary headers.
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/octet-stream") // Correct for binary file upload

	// Log truncated token for security.
	logToken := token
	if len(token) > 10 {
		logToken = token[:10] + "..."
	}
	fmt.Println("Uploading to Sirv with token:", logToken)
	fmt.Println("Upload URL:", uploadURL)
	fmt.Println("File name:", fileName)
	fmt.Println("Upload headers:", req.Header)

	client := &http.Client{
		Timeout: 30 * time.Second, // Increased timeout for file uploads
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Sirv upload request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read upload response body: %w", err)
	}
	fmt.Println("Upload response:", string(respBody))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Sirv upload failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	fmt.Println("File uploaded successfully to Sirv!")
	return nil
}

// cleanFileName removes spaces and potentially other problematic characters from a filename.
func cleanFileName(name string) string {
	// Replace spaces with underscores. You might want to extend this for other characters.
	cleaned := strings.ReplaceAll(name, " ", "_")
	// You could add more sanitization here, e.g., using a regex to allow only alphanumeric, dashes, underscores, and dots.
	return cleaned
}

// UploadImageToSirv handles the HTTP request for image uploads.
func UploadImageToSirv(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form data, with a 10MB limit.
	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err) // Log the actual error
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve the file from the form. The form field name is "file".
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file from form: %v", err) // Log the actual error
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close() // Ensure the uploaded file is closed.

	// Clean the filename to be Sirv-friendly.
	fileName := cleanFileName(handler.Filename)

	// Create a temporary file to store the uploaded content.
	// os.TempDir() provides a cross-platform temporary directory.
	tempFile, err := os.CreateTemp(os.TempDir(), "sirv-upload-*.tmp") // Create a unique temp file
	if err != nil {
		log.Printf("Temp file creation error: %v", err)
		http.Error(w, "Temp file creation error", http.StatusInternalServerError)
		return
	}
	tempPath := tempFile.Name() // Get the full path of the temporary file
	defer os.Remove(tempPath)   // Ensure the temporary file is removed after use
	defer tempFile.Close()      // Ensure the temporary file handle is closed

	// Copy the uploaded file content to the temporary file.
	_, err = io.Copy(tempFile, file)
	if err != nil {
		log.Printf("Error copying file to temp location: %v", err)
		http.Error(w, "Error saving temporary file", http.StatusInternalServerError)
		return
	}

	// Get Sirv authentication token.
	token, err := GetSirvToken()
	if err != nil {
		log.Printf("Failed to get Sirv token: %v", err)
		http.Error(w, "Failed to get Sirv token", http.StatusInternalServerError)
		return
	}

	// Upload the temporary file to Sirv.
	err = UploadToSirv(tempPath, fileName, token)
	if err != nil {
		log.Printf("Error uploading to Sirv: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError) // Return Sirv-specific error message
		return
	}

	// Construct the public URL for the uploaded image.
	imageURL := SirvBaseURL + "Products/" + fileName // Ensure the path matches your upload path
	// Note: SirvBaseURL already includes "/Products/", so you might need to adjust this depending on your Sirv folder structure.
	// If SirvBaseURL is "https://old-souqs.sirv.com/", then for a file in /Products/ you'd use SirvBaseURL + "Products/" + fileName.
	// If SirvBaseURL is "https://old-souqs.sirv.com/Products/", then you'd just use SirvBaseURL + fileName.
	// I've updated SirvBaseURL to "https://old-souqs.sirv.com/" and adjusted the imageURL construction accordingly.

	// Respond with the image URL as JSON.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": imageURL})
	fmt.Println("Image uploaded successfully and URL returned:", imageURL)
}
