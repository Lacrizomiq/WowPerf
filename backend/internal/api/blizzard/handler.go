package blizzard

import (
	"wowperf/internal/api/blizzard/gamedata"
	"wowperf/internal/api/blizzard/profile"
	"wowperf/internal/services/blizzard"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	CharacterProfile            *profile.CharacterProfileHandler
	CharacterMedia              *profile.CharacterMediaHandler
	Equipment                   *profile.EquipmentHandler
	ItemMedia                   *gamedata.ItemMediaHandler
	MythicKeystoneProfile       *profile.MythicKeystoneProfileHandler
	MythicKeystoneSeasonDetails *profile.MythicKeystoneSeasonDetailsHandler
	Specializations             *profile.SpecializationsHandler
	SpellMedia                  *gamedata.SpellMediaHandler
	TalentTreeIndex             *gamedata.TalentTreeIndexHandler
	TalentTree                  *gamedata.TalentTreeHandler
	TalentTreeNodes             *gamedata.TalentTreeNodesHandler
	TalentIndex                 *gamedata.TalentIndexHandler
	TalentByID                  *gamedata.TalentByIDHandler
}

func NewHandler(service *blizzard.Service) (*Handler, error) {
	return &Handler{
		CharacterProfile:            profile.NewCharacterProfileHandler(service),
		CharacterMedia:              profile.NewCharacterMediaHandler(service),
		Equipment:                   profile.NewEquipmentHandler(service),
		ItemMedia:                   gamedata.NewItemMediaHandler(service),
		MythicKeystoneProfile:       profile.NewMythicKeystoneProfileHandler(service),
		MythicKeystoneSeasonDetails: profile.NewMythicKeystoneSeasonDetailsHandler(service),
		Specializations:             profile.NewSpecializationsHandler(service),
		SpellMedia:                  gamedata.NewSpellMediaHandler(service),
		TalentTreeIndex:             gamedata.NewTalentTreeIndexHandler(service),
		TalentTree:                  gamedata.NewTalentTreeHandler(service),
		TalentTreeNodes:             gamedata.NewTalentTreeNodesHandler(service),
		TalentIndex:                 gamedata.NewTalentIndexHandler(service),
		TalentByID:                  gamedata.NewTalentByIDHandler(service),
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {

	// Blizzard Profile API
	r.GET("/blizzard/characters/:realmSlug/:characterName", h.CharacterProfile.GetCharacterProfile)
	r.GET("/blizzard/characters/:realmSlug/:characterName/media", h.CharacterMedia.GetCharacterMedia)

	r.GET("/blizzard/characters/:realmSlug/:characterName/equipment", h.Equipment.GetCharacterEquipment)

	r.GET("/blizzard/characters/:realmSlug/:characterName/mythic-keystone-profile", h.MythicKeystoneProfile.GetCharacterMythicKeystoneProfile)
	r.GET("/blizzard/characters/:realmSlug/:characterName/mythic-keystone-profile/season/:seasonId", h.MythicKeystoneSeasonDetails.GetCharacterMythicKeystoneSeasonDetails)

	r.GET("/blizzard/characters/:realmSlug/:characterName/specializations", h.Specializations.GetCharacterSpecializations)

	// Blizzard Game Data API
	r.GET("/blizzard/data/item/:itemId/media", h.ItemMedia.GetItemMedia)
	r.GET("/blizzard/data/spell/:spellId/media", h.SpellMedia.GetSpellMedia)

	r.GET("/blizzard/data/talent-tree/index", h.TalentTreeIndex.GetTalentTreeIndex)
	r.GET("/blizzard/data/talent-tree/:talentTreeId/playable-specialization/:specId", h.TalentTree.GetTalentTree)
	r.GET("/blizzard/data/talent-tree/:talentTreeId/nodes", h.TalentTreeNodes.GetTalentTreeNodes)
	r.GET("/blizzard/data/talent/index", h.TalentIndex.GetTalentIndex)
	r.GET("/blizzard/data/talent/:talentId", h.TalentByID.GetTalentByID)

}
