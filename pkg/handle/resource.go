package handle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/zulubit/mimi/pkg/dejson"
	"github.com/zulubit/mimi/pkg/render"
)

func GetResource(w http.ResponseWriter, r *http.Request) {

	// Extract the 'slug' from the URL
	vars := mux.Vars(r)
	slug := vars["slug"]

	// Find the page by slug
	pageData, found, err := findpage(slug)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Printf("Error finding page: %v\n", err)
		return
	}

	if !found {
		http.NotFound(w, r)
		fmt.Printf("Page not found: %s\n", slug)
		return
	}

	// Parse and validate the page using dejson
	page, err := dejson.ParsePage([]byte(pageData))
	if err != nil {
		http.Error(w, "Error parsing page JSON", http.StatusInternalServerError)
		fmt.Printf("Error parsing page JSON: %v\n", err)
		return
	}

	// Define SEO and content data
	seo := render.Seo{
		Title:       page.SEO.Title,
		Description: page.SEO.Description,
		Keywords:    strings.Join(page.SEO.Keywords, ", "),
	}

	// content := render.Content{
	// 	Element: "my-element",
	// 	Data:    template.HTML(minifiedJSON.String()), // Prevent escaping
	// }

	contentSlice := []render.Content{}

	for i, c := range page.Data {
		encodedData, err := json.Marshal(c.Data)
		if err != nil {
			return
		}

		if c.Renderable {
			nc := render.Content{
				Element: c.Template,
				Data:    string(encodedData),
				Index:   i,
			}

			contentSlice = append(contentSlice, nc)
		}
	}

	// Render the page using RenderPage
	renderedPage, err := render.RenderPage(seo, contentSlice)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		fmt.Printf("Error rendering page: %v\n", err)
		return
	}

	// Write the rendered HTML
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(renderedPage)); err != nil {
		fmt.Printf("Error writing response: %v\n", err)
	}
}

func findpage(slug string) (string, bool, error) {
	pages, err := mapPagesFromResources()
	if err != nil {
		return "", false, err
	}

	for _, p := range *pages {
		if p.slug == slug {
			return p.data, true, nil
		}
	}

	return "", false, nil
}

func mapPagesFromResources() (*[]Page, error) {
	cdir := "./sitedata/resources/pages"
	dir, err := os.ReadDir(cdir)
	if err != nil {
		return nil, err
	}

	pages := []Page{}

	for _, r := range dir {
		if r.IsDir() {
			continue
		}

		fileName := r.Name()
		baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

		content, err := readPage(path.Join(cdir, fileName))
		if err != nil {
			return nil, err
		}

		pages = append(pages, Page{
			slug: baseName,
			data: *content,
		})
	}

	return &pages, nil
}

func readPage(filePath string) (*string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	data := string(content)
	return &data, nil
}

type Page struct {
	slug string
	data string
}
