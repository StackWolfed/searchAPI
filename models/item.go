package models

// Item is the actual DB representation
type Item struct {
	ID        int64  `json:"id"`
	LowerName string `json:"lowerName"`
	Name      string `json:"name"`
	URL       string `json:"url"`
}

// Items is the collection type
type Items []Item
