package preferences

import (
	"fmt"
	"maps"
)

//go:generate stringer -type=Preference
type Preference int
type SetPreference struct {
	Pref  Preference
	Value interface{}
}
type Preferences map[Preference]interface{}

var prefs Preferences

const (
	JumpToEmptySquare Preference = iota
	SwapCursorOnGridWrap
	SwapCursorOnDirectionChange
	WrapAtEndOfGrid
	WrapOnArrowNavigation
)

var defaultPreferences = Preferences{
	SwapCursorOnDirectionChange: true,
	SwapCursorOnGridWrap:        true,
	WrapOnArrowNavigation:       false,
	WrapAtEndOfGrid:             true,
	JumpToEmptySquare:           true,
}

func Init() {
	prefs = maps.Clone(defaultPreferences)
}

func ListPreferences() []SetPreference {
	preferenceSettings := make([]SetPreference, 0, len(defaultPreferences))

	for key := JumpToEmptySquare; key <= WrapOnArrowNavigation; key++ {
		prefSetting := SetPreference{
			Pref:  key,
			Value: prefs[key],
		}
		preferenceSettings = append(preferenceSettings, prefSetting)
	}

	return preferenceSettings
}

func Get(k Preference) interface{} {
	return prefs[k]
}

func GetBool(k Preference) bool {
	return Get(k).(bool)
}

func Set(k Preference, v interface{}) {
	prefs[k] = v
}

func SetBool(k Preference, v bool) {
	Set(k, v)
}

func (p Preferences) String() string {
	str := "Preferences{"

	for key, val := range p {
		str += fmt.Sprintf("%v: %v, ", key, val)
	}

	str += "}"

	return str
}
