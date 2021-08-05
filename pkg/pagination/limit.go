package pagination

import (
	"net/http"
	"strconv"
)

var defaultLimit = 15

func GetLimitParam(r *http.Request) (limit int) {
	limit = defaultLimit
	if limitQ := r.URL.Query().Get("limit"); limitQ != "" {
		if intLimit, err := strconv.Atoi(limitQ); err == nil {
			if intLimit <= 50 {
				limit = intLimit
			}
		}
	}
	return
}
