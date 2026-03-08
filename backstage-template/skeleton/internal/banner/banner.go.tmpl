package banner

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const bannerText = `██╗  ██╗ ██████╗ ███╗   ███╗███████╗██████╗ ██╗   ██╗███╗   ██╗
██║  ██║██╔═══██╗████╗ ████║██╔════╝██╔══██╗██║   ██║████╗  ██║
███████║██║   ██║██╔████╔██║█████╗  ██████╔╝██║   ██║██╔██╗ ██║
██╔══██║██║   ██║██║╚██╔╝██║██╔══╝  ██╔══██╗██║   ██║██║╚██╗██║
██║  ██║╚██████╔╝██║ ╚═╝ ██║███████╗██║  ██║╚██████╔╝██║ ╚████║
╚═╝  ╚═╝ ╚═════╝ ╚═╝     ╚═╝╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝`

// TODO: Generate ASCII art for ${{ values.serviceName }} at scaffold time or use a generic subtitle
const serviceText = `${{ values.fullName | upper }}`

const glitchChars = "░▒▓█▄▀▐▌╠╣╬═║╗╝╚╔"

const fieldWidth = 70

var (
{%- if values.serviceType == "pitcher" %}
	primaryColor = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6600")).
			Bold(true)
	brightColor = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF9900")).
			Bold(true)
	dimColor = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#994400"))
	accentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFCC00")).
			Bold(true)
	serviceBlockStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF4400")).
				Bold(true)
{%- endif %}
{%- if values.serviceType == "catcher" %}
	primaryColor = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00CC66")).
			Bold(true)
	brightColor = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF99")).
			Bold(true)
	dimColor = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#005522"))
	accentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFCC")).
			Bold(true)
	serviceBlockStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#009966")).
				Bold(true)
{%- endif %}
)

type tickMsg time.Time

type model struct {
	width        int
	frame        int
	glitchPhase  bool
	glitchFrames int
	done         bool
}

func Show() {
	p := tea.NewProgram(initialModel())
	_, _ = p.Run()
	fmt.Println(renderHeader())
}

func initialModel() model {
	return model{
		width:        80,
		glitchPhase:  true,
		glitchFrames: 0,
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		m.done = true
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case tickMsg:
		_ = msg
		m.frame++
		if m.glitchPhase {
			m.glitchFrames++
			if m.glitchFrames >= 10 {
				m.glitchPhase = false
			}
			return m, tickCmd()
		}
		if m.frame >= 40 {
			m.done = true
			return m, tea.Quit
		}
		return m, tickCmd()
	}
	return m, nil
}

func (m model) View() tea.View {
	if m.done {
		return tea.NewView("")
	}

	var b strings.Builder

	var bannerOutput string
	if m.glitchPhase {
		bannerOutput = glitchText(bannerText, m.glitchFrames)
	} else {
		bannerOutput = primaryColor.Render(bannerText)
	}
	b.WriteString(bannerOutput)

	the2 := accentStyle.Render("2")
	if m.glitchPhase && m.glitchFrames < 8 {
		glitchRunes := []rune(glitchChars)
		the2 = accentStyle.Render(string(glitchRunes[rand.IntN(len(glitchRunes))]))
	}
	b.WriteString(the2)
	b.WriteString("\n\n")

	b.WriteString(serviceBlockStyle.Render(serviceText))
	b.WriteString("\n")

	output := applyScanlines(b.String())
	v := tea.NewView(centerText(output, m.width))
	v.AltScreen = true
	return v
}

func glitchText(text string, glitchFrame int) string {
	glitchProbability := float64(10-glitchFrame) / 10.0
	if glitchProbability < 0 {
		glitchProbability = 0
	}
	runes := []rune(text)
	glitchRunes := []rune(glitchChars)
	result := make([]rune, len(runes))
	for i, r := range runes {
		if r == '\n' || r == ' ' {
			result[i] = r
			continue
		}
		if rand.Float64() < glitchProbability {
			result[i] = glitchRunes[rand.IntN(len(glitchRunes))]
		} else {
			result[i] = r
		}
	}
	return brightColor.Render(string(result))
}

func applyScanlines(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if i%2 == 1 {
			lines[i] = dimColor.Render(line)
		}
	}
	return strings.Join(lines, "\n")
}

func centerText(text string, width int) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		visLen := lipgloss.Width(line)
		if visLen < width {
			pad := (width - visLen) / 2
			lines[i] = strings.Repeat(" ", pad) + line
		}
	}
	return strings.Join(lines, "\n")
}

func renderHeader() string {
	var b strings.Builder
	b.WriteString(primaryColor.Render(bannerText))
	b.WriteString(accentStyle.Render("2"))
	b.WriteString("\n")
	b.WriteString(serviceBlockStyle.Render(serviceText))
	b.WriteString("\n")
	return b.String()
}
