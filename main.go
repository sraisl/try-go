package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

var (
	uiRenderer    = lipgloss.NewRenderer(os.Stderr)
	titleStyle    = uiRenderer.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	pathStyle     = uiRenderer.NewStyle().Foreground(lipgloss.Color("244"))
	labelStyle    = uiRenderer.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	queryStyle    = uiRenderer.NewStyle().Foreground(lipgloss.Color("229"))
	selectedStyle = uiRenderer.NewStyle().Bold(true).Foreground(lipgloss.Color("229")).Background(lipgloss.Color("62")).Padding(0, 1)
	itemStyle     = uiRenderer.NewStyle().Foreground(lipgloss.Color("250")).PaddingLeft(2)
	createStyle   = uiRenderer.NewStyle().Bold(true).Foreground(lipgloss.Color("48")).PaddingLeft(2)
	emptyStyle    = uiRenderer.NewStyle().Foreground(lipgloss.Color("214"))
	warningStyle  = uiRenderer.NewStyle().Bold(true).Foreground(lipgloss.Color("203"))
	helpStyle     = uiRenderer.NewStyle().Foreground(lipgloss.Color("241"))
)

type model struct {
	root     string
	path     string
	allDirs  []string
	dirs     []string
	query    string
	cursor   int
	confirm  string
	errMsg   string
	selected bool
}

func (m model) View() string {
	s := titleStyle.Render("tries") + "\n"
	s += labelStyle.Render("Aktuelles Verzeichnis") + "\n"
	s += pathStyle.Render(m.path) + "\n\n"
	s += labelStyle.Render("Suche") + " " + queryStyle.Render(m.query) + "\n\n"

	if len(m.dirs) == 0 && m.query != "" {
		s += emptyStyle.Render("Keine Treffer.") + "\n"
		if name := m.createDirName(); name != "" {
			s += createStyle.Render("+ "+name+" anlegen") + "\n"
		} else {
			s += emptyStyle.Render("Bitte einen Namen mit Buchstaben oder Zahlen eingeben.") + "\n"
		}
	} else {
		s += labelStyle.Render("Unterverzeichnisse") + "\n\n"

		for i, dir := range m.dirs {
			if i == m.cursor {
				s += selectedStyle.Render("> "+dir) + "\n"
			} else {
				s += itemStyle.Render(dir) + "\n"
			}
		}
	}

	if m.confirm != "" {
		s += "\n" + warningStyle.Render("Backspace erneut drücken, um "+m.confirm+" zu löschen.") + "\n"
	}

	if m.errMsg != "" {
		s += "\n" + warningStyle.Render(m.errMsg) + "\n"
	}

	s += "\n" + helpStyle.Render("Tippen zum Suchen · ↑/↓ bewegen · Enter öffnen/anlegen · Backspace bearbeiten/löschen · Esc beenden") + "\n"

	return s
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			if m.confirm != "" {
				m.confirm = ""
				m.errMsg = ""
				return m, nil
			}

			return m, tea.Quit

		case tea.KeyRunes:
			if len(msg.Runes) == 1 && msg.Runes[0] == 'q' && m.query == "" {
				return m, tea.Quit
			}

			m.query += string(msg.Runes)
			m.confirm = ""
			m.errMsg = ""
			m.applyFilter()

		case tea.KeyBackspace:
			if m.query != "" {
				runes := []rune(m.query)
				m.query = string(runes[:len(runes)-1])
				m.confirm = ""
				m.errMsg = ""
				m.applyFilter()
				return m, nil
			}

			if len(m.dirs) == 0 {
				return m, nil
			}

			selectedDir := m.dirs[m.cursor]
			if m.confirm != selectedDir {
				m.confirm = selectedDir
				m.errMsg = ""
				return m, nil
			}

			if err := removeDir(m.root, selectedDir); err != nil {
				m.errMsg = "Konnte " + selectedDir + " nicht löschen: " + err.Error()
				m.confirm = ""
				return m, nil
			}

			m.allDirs = removeDirName(m.allDirs, selectedDir)
			m.confirm = ""
			m.errMsg = ""
			m.applyFilter()

		case tea.KeyDown:
			if m.cursor < len(m.dirs)-1 {
				m.cursor++
				m.confirm = ""
				m.errMsg = ""
			}

		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
				m.confirm = ""
				m.errMsg = ""
			}

		case tea.KeyEnter:
			if len(m.dirs) == 0 || m.path != m.root {
				if m.query == "" || m.path != m.root {
					return m, nil
				}

				newDirName := m.createDirName()
				if newDirName == "" {
					return m, nil
				}

				newPath := filepath.Join(m.root, newDirName)
				if err := os.Mkdir(newPath, 0755); err != nil {
					return m, nil
				}

				m.path = newPath
				m.selected = true
				return m, tea.Quit
			}

			selectedDir := m.dirs[m.cursor]
			m.path = filepath.Join(m.path, selectedDir)
			m.selected = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *model) applyFilter() {
	m.dirs = filterDirs(m.allDirs, m.query)
	m.cursor = 0
	if len(m.dirs) == 0 {
		m.cursor = 0
	}
}

func (m model) createDirName() string {
	name := sanitizeDirName(m.query)
	if name == "" {
		return ""
	}

	return time.Now().Format("2006-01-02") + "-" + name
}

func filterDirs(dirs []string, query string) []string {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return dirs
	}

	terms := strings.Fields(query)
	var filtered []string

	for _, dir := range dirs {
		haystack := strings.ToLower(dir)
		matches := true

		for _, term := range terms {
			if !strings.Contains(haystack, term) {
				matches = false
				break
			}
		}

		if matches {
			filtered = append(filtered, dir)
		}
	}

	return filtered
}

func sanitizeDirName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))

	var b strings.Builder
	lastWasDash := false

	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			lastWasDash = false
			continue
		}

		if !lastWasDash {
			b.WriteRune('-')
			lastWasDash = true
		}
	}

	return strings.Trim(b.String(), "-")
}

func removeDir(root, name string) error {
	path := filepath.Join(root, name)
	if filepath.Dir(path) != root {
		return fmt.Errorf("ungültiger Pfad")
	}

	return os.RemoveAll(path)
}

func removeDirName(dirs []string, name string) []string {
	filtered := dirs[:0]

	for _, dir := range dirs {
		if dir != name {
			filtered = append(filtered, dir)
		}
	}

	return filtered
}

func readDirs(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var dirs []string

	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}

	return dirs, nil
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	path := filepath.Join(home, "src", "tries")
	dirs, err := readDirs(path)
	if err != nil {
		panic(err)
	}

	m := model{
		root:    path,
		path:    path,
		allDirs: dirs,
		dirs:    dirs,
	}
	_ = m

	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))

	finalModel, err := p.Run()
	if err != nil {
		panic(err)
	}

	if m, ok := finalModel.(model); ok && m.selected {
		fmt.Println(m.path)
	}
}
