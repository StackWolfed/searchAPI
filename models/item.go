package models

// Item is the actual DB representation
type Item struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Items is the collection type
type Items []Item
