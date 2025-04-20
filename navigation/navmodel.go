package navigation

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mfulz/chocolate"
)

type Navigator interface {
	id() string
	bar() string
	name() string
	next() string
	prev() string
	autoFocus() bool
	navKeysFocus() bool
	keysMustFocus() bool
	hasFocus() bool
	setFocus(bool)
	Update(tea.Msg) (tea.Model, tea.Cmd)
}

type navigator[T any] struct {
	model T

	mbar  string
	mname string
	mnext string
	mprev string

	autofocus     bool
	navKeysfocus  bool
	keysMustfocus bool

	canfocus bool
	hasfocus bool

	updateHandler func(tea.Msg) (tea.Model, tea.Cmd)
}

func (n *navigator[T]) id() string          { return fmt.Sprintf("%s:%s", n.mbar, n.mname) }
func (n *navigator[T]) bar() string         { return n.mbar }
func (n *navigator[T]) name() string        { return n.mname }
func (n *navigator[T]) next() string        { return n.mnext }
func (n *navigator[T]) prev() string        { return n.mprev }
func (n *navigator[T]) autoFocus() bool     { return n.autofocus }
func (n *navigator[T]) navKeysFocus() bool  { return n.navKeysfocus }
func (n *navigator[T]) keysMustFocus() bool { return n.keysMustfocus }
func (n *navigator[T]) hasFocus() bool      { return n.hasfocus }
func (n *navigator[T]) setFocus(v bool) {
	if n.canfocus {
		n.hasfocus = v
	} else {
		n.hasfocus = false
	}
}
func (n *navigator[T]) Model() T { return n.model }
func (n *navigator[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if n.updateHandler == nil {
		return nil, nil
	}
	return n.updateHandler(msg)
}

func NewChocolateNavigator[T *chocolate.Chocolate](
	choc T, bar, name,
	next, prev string,
	autoFocus, navKeysFocus, keysMustFocus bool,
) *navigator[*NavigationModel] {
	model := NewNavigationModel(choc)

	return &navigator[*NavigationModel]{
		model:         model,
		mbar:          bar,
		mname:         name,
		mnext:         next,
		mprev:         prev,
		autofocus:     autoFocus,
		navKeysfocus:  navKeysFocus,
		keysMustfocus: keysMustFocus,
		canfocus:      true,
		updateHandler: model.Update,
	}
}

func NewTeaModelNavigator[T tea.Model](
	model T, bar, name,
	next, prev string,
	autoFocus, navKeysFocus, keysMustFocus bool,
) *navigator[T] {
	return &navigator[T]{
		model:         model,
		mbar:          bar,
		mname:         name,
		mnext:         next,
		mprev:         prev,
		autofocus:     autoFocus,
		navKeysfocus:  navKeysFocus,
		keysMustfocus: keysMustFocus,
		canfocus:      false,
		updateHandler: model.Update,
	}
}

type NavigationModel struct {
	choc *chocolate.Chocolate

	KeyMap *KeyMap

	navs    map[string]Navigator
	current Navigator
	focused bool
}

func (nm *NavigationModel) handleSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, nm.KeyMap.Select):
			if !nm.focused {
				nm.focused = true
			}
			nm.current.setFocus(true)
			nm.choc.SelectStyle(chocolate.TS_FOCUSED, nm.current.bar())
			return nm, nil
		}
	}

	return nil, func() tea.Msg { return msg }
}

func (nm *NavigationModel) handleLeave(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, nm.KeyMap.Leave) ||
			key.Matches(msg, nm.KeyMap.Quit):
			if nm.focused {
				nm.focused = false
			}
			nm.choc.SelectStyle(chocolate.TS_SELECTED, nm.current.bar())
			nm.current.setFocus(false)
			return nm, nil
		}
	}

	return nil, func() tea.Msg { return msg }
}

func (nm *NavigationModel) handleChangeCurrent(msg tea.Msg) (tea.Model, tea.Cmd) {
	var targetNav string

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, nm.KeyMap.Next):
			targetNav = nm.current.next()
		case key.Matches(msg, nm.KeyMap.Prev):
			targetNav = nm.current.prev()
		}
		if len(targetNav) > 0 {
			nm.setCurrentNav(targetNav)
		}

		return nm, nil
	}

	return nil, func() tea.Msg { return msg }
}

func (nm *NavigationModel) setCurrentNav(id string) error {
	nav, ok := nm.navs[id]
	if !ok {
		return fmt.Errorf("no navigator with id '%s'", id)
	}

	if nm.current == nav {
		return nil
	}

	if nm.current != nil {
		nm.choc.SelectStyle(chocolate.TS_DEFAULT, nm.current.bar())
		nm.current.setFocus(false)
	}
	nm.current = nav
	nm.focused = nav.autoFocus()
	nm.current.setFocus(nm.focused)

	if nm.focused {
		nm.choc.SelectStyle(chocolate.TS_FOCUSED, nm.current.bar())
	} else {
		nm.choc.SelectStyle(chocolate.TS_SELECTED, nm.current.bar())
	}

	return nil
}

func (nm *NavigationModel) Init() tea.Cmd { return nil }
func (nm *NavigationModel) View() string  { return nm.choc.View() }

func (nm *NavigationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if nm.current == nil {
		return nm, nil
	}

	nm.KeyMap.Next.SetEnabled(true)
	nm.KeyMap.Prev.SetEnabled(true)

	if nm.focused {
		if nm.current.hasFocus() {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch {
				case key.Matches(msg, nm.KeyMap.Leave) ||
					key.Matches(msg, nm.KeyMap.Quit):
					_, cmd = nm.current.Update(msg)
					cmds = append(cmds, cmd)
					nm.choc.SelectStyle(chocolate.TS_DEFAULT, nm.current.bar())
					// return nm, tea.Batch(cmds...)
				case key.Matches(msg, nm.KeyMap.Next) ||
					key.Matches(msg, nm.KeyMap.Prev):
					_, cmd = nm.current.Update(msg)
					cmds = append(cmds, cmd)
					return nm, tea.Batch(cmds...)
				case key.Matches(msg, nm.KeyMap.Select):
					_, cmd = nm.current.Update(msg)
					cmds = append(cmds, cmd)
					return nm, tea.Batch(cmds...)
				}
			}
		} else {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch {
				case key.Matches(msg, nm.KeyMap.Select):
					_, cmd = nm.current.Update(msg)
					cmds = append(cmds, cmd)
				}
			}
		}
		if !nm.current.navKeysFocus() {
			nm.KeyMap.Next.SetEnabled(false)
			nm.KeyMap.Prev.SetEnabled(false)
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		nm.choc.Resize(msg.Width, msg.Height)
		return nm, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, nm.KeyMap.Quit):
			if nm.current.hasFocus() {
				nm.current.setFocus(false)
				return nm, nil
			}
			if nm.focused {
				return nm.handleLeave(msg)
			}
			return nm, tea.Quit
		case key.Matches(msg, nm.KeyMap.Select):
			return nm.handleSelect(msg)
		case key.Matches(msg, nm.KeyMap.Leave):
			if nm.current.hasFocus() {
				nm.current.setFocus(false)
				return nm, nil
			}
			return nm.handleLeave(msg)
		case key.Matches(msg, nm.KeyMap.Next) ||
			key.Matches(msg, nm.KeyMap.Prev):
			return nm.handleChangeCurrent(msg)
		default:
			if !nm.current.keysMustFocus() {
				_, cmd = nm.current.Update(msg)
				cmds = append(cmds, cmd)
			} else {
				if nm.focused {
					_, cmd = nm.current.Update(msg)
					cmds = append(cmds, cmd)
				}
			}
		}
	default:
		_, cmd = nm.current.Update(msg)
		cmds = append(cmds, cmd)
	}

	return nm, tea.Batch(cmds...)
}

func (nm *NavigationModel) AddNavigator(nav Navigator, overwrites ...bool) error {
	if nm.navs == nil {
		nm.navs = make(map[string]Navigator)
	}
	overwrite := len(overwrites) > 0

	if _, ok := nm.navs[nav.id()]; ok && !overwrite {
		return fmt.Errorf("navigator with id '%s' already exists", nav.id())
	}

	nm.navs[nav.id()] = nav
	if nm.current == nil {
		nm.current = nav
		if nav.autoFocus() {
			nm.focused = true
		}
	}

	return nil
}

func NewNavigationModel(choc *chocolate.Chocolate) *NavigationModel {
	return &NavigationModel{
		choc:   choc,
		KeyMap: DefaultKeyMap(),
	}
}
