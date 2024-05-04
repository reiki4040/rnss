package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	InstanceId string
	Name       string
	Others     string
}

func (i item) Title() string {
	return i.Name
}
func (i item) Description() string {
	return fmt.Sprintf("%s\t%s", i.InstanceId, i.Others)
}
func (i item) FilterValue() string {
	return fmt.Sprintf("%s\t%s\t%s", i.Name, i.InstanceId, i.Others)
}

func NewSelectionCUI(ec2list []string, filter []rune) (*model, error) {
	items := make([]list.Item, 0, len(ec2list))
	for _, l := range ec2list {
		splited := strings.SplitN(l, "\t", 3)
		item := item{}
		switch len(splited) {
		case 0:
			continue
		default:
			fallthrough
		case 3:
			item.Others = splited[2]
			fallthrough
		case 2:
			item.Name = splited[1]
			fallthrough
		case 1:
			item.InstanceId = splited[0]
		}
		items = append(items, item)
	}

	return &model{
		list: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}, nil
}

type model struct {
	list list.Model

	selectedItem item
	selected     bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// default cursor, filter paging implemented in bubble/list
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.selected = true
				m.selectedItem = item(i)
			}

			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}
