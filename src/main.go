package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

// ResponseData defines the structure of our JSON response
type ResponseData struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type GetData struct {
	Name string `json:"name"`
}

type UpdateData struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func isProductionMode() bool {
	var value, exists = os.LookupEnv("WOSP_MODE")
	MODE := "production" // production
	if exists {
		MODE = value
	}
	return MODE == "production"
}

func notAcceptable(w http.ResponseWriter) {
	// common response to not allow testing in production

	response := ResponseData{
		Message: "MODE==production, no testing allowed",
		Status:  "error",
	}

	// Set headers and encode response to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotAcceptable)
	json.NewEncoder(w).Encode(response)
}

func recursiveCreateTestSecret(count int) string {
	if isProductionMode() {
		// dont allow test in production
		return string('0')
	} else {
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
}

func handleTestCreateSecret(w http.ResponseWriter, r *http.Request) {
	if isProductionMode() {
		notAcceptable(w)
	} else {
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
}

func handleTestGetSecretValue(w http.ResponseWriter, r *http.Request) {
	if isProductionMode() {
		notAcceptable(w)
	} else {
		ctx := context.TODO()

		// Strictly enforce the POST method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read put parameters (e.g., name, value)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var data GetData
		err = json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		getSecretValueInput := &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(data.Name),
		}

		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			panic(err)
		}

		client := secretsmanager.NewFromConfig(cfg)

		response, err := client.GetSecretValue(ctx, getSecretValueInput)
		if err != nil {
			log.Println("failed to get secret value, %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(
				ResponseData{
					Message: fmt.Sprintf("failed to get secret value, %v", err),
					Status:  "error",
				},
			)
		} else {
			// Set headers and encode response to JSON
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}
	}
}

func handleListSecrets(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	// Strictly enforce the GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read query parameters (e.g., /api?contains=Golang)
	contains := r.URL.Query().Get("contains")
	max_str := r.URL.Query().Get("max")
	max_parsed, err := strconv.ParseInt(max_str, 10, 32)
	maxResults := int32(20)
	if err == nil {
		maxResults = int32(max_parsed)
	}
	listSecretsInput := &secretsmanager.ListSecretsInput{
		MaxResults: &maxResults,
	}
	if contains != "" {
		listSecretsInput = &secretsmanager.ListSecretsInput{
			Filters: []types.Filter{
				{
					Key:    types.FilterNameStringTypeName,
					Values: []string{contains},
				},
			},
			MaxResults: &maxResults,
		}
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	client := secretsmanager.NewFromConfig(cfg)

	response, err := client.ListSecrets(ctx, listSecretsInput)
	if err != nil {
		log.Println("failed to list secrets, %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			ResponseData{
				Message: fmt.Sprintf("failed to list secrets, %v", err),
				Status:  "error",
			},
		)
	} else {
		// Set headers and encode response to JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func handlePutSecretValue(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	// Strictly enforce the PUT method
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read put parameters (e.g., name, value)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data UpdateData
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	putSecretValueInput := &secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(data.Name),
		SecretString: aws.String(data.Value),
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	client := secretsmanager.NewFromConfig(cfg)

	response, err := client.PutSecretValue(ctx, putSecretValueInput)
	if err != nil {
		log.Println("failed to put secret value, %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			ResponseData{
				Message: fmt.Sprintf("failed to put secret value, %v", err),
				Status:  "error",
			},
		)
	} else {
		// Set headers and encode response to JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	// Register the route and handler function
	http.HandleFunc("/api/aws/secretsmanager/listsecrets", handleListSecrets)
	http.HandleFunc("/api/aws/secretsmanager/putsecretvalue", handlePutSecretValue)

	http.HandleFunc("/test/aws/create", handleTestCreateSecret)
	http.HandleFunc("/test/aws/get", handleTestGetSecretValue)

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

	log.Println("Server starting on 127.0.0.1:8888...")
	if err := http.ListenAndServe("127.0.0.1:8888", nil); err != nil {
		log.Fatal(err)
	}
}
