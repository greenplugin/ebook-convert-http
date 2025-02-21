package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const tempDir = "./temp"

type Convert struct {
}

func (c *Convert) Register(mux *http.ServeMux) {
	mux.HandleFunc("/convert", c.handleConversion)
}

func NewConvert() *Convert {
	// Temporary directory for saving files
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		if err := os.Mkdir(tempDir, 0755); err != nil {
			log.Fatalf("Error creating temporary directory: %s", err)
		}
	}

	return &Convert{}
}

// Route handler for processing the file
func (c *Convert) handleConversion(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are supported", http.StatusMethodNotAllowed)
		return
	}

	// Read the conversion format from the "convert-to" header
	convertTo := r.FormValue("convert-to")
	if convertTo == "" {
		http.Error(w, "Target format not specified (convert-to parameter)", http.StatusBadRequest)
		return
	}

	// Read the uploaded file
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Generate a temporary file for input data
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		http.Error(w, "Error: file must have an extension", http.StatusBadRequest)
		return
	}
	inputFile, err := os.CreateTemp(tempDir, fmt.Sprintf("input-*%s", ext))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating temporary file: %v", err), http.StatusInternalServerError)
		return
	}
	defer os.Remove(inputFile.Name()) // Delete the file after completion

	log.Printf("File uploaded: %s", inputFile.Name())

	// Copy the contents of the input file
	if _, err := io.Copy(inputFile, file); err != nil {
		http.Error(w, fmt.Sprintf("Error writing file: %v", err), http.StatusInternalServerError)
		return
	}

	// Close the file for conversion
	if err := inputFile.Close(); err != nil {
		http.Error(w, fmt.Sprintf("Error closing file: %v", err), http.StatusInternalServerError)
		return
	}

	// Generate a temporary file for output data
	outputFile, err := os.CreateTemp(tempDir, fmt.Sprintf("output-*.%s", convertTo))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating output file: %v", err), http.StatusInternalServerError)
		return
	}
	defer os.Remove(outputFile.Name()) // Delete the file after completion
	outputFile.Close()

	cmdArgs := []string{"--no-sandbox", inputFile.Name(), outputFile.Name()}

	// Get additional parameters from the request
	options := r.FormValue("convert-options")
	if options != "" {
		// Split the options string in "key=value" format by spaces
		// Example: "margin-top=12 margin-right=10"
		for _, opt := range strings.Split(options, " ") {
			if opt != "" {
				cmdArgs = append(cmdArgs, "--"+opt)
			}
		}
	}
	// Execute the `ebook-convert` command to convert the file
	cmd := exec.Command("ebook-convert", cmdArgs...)
	cmd.Stderr = os.Stderr // Перенаправляем ошибки утилиты в вывод
	if err := cmd.Run(); err != nil {
		http.Error(w, fmt.Sprintf("Error converting file: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("File successfully converted: %s", outputFile.Name())

	// Prepare the file for sending to the client
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.%s", strings.TrimSuffix(header.Filename, ext), convertTo))
	resultFile, err := os.Open(outputFile.Name())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading output file: %v", err), http.StatusInternalServerError)
		return
	}
	defer resultFile.Close()

	// Send the file in the response
	if _, err := io.Copy(w, resultFile); err != nil {
		http.Error(w, fmt.Sprintf("Error sending file: %v", err), http.StatusInternalServerError)
		return
	}
}
