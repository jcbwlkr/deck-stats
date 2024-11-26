package magic

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

type Color string

const (
	White Color = "w"
	Blue        = "u"
	Black       = "b"
	Red         = "r"
	Green       = "g"
)

type ColorIdentity []Color

var (
	MonoWhite     = ColorIdentity{White}
	MonoBlue      = ColorIdentity{Blue}
	MonoBlack     = ColorIdentity{Black}
	MonoRed       = ColorIdentity{Red}
	MonoGreen     = ColorIdentity{Green}
	Azorius       = ColorIdentity{White, Blue}
	Orzhov        = ColorIdentity{White, Black}
	Boros         = ColorIdentity{White, Red}
	Selesnya      = ColorIdentity{White, Green}
	Dimir         = ColorIdentity{Blue, Black}
	Simic         = ColorIdentity{Blue, Green}
	Izzet         = ColorIdentity{Blue, Red}
	Rakdos        = ColorIdentity{Black, Red}
	Golgari       = ColorIdentity{Black, Green}
	Gruul         = ColorIdentity{Red, Green}
	Esper         = ColorIdentity{White, Blue, Black}
	Bant          = ColorIdentity{White, Blue, Green}
	Jeskai        = ColorIdentity{White, Blue, Red}
	Mardu         = ColorIdentity{White, Black, Red}
	Abzan         = ColorIdentity{White, Black, Green}
	Naya          = ColorIdentity{White, Red, Green}
	Grixis        = ColorIdentity{Blue, Black, Red}
	Sultai        = ColorIdentity{Blue, Black, Green}
	Temur         = ColorIdentity{Blue, Red, Green}
	Jund          = ColorIdentity{Black, Red, Green}
	YoreTiller    = ColorIdentity{White, Blue, Black, Red}
	WitchMaw      = ColorIdentity{White, Blue, Black, Green}
	InkTreader    = ColorIdentity{White, Blue, Red, Green}
	DuneBlackrood = ColorIdentity{White, Black, Red, Green}
	GlintEye      = ColorIdentity{Blue, Black, Red, Green}
	WUBRG         = ColorIdentity{White, Blue, Black, Red, Green}
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

// Value implements the driver.Valuer interface
func (id ColorIdentity) Value() (driver.Value, error) {
	s := make([]string, 0, len(id))
	for _, c := range id {
		s = append(s, string(c))
	}
	return pq.StringArray(s), nil
}

// Scan implements the sql.Scanner interface,
func (id *ColorIdentity) Scan(src interface{}) error {
	var s pq.StringArray
	if err := s.Scan(src); err != nil {
		return err
	}

	tmp := make(ColorIdentity, 0, len(s))
	for _, c := range s {
		tmp = append(tmp, Color(c))
	}
	*id = tmp
	return nil
}

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
