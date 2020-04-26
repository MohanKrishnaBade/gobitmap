package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/robtec/newsapi/api"
	"net/http"
	"os"
)

type counter struct {
	count int
}

var requestCount *counter

func init() {
	requestCount = &counter{}
}
func News(c *gin.Context) {
	res, _ := getNews(c)
	c.JSON(http.StatusOK, res)
}

func getNews(c *gin.Context) (*api.Response, error) {

	//get Query parameters
	searchQuery := c.Request.URL.Query().Get("searchKey")
	source := c.Request.URL.Query().Get("source")

	// prepare redis Key and get the info associated with given key.
	redisKey := searchQuery + source
	data, err := Get(redisKey)

	response := &api.Response{}
	if err == redis.Nil {
		httpClient := http.Client{}
		client, err := api.New(&httpClient, os.Getenv("NEWS_SECRET_KEY"), os.Getenv("NEW_APP_URL"))
		if err != nil {
			panic(err)
		}

		//Different options
		moreOpts := api.Options{Language: "en", Q: searchQuery, SortBy: "popularity", PageSize: 50, Sources: source}

		// Get Everything with options from above
		data, err := client.Everything(moreOpts)
		requestCount.count++
		if err == nil {
			redisData, _ := json.Marshal(data)
			if error1 := CreateUnless(redisKey, string(redisData)); error1 != nil {
				return data, error1
			}
		}
		fmt.Println(requestCount.count)
		return data, err
	} else {
		if err := json.Unmarshal([]byte(data), response); err != nil {
			return response, err
		}
	}
	return response, nil
}
