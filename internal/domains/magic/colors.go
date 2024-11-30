package magic

import (
	"database/sql/driver"

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

////////////////////////////////////////////////////////////////////////////////
// DB Methods for Storing
////////////////////////////////////////////////////////////////////////////////

// Value implements the driver.Valuer interface
func (id ColorIdentity) Value() (driver.Value, error) {
	s := make([]string, 0, len(id))
	for _, c := range id {
		s = append(s, string(c))
	}
	return pq.StringArray(s).Value()
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
