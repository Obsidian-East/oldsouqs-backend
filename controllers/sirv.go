package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	SirvClientID     = "YJSw6mQ8yagO4n37YEXPhKto3kE"
	SirvClientSecret = "i0G1wKuzM+qa7VLV3PCaZJjwyONW+J4bdZNoCM+WUgpSdFktUZNR3SqDDLFUxtvrm0/HVLOxlPRwORLl9L70xg=="
	SirvBaseURL      = "https://old-souqs.sirv.com/Products/"
)

func getSirvToken() (string, error) {
	data := url.Values{}
	data.Set("clientId", "YJSw6mQ8yagO4n37YEXPhKto3kE")
	data.Set("clientSecret", "i0G1wKuzM+qa7VLV3PCaZJjwyONW+J4bdZNoCM+WUgpSdFktUZNR3SqDDLFUxtvrm0/HVLOxlPRwORLl9L70xg==")

	resp, err := http.PostForm("https://api.sirv.com/v2/token", data)
	if err != nil {
		return "", fmt.Errorf("token request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("Sirv raw response:", string(bodyBytes)) // 🔍 debug output

	// Attempt to decode token
	var result struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	fmt.Println("Sirv token:", result.Token) // ✅ verify it's not empty

	return result.Token, nil
}

func uploadToSirv(filePath, fileName, token string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	uploadURL := "https://api.sirv.com/v2/files/upload?filename=/Products/" + fileName

	req, err := http.NewRequest("POST", uploadURL, file)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	fmt.Println("Sirv token:", token)
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Sirv upload failed: %s", string(body))
	}

	return nil
}

func cleanFileName(name string) string {
	return strings.ReplaceAll(name, " ", "_")
}

func UploadImageToSirv(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := cleanFileName(handler.Filename)
	tempPath := "/tmp/" + fileName

	// Save file temporarily
	tempFile, err := os.Create(tempPath)
	if err != nil {
		http.Error(w, "Temp file creation error", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()
	io.Copy(tempFile, file)

	// Get Sirv token
	token, err := getSirvToken()
	if err != nil {
		http.Error(w, "Failed to get Sirv token", http.StatusInternalServerError)
		return
	}

	// Upload to Sirv
	err = uploadToSirv(tempPath, fileName, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Clean up temp file
	os.Remove(tempPath)

	// Return the image URL
	imageURL := SirvBaseURL + fileName
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": imageURL})
}
