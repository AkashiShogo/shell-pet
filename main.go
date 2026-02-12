package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Pet はゲームの中心となるデータ構造
type Pet struct {
	Name  string `json:"name"`  // ペットの名前
	Level int    `json:"level"` // レベル（未使用だが将来の拡張用）
	Exp   int    `json:"exp"`   // 経験値（未使用だが将来の拡張用）
	Stage int    `json:"stage"` // 0: Egg（卵）, 1: Baby（幼年期）

	// 4大ステータス（コマンド種別に対応）
	STR int `json:"str"` // 攻撃力 (Git/Build系)
	VIT int `json:"vit"` // 防御力 (Docker/Infra系)
	INT int `json:"int"` // 知力 (Editor/Code系)
	AGI int `json:"agi"` // 素早さ (Shell/Net系)

	// 生存パラメータ（たまごっち要素）
	Hunger int `json:"hunger"` // 満腹度 (0-100)
	Bugs   int `json:"bugs"`   // 汚れ/バグ (0-100)
}

const (
	saveFile = "pet.json" // 保存ファイル名
)

// model はBubble Teaアプリケーションの状態を保持
type model struct {
	pet Pet
}

// Init は初期化処理（今回は不要だがBubble Teaのインターフェース要件）
func (m model) Init() tea.Cmd {
	return nil
}

// Update はキー入力などのメッセージを受け取り、状態を更新する
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			// 終了時にペット状態を保存
			savePet(m.pet)
			return m, tea.Quit

		case "f":
			// Feed（餌やり）処理
			m.pet = feedPet(m.pet)

		case "c":
			// Clean（掃除）処理
			m.pet.Bugs = 0
		}
	}
	return m, nil
}

// View は現在の状態を文字列としてレンダリング
func (m model) View() string {
	// AAアート取得
	art := getPetArt(m.pet)

	// ステータス表示を構築
	stats := buildStats(m.pet)

	// 汚れている場合の背景色
	bgColor := ""
	if m.pet.Bugs > 50 {
		bgColor = "#3a3a00" // 黄色がかった暗い背景
	}

	// レイアウト: 左にAA、右にステータス
	leftStyle := lipgloss.NewStyle().
		Width(40).
		Align(lipgloss.Center).
		Background(lipgloss.Color(bgColor))

	rightStyle := lipgloss.NewStyle().
		Width(40).
		Padding(1).
		Background(lipgloss.Color(bgColor))

	left := leftStyle.Render(art)
	right := rightStyle.Render(stats)

	content := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	// 操作説明
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("\n[f] Feed  [c] Clean  [q] Quit")

	return content + help
}

// feedPet は餌やり処理（ランダムでステータス上昇）
func feedPet(p Pet) Pet {
	// Hungerを回復（上限100）
	p.Hunger += 30
	if p.Hunger > 100 {
		p.Hunger = 100
	}

	// ランダムに1つのステータスを上昇
	statType := rand.Intn(4)
	switch statType {
	case 0:
		p.STR += rand.Intn(3) + 1
	case 1:
		p.VIT += rand.Intn(3) + 1
	case 2:
		p.INT += rand.Intn(3) + 1
	case 3:
		p.AGI += rand.Intn(3) + 1
	}

	// 一定確率でBugs（汚れ）が発生
	if rand.Float32() < 0.4 {
		p.Bugs += rand.Intn(15) + 5
		if p.Bugs > 100 {
			p.Bugs = 100
		}
	}

	// 累積ステータスで進化判定（STR+VIT+INT+AGI > 20で卵からBabyへ）
	totalStats := p.STR + p.VIT + p.INT + p.AGI
	if p.Stage == 0 && totalStats > 20 {
		p.Stage = 1
	}

	return p
}

// getPetArt はステージに応じたAAアートを返す
func getPetArt(p Pet) string {
	switch p.Stage {
	case 0: // Egg（卵）
		return `
    ___
   /   \
  |  o  |
   \___/
`
	case 1: // Baby（幼年期）
		// Bugsが多いと表情が変わる
		if p.Bugs > 50 {
			return `
      ___
     /   \
    | T T |
    |  ~  |
     \___/
    /|   |\
   / |   | \
     dirty!
`
		}
		return `
      ___
     /   \
    | ^ ^ |
    |  v  |
     \___/
    /|   |\
   / |   | \
`
	default:
		return "???"
	}
}

// buildStats はステータス情報を整形して返す
func buildStats(p Pet) string {
	// ステージ名
	stageName := "Egg"
	if p.Stage == 1 {
		stageName = "Baby"
	}

	// スタイル定義
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	barStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86"))

	// ステータスバー生成
	hungerBar := createBar(p.Hunger, 100, 20)
	bugsBar := createBar(p.Bugs, 100, 20)

	return fmt.Sprintf(`%s

%s
Level: %d  Stage: %s

【4大ステータス】
STR: %d  (Git/Build系)
VIT: %d  (Docker/Infra系)
INT: %d  (Editor/Code系)
AGI: %d  (Shell/Net系)

【生存パラメータ】
Hunger: %s
Bugs:   %s
`,
		titleStyle.Render("🥚 "+p.Name),
		lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("━━━━━━━━━━━━━━━━━"),
		p.Level, stageName,
		p.STR,
		p.VIT,
		p.INT,
		p.AGI,
		barStyle.Render(hungerBar),
		barStyle.Render(bugsBar),
	)
}

// createBar は値に応じたプログレスバーを生成
func createBar(current, max, width int) string {
	filled := int(float64(current) / float64(max) * float64(width))
	if filled > width {
		filled = width
	}
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return fmt.Sprintf("%s %d/%d", bar, current, max)
}

// loadPet は保存ファイルからペット情報を読み込む
func loadPet() Pet {
	data, err := os.ReadFile(saveFile)
	if err != nil {
		// ファイルが存在しない場合は新規ペット作成
		return Pet{
			Name:   "ShellPet",
			Level:  1,
			Exp:    0,
			Stage:  0,
			STR:    0,
			VIT:    0,
			INT:    0,
			AGI:    0,
			Hunger: 50,
			Bugs:   0,
		}
	}

	var pet Pet
	if err := json.Unmarshal(data, &pet); err != nil {
		// JSONパースエラーの場合も新規ペット作成
		return Pet{
			Name:   "ShellPet",
			Level:  1,
			Exp:    0,
			Stage:  0,
			STR:    0,
			VIT:    0,
			INT:    0,
			AGI:    0,
			Hunger: 50,
			Bugs:   0,
		}
	}

	return pet
}

// savePet はペット情報をJSONファイルに保存
func savePet(pet Pet) {
	data, err := json.MarshalIndent(pet, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Save error: %v\n", err)
		return
	}

	if err := os.WriteFile(saveFile, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
	}
}

func main() {
	// 乱数シード初期化
	rand.Seed(time.Now().UnixNano())

	// ペット情報読み込み
	pet := loadPet()

	// Bubble Teaアプリケーション起動
	p := tea.NewProgram(model{pet: pet})
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
