package home

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stonksdex/externalapi/api/common"
	"github.com/stonksdex/externalapi/codes"
	"github.com/stonksdex/externalapi/fullindex"
	"github.com/stonksdex/externalapi/model"
	"github.com/stonksdex/externalapi/system"
)

func Search(c *gin.Context) {
	res := common.Response{}
	res.Timestamp = time.Now().Unix()
	res.Code = codes.CODE_SUCCESS
	res.Msg = "success"

	searchKey := c.Param("key")
	sizeStr := c.DefaultQuery("size", "50")

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		size = 50
	}

	if len(searchKey) == 0 {
		c.JSON(http.StatusOK, res)
		return
	}

	result, err := fullindex.SearchIndex(searchKey, size, true)
	if err != nil {
		res.Code = codes.CODE_ERR_UNKNOWN
		res.Msg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusOK, res)
		return
	}

	catIds := retrieveIds(result)
	tokenIds := catIds["token"]
	if len(tokenIds) > 0 {
		db := system.GetDb()
		var tokenMetas []model.TokenMeta
		db.Model(&model.TokenMeta{}).Where(" id in ? ", tokenIds).Find(&tokenMetas)
		res.Data = tokenMetas
	}

	c.JSON(http.StatusOK, res)
}

func retrieveIds(result []fullindex.SearchResult) map[string][]string {
	ids := sync.Map{}
	mutexes := sync.Map{}

	var wg sync.WaitGroup
	for _, v := range result {
		wg.Add(1)
		go func(v fullindex.SearchResult) {
			defer wg.Done()

			a := strings.Split(v.ID, "-")
			if len(a) == 2 {
				lockInterface, _ := mutexes.LoadOrStore(a[0], &sync.Mutex{})
				mu := lockInterface.(*sync.Mutex)
				mu.Lock()

				val, _ := ids.Load(a[0])
				var list []string
				if val != nil {
					list = val.([]string)
				}
				list = append(list, a[1])

				ids.Store(a[0], list)

				mu.Unlock()
			}
		}(v)
	}

	wg.Wait()

	finalResult := make(map[string][]string)
	ids.Range(func(key, value interface{}) bool {
		finalResult[key.(string)] = value.([]string)
		return true
	})

	return finalResult
}
