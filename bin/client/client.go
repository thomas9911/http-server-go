package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"example/http-server-go/pkgs/types"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/spf13/cobra"
)

var URL = "http://localhost:8080"

type RequestConfig struct {
	RootUrl   string
	client    http.Client
	SecretKey string
}

func (config *RequestConfig) UpdateFromEnv() {
	if value, ok := os.LookupEnv("SERVER_ROOT_URL"); ok {
		fmt.Println(value)
		config.RootUrl = value
	}

	if value, ok := os.LookupEnv("SERVER_SECRET_KEY"); ok {
		config.SecretKey = value
	}
}

func generateNewAlbum() types.Album {
	album := types.Album{
		ID:     strconv.Itoa(gofakeit.Number(3, 999999)),
		Title:  gofakeit.ProductName(),
		Artist: gofakeit.Name(),
		Price:  gofakeit.Price(5, 60),
	}
	return album
}

func postAlbum(requestConfig RequestConfig, album types.Album) (types.Album, error) {
	jsonData, err := json.Marshal(album)
	if err != nil {
		log.Fatal("Error while trying to Marshal the data")
		return album, err
	}

	request, err := http.NewRequest("POST", requestConfig.RootUrl+"/albums", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error while creating request")
		return album, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", requestConfig.SecretKey)

	resp, err := requestConfig.client.Do(request)

	if err != nil {
		log.Fatal("Error while trying to Post the data")
		return album, err
	}

	if resp.StatusCode != http.StatusCreated {
		var errorMessage types.ErrorMessage
		err = json.NewDecoder(resp.Body).Decode(&errorMessage)

		if err != nil {
			return album, errors.New("error while trying to get albums")
		}
		return album, fmt.Errorf("%s", errorMessage.Error)
	}

	var target types.Album

	err = json.NewDecoder(resp.Body).Decode(&target)

	if err != nil {
		return types.Album{}, err
	}

	return target, nil
}

func postAlbumJson(requestConfig RequestConfig, album types.Album) (string, error) {
	album, err := postAlbum(requestConfig, album)

	if err != nil {
		return "", err
	}

	b, err := json.MarshalIndent(album, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func getAlbums(requestConfig RequestConfig) ([]types.Album, error) {
	request, err := http.NewRequest("GET", requestConfig.RootUrl+"/albums", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", requestConfig.SecretKey)

	resp, err := requestConfig.client.Do(request)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorMessage types.ErrorMessage
		err = json.NewDecoder(resp.Body).Decode(&errorMessage)

		if err != nil {
			return nil, errors.New("error while trying to get albums")
		}
		return nil, fmt.Errorf("%s", errorMessage.Error)
	}

	var target []types.Album

	err = json.NewDecoder(resp.Body).Decode(&target)

	if err != nil {
		return nil, err
	}

	return target, nil
}

func getAlbumsJson(requestConfig RequestConfig) (string, error) {
	albums, err := getAlbums(requestConfig)

	if err != nil {
		return "", err
	}

	b, err := json.MarshalIndent(albums, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func getAlbumsCsv(requestConfig RequestConfig) (string, error) {
	albums, err := getAlbums(requestConfig)
	if err != nil {
		return "", err
	}

	var csvData bytes.Buffer
	writer := csv.NewWriter(&csvData)

	header := types.AllAlbumKeys
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

	// cfg := RequestConfig{RootUrl: URL, SecretKey: "MAGICSTRING"}
	cfg := RequestConfig{}

	rootCmd := &cobra.Command{
		Use: "client",
	}
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list album",
		Run: func(cmd *cobra.Command, args []string) {
			// cfg := RequestConfig{RootUrl: URL, SecretKey: "MAGICSTRING"}
			cfg.UpdateFromEnv()

			var b string
			var err error

			if ofCSV {
				b, err = getAlbumsCsv(cfg)
			} else {
				b, err = getAlbumsJson(cfg)
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
			// cfg := RequestConfig{RootUrl: URL, SecretKey: "MAGICSTRING"}
			cfg.UpdateFromEnv()

			b, err := postAlbumJson(cfg, generateNewAlbum())

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
