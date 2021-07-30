package pagination

import (
	"net/http"
	"reflect"
	"strconv"
)

var itemsPerPage = 15

type PaginatorData struct {
	CurrentPage int         `json:"current_page"`
	Total       int         `json:"total"`
	PerPage     int         `json:"per_page"`
	Data        interface{} `json:"data"`
}

func GetPageParam(r *http.Request) (page int) {
	page = 1
	if pageQ := r.URL.Query().Get("page"); pageQ != "" {
		if intPage, err := strconv.Atoi(pageQ); err == nil {
			page = intPage
		}
	}
	return
}

// Create an array of data with limited items
func Paginator(data interface{}, page int) (results PaginatorData) {

	start := (page - 1) * itemsPerPage
	stop := start + itemsPerPage

	length := reflect.ValueOf(data).Len()

	results = PaginatorData{
		CurrentPage: page,
		Total:       0,
		PerPage:     itemsPerPage,
		Data:        []interface{}{},
	}

	if start > length {
		return
	}

	if stop > length {
		stop = length
	}

	if length != 0 {
		results.Data = reflect.ValueOf(data).Slice(start, stop).Interface()
	}

	results.Total = length
	return
}
