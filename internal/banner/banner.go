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

const serviceText = ` ██████╗ ██████╗ ██████╗ ███████╗       ██████╗ █████╗ ████████╗ ██████╗██╗  ██╗███████╗██████╗
██╔════╝██╔═══██╗██╔══██╗██╔════╝      ██╔════╝██╔══██╗╚══██╔══╝██╔════╝██║  ██║██╔════╝██╔══██╗
██║     ██║   ██║██████╔╝█████╗  █████╗██║     ███████║   ██║   ██║     ███████║█████╗  ██████╔╝
██║     ██║   ██║██╔══██╗██╔══╝  ╚════╝██║     ██╔══██║   ██║   ██║     ██╔══██║██╔══╝  ██╔══██╗
╚██████╗╚██████╔╝██║  ██║███████╗       ╚██████╗██║  ██║   ██║   ╚██████╗██║  ██║███████╗██║  ██║
 ╚═════╝ ╚═════╝ ╚═╝  ╚═╝╚══════╝        ╚═════╝╚═╝  ╚═╝   ╚═╝    ╚═════╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝`

// Glitch characters used during the intro phase.
const glitchChars = "░▒▓█▄▀▐▌╠╣╬═║╗╝╚╔"

// Mitt art for the animation (catching a ball).
var mittArt = []string{
	"    ___     ",
	"  / _ _ \\   ",
	" | (o o) |  ",
	"  \\ === /   ",
	"   \\___/    ",
}

const fieldWidth = 70

var (
	// Green gradient — catcher style
	greenHot = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00CC66")).
			Bold(true)

	greenBright = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF99")).
			Bold(true)

	dimGreen = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#005522"))

	// The "2" in a contrasting cyan
	accentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFCC")).
			Bold(true)

	serviceBlockStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#009966")).
				Bold(true)

	mittStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true)
)

type tickMsg time.Time

type model struct {
	width        int
	frame        int
	glitchPhase  bool
	glitchFrames int
	done         bool
}

// Show displays the animated banner for a brief duration, then prints
// the final frame as a persistent header before returning.
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
		bannerOutput = greenHot.Render(bannerText)
	}

	b.WriteString(bannerOutput)

	the2 := accentStyle.Render("2")
	if m.glitchPhase && m.glitchFrames < 8 {
		glitchRunes := []rune(glitchChars)
		the2 = accentStyle.Render(string(glitchRunes[rand.IntN(len(glitchRunes))]))
	}
	b.WriteString(the2)
	b.WriteString("\n\n")

	// Mitt animation — sliding across from right to left (catching)
	if !m.glitchPhase {
		mittPos := fieldWidth - (m.frame*3)%fieldWidth
		for _, artLine := range mittArt {
			pad := strings.Repeat(" ", mittPos)
			b.WriteString(mittStyle.Render(pad + artLine))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	} else {
		b.WriteString("\n\n")
	}

	var svcOutput string
	if m.glitchPhase {
		svcOutput = glitchText(serviceText, m.glitchFrames)
	} else {
		svcOutput = serviceBlockStyle.Render(serviceText)
	}
	b.WriteString(svcOutput)
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

	return greenBright.Render(string(result))
}

func applyScanlines(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if i%2 == 1 {
			lines[i] = dimGreen.Render(line)
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
	b.WriteString(greenHot.Render(bannerText))
	b.WriteString(accentStyle.Render("2"))
	b.WriteString("\n")
	b.WriteString(serviceBlockStyle.Render(serviceText))
	b.WriteString("\n")
	return b.String()
}
