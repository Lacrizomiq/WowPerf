package blizzard

import (
	"wowperf/internal/api/blizzard/gamedata"
	"wowperf/internal/api/blizzard/profile"
	"wowperf/internal/services/blizzard"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	MythicKeystoneAffixIndex    *gamedata.MythicKeystoneAffixIndexHandler
	MythicKeystoneAffixByID     *gamedata.MythicKeystoneAffixByIDHandler
	MythicKeystoneAffixMedia    *gamedata.MythicKeystoneAffixMediaHandler
	MythicKeystoneIndex         *gamedata.MythicKeystoneIndexHandler
	MythicKeystoneDungeonsIndex *gamedata.MythicKeystoneDungeonsIndexHandler
	MythicKeystoneByID          *gamedata.MythicKeystoneByIDHandler
	MythicKeystonePeriodsIndex  *gamedata.MythicKeystonePeriodsIndexHandler
	MythicKeystonePeriodByID    *gamedata.MythicKeystonePeriodByIDHandler
	MythicKeystoneSeasonsIndex  *gamedata.MythicKeystoneSeasonsIndexHandler
	MythicKeystoneSeasonByID    *gamedata.MythicKeystoneSeasonByIDHandler
	JournalInstanceIndex        *gamedata.JournalInstanceIndexHandler
	JournalInstanceByID         *gamedata.JournalInstanceByIDHandler
	JournalInstanceMedia        *gamedata.JournalInstanceMediaHandler
	GetSeasonDungeons           *profile.GetSeasonDungeonsHandler
}

func NewHandler(service *blizzard.Service, db *gorm.DB) *Handler {
	return &Handler{
		CharacterProfile:            profile.NewCharacterProfileHandler(service),
		CharacterMedia:              profile.NewCharacterMediaHandler(service),
		Equipment:                   profile.NewEquipmentHandler(service),
		ItemMedia:                   gamedata.NewItemMediaHandler(service),
		MythicKeystoneProfile:       profile.NewMythicKeystoneProfileHandler(service),
		MythicKeystoneSeasonDetails: profile.NewMythicKeystoneSeasonDetailsHandler(service, db),
		Specializations:             profile.NewSpecializationsHandler(service, db),
		SpellMedia:                  gamedata.NewSpellMediaHandler(service),
		TalentTreeIndex:             gamedata.NewTalentTreeIndexHandler(service),
		TalentTree:                  gamedata.NewTalentTreeHandler(db),
		TalentTreeNodes:             gamedata.NewTalentTreeNodesHandler(service),
		TalentIndex:                 gamedata.NewTalentIndexHandler(service),
		TalentByID:                  gamedata.NewTalentByIDHandler(service),
		MythicKeystoneAffixIndex:    gamedata.NewMythicKeystoneAffixIndexHandler(service),
		MythicKeystoneAffixByID:     gamedata.NewMythicKeystoneAffixByIDHandler(service),
		MythicKeystoneAffixMedia:    gamedata.NewMythicKeystoneAffixMediaHandler(service),
		MythicKeystoneIndex:         gamedata.NewMythicKeystoneIndexHandler(service),
		MythicKeystoneDungeonsIndex: gamedata.NewMythicKeystoneDungeonsIndexHandler(service),
		MythicKeystoneByID:          gamedata.NewMythicKeystoneByIDHandler(service),
		MythicKeystonePeriodsIndex:  gamedata.NewMythicKeystonePeriodsIndexHandler(service),
		MythicKeystonePeriodByID:    gamedata.NewMythicKeystonePeriodByIDHandler(service),
		MythicKeystoneSeasonsIndex:  gamedata.NewMythicKeystoneSeasonsIndexHandler(service),
		MythicKeystoneSeasonByID:    gamedata.NewMythicKeystoneSeasonByIDHandler(service),
		JournalInstanceIndex:        gamedata.NewJournalInstanceIndexHandler(service),
		JournalInstanceByID:         gamedata.NewJournalInstanceByIDHandler(service),
		JournalInstanceMedia:        gamedata.NewJournalInstanceMediaHandler(service),
		GetSeasonDungeons:           profile.NewGetSeasonDungeonsHandler(service, db),
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {

	// Blizzard Profile API
	r.GET("/blizzard/characters/:realmSlug/:characterName", h.CharacterProfile.GetCharacterProfile)
	r.GET("/blizzard/characters/:realmSlug/:characterName/media", h.CharacterMedia.GetCharacterMedia)

	r.GET("/blizzard/characters/:realmSlug/:characterName/equipment", h.Equipment.GetCharacterEquipment)

	r.GET("/blizzard/characters/:realmSlug/:characterName/mythic-keystone-profile", h.MythicKeystoneProfile.GetCharacterMythicKeystoneProfile)
	r.GET("/blizzard/characters/:realmSlug/:characterName/mythic-keystone-profile/season/:seasonId", h.MythicKeystoneSeasonDetails.GetCharacterMythicKeystoneSeasonBestRuns)

	r.GET("/blizzard/characters/:realmSlug/:characterName/specializations", h.Specializations.GetCharacterSpecializations)

	// Blizzard Game Data API
	r.GET("/blizzard/data/item/:itemId/media", h.ItemMedia.GetItemMedia)
	r.GET("/blizzard/data/spell/:spellId/media", h.SpellMedia.GetSpellMedia)

	r.GET("/blizzard/data/talent-tree/index", h.TalentTreeIndex.GetTalentTreeIndex)
	r.GET("/blizzard/data/talent-tree/:talentTreeId/playable-specialization/:specId", h.TalentTree.GetTalentTree)
	r.GET("/blizzard/data/talent-tree/:talentTreeId/nodes", h.TalentTreeNodes.GetTalentTreeNodes)
	r.GET("/blizzard/data/talent/index", h.TalentIndex.GetTalentIndex)
	r.GET("/blizzard/data/talent/:talentId", h.TalentByID.GetTalentByID)

	r.GET("/blizzard/data/mythic-keystone-affix/index", h.MythicKeystoneAffixIndex.GetMythicKeystoneAffixIndex)
	r.GET("/blizzard/data/mythic-keystone-affix/:affixId", h.MythicKeystoneAffixByID.GetMythicKeystoneAffixByID)
	r.GET("/blizzard/data/mythic-keystone-affix/:affixId/media", h.MythicKeystoneAffixMedia.GetMythicKeystoneAffixMedia)

	r.GET("/blizzard/data/mythic-keystone/index", h.MythicKeystoneIndex.GetMythicKeystoneIndex)
	r.GET("/blizzard/data/mythic-keystone/dungeon/index", h.MythicKeystoneDungeonsIndex.GetMythicKeystoneDungeonsIndex)
	r.GET("/blizzard/data/mythic-keystone/:mythicKeystoneId", h.MythicKeystoneByID.GetMythicKeystoneByID)
	r.GET("/blizzard/data/mythic-keystone/period/index", h.MythicKeystonePeriodsIndex.GetMythicKeystonePeriodsIndex)
	r.GET("/blizzard/data/mythic-keystone/period/:periodId", h.MythicKeystonePeriodByID.GetMythicKeystonePeriodByID)
	r.GET("/blizzard/data/mythic-keystone/season/index", h.MythicKeystoneSeasonsIndex.GetMythicKeystoneSeasonsIndex)
	r.GET("/blizzard/data/mythic-keystone/season/:seasonId", h.MythicKeystoneSeasonByID.GetMythicKeystoneSeasonByID)

	r.GET("/blizzard/data/journal-instance/index", h.JournalInstanceIndex.GetJournalInstanceIndex)
	r.GET("/blizzard/data/journal-instance/:instanceId", h.JournalInstanceByID.GetJournalInstanceByID)
	r.GET("/blizzard/data/journal-instance/:instanceId/media", h.JournalInstanceMedia.GetJournalInstanceMedia)

	r.GET("data/mythic-keystone/season/:seasonSlug/dungeons", h.GetSeasonDungeons.GetSeasonDungeons)
}
