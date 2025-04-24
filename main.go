package main

import (
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("Missing arg. Usage: go run main.go <url>")
		return
	}

	url := os.Args[1]

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating request:\n", err)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:\n", err)
	}
	defer res.Body.Close()

	doc, err := html.Parse(res.Body)
	if err != nil {
		log.Fatal("Error parsing the html response:\n", err)
	}

	jsonData := traverseParsedHtml(doc)

	err = writeStringToFile(jsonData)
	if err != nil {
		log.Fatal("Error creating file:\n", err)
	}

}

func writeStringToFile(jsonData string) error {
	out, err := os.Create("output.json")
	if err != nil {
		return err
	}
	defer out.Close()

	out.WriteString(jsonData)

	return nil
}

func traverseParsedHtml(node *html.Node) string {
	htmlNode := node.FirstChild.NextSibling

	bodyNode := htmlNode.FirstChild.NextSibling

	for c := bodyNode.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == "script" && attrContains(c.Attr, "id", "__NEXT_DATA__") {
			return c.FirstChild.Data
		}
	}

	return ""
}

func attrContains(attrs []html.Attribute, key string, value string) bool {
	for _, attr := range attrs {
		if attr.Key == key && attr.Val == value {
			return true
		}
	}
	return false
}
