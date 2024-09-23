package main

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"

	"example/http-server-go/pkgs/types"
)

var albums = []types.Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	album, found := findAlbumByField("ID", id)
	if !found {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "album not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, album)
}

func findAlbumByField(field string, value string) (types.Album, bool) {
	if !slices.Contains(types.AllowedAlbumFields, field) {
		return types.Album{}, false
	}

	album_index, found := slices.BinarySearchFunc(
		albums,
		value,
		func(a types.Album, getByData string) int {
			// Field is check beforehand
			data, _ := a.GetByField(field)
			return strings.Compare(data, getByData)
		},
	)
	if !found {
		return types.Album{}, false
	}

	return albums[album_index], found
}

func postAlbums(c *gin.Context) {
	var newAlbum types.Album

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	albums = append(albums, newAlbum)
	sortFunc := func(a, b types.Album) int { return strings.Compare(a.ID, b.ID) }

	if !slices.IsSortedFunc(albums, sortFunc) {
		slices.SortFunc(albums, sortFunc)
	}
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func AuthRequired(context *gin.Context) {
	auth := context.Request.Header.Get("Authorization")
	if auth != "MAGICSTRING" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, types.ErrorMessage{Error: "Invalid Authorization header"})
		return
	}

	context.Next()
}

func main() {
	router := gin.Default()
	authorized := router.Group("/")
	authorized.Use(AuthRequired)
	{
		authorized.GET("/albums", getAlbums)
		authorized.GET("/albums/:id", getAlbumByID)
		authorized.POST("/albums", postAlbums)
	}

	router.Run("localhost:8080")
}
