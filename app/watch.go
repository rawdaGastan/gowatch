package app

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
	"unicode"

	"github.com/rivo/tview"
)

// TODO: handle errors (panics).
type App struct {
	Interval float64

	Diff    bool
	Exec    bool
	NoTitle bool
	ChgExit bool
	ErrExit bool
	Beep    bool

	Cmd  string
	Args []string

	UpdateCmd  string
	UpdateArgs []string

	view *tview.Application
}

func (a *App) output() string {
	var b bytes.Buffer
	var cmd *exec.Cmd

	if a.Exec {
		cmd = exec.Command(a.Cmd, a.Args...)
	} else {
		cmd = exec.Command("bash", "-c", strings.Join(append([]string{a.Cmd}, a.Args...), " "))
	}

	cmd.Stderr = &b
	cmd.Stdout = &b

	err := cmd.Run()
	if cmd.ProcessState.ExitCode() != 0 {
		// if a.Beep {

		// }

		if a.ErrExit {
			a.view.Stop()
		}
	}

	if err != nil {
		return err.Error()
	}

	return b.String()
}

func (a *App) Run(ctx context.Context) {
	a.view = tview.NewApplication()

	currentResult := a.output()

	go func() {
		for {
			result := a.output()
			if currentResult != result {
				if len(strings.TrimSpace(a.UpdateCmd)) > 0 {
					var cmd *exec.Cmd
					if a.Exec {
						cmd = exec.Command(a.UpdateCmd, a.UpdateArgs...)
					} else {
						cmd = exec.Command("bash", "-c", strings.Join(append([]string{a.UpdateCmd}, a.UpdateArgs...), " "))
					}
					err := cmd.Run()
					if err != nil {
						panic(err)
					}
				}

				if a.ChgExit {
					a.view.Stop()
				}
			}

			a.view.SetRoot(a.snapshot(currentResult, result), true).Draw()
			currentResult = result

			select {
			case <-ctx.Done():
				panic(ctx.Err())
			case <-time.After(time.Duration(a.Interval * float64(time.Second))):
			}
		}
	}()

	snapshot := a.snapshot(currentResult, currentResult)
	if err := a.view.SetRoot(snapshot, true).Run(); err != nil {
		panic(err)
	}
}

func (a *App) snapshot(old, data string) *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	if !a.NoTitle {
		flex.AddItem(a.header(), 0, 1, false)
	}
	flex.AddItem(a.body(old, data), 0, 3, false)
	return flex
}

func (a *App) body(old, data string) *tview.Flex {
	updatedData := a.highlight(old, data)

	body := tview.NewTextView().SetDynamicColors(true)
	body.SetText(updatedData)

	snapshot := tview.NewFlex().SetDirection(tview.FlexRow)
	snapshot.AddItem(body, 0, 1, false)

	return snapshot
}

func (a *App) highlight(oldData, newData string) string {
	if a.Diff {
		oldLines := strings.Split(oldData, "\n")
		newLines := strings.Split(newData, "\n")
		var updatedLines []string

		// make new and old have equal lines
		if len(oldLines) > len(newLines) {
			// TODO: the loop failed with mince condition
			diffLines := len(oldLines) - len(newLines)
			for i := 0; i < diffLines; i++ {
				newLines = append(newLines, "")
			}
		} else {
			diffLines := len(newLines) - len(oldLines)
			for i := 0; i < diffLines; i++ {
				oldLines = append(oldLines, "")
			}
		}

		for x, newLine := range newLines {
			oldLine := oldLines[x]
			updatedLine := a.diffLine(oldLine, newLine)
			updatedLines = append(updatedLines, updatedLine)
		}

		return strings.Join(updatedLines, "\n")
	}

	return newData
}

func (a *App) diffLine(oldLine, newLine string) string {
	updatedLine := newLine
	var updatedI int

	// TODO: range loop diff results
	for i := 0; i < len(newLine); i++ {
		if unicode.IsSpace(rune(newLine[i])) {
			updatedI++
			continue
		}

		if (i < len(oldLine) && oldLine[i] != newLine[i]) ||
			i >= len(oldLine) {
			updatedLine = fmt.Sprintf("%s[black:white]%s[white:black]%s", updatedLine[:updatedI], string(updatedLine[updatedI]), updatedLine[updatedI+1:])
			updatedI += 26
		}
		updatedI++
	}

	if len(oldLine) > len(newLine) {
		for i := 0; i < len(oldLine)-len(newLine); i++ {
			if !unicode.IsSpace(rune(oldLine[i])) {
				updatedLine += "[black:white] [white:black]"
				continue
			}
			updatedLine += string(oldLine[i])
		}
	}

	return updatedLine
}

func (a *App) header() *tview.Flex {
	every := tview.NewTextView()
	every.SetBorder(true).SetTitle("Every")
	every.SetTitleAlign(tview.AlignTop)
	every.SetText(fmt.Sprintf("%v s", a.Interval))

	cmd := tview.NewTextView()
	cmd.SetBorder(true).SetTitle("Command")
	cmd.SetTitleAlign(tview.AlignTop)
	cmd.SetText(fmt.Sprintf("%s %s", a.Cmd, strings.Join(a.Args, " ")))

	update := tview.NewTextView()
	update.SetBorder(true).SetTitle("Update")
	update.SetTitleAlign(tview.AlignTop)
	update.SetText(fmt.Sprintf("%s %s", a.UpdateCmd, strings.Join(a.UpdateArgs, " ")))

	timer := tview.NewTextView()
	timer.SetBorder(true).SetTitle("Time")
	timer.SetTitleAlign(tview.AlignTop)
	timer.SetText(time.Now().Format(time.ANSIC))

	header := tview.NewFlex().SetDirection(tview.FlexColumn)
	header.AddItem(every, 0, 1, false)
	header.AddItem(update, 0, 1, false)
	header.AddItem(cmd, 0, 1, false)
	header.AddItem(timer, 0, 1, false)

	return header
}
