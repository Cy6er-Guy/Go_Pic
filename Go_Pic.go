package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/Jeffail/gabs/v2"
)

const (
	apiKey            = "YOUR_API_KEY"
	searchURL         = "https://api.pexels.com/v1/search"
	downloadFolder    = "downloaded_images"
	serverPort        = ":5050"
	maxResultsPerPage = 3
)

func main() {
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/search", searchHandler)
	fmt.Printf("Serving on http://localhost%s\n", serverPort)
	http.ListenAndServe(serverPort, nil)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, `
        <!DOCTYPE html>
        <html>
        <head>
            <title>Image Keyword Search</title>
        </head>
        <body>
            <form action="/search" method="post">
                <textarea name="text" rows="10" cols="50"></textarea><br>
                <input type="submit" value="Submit">
            </form>
        </body>
        </html>
    `)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, "Searching for images...\n")

	os.Mkdir(downloadFolder, 0755)
	pairs := splitIntoPairs(text)

	for _, pair := range pairs {
		if len(pair) != 2 {
			continue
		}
		query := fmt.Sprintf("%s %s", pair[0], pair[1])
		searchImages(query, w)
	}

	fmt.Fprintln(w, "Image search complete.")
}

func splitIntoPairs(text string) [][]string {
	// Split text into pairs of words
	words := strings.Fields(text)
	var pairs [][]string
	for i := 0; i < len(words)-1; i++ {
		pairs = append(pairs, []string{words[i], words[i+1]})
	}
	return pairs
}

func searchImages(query string, w http.ResponseWriter) {
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		fmt.Fprintf(w, "Error creating request: %v\n", err)
		return
	}

	q := req.URL.Query()
	q.Add("query", query)
	q.Add("per_page", fmt.Sprintf("%d", maxResultsPerPage))
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(w, "Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(w, "Error response from API: %v\n", resp.Status)
		return
	}

	jsonParsed, err := gabs.ParseJSONBuffer(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "Error parsing JSON: %v\n", err)
		return
	}

	photos := jsonParsed.Path("photos").Children()
	for _, photo := range photos {
		photoID := int(photo.Path("id").Data().(float64))
		photographer := photo.Path("photographer").Data().(string)
		photoURL := photo.Path("src.original").Data().(string)

		fmt.Fprintf(w, "Photo ID: %d\n", photoID)
		fmt.Fprintf(w, "Photographer: %s\n", photographer)
		fmt.Fprintf(w, "URL: %s\n\n", photoURL)

		downloadImage(photoURL, query, photoID, w)
	}
}

func downloadImage(url, query string, photoID int, w http.ResponseWriter) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(w, "Error downloading image %d: %v\n", photoID, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(w, "Error response when downloading image %d: %v\n", photoID, resp.Status)
		return
	}

	// Extract file extension from URL
	ext := getFileExtension(url)
	if ext == "" {
		ext = ".jpg" // Default extension if unable to determine
	}

	// Replace invalid characters in query with underscores
	cleanQuery := sanitizeFilename(query)

	filePath := fmt.Sprintf("%s/%s_%d%s", downloadFolder, cleanQuery, photoID, ext)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Fprintf(w, "Error creating file for image %d: %v\n", photoID, err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Fprintf(w, "Error saving image %d: %v\n", photoID, err)
		return
	}

	fmt.Fprintf(w, "Image downloaded and saved to %s\n\n", filePath)
}

func getFileExtension(url string) string {
	parts := strings.Split(url, ".")
	if len(parts) > 1 {
		return "." + parts[len(parts)-1]
	}
	return ""
}

func sanitizeFilename(filename string) string {
	// Replace characters that are invalid in filenames with underscores
	invalidChars := regexp.MustCompile(`[^\w\-. ]`)
	return invalidChars.ReplaceAllString(filename, "_")
}
