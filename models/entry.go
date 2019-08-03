package models

// Entry is the item with the minimum requirements: Name and URL
type Entry struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
