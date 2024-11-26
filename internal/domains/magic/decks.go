package magic

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Deck struct {
	ID            string        `db:"id" json:"id"`
	UserID        string        `db:"user_id" json:"user_id"`
	Service       string        `db:"service" json:"service"`
	ServiceID     string        `db:"service_id" json:"service_id"`
	Name          string        `db:"name" json:"name"`
	Format        string        `db:"format" json:"format"`
	URL           string        `db:"url" json:"url"`
	ColorIdentity ColorIdentity `db:"color_identity" json:"color_identity"`
	Folder        Folder        `db:"folder" json:"folder"`
	Archetypes    Archetypes    `db:"archetypes" json:"archetypes"`
	Leaders       Leaders       `db:"leaders" json:"leaders"`
	UpdatedAt     time.Time     `db:"updated_at" json:"updated_at"`
	RefreshedAt   time.Time     `db:"refreshed_at" json:"refreshed_at"`
}

type Leaders struct {
	Commanders      []Card `json:"commanders,omitempty"`
	Companion       *Card  `json:"companion,omitempty"`
	Oathbreakers    []Card `json:"oathbreakers,omitempty"`
	SignatureSpells []Card `json:"signature_spells,omitempty"`
}

type Card struct {
	ID   string // Scryfall ID
	Name string
}

type Folder struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Archetypes []Archetype

type Archetype struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

////////////////////////////////////////////////////////////////////////////////
// DB Methods for Storing
////////////////////////////////////////////////////////////////////////////////

// Scan implements the Scanner interface for Folder.
func (f *Folder) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	b, ok := v.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, f)
}

// Value implements the Valuer interface for Folder.
func (f Folder) Value() (driver.Value, error) {
	return json.Marshal(f)
}

// Scan implements the Scanner interface for Archetypes.
func (a *Archetypes) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	b, ok := v.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, a)
}

// Value implements the Valuer interface for Archetypes.
func (a Archetypes) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan implements the Scanner interface for Leaders.
func (a *Leaders) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	b, ok := v.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, a)
}

// Value implements the Valuer interface for Leaders.
func (a Leaders) Value() (driver.Value, error) {
	return json.Marshal(a)
}
