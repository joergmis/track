package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joergmis/track"
	"github.com/joergmis/track/clockodo"
	"github.com/spf13/viper"
)

const (
	customerViewFocused = iota
	projectViewFocused
	activitiesViewFocused
)

type customer struct {
	name string
}

func (c customer) FilterValue() string {
	return c.name
}
func (c customer) Title() string {
	return c.name
}
func (c customer) Description() string {
	return ""
}

type project struct {
	name        string
	description string
}

func (p project) FilterValue() string {
	return p.name
}
func (p project) Title() string {
	return p.name
}
func (p project) Description() string {
	return p.description
}

type model struct {
	focused int

	customerView   list.Model
	projectView    list.Model
	activitiesView list.Model

	customerRepository track.CustomerRepository
	projectRepository  track.ProjectRepository

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

func (m *model) selectCustomer() {
	customer, ok := m.customerView.SelectedItem().(customer)
	if !ok {
		tea.Quit()
		log.Fatal("failed to cast to customer")
	}

	projects, err := m.projectRepository.GetCustomerProjects(track.Customer{Name: customer.Title()})
	if err != nil {
		tea.Quit()
		log.Fatal(err)
	}

	items := []list.Item{}

	for _, p := range projects {
		items = append(items, project{
			name:        p.Name,
			description: fmt.Sprintf("completed: %v, active: %v", p.Completed, p.Active),
		})
	}

	m.projectView.SetItems(items)
	m.projectView.ResetFilter()
	m.focused = projectViewFocused
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

		case "enter":
			m.selectCustomer()

		case "left", "h":
			m.deselectCustomer()
		}
	}

	var cmd tea.Cmd

	switch m.focused {
	case customerViewFocused:
		m.customerView, cmd = m.customerView.Update(msg)
	case projectViewFocused:
		m.projectView, cmd = m.projectView.Update(msg)
	case activitiesViewFocused:
		m.activitiesView, cmd = m.activitiesView.Update(msg)
	}

	return m, cmd
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
	config := clockodo.Config{
		EmailAddress: viper.GetString("clockodo.email"),
		ApiToken:     viper.GetString("clockodo.token"),
	}

	customerRepository, err := clockodo.NewCustomerRepository(config)
	if err != nil {
		log.Fatal(err)
	}

	customers, err := customerRepository.GetAllCustomers()
	if err != nil {
		log.Fatal(err)
	}

	items := []list.Item{}

	for _, c := range customers {
		items = append(items, customer{
			name: c.Name,
		})
	}

	customerList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	customerList.Title = "Select project"
	customerList.SetItems(items)
	customerList.SetShowHelp(false)

	projecRepository, err := clockodo.NewProjectRepository(config)
	if err != nil {
		log.Fatal(err)
	}

	_, err = projecRepository.GetAllProjects()
	if err != nil {
		log.Fatal(err)
	}

	projectList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	projectList.Title = "Select customer"
	projectList.SetItems([]list.Item{})
	projectList.SetShowHelp(false)

	return &model{
		focused:            customerViewFocused,
		customerView:       customerList,
		projectView:        projectList,
		customerRepository: customerRepository,
		projectRepository:  projecRepository,
		loading:            true,
	}
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
