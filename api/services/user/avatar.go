package user

import (
	"fmt"
	"log"
	"os"

	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/homepi/homepi/api/app/components/strings"
	"github.com/homepi/homepi/api/db/models"
)

func (s *Service) UpdateAvatar(ctx *gin.Context) {

	var (
		user = ctx.MustGet("user").(*models.User)
	)

	_, file, err := ctx.Request.FormFile("avatar")
	if err != nil {
		log.Println(err)
		return
	}

	randomAvatarDigits := strings.RandomDigits(20)

	if err = ctx.SaveUploadedFile(file, fmt.Sprintf("./public/uploads/avatars/%s.png", randomAvatarDigits)); err != nil {

		ctx.JSON(respond.Default.
			SetStatusText("Failed!").
			SetStatusCode(400).
			RespondWithMessage("Upload failed. Please try again later!"))
		return
	}

	if user.Avatar != "default" {
		_ = os.Remove(fmt.Sprintf("./public/uploads/avatars/%s.png", user.Avatar))
	}

	s.db.Model(&user).Update("avatar", randomAvatarDigits)

	ctx.JSON(respond.Default.Succeed(user))
	return
}
