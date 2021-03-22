package accessory

import (
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/homepi/homepi/api/app/components/gpio"
)

func (s *Service) GetGpioPins(ctx *gin.Context) {
	ctx.JSON(respond.Default.Succeed(gpio.GetPins(s.db)))
	return
}
