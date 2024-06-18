

```markdown
# Image Keyword Search and Download Program

This program allows you to search for images on Pexels based on keywords extracted from text input. It then downloads these images to your local machine.

## Prerequisites

Before running the program, make sure you have the following:

- Go installed on your system: [Go Installation Guide](https://golang.org/doc/install)
- Pexels API key: You can obtain one by signing up on [Pexels](https://www.pexels.com/api/documentation/) (replace `YOUR_API_KEY` in the code with your actual API key)

## Setup

1. Clone this repository:
   ```bash
   git clone https://github.com/Cy6er-Guy/Go_Pic.git
   cd Go_Pic
   ```

2. Install dependencies:
   ```bash
   go get github.com/Jeffail/gabs/v2
   ```

3. Set your Pexels API key:
   - Open `Go_Pic.go` and replace `YOUR_API_KEY` with your actual Pexels API key.

## Usage

1. Run the server:
   ```bash
   go run Go_Pic.go
   ```

2. Access the web interface:
   - Open your web browser and go to `http://localhost:5050`
   - Enter your text in the provided textarea and click "Submit".

3. Results:
   - The program will analyze the text, search for keywords in pairs of words, and download corresponding images from Pexels.
   - Images will be saved in the `downloaded_images` folder.

## Notes

- The program searches for images based on pairs of consecutive words extracted from the input text.
- Each pair of words is treated as a search query.
- Images are downloaded and saved with filenames based on the search query.

