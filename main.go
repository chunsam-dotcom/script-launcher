package main

import (
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
		return
	}

	baseDir := filepath.Join(home, "system")
	var allFiles, listData []string

	entries, _ := os.ReadDir(baseDir)
	for _, entry := range entries {
		if !entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			allFiles = append(allFiles, entry.Name())
		}
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

	run := func() {
		if selectedIdx < 0 || selectedIdx >= len(listData) {
			return
		}
		scriptName := listData[selectedIdx]
		scriptPath := filepath.Join(baseDir, scriptName)

		var cmd *exec.Cmd

		if filepath.Ext(scriptName) == "" {
			cmd = exec.Command(scriptPath)
			cmd.Dir = baseDir
			if err := cmd.Start(); err != nil {
				dialog.ShowError(err, myWindow)
			}
		} else {
			switch runtime.GOOS {
			case "darwin":
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
						set custom title of tab (count of tabs of front window) of front window to "%s"
					end tell`, scriptPath, scriptPath, scriptName)
				cmd = exec.Command("osascript", "-e", osa)
			case "linux":
				cmd = exec.Command("x-terminal-emulator", "-e", "bash", "-c", fmt.Sprintf("bash %s; exec bash", scriptPath))
			default:
				return
			}
			cmd.Dir = baseDir
			if err := cmd.Start(); err != nil {
				dialog.ShowError(err, myWindow)
			}
		}

		myWindow.Canvas().Focus(input)
		input.SetText("")
	}

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

	list.Select(0)
	myWindow.SetContent(container.NewBorder(input, nil, nil, nil, list))
	myWindow.Canvas().Focus(input)
	myWindow.ShowAndRun()
}
