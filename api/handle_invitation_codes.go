package api

import (
	"crypto/rand"
	"encoding/hex"
	"novaro-server/auth"
	"novaro-server/config"
	"novaro-server/model"

	"github.com/gin-gonic/gin"
)

type InvitationCodesApi struct {
}

// MakeInvitationCodes godoc
// @Summary Generate new invitation codes
// @Description Generate new invitation codes and save to the system
// @Accept json
// @Produce json
// @Success 200 " Successfully generated and saved invitation codes"
// @Failure 500 " Error generating and adding invitation codes"
// @Router /v1/api/invitation/codes/add [post]
func (InvitationCodesApi) MakeInvitationCodes(c *gin.Context) {
	code, err := makeInvitationCode(config.InvitatioCodeLength)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	currentUser := auth.CurrentUser(c)
	if err := model.SaveInvitationCodes(currentUser, code); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Successfully added invitation codes"})
}

func makeInvitationCode(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
