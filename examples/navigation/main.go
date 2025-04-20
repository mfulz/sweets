package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mfulz/chocolate"
	"github.com/mfulz/sweets/navigation"
)

var layout string = `
[
  {
    "_comment": "set width of tuimodel1 20% of the parent's width",
    "source": "super",
    "source_attribute": "width",
    "target": "tuimodel1",
    "target_attribute": "width",
    "relation": "eq",
    "multiplier": 0.2,
    "strength": "required"
  },
  {
    "_comment": "set height of tuimodel1 equal to height of parent",
    "source": "super",
    "source_attribute": "height",
    "target": "tuimodel1",
    "target_attribute": "height",
    "relation": "eq",
    "multiplier": 1.0
  },
  {
    "_comment": "place tuimodel2 right to tuimodel1 and make it required",
    "source": "tuimodel1",
    "source_attribute": "xend",
    "target": "tuimodel2",
    "target_attribute": "xstart",
    "relation": "eq",
    "multiplier": 1.0,
    "strength": "required"
  },
  {
    "_comment": "set width of tuimodel2 equal to width of parent. (the bias algorithm will shrink it",
    "_comment": "to fill the remaining width after sizing tuimodel1. This will work as the default strength is medium",
    "_comment": "and tuimodel1 uses a higher one",
    "source": "super",
    "source_attribute": "width",
    "target": "tuimodel2",
    "target_attribute": "width",
    "relation": "eq",
    "multiplier": 1.0
  },
  {
    "_comment": "set width of tuimodel3 equal to width of tuimodel2 and make it required.",
    "source": "tuimodel2",
    "source_attribute": "width",
    "target": "tuimodel3",
    "target_attribute": "width",
    "relation": "eq",
    "multiplier": 1.0,
    "strength": "required"
  },
  {
    "_comment": "set height of tuimodel2 equal to height of parent",
    "source": "super",
    "source_attribute": "height",
    "target": "tuimodel2",
    "target_attribute": "height",
    "relation": "eq",
    "multiplier": 1.0
  },
  {
    "_comment": "set height of tuimodel3 equal to height of parent",
    "source": "super",
    "source_attribute": "height",
    "target": "tuimodel3",
    "target_attribute": "height",
    "relation": "eq",
    "multiplier": 1.0
  },
  {
    "_comment": "place tuimodel3 right of tuimodel2 and make it required",
    "source": "tuimodel2",
    "source_attribute": "xend",
    "target": "tuimodel3",
    "target_attribute": "xstart",
    "relation": "eq",
    "multiplier": 1.0,
    "strength": "required"
  }
]
`

// styles used for flavour
var (
	viewStyle = lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color("246")).Background(lipgloss.Color("232")).
			BorderForeground(lipgloss.Color("246")).BorderBackground(lipgloss.Color("232")).
			Border(lipgloss.RoundedBorder())
	selectedStyle = lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color("15")).Background(lipgloss.Color("237")).
			BorderForeground(lipgloss.Color("15")).BorderBackground(lipgloss.Color("237")).
			Border(lipgloss.RoundedBorder())
	focusedStyle = lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color("196")).Background(lipgloss.Color("232")).
			BorderForeground(lipgloss.Color("15")).BorderBackground(lipgloss.Color("232")).
			Border(lipgloss.RoundedBorder())
)

type textModel struct {
	text string
	alt  string

	cur *string
}

func (t *textModel) Init() tea.Cmd { return nil }
func (t *textModel) View() string  { return *t.cur }
func (t *textModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "t":
			if t.cur == &t.alt {
				t.cur = &t.text
			} else {
				t.cur = &t.alt
			}
		case "a":
			t.cur = &t.alt
		case "s":
			t.cur = &t.text
		}
	}

	return t, nil
}

func main() {
	theme := chocolate.NewChocolateFlavour(
		chocolate.WithDefaults(
			&viewStyle,
			&selectedStyle,
			&focusedStyle,
		),
	)

	choc := chocolate.NewChocolate(chocolate.WithFlavour(theme))
	if err := choc.FromJson([]byte(layout)); err != nil {
		panic(err)
	}

	sub := choc.MakeChocolate(
		"tui3", "tuimodel3",
		true,
	)
	if err := sub.FromJson([]byte(layout)); err != nil {
		panic(err)
	}

	tui21 := &textModel{
		text: "Tui21 Text",
		alt:  "Tui21 Alt Text",
	}
	tui21.cur = &tui21.text

	tui22 := &textModel{
		text: "Tui22 Text",
		alt:  "Tui22 Alt Text",
	}
	tui22.cur = &tui22.text

	tui23 := &textModel{
		text: "Tui23 Text",
		alt:  "Tui23 Alt Text",
	}
	tui23.cur = &tui23.text

	sub.AddTeaModelBarModel(tui21, "tui1", "tuimodel1", true)
	sub.AddTeaModelBarModel(tui22, "tui2", "tuimodel2", true)
	sub.AddTeaModelBarModel(tui23, "tui3", "tuimodel3", true)

	tui1 := &textModel{
		text: "Tui1 Text",
		alt:  "Tui1 Alt Text",
	}
	tui1.cur = &tui1.text

	tui2 := &textModel{
		text: "Tui2 Text",
		alt:  "Tui2 Alt Text",
	}
	tui2.cur = &tui2.text

	// tui3 := &textModel{
	// 	text: "Tui3 Text",
	// 	alt:  "Tui3 Alt Text",
	// }
	// tui3.cur = &tui3.text

	choc.AddTeaModelBarModel(tui1, "tui1", "tuimodel1", true)
	choc.AddTeaModelBarModel(tui2, "tui2", "tuimodel2", true)
	// choc.AddTeaModelBarModel(tui3, "tui3", "tuimodel3", true)

	nav2 := navigation.NewChocolateNavigator(
		sub, "tuimodel3", "tui3",
		"tuimodel1:tui1", "tuimodel2:tui2",
		true, false, true,
	)
	nav2.Model().AddNavigator(
		navigation.NewTeaModelNavigator(
			tui21, "tuimodel1", "tui1",
			"tuimodel2:tui2", "tuimodel3:tui3",
			false, false, true),
	)
	nav2.Model().AddNavigator(
		navigation.NewTeaModelNavigator(
			tui22, "tuimodel2", "tui2",
			"tuimodel3:tui3", "tuimodel1:tui1",
			false, false, true),
	)
	nav2.Model().AddNavigator(
		navigation.NewTeaModelNavigator(
			tui23, "tuimodel3", "tui3",
			"tuimodel1:tui1", "tuimodel2:tui2",
			false, false, true),
	)
	// nav2.AddNavigationModel(
	// 	"",
	// 	tui21, "tuimodel1", "tui21",
	// 	"tui22", "tui23",
	// 	false, false, true,
	// )
	// nav2.AddNavigationModel(
	// 	"",
	// 	tui22, "tuimodel2", "tui22",
	// 	"tui23", "tui21",
	// 	false, false, false,
	// )
	// nav2.AddNavigationModel(
	// 	"",
	// 	tui23, "tuimodel3", "tui23",
	// 	"tui21", "tui22",
	// 	false, false, false,
	// )

	nav := navigation.NewNavigationModel(choc)
	nav.AddNavigator(
		navigation.NewTeaModelNavigator(
			tui1, "tuimodel1", "tui1",
			"tuimodel2:tui2", "tuimodel3:tui3",
			false, false, true),
	)
	nav.AddNavigator(
		navigation.NewTeaModelNavigator(
			tui2, "tuimodel2", "tui2",
			"tuimodel3:tui3", "tuimodel1:tui1",
			false, true, true),
	)
	nav.AddNavigator(
		nav2,
	)
	// 		nav.AddNavigationModel(
	// 	"",
	// 	tui2, "tuimodel2", "tui2",
	// 	"tui1", "nav2",
	// 	false, true, true,
	// )
	// nav.AddNavigationModel(
	// 	"",
	// 	nav2, "tuimodel3", "tui3",
	// 	"tui1", "tui2",
	// 	true, false, false,
	// )
	// nav.AddNavigationModel(
	// 	"nav2",
	// 	nav2, "tuimodel3", "nav2",
	// 	"tui2", "tui1",
	// 	false, false, true,
	// )
	p := tea.NewProgram(nav)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
