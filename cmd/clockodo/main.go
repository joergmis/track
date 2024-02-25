package main

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joergmis/track"
	"github.com/joergmis/track/local"
	"github.com/spf13/viper"
)

const (
	customerViewFocused = iota
	projectViewFocused
	activitiesViewFocused
)

type customer struct {
	name        string
	description string
}

func (c customer) FilterValue() string {
	return c.name + " " + c.description
}
func (c customer) Title() string {
	return c.name
}
func (c customer) Description() string {
	return c.description
}

type model struct {
	focused int

	customerView   list.Model
	projectView    list.Model
	activitiesView list.Model

	repository track.Repository

	loading bool
	quiting bool
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) deselectCustomer() {
	items := []list.Item{}
	m.projectView.SetItems(items)
	m.focused = customerViewFocused
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.loading {
			m.loading = false
			m.projectView.SetWidth(msg.Width)
			m.projectView.SetHeight(msg.Height)
			m.customerView.SetWidth(msg.Width)
			m.customerView.SetHeight(msg.Height)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quiting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	cmds := []tea.Cmd{}

	switch m.focused {
	case customerViewFocused:
		m.projectView, cmd = m.projectView.Update(msg)
		cmds = append(cmds, cmd)
		m.customerView, cmd = m.customerView.Update(msg)
		cmds = append(cmds, cmd)

	case projectViewFocused:
		m.projectView, cmd = m.projectView.Update(msg)
		cmds = append(cmds, cmd)

	case activitiesViewFocused:
		m.activitiesView, cmd = m.activitiesView.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	if m.loading {
		return "loading..."
	}

	if m.quiting {
		return ""
	}

	switch m.focused {
	case customerViewFocused:
		fallthrough
	case projectViewFocused:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.customerView.View(),
			m.projectView.View(),
		)
	case activitiesViewFocused:
		fallthrough
	default:
		return m.activitiesView.View()
	}
}

func newModel() *model {
	// config := clockodo.Config{
	// 	EmailAddress: viper.GetString("clockodo.email"),
	// 	ApiToken:     viper.GetString("clockodo.token"),
	// }

	repo := local.NewRepository()

	customers, err := repo.GetAllCustomers()
	if err != nil {
		log.Fatal(err)
	}

	items := []list.Item{}

	for _, c := range customers {
		projects, err := repo.GetCustomerProjects(c.Name)
		if err != nil {
			log.Fatal(err)
		}

		for _, p := range projects {
			items = append(items, customer{
				name:        c.Name,
				description: p.Name,
			})
		}
	}

	m := &model{}

	customerList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	customerList.Title = "Select customer"
	customerList.SetItems(items)
	customerList.SetShowHelp(false)

	projectList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	projectList.Title = "Select project"
	projectList.SetItems([]list.Item{})
	projectList.SetShowHelp(false)

	m.focused = customerViewFocused
	m.customerView = customerList
	m.projectView = projectList
	m.repository = repo
	m.loading = true

	return m
}

func main() {
	m := newModel()

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	viper.SetConfigName("track")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("no config file found")
		}
	}

	if !viper.IsSet("clockodo.email") {
		log.Fatal("no email address set in configuration")
	}

	if !viper.IsSet("clockodo.token") {
		log.Fatal("no apikey set in configuration")
	}
}
