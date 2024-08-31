package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// POST endpoint to generate PDF from HTML string
	r.POST("/generate-pdf", func(c *gin.Context) {
		log.Println("Received request to generate PDF")

		// Start timing the process
		startTime := time.Now()

		// Read HTML content from request body
		htmlBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("Failed to read request body:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}

		log.Println("HTML content read successfully")

		// Prepare the request to Gotenberg
		requestBody := new(bytes.Buffer)
		writer := multipart.NewWriter(requestBody)

		// Add the HTML file part
		part, err := writer.CreateFormFile("files", "index.html")
		if err != nil {
			log.Println("Failed to create form file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create form file"})
			return
		}

		part.Write(htmlBytes)
		writer.Close()

		// Send the request to Gotenberg
		req, err := http.NewRequest("POST", "http://localhost:3001/forms/chromium/convert/html", requestBody)
		if err != nil {
			log.Println("Failed to create request to Gotenberg:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request to Gotenberg"})
			return
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Failed to send request to Gotenberg:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to Gotenberg"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Println("Gotenberg returned non-OK status:", resp.Status)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gotenberg returned non-OK status"})
			return
		}

		// Read the PDF response from Gotenberg
		pdfBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Failed to read PDF response:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read PDF response"})
			return
		}

		// Calculate the duration
		duration := time.Since(startTime)
		log.Printf("PDF rendering completed in %v\n", duration)

		// Send the PDF file back to the frontend
		c.Header("Content-Disposition", "attachment; filename=generated.pdf")
		c.Data(http.StatusOK, "application/pdf", pdfBytes)
	})

	// Start the server
	r.Run(":3000")
}
