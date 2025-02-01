package preferences

import "maps"

type Preference int
var prefs map[Preference]interface{}

const (
	SwapCursorOnDirectionChange Preference = iota
	SwapCursorOnGridWrap
	WrapAtEndOfGrid
)

var defaultPreferences = map[Preference]interface{}{
	SwapCursorOnDirectionChange: true,
	SwapCursorOnGridWrap: true,
	WrapAtEndOfGrid: true,
}

func Init() {
	prefs = maps.Clone(defaultPreferences)
}

func Get(k Preference) interface{} {
	return prefs[k]
}

func GetBool(k Preference) bool {
	return Get(k).(bool)
}