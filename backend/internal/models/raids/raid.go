package raids

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Raid struct {
	*gorm.Model `json:"-"`
	ID          uint   `gorm:"primaryKey"`
	Slug        string `gorm:"uniqueIndex"`
	Name        string
	ShortName   string
	Expansion   string
	MediaURL    string
	Icon        string
	Starts      StartEndMap `gorm:"type:jsonb"`
	Ends        StartEndMap `gorm:"type:jsonb"`
	Encounters  Encounters  `gorm:"type:jsonb"`
}

type StartEndMap map[string]time.Time

func (m StartEndMap) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *StartEndMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, &m)
}

type Encounters []Encounter

func (e Encounters) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *Encounters) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, &e)
}

type Encounter struct {
	ID   uint   `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}
