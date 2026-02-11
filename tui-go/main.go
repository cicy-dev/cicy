package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const API_URL = "http://localhost:13001"

type message struct {
	question string
	answer   string
	elapsed  float64
}

type model struct {
	input    textinput.Model
	spinner  spinner.Model
	messages []message
	loading  bool
	showHelp bool
	err      string
}

type responseMsg struct {
	text    string
	elapsed float64
}

var (
	titleColor  = lipgloss.Color("#7aa2f7")
	userColor   = lipgloss.Color("#f7768e")
	aiColor     = lipgloss.Color("#9ece6a")
	timeColor   = lipgloss.Color("#565f89")
	helpColor   = lipgloss.Color("#bb9af7")
	errorColor  = lipgloss.Color("#f7768e")
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func initialModel() model {
	// 输入框
	in := textinput.New()
	in.Placeholder = ""
	in.Focus()
	in.Width = 60
	in.Prompt = "> "
	in.PromptStyle = lipgloss.NewStyle().Foreground(userColor)

	// Loading spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(aiColor)

	return model{
		input:    in,
		spinner:  s,
		messages: []message{},
		loading:  false,
		showHelp: false,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// 帮助界面按键处理
		if m.showHelp {
			if msg.Type == tea.KeyEsc || msg.String() == "q" {
				m.showHelp = false
			}
			return m, nil
		}

		switch msg.Type {
		case tea.KeyEnter:
			if m.loading {
				return m, nil
			}

			value := strings.TrimSpace(m.input.Value())
			if value == "" {
				return m, nil
			}

			// 处理命令
			if strings.HasPrefix(value, "/") {
				return m.handleCommand(value)
			}

			// 发送普通消息
			m.input.SetValue("")
			m.loading = true
			m.err = ""

			return m, tea.Batch(
				m.spinner.Tick,
				sendRequest(value),
			)

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case responseMsg:
		// 收到响应
		m.loading = false
		m.messages = append(m.messages, message{
			question: m.input.Value(),
			answer:   msg.text,
			elapsed:  msg.elapsed,
		})
		return m, nil

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) handleCommand(cmd string) (tea.Model, tea.Cmd) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return m, nil
	}

	command := strings.ToLower(parts[0])

	switch command {
	case "/quit", "/q":
		return m, tea.Quit

	case "/clear", "/c":
		m.messages = []message{}
		m.err = ""
		m.input.SetValue("")

	case "/list", "/l":
		// 显示所有消息（已经在界面上）
		m.input.SetValue("")

	case "/help", "/h":
		m.showHelp = true
		m.input.SetValue("")

	default:
		m.err = fmt.Sprintf("Unknown command: %s", command)
		m.input.SetValue("")
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	// 帮助界面
	if m.showHelp {
		return m.renderHelp()
	}

	// 标题
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(titleColor).
		Render("◇ CICY")
	b.WriteString(title)
	b.WriteString("\n\n")

	// 显示历史对话（最近 5 条）
	start := 0
	if len(m.messages) > 5 {
		start = len(m.messages) - 5
	}

	for i := start; i < len(m.messages); i++ {
		msg := m.messages[i]

		// 用户消息
		userMsg := lipgloss.NewStyle().
			Foreground(userColor).
			Render(msg.question)
		b.WriteString(userMsg)
		b.WriteString("\n")

		// AI 回复
		aiMsg := lipgloss.NewStyle().
			Foreground(aiColor).
			Render(msg.answer)
		b.WriteString(aiMsg)
		b.WriteString("\n")

		// 完成耗时
		timeMsg := lipgloss.NewStyle().
			Foreground(timeColor).
			Render(fmt.Sprintf(" - Completed in %.2fs", msg.elapsed))
		b.WriteString(timeMsg)
		b.WriteString("\n\n")
	}

	// 错误信息
	if m.err != "" {
		errMsg := lipgloss.NewStyle().
			Foreground(errorColor).
			Render(fmt.Sprintf("Error: %s", m.err))
		b.WriteString(errMsg)
		b.WriteString("\n\n")
	}

	// 输入框
	if m.loading {
		// Loading 状态
		loadingPrompt := lipgloss.NewStyle().
			Foreground(aiColor).
			Render(m.spinner.View() + " > ")
		b.WriteString(loadingPrompt)
	} else {
		// 正常输入
		b.WriteString(m.input.View())
	}

	b.WriteString("\n")

	// 提示信息
	hint := lipgloss.NewStyle().
		Foreground(timeColor).
		Render("Type /help for commands • Ctrl+C to quit")
	b.WriteString("\n")
	b.WriteString(hint)

	return b.String()
}

func (m model) renderHelp() string {
	var b strings.Builder

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(titleColor).
		Render("◇ CICY - Help")
	b.WriteString(title)
	b.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().Foreground(helpColor)
	cmdStyle := lipgloss.NewStyle().Foreground(aiColor).Bold(true)

	commands := []struct {
		cmd  string
		desc string
	}{
		{"/help, /h", "显示此帮助信息"},
		{"/quit, /q", "退出程序"},
		{"/clear, /c", "清空消息历史"},
		{"/list, /l", "显示所有消息"},
		{"Ctrl+C", "退出程序"},
		{"Esc", "退出帮助/退出程序"},
	}

	b.WriteString(helpStyle.Render("可用命令："))
	b.WriteString("\n\n")

	for _, cmd := range commands {
		b.WriteString("  ")
		b.WriteString(cmdStyle.Render(cmd.cmd))
		b.WriteString("  -  ")
		b.WriteString(cmd.desc)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	hint := lipgloss.NewStyle().
		Foreground(timeColor).
		Render("按 Esc 或 q 返回")
	b.WriteString(hint)

	return b.String()
}

func sendRequest(message string) tea.Cmd {
	return func() tea.Msg {
		startTime := time.Now()

		data := map[string]string{"message": message}
		jsonData, _ := json.Marshal(data)

		resp, err := http.Post(
			API_URL+"/message",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			elapsed := time.Since(startTime).Seconds()
			return responseMsg{
				text:    fmt.Sprintf("Error: %v", err),
				elapsed: elapsed,
			}
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		elapsed := time.Since(startTime).Seconds()

		if msg, ok := result["message"].(string); ok {
			return responseMsg{
				text:    msg,
				elapsed: elapsed,
			}
		}

		return responseMsg{
			text:    "Message received!",
			elapsed: elapsed,
		}
	}
}
