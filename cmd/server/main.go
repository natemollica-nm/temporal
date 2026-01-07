package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/natemollica-nm/temporal/internal/config"
	"github.com/natemollica-nm/temporal/internal/metrics"
	"github.com/natemollica-nm/temporal/pkg/temporal/shared"
	"github.com/natemollica-nm/temporal/pkg/temporal/workflows/basic"
	"go.temporal.io/sdk/client"

	sdktally "go.temporal.io/sdk/contrib/tally"
	sdklog "go.temporal.io/sdk/log"
)

var temporalClient client.Client

// Initialize Temporal Client
func initializeTemporal() error {
	// Load configuration
	cfg := config.Default()
	
	// Create logger - use nil for default logger
	var logger sdklog.Logger
	
	// Initialize metrics
	if err := metrics.Initialize(cfg.Metrics, logger); err != nil {
		return fmt.Errorf("failed to initialize metrics: %w", err)
	}

	var err error
	temporalClient, err = client.Dial(client.Options{
		HostPort:       cfg.Temporal.HostPort,
		Namespace:      cfg.Temporal.Namespace,
		MetricsHandler: sdktally.NewMetricsHandler(metrics.GetScope()),
		Logger:         logger,
	})
	return err
}

// Start the Temporal Workflow
func startWorkflow(name string) (string, error) {
	workflowID := "getAddressFromIP-" + uuid.New().String()
	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: shared.TaskQueueName,
	}

	we, err := temporalClient.ExecuteWorkflow(context.Background(), options, basic.GetAddressFromIP, name)
	if err != nil {
		return "", err
	}

	var result string
	err = we.Get(context.Background(), &result)
	return result, err
}

// Handle HTMX form submission
func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<p class="error">Name is required</p>`)
		return
	}

	result, err := startWorkflow(name)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<p class="error">Error: %s</p>`, html.EscapeString(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<p class="success">%s</p>`, html.EscapeString(result))
}

// Handle API request
func handleAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	var requestData struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	requestData.Name = strings.TrimSpace(requestData.Name)
	if requestData.Name == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Name is required"})
		return
	}

	result, err := startWorkflow(requestData.Name)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": result})
}

// Serve static files with proper MIME types
func serveStaticFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	// Security: prevent directory traversal
	if strings.Contains(path, "..") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	filePath := filepath.Join("web/static", path)

	// Set appropriate MIME type
	switch filepath.Ext(path) {
	case ".html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	http.ServeFile(w, r, filePath)
}

func main() {
	err := initializeTemporal()
	if err != nil {
		log.Fatalf("Failed to initialize Temporal client: %v", err)
	}

	http.HandleFunc("/submit", handleSubmit)
	http.HandleFunc("/api", handleAPI)
	http.HandleFunc("/", serveStaticFiles)

	port := 4000
	fmt.Printf("Server running on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
