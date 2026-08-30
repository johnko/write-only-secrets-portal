package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"log"
	"net/http"
)

// ResponseData defines the structure of our JSON response
type ResponseData struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func recursiveCreateTestSecret(count int) string {
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	client := secretsmanager.NewFromConfig(cfg)

	secretName := fmt.Sprintf("my-app-test-creds-%d", count)
	secretString := `{"username":"test","password":"test"}`
	input := &secretsmanager.CreateSecretInput{
		Name:         aws.String(secretName),
		SecretString: aws.String(secretString),
		Description:  aws.String("Database credentials created via Go SDK"),
	}
	_, err = client.CreateSecret(ctx, input)
	if err != nil {
		count = count + 1
		// log.Println("failed to create secret, %v", err)
		return recursiveCreateTestSecret(count)
	} else {
		// log.Println(result)
		return secretName
	}
}

func handleTestCreateTestSecretRequest(w http.ResponseWriter, r *http.Request) {
	// Strictly enforce the POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	secretName := recursiveCreateTestSecret(0)

	// Prepare data structure
	response := ResponseData{
		Message: fmt.Sprintf("Created test secret %s", secretName),
		Status:  "success",
	}

	// Set headers and encode response to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	// Strictly enforce the GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read query parameters (e.g., /api?contains=Golang)
	contains := r.URL.Query().Get("contains")
	listSecretsInput := &secretsmanager.ListSecretsInput{}
	if contains != "" {
		listSecretsInput = &secretsmanager.ListSecretsInput{
			Filters: []types.Filter{
				{
					Key:    types.FilterNameStringTypeName,
					Values: []string{contains},
				},
			},
		}
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	client := secretsmanager.NewFromConfig(cfg)

	response, err := client.ListSecrets(ctx, listSecretsInput)
	if err != nil {
		log.Fatalf("failed to list secrets, %v", err)
	}

	// Set headers and encode response to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Register the route and handler function
	http.HandleFunc("/api/aws/secretsmanager/listsecrets", handleGetRequest)
	http.HandleFunc("/test/aws", handleTestCreateTestSecretRequest)

	// Create a file server handler pointing to your local directory
	fileServer := http.FileServer(http.Dir("./aws"))

	// Strip the URL prefix so the file server looks for files correctly
	http.Handle("/aws/", http.StripPrefix("/aws/", fileServer))

	http.HandleFunc("/aws", func(w http.ResponseWriter, r *http.Request) {
		// Target URL can be relative to domain ("/new-path") or absolute ("https://example.com")
		// Redirects the client temporarily to a new URL using HTTP 302 via http.StatusFound
		http.Redirect(w, r, "/aws/", http.StatusFound)
	})

	// Default path redirects to web page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/aws/", http.StatusFound)
	})

	log.Println("Server starting on :8888...")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal(err)
	}
}
