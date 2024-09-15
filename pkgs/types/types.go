package types

import "fmt"

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var AllowedAlbumFields = []string{
	"ID",
	"Title",
	"Artist",
	"Price",
}

func (a *Album) GetByField(field string) (string, bool) {
	switch field {
	case "ID":
		return a.ID, true
	case "Title":
		return a.Title, true
	case "Artist":
		return a.Artist, true
	case "Price":
		return fmt.Sprintf("%.2f", a.Price), true
	default:
		return "", false
	}
}
