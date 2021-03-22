package pagination

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

var defaultLimit = 15

func GetLimitParam(ctx *gin.Context) (limit int) {
	limit = defaultLimit
	if limitQ := ctx.Query("limit"); limitQ != "" {
		if intLimit, err := strconv.Atoi(limitQ); err == nil {
			if intLimit <= 50 {
				limit = intLimit
			}
		}
	}
	return
}
