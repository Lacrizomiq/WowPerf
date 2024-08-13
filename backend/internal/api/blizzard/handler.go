package blizzard

import (
	"wowperf/internal/api/blizzard/profile"
	"wowperf/internal/services/blizzard"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	CharacterProfile *profile.CharacterProfileHandler
	Equipment        *profile.EquipmentHandler
}

func NewHandler(service *blizzard.Service) (*Handler, error) {
	return &Handler{
		CharacterProfile: profile.NewCharacterProfileHandler(service),
		Equipment:        profile.NewEquipmentHandler(service),
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/blizzard/characters/:realmSlug/:characterName", h.CharacterProfile.GetCharacterProfile)
	r.GET("/blizzard/characters/:realmSlug/:characterName/equipment", h.Equipment.GetCharacterEquipment)
}
