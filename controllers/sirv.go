package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	SirvClientID     = "YJSw6mQ8yagO4n37YEXPhKto3kE"
	SirvClientSecret = "i0G1wKuzM+qa7VLV3PCaZJjwyONW+J4bdZNoCM+WUgpSdFktUZNR3SqDDLFUxtvrm0/HVLOxlPRwORLl9L70xg=="
	SirvBaseURL      = "https://old-souqs.sirv.com/Products/"
)

func getSirvToken() (string, error) {
	client := &http.Client{}

	data := "clientId=YJSw6mQ8yagO4n37YEXPhKto3kE&clientSecret=i0G1wKuzM+qa7VLV3PCaZJjwyONW+J4bdZNoCM+WUgpSdFktUZNR3SqDDLFUxtvrm0/HVLOxlPRwORLl9L70xg=="
	req, err := http.NewRequest("POST", "https://api.sirv.com/v2/token", strings.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // ✅ MUST HAVE

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("Sirv raw response:", string(bodyBytes)) // 🔍 log response

	var result struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	fmt.Println("Sirv token:", result.Token)
	return result.Token, nil
}

func uploadToSirv(filePath, fileName, token string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	url := "https://api.sirv.com/v2/files/upload?filename=/Products/" + fileName
	req, err := http.NewRequest("POST", url, file)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token) // ✅ Correct casing
	req.Header.Set("Content-Type", "application/octet-stream")

	fmt.Println("Uploading to Sirv with token:", token[:10]+"...")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Upload response:", string(respBody)) // 🔍 See Sirv response

	if resp.StatusCode != 200 {
		return fmt.Errorf("Sirv upload failed: %s", string(respBody))
	}
	fmt.Println("Token:", token)
	fmt.Println("Upload URL:", url)
	fmt.Println("File name:", fileName)
	fmt.Println("Upload headers:", req.Header)

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
