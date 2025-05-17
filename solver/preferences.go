package solver

import (
	"fmt"

	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2/list"
	"github.com/tylerwgrass/cruciterm/logger"
	prefs "github.com/tylerwgrass/cruciterm/preferences"
)

var activePreferenceIndex int = 0

type preferencesModel struct {
	preferences []prefs.SetPreference
	preferencesList *list.List 
}

func preferencesEnumerator(l list.Items, i int) string {
	if i == activePreferenceIndex {
		return "⮕ "
	}
	return ""
}

func initPreferencesModel() preferencesModel {
	preferences := prefs.ListPreferences()
	preferencesList := getPreferencesList(preferences)
		
	return preferencesModel{
		preferences: preferences, 
		preferencesList: preferencesList,
	} 
	
}

func (m preferencesModel) Init() tea.Cmd {
	return nil
}

func (m preferencesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
			case key.Matches(msg, keys.Up):
				activePreferenceIndex = max(activePreferenceIndex - 1, 0);
			case key.Matches(msg, keys.Down):
				activePreferenceIndex = min(activePreferenceIndex + 1, len(m.preferences) - 1);
		}
	}

	m.preferencesList = getPreferencesList(m.preferences)
	return m, nil
}

func getPreferencesList(preferences []prefs.SetPreference) *list.List {
	logger.Debug(fmt.Sprintf("active index: %v", activePreferenceIndex))
	preferencesList := list.New().
		Enumerator(preferencesEnumerator)
	
	for _, setPref := range(preferences) {
		if setPref.Value == true {
			preferencesList.Item(setPref.Pref.String() + " ✓")
		} else {
			preferencesList.Item(setPref.Pref.String() + " x")
		}
	}
	return preferencesList
}

func (m preferencesModel) View() string {
	return m.preferencesList.String() 
}