package main

import (
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// 테마 설정: 리스트 선택 시 하이라이트 가독성 향상
type customTheme struct{ fyne.Theme }

func (m customTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if n == theme.ColorNameFocus {
		return color.Transparent
	}
	if n == theme.ColorNameSelection {
		return color.NRGBA{R: 200, G: 200, B: 200, A: 100}
	}
	return theme.DefaultTheme().Color(n, v)
}
func (m customTheme) Font(s fyne.TextStyle) fyne.Resource     { return theme.DefaultTheme().Font(s) }
func (m customTheme) Icon(n fyne.ThemeIconName) fyne.Resource { return theme.DefaultTheme().Icon(n) }
func (m customTheme) Size(n fyne.ThemeSizeName) float32       { return theme.DefaultTheme().Size(n) }

// 커스텀 입력창: 위/아래 방향키로 리스트 제어
type smartEntry struct {
	widget.Entry
	onUp, onDown, onEnter func()
}

func newSmartEntry() *smartEntry {
	e := &smartEntry{}
	e.ExtendBaseWidget(e)
	return e
}

func (e *smartEntry) TypedKey(k *fyne.KeyEvent) {
	switch k.Name {
	case fyne.KeyUp:
		if e.onUp != nil {
			e.onUp()
		}
	case fyne.KeyDown:
		if e.onDown != nil {
			e.onDown()
		}
	case fyne.KeyReturn:
		if e.onEnter != nil {
			e.onEnter()
		}
	default:
		e.Entry.TypedKey(k)
	}
}

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(&customTheme{})
	myWindow := myApp.NewWindow("Script Launcher")
	myWindow.Resize(fyne.NewSize(500, 600))

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("홈 디렉토리를 찾을 수 없습니다.")
		return
	}

	baseDir := home + "/system"
	var allFiles, listData []string
	files, _ := filepath.Glob(filepath.Join(baseDir, "*.sh"))
	for _, f := range files {
		allFiles = append(allFiles, filepath.Base(f))
	}
	listData = allFiles

	var selectedIdx int = 0
	list := widget.NewList(
		func() int { return len(listData) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(listData[id])
		},
	)
	list.OnSelected = func(id widget.ListItemID) { selectedIdx = id }

	input := newSmartEntry()
	input.SetPlaceHolder("Search & Run (Enter)...")

	// --- 핵심 실행 함수 (AppleScript 활용) ---
	run := func() {
		if selectedIdx >= 0 && selectedIdx < len(listData) {
			scriptName := listData[selectedIdx]
			scriptPath := filepath.Join(baseDir, scriptName)

			// AppleScript 로직:
			// 1. 터미널 활성화
			// 2. 창이 없으면 새로 생성, 있으면 새 탭(Cmd+T) 생성
			// 3. 해당 탭의 테마를 'Homebrew'(초록색)로 변경하고 제목 설정
			osa := fmt.Sprintf(`
				tell application "Terminal"
					activate
					if (count of windows) is 0 then
						do script "%s"
					else
						tell application "System Events" to tell process "Terminal" to keystroke "t" using command down
						delay 0.1
						do script "%s" in front window
					end if
					-- 테마 및 제목 설정
					set custom title of tab (count of tabs of front window) of front window to "%s"
				end tell`, scriptPath, scriptPath, scriptName)

			cmd := exec.Command("osascript", "-e", osa)
			cmd.Dir = baseDir

			// 실행 및 에러 처리 (기존 로직 유지)
			err := cmd.Run()
			if err != nil {
				dialog.ShowError(err, myWindow)
			}

			// 후처리: 포커스 유지
			myWindow.Canvas().Focus(input)
			input.SetText("") // 입력창 내용을 ""(빈 칸)으로 변경
		}
	}

	// --- 이벤트 바인딩 ---
	input.onUp = func() {
		if selectedIdx > 0 {
			selectedIdx--
			list.Select(selectedIdx)
			list.ScrollTo(selectedIdx)
		}
	}
	input.onDown = func() {
		if selectedIdx < len(listData)-1 {
			selectedIdx++
			list.Select(selectedIdx)
			list.ScrollTo(selectedIdx)
		}
	}
	input.onEnter = run

	input.OnChanged = func(s string) {
		listData = nil
		search := strings.ToLower(s)
		for _, name := range allFiles {
			if strings.Contains(strings.ToLower(name), search) {
				listData = append(listData, name)
			}
		}
		list.Refresh()
		if len(listData) > 0 {
			selectedIdx = 0
			list.Select(0)
		} else {
			selectedIdx = -1
			list.UnselectAll()
		}
	}

	// 초기 설정
	list.Select(0)
	myWindow.SetContent(container.NewBorder(input, nil, nil, nil, list))
	myWindow.Canvas().Focus(input)
	myWindow.ShowAndRun()
}
