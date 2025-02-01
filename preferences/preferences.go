package preferences

type Preference int

const (
	SwapCursorOnDirectionChange Preference = iota
)

var prefs map[Preference]interface{}

func Init() {
	prefs = getDefaultPreferences()
}

func getDefaultPreferences() map[Preference]interface{} {
	return map[Preference]interface{}{
		SwapCursorOnDirectionChange: false,
	}
}

func Get(k Preference) interface{} {
	return prefs[k]
}

func GetBool(k Preference) bool {
	return Get(k).(bool)
}