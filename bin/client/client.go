package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"example/http-server-go/pkgs/types"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/spf13/cobra"
)

var URL = "http://localhost:8080"

func generateNewAlbum() types.Album {
	album := types.Album{
		ID:     strconv.Itoa(gofakeit.Number(3, 999999)),
		Title:  gofakeit.ProductName(),
		Artist: gofakeit.Name(),
		Price:  gofakeit.Price(5, 60),
	}
	return album
}

func postAlbum(rootUrl string, album types.Album) (types.Album, error) {
	jsonData, err := json.Marshal(album)
	if err != nil {
		log.Fatal("Error while trying to Marshal the data")
		return album, err
	}

	resp, err := http.Post(rootUrl+"/albums", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error while trying to Post the data")
		return album, err
	}

	var target types.Album

	err = json.NewDecoder(resp.Body).Decode(&target)

	if err != nil {
		return types.Album{}, err
	}

	return target, nil
}

func postAlbumJson(rootUrl string, album types.Album) (string, error) {
	album, err := postAlbum(rootUrl, album)

	if err != nil {
		return "", err
	}

	b, err := json.MarshalIndent(album, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func getAlbums(rootUrl string) ([]types.Album, error) {
	resp, err := http.Get(rootUrl + "/albums")

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var target []types.Album

	err = json.NewDecoder(resp.Body).Decode(&target)

	if err != nil {
		return nil, err
	}

	return target, nil
}

func getAlbumsJson(rootUrl string) (string, error) {
	albums, err := getAlbums(rootUrl)

	if err != nil {
		return "", err
	}

	b, err := json.MarshalIndent(albums, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func getAlbumsCsv(rootUrl string) (string, error) {
	albums, err := getAlbums(rootUrl)
	if err != nil {
		return "", err
	}

	var csvData bytes.Buffer
	writer := csv.NewWriter(&csvData)

	header := []string{
		"ID",
		"Title",
		"Artist",
		"Price",
	}
	writer.Write(header)

	for _, album := range albums {
		id, _ := album.GetByField("ID")
		title, _ := album.GetByField("Title")
		artist, _ := album.GetByField("Artist")
		price, _ := album.GetByField("Price")
		row := []string{
			id,
			title,
			artist,
			price,
		}
		writer.Write(row)
	}
	writer.Flush()

	return csvData.String(), nil
}

func main() {
	var ofJson bool
	var ofCSV bool

	rootCmd := &cobra.Command{
		Use: "client",
	}
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list album",
		Run: func(cmd *cobra.Command, args []string) {

			var b string
			var err error

			if ofCSV {
				b, err = getAlbumsCsv(URL)
			} else {
				b, err = getAlbumsJson(URL)
			}

			if err != nil {
				log.Fatal(err)
				return
			}

			fmt.Println(b)
		},
	}

	listCmd.Flags().BoolVar(&ofJson, "json", false, "Output in JSON")
	listCmd.Flags().BoolVar(&ofCSV, "csv", false, "Output in CSV")
	listCmd.MarkFlagsMutuallyExclusive("json", "csv")

	newCmd := &cobra.Command{
		Use:   "new",
		Short: "new album",
		Run: func(cmd *cobra.Command, args []string) {
			b, err := postAlbumJson(URL, generateNewAlbum())

			if err != nil {
				log.Fatal(err)
				return
			}

			fmt.Println(b)
		},
	}
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.Execute()
}
