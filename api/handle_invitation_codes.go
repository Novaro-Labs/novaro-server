package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/auth"
	"novaro-server/service"
)

type InvitationCodesApi struct {
	service *service.InvitationCodesService
}

func NewInvitationCodesApi() *InvitationCodesApi {
	return &InvitationCodesApi{
		service: service.NewInvitationCodesService(),
	}
}

// MakeInvitationCodes godoc
// @Summary Generate new invitation codes
// @Description Generate new invitation codes and save to the system
// @Accept json
// @Produce json
// @Success 200 " Successfully generated and saved invitation codes"
// @Failure 500 " Error generating and adding invitation codes"
// @Router /v1/api/invitation/codes/add [post]
func (api *InvitationCodesApi) MakeInvitationCodes(c *gin.Context) {
	currentUser := auth.CurrentUser(c)
	if currentUser == nil {
		c.JSON(401, gin.H{"error": "please login"})
		return
	}
	if code, expiresAt, err := api.service.MakeInvitationCodes(*currentUser); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{
			"message":   "Successfully added invitation codes",
			"code":      *code,
			"expiresAt": *expiresAt,
		})
	}
}
