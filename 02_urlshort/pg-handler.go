package urlshort

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/joho/godotenv"
)

// UrlMap is the db schema for url mappings
type UrlMap struct {
	tableName struct{} `pg:"url_map"`
    Id     	int64
    Name  	string
    Url 	string
}

func (urlMap UrlMap) String() string {
    return fmt.Sprintf("UrlMap<%d %s %v>", urlMap.Id, urlMap.Name, urlMap.Url)
}

func readURLsFromDb() []UrlMap {
	_ = godotenv.Load()

    db := pg.Connect(&pg.Options{
        User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Database: os.Getenv("DB_NAME"),
    })
    defer db.Close()

    var urls []UrlMap

    err := db.Model(&urls).Select()
    if err != nil {
        panic(err)
    }

    return urls
}

func buildPathMap(urlMap []UrlMap) PathMap {
	pathMap := make(PathMap)
	for _, item := range urlMap {
		pathMap[item.Name] = item.Url
	}
	return pathMap
}

// PgHandler is a URL map handler reading data from a Postgres database
func PgHandler (fallback http.Handler) http.HandlerFunc {
	dbEntries := readURLsFromDb()

	pathMap := buildPathMap(dbEntries)

	return MapHandler(pathMap, fallback)
}
