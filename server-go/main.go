package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	mathrand "math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const VERSION = "1.0.1"
const PROTOCOL_VERSION = "2024-11-05"

// 颜色主题 (Tokyo Night)
var (
	primaryColor   = lipgloss.Color("#7aa2f7")
	successColor   = lipgloss.Color("#9ece6a")
	errorColor     = lipgloss.Color("#f7768e")
	mutedColor     = lipgloss.Color("#565f89")
	backgroundColor = lipgloss.Color("#1a1b26")
)

// 样式
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#c0caf5")).
			Padding(0, 1)

	statusStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	inputStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)
)

// 消息存储
type Message struct {
	Type      string    `json:"type"`
	Text      string    `json:"text,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	ID        int       `json:"id"`
}

type Image struct {
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	MimeType  string    `json:"mimeType"`
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
	ID        int       `json:"id"`
}

var (
	messages []Message
	images   []Image
	msgMutex sync.RWMutex
	authToken string
)

// MCP 工具定义
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

var tools = []Tool{
	{
		Name:        "send_message",
		Description: "Send a text message to the server",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"message": map[string]interface{}{
					"type":        "string",
					"description": "The message to send",
				},
			},
			"required": []string{"message"},
		},
	},
	{
		Name:        "get_messages",
		Description: "Get all messages from the server",
		InputSchema: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
	},
	{
		Name:        "clear_messages",
		Description: "Clear all messages",
		InputSchema: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
	},
}

// JSON-RPC 结构
type JSONRPCRequest struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      interface{}            `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 生成随机 token
func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// 加载或生成 token
func loadOrGenerateToken() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("无法获取用户目录: %v", err)
		return generateToken()
	}

	dataDir := filepath.Join(homeDir, "data")
	tokenFile := filepath.Join(dataDir, "cicy-server.txt")

	// 读取现有 token
	if data, err := os.ReadFile(tokenFile); err == nil {
		token := strings.TrimSpace(string(data))
		if token != "" {
			log.Printf("已加载 token: %s", tokenFile)
			return token
		}
	}

	// 生成新 token
	token := generateToken()
	
	// 创建目录
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Printf("无法创建目录: %v", err)
		return token
	}

	// 保存 token
	if err := os.WriteFile(tokenFile, []byte(token), 0600); err != nil {
		log.Printf("无法保存 token: %v", err)
	} else {
		log.Printf("已生成新 token: %s", tokenFile)
	}

	return token
}

// 验证 token 中间件
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			token = r.Header.Get("X-Auth-Token")
		}
		
		// 移除 "Bearer " 前缀
		token = strings.TrimPrefix(token, "Bearer ")
		
		if token != authToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		next(w, r)
	}
}

// API 消息结构
type APIMessage struct {
	Type    string                   `json:"type"` // "text" or "image"
	Text    string                   `json:"text,omitempty"`
	URL     string                   `json:"url,omitempty"`
	Data    string                   `json:"data,omitempty"` // base64
	Content []map[string]interface{} `json:"content,omitempty"` // MCP 格式
}

// 全局 program 变量，用于发送消息到 TUI
var tuiProgram *tea.Program

// 图片消息结构
type imageMsg struct {
	path string
	size string
}

// API 处理器
func apiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg APIMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	msgMutex.Lock()
	defer msgMutex.Unlock()

	// 处理 MCP content 数组格式
	if len(msg.Content) > 0 {
		for _, item := range msg.Content {
			itemType, _ := item["type"].(string)
			
			switch itemType {
			case "text":
				text, _ := item["text"].(string)
				if text != "" {
					messages = append(messages, Message{
						Type:      "text",
						Text:      text,
						Timestamp: time.Now(),
						ID:        len(messages) + 1,
					})
					log.Printf("📝 收到文本消息: %s", text)
					
					// 发送消息到 TUI
					if tuiProgram != nil {
						tuiProgram.Send(newMessageMsg{text: text})
					}
				}
				
			case "image":
				imageURL, _ := item["url"].(string)
				imageData, _ := item["data"].(string)
				
				if imageURL == "" && imageData == "" {
					continue
				}
				
				finalData := imageData
				imageSize := len(imageData)
				
				if imageURL != "" {
					// 从 URL 下载图片
					resp, err := http.Get(imageURL)
					if err != nil {
						log.Printf("❌ 下载图片失败: %v", err)
						continue
					}
					defer resp.Body.Close()

					var buf []byte
					buf = make([]byte, resp.ContentLength)
					resp.Body.Read(buf)
					finalData = base64.StdEncoding.EncodeToString(buf)
					imageSize = len(buf)
				} else {
					decoded, _ := base64.StdEncoding.DecodeString(imageData)
					imageSize = len(decoded)
				}

				images = append(images, Image{
					Type:      "image",
					Name:      fmt.Sprintf("image_%d", len(images)+1),
					MimeType:  "image/png",
					Data:      finalData,
					Timestamp: time.Now(),
					ID:        len(images) + 1,
				})
				
				sizeStr := formatSize(imageSize)
				log.Printf("🖼️  收到图片消息 (大小: %s)", sizeStr)
				
				// 保存图片到文件
				imagePath, err := saveImageToFile(finalData)
				if err != nil {
					log.Printf("❌ 保存图片失败: %v", err)
					continue
				}
				
				// 发送图片消息到 TUI
				if tuiProgram != nil {
					tuiProgram.Send(imageMsg{path: imagePath, size: sizeStr})
				}
			}
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Message received",
		})
		return
	}

	// 处理旧格式（单个消息）
	switch msg.Type {
	case "text":
		if msg.Text == "" {
			http.Error(w, "Text is required", http.StatusBadRequest)
			return
		}
		messages = append(messages, Message{
			Type:      "text",
			Text:      msg.Text,
			Timestamp: time.Now(),
			ID:        len(messages) + 1,
		})
		log.Printf("📝 收到文本消息: %s", msg.Text)
		
		// 发送消息到 TUI
		if tuiProgram != nil {
			tuiProgram.Send(newMessageMsg{text: msg.Text})
		}

	case "image":
		if msg.URL == "" && msg.Data == "" {
			http.Error(w, "URL or Data is required", http.StatusBadRequest)
			return
		}
		
		imageData := msg.Data
		imageSize := len(msg.Data)
		
		if msg.URL != "" {
			resp, err := http.Get(msg.URL)
			if err != nil {
				http.Error(w, "Failed to download image", http.StatusBadRequest)
				return
			}
			defer resp.Body.Close()

			var buf []byte
			buf = make([]byte, resp.ContentLength)
			resp.Body.Read(buf)
			imageData = base64.StdEncoding.EncodeToString(buf)
			imageSize = len(buf)
		} else {
			decoded, _ := base64.StdEncoding.DecodeString(msg.Data)
			imageSize = len(decoded)
		}

		images = append(images, Image{
			Type:      "image",
			Name:      fmt.Sprintf("image_%d", len(images)+1),
			MimeType:  "image/png",
			Data:      imageData,
			Timestamp: time.Now(),
			ID:        len(images) + 1,
		})
		
		sizeStr := formatSize(imageSize)
		log.Printf("🖼️  收到图片消息 (大小: %s)", sizeStr)
		
		// 保存图片到文件
		imagePath, err := saveImageToFile(imageData)
		if err != nil {
			log.Printf("❌ 保存图片失败: %v", err)
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
		
		// 发送图片消息到 TUI
		if tuiProgram != nil {
			tuiProgram.Send(imageMsg{path: imagePath, size: sizeStr})
		}

	default:
		http.Error(w, "Invalid type", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Message received",
	})
}

// 保存图片到临时文件
func saveImageToFile(base64Data string) (string, error) {
	// 解码 base64
	decoded, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", err
	}
	
	// 创建目录
	homeDir, _ := os.UserHomeDir()
	imageDir := filepath.Join(homeDir, "Desktop", "images")
	os.MkdirAll(imageDir, 0755)
	
	// 生成文件名（使用时间戳）
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("image_%s.png", timestamp)
	filepath := filepath.Join(imageDir, filename)
	
	// 保存文件
	err = os.WriteFile(filepath, decoded, 0644)
	if err != nil {
		return "", err
	}
	
	return filepath, nil
}

// 格式化文件大小
func formatSize(size int) string {
	if size > 1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
	}
	if size > 1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	}
	return fmt.Sprintf("%d bytes", size)
}

// 在终端显示图片（iTerm2 内联图片协议）
func displayImageInTerminal(base64Data string) {
	// iTerm2 图片协议格式
	// ESC ] 1337 ; File=inline=1;width=auto;height=auto : <base64> BEL
	fmt.Printf("\033]1337;File=inline=1;width=40;height=auto:%s\a\n", base64Data)
}
var responses = []string{
	"好",
	"收到",
	"了解",
	"明白",
	"知道了",
	"没问题",
	"好的好的",
	"收到了",
	"明白了",
	"我知道了",
	"好的我明白",
	"收到你的消息",
	"了解了解",
	"明白了会处理",
	"好的马上开始",
	"收到了我会认真处理",
	"明白了这个任务我会仔细完成",
	"好的我已经收到你的指示了",
	"了解我会立即开始处理这个任务",
	"收到了我会按照要求来做请放心",
	"明白了我会认真完成这个任务完成后会及时汇报结果",
}

// 多行回复
var multiLineResponses = [][]string{
	{"hi", "how are you"},
	{"好的", "收到了"},
	{"明白", "马上处理"},
	{"了解", "我会做好的"},
	{"收到", "正在处理中"},
}

func getRandomResponse() string {
	// 30% 概率返回多行
	if mathrand.Float32() < 0.3 {
		lines := multiLineResponses[mathrand.Intn(len(multiLineResponses))]
		return strings.Join(lines, "\n")
	}
	return responses[mathrand.Intn(len(responses))]
}

// TUI Model
type model struct {
	input        string
	messages     []string
	pendingImage string // 待打开的图片路径
	loading      bool
	loadingDots  int
	startTime    time.Time
	width        int
	height       int
	serverPort   int
	ctrlCCount   int
	lastCtrlC    time.Time
	sshMode      bool
	sshHosts     []string
	sshSelected  int
	sshConnected string // 已连接的 SSH 主机
}

type tickMsg time.Time
type responseMsg struct {
	text     string
	duration time.Duration
}
type newMessageMsg struct {
	text string
}

func initialModel(port int) model {
	// ASCII Logo - 更宽更大
	logo := []string{
		"",
		"",
		"      ██████╗ ██╗ ██████╗ ██╗   ██╗",
		"     ██╔════╝ ██║██╔════╝ ╚██╗ ██╔╝",
		"     ██║      ██║██║       ╚████╔╝ ",
		"     ██║      ██║██║        ╚██╔╝  ",
		"     ╚██████╗ ██║╚██████╗    ██║   ",
		"      ╚═════╝ ╚═╝ ╚═════╝    ╚═╝   ",
		"",
		"        MCP Message Communication System",
		"        ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━",
		"",
	}
	
	messages := make([]string, len(logo))
	copy(messages, logo)
	
	if port != 0 {
		messages = append(messages, fmt.Sprintf("        🚀 服务器已启动 (端口: %d)", port))
		messages = append(messages, "")
	}
	
	return model{
		messages:   messages,
		serverPort: port,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// SSH 模式下的按键处理
		if m.sshMode {
			switch msg.String() {
			case "esc", "q":
				m.sshMode = false
				m.input = ""
				return m, nil
			
			case "up", "k":
				if m.sshSelected > 0 {
					m.sshSelected--
				}
				return m, nil
			
			case "down", "j":
				if m.sshSelected < len(m.sshHosts)-1 {
					m.sshSelected++
				}
				return m, nil
			
			case "enter":
				if m.sshSelected < len(m.sshHosts) {
					selected := m.sshHosts[m.sshSelected]
					m.sshConnected = selected
					m.messages = append(m.messages, fmt.Sprintf("✓ 已连接到: %s", selected))
					m.sshMode = false
					m.input = ""
				} else {
					// 调试：显示索引信息
					m.messages = append(m.messages, fmt.Sprintf("错误: 索引 %d >= 长度 %d", m.sshSelected, len(m.sshHosts)))
				}
				return m, nil
			}
			return m, nil
		}

		// 正常模式下的按键处理
		switch msg.String() {
		case "ctrl+c":
			now := time.Now()
			// 如果距离上次 Ctrl+C 超过 2 秒，重置计数
			if now.Sub(m.lastCtrlC) > 2*time.Second {
				m.ctrlCCount = 0
			}
			
			m.ctrlCCount++
			m.lastCtrlC = now
			
			if m.ctrlCCount >= 2 {
				return m, tea.Quit
			}
			
			// 第一次按 Ctrl+C，显示提示
			m.messages = append(m.messages, statusStyle.Render("  再按一次 Ctrl+C 退出"))
			return m, nil

		case "esc":
			return m, tea.Quit
		
		case "o":
			// 打开待查看的图片
			if m.pendingImage != "" {
				go openImage(m.pendingImage)
				m.messages = append(m.messages, statusStyle.Render("  ✓ 已打开图片"))
				m.pendingImage = ""
			}
			return m, nil

		case "enter":
			if m.input == "" {
				return m, nil
			}

			// 处理 /ssh 命令
			if m.input == "/ssh" {
				hosts := getSSHHosts()
				if len(hosts) == 0 {
					m.messages = append(m.messages, "  未找到 SSH 配置")
					m.input = ""
				} else {
					m.sshMode = true
					m.sshHosts = hosts
					m.sshSelected = 0
					m.input = ""
				}
				return m, nil
			}

			// 处理 /exit 命令（断开 SSH）
			if m.input == "/exit" && m.sshConnected != "" {
				m.messages = append(m.messages, fmt.Sprintf("✓ 已断开: %s", m.sshConnected))
				m.sshConnected = ""
				m.input = ""
				return m, nil
			}

			// 如果已连接 SSH，转发命令
			if m.sshConnected != "" {
				m.messages = append(m.messages, fmt.Sprintf("$ %s", m.input))
				m.loading = true
				m.startTime = time.Now()

				cmd := m.input
				host := m.sshConnected
				m.input = ""

				return m, tea.Batch(
					tickCmd(),
					func() tea.Msg {
						// 执行 SSH 命令
						output := executeSSHCommand(host, cmd)
						duration := time.Since(m.startTime)
						return responseMsg{text: output, duration: duration}
					},
				)
			}

			// 检查服务器是否启动
			if m.serverPort == 0 {
				m.messages = append(m.messages, fmt.Sprintf("你: %s", m.input))
				m.messages = append(m.messages, "❌ 错误: 服务器未启动，无法发送消息")
				m.input = ""
				return m, nil
			}

			// 发送消息
			m.messages = append(m.messages, fmt.Sprintf("你: %s", m.input))
			m.loading = true
			m.startTime = time.Now()

			input := m.input
			m.input = ""

			return m, tea.Batch(
				tickCmd(),
				func() tea.Msg {
					// 调用本地 API
					resp := sendMessageToServer(input, m.serverPort)
					duration := time.Since(m.startTime)
					return responseMsg{text: resp, duration: duration}
				},
			)

		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}

		default:
			if len(msg.String()) == 1 {
				m.input += msg.String()
			}
		}

	case tickMsg:
		if m.loading {
			m.loadingDots = (m.loadingDots + 1) % 4
			return m, tickCmd()
		}

	case responseMsg:
		m.loading = false
		// 处理多行回复
		lines := strings.Split(msg.text, "\n")
		for _, line := range lines {
			m.messages = append(m.messages, fmt.Sprintf("✓ %s", line))
		}
		// 显示秒数，保留2位小数
		seconds := float64(msg.duration.Milliseconds()) / 1000.0
		m.messages = append(m.messages, statusStyle.Render(fmt.Sprintf("  - %.2f", seconds)))
		return m, nil
	
	case newMessageMsg:
		// 从 API 收到的新消息
		m.messages = append(m.messages, fmt.Sprintf("📨 %s", msg.text))
		return m, nil
	
	case imageMsg:
		// 从 API 收到的图片消息
		m.pendingImage = msg.path
		m.messages = append(m.messages, fmt.Sprintf("🖼️  收到图片 (%s)", msg.size))
		m.messages = append(m.messages, statusStyle.Render("  按 'o' 打开图片"))
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m model) View() string {
	// SSH 选择模式
	if m.sshMode {
		// 标题（放在边框内）
		title := lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Render("选择 SSH 主机")
		
		// 主机列表
		var items []string
		for i, host := range m.sshHosts {
			if i == m.sshSelected {
				// 选中项 - 绿色背景 + 黑色文字 + 箭头
				item := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#000000")).
					Background(successColor).
					Bold(true).
					Padding(0, 1).
					Render(fmt.Sprintf("▶ %s", host))
				items = append(items, item)
			} else {
				// 未选中项
				item := lipgloss.NewStyle().
					Foreground(mutedColor).
					Render(fmt.Sprintf("  %s", host))
				items = append(items, item)
			}
		}
		
		// 组合内容：标题 + 空行 + 列表
		content := title + "\n\n" + strings.Join(items, "\n")
		
		// 帮助信息
		help := statusStyle.Render("↑/↓: 选择 | Enter: 确认 | ESC: 取消")
		
		// 创建边框 - 固定宽度 50
		boxStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Width(50).
			Padding(1, 2)
		
		box := boxStyle.Render(content)
		
		// 居中显示
		centered := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box+"\n\n"+help)
		
		return centered
	}

	// 正常模式
	// 标题
	title := titleStyle.Render("CICY - MCP 消息系统")

	// 计算可用于消息显示的行数
	// 标题(1行) + 空行(1行) + 输入框(3行) + 帮助(1行) + 空行(2行) = 8行
	availableLines := m.height - 8
	if availableLines < 5 {
		availableLines = 5
	}

	// 消息列表（根据可用行数动态调整）
	// 对于 logo 行，使用居中样式
	msgList := ""
	start := 0
	if len(m.messages) > availableLines {
		start = len(m.messages) - availableLines
	}
	
	// 创建居中样式
	centerStyle := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center)
	
	for i, msg := range m.messages[start:] {
		actualIndex := start + i
		// 前 12 行是 logo 和启动信息，需要居中
		if actualIndex < 12 {
			msgList += centerStyle.Render(msg) + "\n"
		} else {
			msgList += messageStyle.Render(msg) + "\n"
		}
	}

	// Loading 动画
	loadingText := ""
	if m.loading {
		dots := ""
		for i := 0; i < m.loadingDots; i++ {
			dots += "."
		}
		loadingText = statusStyle.Render(fmt.Sprintf("  发送中%s", dots)) + "\n"
	}

	// 输入框（固定在底部，宽度占满窗口）
	prompt := ">"
	if m.sshConnected != "" {
		prompt = fmt.Sprintf("[%s]>", m.sshConnected)
	}
	inputContent := fmt.Sprintf("%s %s█", prompt, m.input)
	
	// 计算输入框宽度（窗口宽度 - 4，留出边距）
	inputWidth := m.width - 4
	if inputWidth < 20 {
		inputWidth = 20
	}
	
	inputBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Width(inputWidth).
		Padding(0, 1).
		Render(inputContent)

	// 帮助
	helpText := "Ctrl+C 两次退出 | ESC 退出"
	if m.pendingImage != "" {
		helpText = "按 'o' 打开图片 | " + helpText
	} else if m.sshConnected != "" {
		helpText = "/exit 断开SSH | " + helpText
	}
	help := statusStyle.Render("  " + helpText)

	return fmt.Sprintf("%s\n\n%s%s\n%s\n%s",
		title,
		msgList,
		loadingText,
		inputBox,
		help,
	)
}

// HTTP 服务器
func startServer(port int) (chan bool, error) {
	ready := make(chan bool)
	
	// 加载或生成 token
	authToken = loadOrGenerateToken()
	
	// 先检查端口是否可用
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("端口 %d 已被占用或无法使用", port)
	}
	
	http.HandleFunc("/mcp", mcpHandler)
	http.HandleFunc("/message", messageHandler)
	http.HandleFunc("/messages", messagesHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/message", authMiddleware(apiHandler))
	
	go func() {
		log.Printf("MCP Server listening on http://localhost%s\n", addr)
		log.Printf("API Endpoint: POST /api/message (需要 token 认证)\n")
		ready <- true
		if err := http.Serve(listener, nil); err != nil {
			log.Fatal(err)
		}
	}()
	
	return ready, nil
}

func mcpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, req.ID, -32700, "Parse error")
		return
	}

	switch req.Method {
	case "initialize":
		sendResponse(w, req.ID, map[string]interface{}{
			"protocolVersion": PROTOCOL_VERSION,
			"capabilities": map[string]interface{}{
				"tools": map[string]bool{"listChanged": true},
			},
			"serverInfo": map[string]string{
				"name":    "cicy-go-server",
				"version": VERSION,
			},
		})

	case "tools/list":
		sendResponse(w, req.ID, map[string]interface{}{
			"tools": tools,
		})

	case "tools/call":
		handleToolCall(w, req)

	default:
		sendError(w, req.ID, -32601, fmt.Sprintf("Method not found: %s", req.Method))
	}
}

func handleToolCall(w http.ResponseWriter, req JSONRPCRequest) {
	name, _ := req.Params["name"].(string)
	args, _ := req.Params["arguments"].(map[string]interface{})

	switch name {
	case "send_message":
		message, _ := args["message"].(string)
		if message == "" {
			sendError(w, req.ID, -32602, "Invalid params: message required")
			return
		}

		msgMutex.Lock()
		msg := Message{
			Type:      "text",
			Text:      message,
			Timestamp: time.Now(),
			ID:        len(messages) + 1,
		}
		messages = append(messages, msg)
		msgMutex.Unlock()

		reply := getRandomResponse()
		sendResponse(w, req.ID, map[string]interface{}{
			"content": []map[string]string{
				{"type": "text", "text": reply},
			},
			"isError": false,
		})

	case "get_messages":
		msgMutex.RLock()
		allMessages := append([]Message{}, messages...)
		msgMutex.RUnlock()

		data, _ := json.MarshalIndent(allMessages, "", "  ")
		sendResponse(w, req.ID, map[string]interface{}{
			"content": []map[string]string{
				{"type": "text", "text": string(data)},
			},
			"isError": false,
		})

	case "clear_messages":
		msgMutex.Lock()
		messages = []Message{}
		images = []Image{}
		msgMutex.Unlock()

		sendResponse(w, req.ID, map[string]interface{}{
			"content": []map[string]string{
				{"type": "text", "text": "All messages cleared"},
			},
			"isError": false,
		})

	default:
		sendError(w, req.ID, -32601, fmt.Sprintf("Tool not found: %s", name))
	}
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, _ := body["message"].(string)
	if message == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "No message provided",
		})
		return
	}

	msgMutex.Lock()
	msg := Message{
		Type:      "text",
		Text:      message,
		Timestamp: time.Now(),
		ID:        len(messages) + 1,
	}
	messages = append(messages, msg)
	msgMutex.Unlock()

	reply := getRandomResponse()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": reply,
	})
}

func messagesHandler(w http.ResponseWriter, r *http.Request) {
	msgMutex.RLock()
	defer msgMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	msgMutex.RLock()
	msgCount := len(messages)
	imgCount := len(images)
	msgMutex.RUnlock()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "ok",
		"protocol": "mcp",
		"version":  PROTOCOL_VERSION,
		"messages": msgCount,
		"images":   imgCount,
	})
}

func sendResponse(w http.ResponseWriter, id interface{}, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	})
}

func sendError(w http.ResponseWriter, id interface{}, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
		},
	})
}

// 发送消息到本地服务器
func sendMessageToServer(message string, port int) string {
	url := fmt.Sprintf("http://localhost:%d/message", port)
	body := map[string]string{"message": message}
	data, _ := json.Marshal(body)

	resp, err := http.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return "错误: 无法连接到服务器"
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if msg, ok := result["message"].(string); ok {
		return msg
	}
	return "收到"
}

// 打开图片文件
func openImage(path string) {
	cmd := exec.Command("open", path)
	if err := cmd.Run(); err != nil {
		log.Printf("❌ 打开图片失败: %v", err)
	}
}

// 读取 SSH 配置文件中的主机名
func getSSHHosts() []string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	configPath := filepath.Join(homeDir, ".ssh", "config")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil
	}

	var hosts []string
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Host ") {
			host := strings.TrimPrefix(line, "Host ")
			host = strings.TrimSpace(host)
			// 跳过通配符
			if !strings.Contains(host, "*") && host != "" {
				hosts = append(hosts, host)
			}
		}
	}

	return hosts
}

// 执行 SSH 命令
func executeSSHCommand(host, command string) string {
	cmd := exec.Command("ssh", host, command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("错误: %v", err)
	}
	return strings.TrimSpace(string(output))
}

func main() {
	// 命令行参数
	helpFlag := flag.Bool("help", false, "显示帮助信息")
	versionFlag := flag.Bool("version", false, "显示版本号")
	portFlag := flag.Int("port", 13001, "服务器端口")
	flag.BoolVar(helpFlag, "h", false, "显示帮助信息")
	flag.BoolVar(versionFlag, "v", false, "显示版本号")
	flag.IntVar(portFlag, "p", 13001, "服务器端口")
	flag.Parse()

	if *helpFlag {
		fmt.Printf(`
CICY - MCP Message Communication System v%s (Go Edition)

用法 (Usage):
  cicy-go [选项]

选项 (Options):
  -h, --help       显示帮助信息
  -v, --version    显示版本号
  -p, --port PORT  指定端口 (默认: 13001)

功能 (Features):
  • 单进程运行 TUI 客户端 + MCP 服务器
  • 高性能 Go 实现
  • 内存占用低 (~15MB)
  • 启动速度快 (<10ms)

快捷键 (Shortcuts):
  Enter      发送消息
  Ctrl+C     退出
  ESC        退出

`, VERSION)
		os.Exit(0)
	}

	if *versionFlag {
		fmt.Printf("cicy-go v%s\n", VERSION)
		os.Exit(0)
	}

	// 启动 HTTP 服务器
	serverPort := *portFlag
	ready, err := startServer(serverPort)
	if err != nil {
		// 端口被占用，只显示警告，不启动服务器
		fmt.Printf("⚠️  警告: 端口 %d 已被占用，服务器未启动\n", serverPort)
		fmt.Printf("提示: TUI 将继续运行，但无法发送消息\n\n")
		time.Sleep(2 * time.Second) // 让用户看到警告
		serverPort = 0 // 标记服务器未启动
	} else {
		<-ready // 等待服务器就绪
	}

	// 启动 TUI
	p := tea.NewProgram(initialModel(serverPort), tea.WithAltScreen())
	tuiProgram = p // 保存全局引用
	
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
