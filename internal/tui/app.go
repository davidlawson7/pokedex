package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Messages for screen transitions
type switchToDetailMsg struct {
	pokemonID uint16
}

type switchToSearchMsg struct{}

// screen identifies which screen is active.
type screen int

const (
	screenSearch screen = iota
	screenDetail
)

// AppModel is the root Bubble Tea model that routes between screens.
type AppModel struct {
	current screen
	search  SearchModel
	detail  DetailModel
	width   int
	height  int
}

// NewAppModel creates the root model with the search screen active.
func NewAppModel() AppModel {
	return AppModel{
		current: screenSearch,
		search:  NewSearchModel(),
	}
}

func (a AppModel) Init() tea.Cmd {
	return a.search.Init()
}

func (a AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

	case switchToDetailMsg:
		a.detail = NewDetailModel(msg.pokemonID, a.width, a.height)
		a.current = screenDetail
		return a, a.detail.Init()

	case switchToSearchMsg:
		a.current = screenSearch
		return a, nil
	}

	switch a.current {
	case screenSearch:
		m, cmd := a.search.Update(msg)
		a.search = m.(SearchModel)
		return a, cmd
	case screenDetail:
		m, cmd := a.detail.Update(msg)
		a.detail = m.(DetailModel)
		return a, cmd
	}
	return a, nil
}

func (a AppModel) View() string {
	switch a.current {
	case screenDetail:
		return a.detail.View()
	default:
		return a.search.View()
	}
}

// Run starts the Bubble Tea application.
func Run() error {
	p := tea.NewProgram(NewAppModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
