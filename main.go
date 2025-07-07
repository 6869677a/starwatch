package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/b0ch3nski/go-starlink/starlink"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const version = "v1.0.0 alpha"

var (
	logToCSV    = false
	showVersion = false
	csvWriter   *csv.Writer
	csvFile     *os.File
	eventLog    = []string{}
)

func formatDuration(seconds float32) string {
	if seconds >= 3600 {
		h := int(seconds) / 3600
		m := (int(seconds) % 3600) / 60
		return fmt.Sprintf("%dh %dm", h, m)
	} else if seconds >= 60 {
		m := int(seconds) / 60
		s := int(seconds) % 60
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%.0f s", seconds)
}

func main() {
	flag.BoolVar(&logToCSV, "log", false, "Enable CSV logging to starwatch-log.csv")
	flag.BoolVar(&showVersion, "version", false, "Show StarWatch version and exit")
	flag.Parse()

	if showVersion {
		fmt.Printf("StarWatch – Shellz Code Solutions by w0rmer\nVersion: %s\n", version)
		return
	}

	if logToCSV {
		var err error
		csvFile, err = os.Create("starwatch-log.csv")
		if err != nil {
			fmt.Println("Error creating log file:", err)
			return
		}
		defer csvFile.Close()
		csvWriter = csv.NewWriter(csvFile)
		defer csvWriter.Flush()
		csvWriter.Write([]string{
			"Timestamp", "Uptime(s)", "POP_Latency(ms)", "Downlink(kbps)", "Uplink(kbps)",
			"Obstruct(%)", "Tilt", "Rotation",
		})
	}

	app := tview.NewApplication()
	title := fmt.Sprintf("🌌 StarWatch %s by w0rmer", version)

	mainView := tview.NewTextView().SetDynamicColors(true).SetWrap(false).SetScrollable(true)
	mainView.SetBorder(true).SetTitle(title)

	logView := tview.NewTextView().SetDynamicColors(true).SetScrollable(true)
	logView.SetBorder(true).SetTitle("📜 Event Log")

	// FLEX must be declared before helpModal can reference it
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(mainView, 0, 4, false).
		AddItem(logView, 0, 1, false)

	// HELP MODAL
	helpText := `[green]Legend & Help:
[white]Tilt: Dish vertical aim (elevation)
Rotation: Dish horizontal aim (azimuth)
Signal Quality: Above/Below usable signal threshold
Obstruct: Percent of time signal is blocked

[gray]Keys:
Q = Quit
L = Toggle CSV Logging
H = Toggle Help`
	helpModal := tview.NewModal().
		SetText(helpText).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(flex, true)
		})

	// WELCOME POPUP
	welcome := tview.NewModal().
		SetText("Welcome to StarWatch by w0rmer!\nPress 'H' at any time for help.").
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(flex, true)
		})


	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q', 'Q':
			app.Stop()
		case 'l', 'L':
			logToCSV = !logToCSV
			appendLog(logView, fmt.Sprintf("[yellow]CSV logging %v", logToCSV))
		case 'h', 'H':
			app.SetRoot(helpModal, false).SetFocus(helpModal)
		}
		return event
	})

	go func() {
		ctx := context.Background()
		client, err := starlink.NewClient(ctx, starlink.DefaultDishyAddr)
		if err != nil {
			panic(err)
		}

		for {
			status, err := client.Status(ctx)
			if err != nil {
				app.QueueUpdateDraw(func() {
					mainView.SetText("[red]DISH UNREACHABLE")
					appendLog(logView, "[red]Dish appears to be offline!")
				})
				time.Sleep(5 * time.Second)
				continue
			}

			uptime := float64(status.DeviceState.GetUptimeS()) / 3600
			dishStatus := "OFFLINE"
			if status.DeviceState.GetUptimeS() > 0 {
				dishStatus = "ONLINE"
			}

			signalStatus := "[red]⚠️ Below Noise Floor"
			if status.GetIsSnrAboveNoiseFloor() {
				signalStatus = "[green]✅ Above Noise Floor"
			}

			currentTilt := status.AlignmentStats.GetBoresightElevationDeg()
			targetTilt := status.AlignmentStats.GetDesiredBoresightElevationDeg()
			tiltDiff := float64(targetTilt - currentTilt)
			tiltDir := "↑"
			if tiltDiff < 0 {
				tiltDir = "↓"
			}
			tiltMag := abs(tiltDiff)

			currentRot := status.AlignmentStats.GetBoresightAzimuthDeg()
			targetRot := status.AlignmentStats.GetDesiredBoresightAzimuthDeg()
			rotDiff := float64(targetRot - currentRot)
			rotDir := "↻ (R)"
			if rotDiff < 0 {
				rotDir = "↺ (L)"
			}
			rotMag := abs(rotDiff)

			alignmentStatus := "[red]NEEDS ADJUSTMENT"
			if tiltMag < 2 && rotMag < 2 {
				alignmentStatus = "[green]OKAY"
			} else if tiltMag < 5 && rotMag < 5 {
				alignmentStatus = "[yellow]MARGINAL"
			}

			obsPercent := status.ObstructionStats.GetFractionObstructed() * 100
			alert := ""
			if obsPercent > 3.0 {
				alert = fmt.Sprintf("[red]⚠ HIGH OBSTRUCTION: %.2f%%", obsPercent)
				appendLog(logView, alert)
			}

			now := time.Now().Format("15:04:05 MST 01/02")
			output := fmt.Sprintf(`%s
[yellow]Dish State: [white]%s   [yellow]Uptime: [white]%.2f hrs   [yellow]Boot Count: [white]%d   [yellow]Signal Quality: [white]%s

[green]POP Latency : [white]%.2f ms
[green]Downlink    : [white]%.2f kbps
[green]Uplink      : [white]%.2f kbps
[green]Obstruct %%   : [white]%.2f%%
[green]Obstruction Data Time: [white]%.0f s   [green]Samples: [white]%d
[green]Recent Obstruction Duration: [white]%s   [green]Avg Interval Between Obstructions: [white]%s

[cyan]Alignment     : %s
[cyan]Tilt          : [white]%.2f° → %.2f° (%.2f°%s)
[cyan]Rotation      : [white]%.2f° → %.2f° (%.2f°%s)
[cyan]Uncertainty   : [white]%.2f°

[blue]GPS Valid     : [white]%v   [blue]Sats: [white]%d
[blue]Connected Router : [white]%s

[gray]HW: %s | SW: %s | Build ID: %.10s
[gray]Logging: %v | Time: %s
[darkgray]Q=Quit  L=Toggle Logging  H=Help`,
				alert,
				dishStatus, uptime, status.DeviceInfo.GetBootcount(), signalStatus,
				status.GetPopPingLatencyMs(), status.GetDownlinkThroughputBps()/1000,
				status.GetUplinkThroughputBps()/1000, obsPercent,
				status.ObstructionStats.GetValidS(), status.ObstructionStats.GetPatchesValid(),
				formatDuration(status.ObstructionStats.GetTimeObstructed()),
				formatDuration(status.ObstructionStats.GetAvgProlongedObstructionIntervalS()),
				alignmentStatus, currentTilt, targetTilt, tiltMag, tiltDir,
				currentRot, targetRot, rotMag, rotDir,
				status.AlignmentStats.GetAttitudeUncertaintyDeg(),
				status.GpsStats.GetGpsValid(), status.GpsStats.GetGpsSats(),
				status.GetConnectedRouters(),
				status.DeviceInfo.GetHardwareVersion(),
				status.DeviceInfo.GetSoftwareVersion(),
				status.DeviceInfo.GetBuildId(),
				logToCSV, now,
			)

			app.QueueUpdateDraw(func() {
				mainView.SetText(output)
			})

			if logToCSV && csvWriter != nil {
				_ = csvWriter.Write([]string{
					time.Now().Format(time.RFC3339),
					strconv.Itoa(int(status.DeviceState.GetUptimeS())),
					fmt.Sprintf("%.2f", status.GetPopPingLatencyMs()),
					fmt.Sprintf("%.2f", status.GetDownlinkThroughputBps()/1000),
					fmt.Sprintf("%.2f", status.GetUplinkThroughputBps()/1000),
					fmt.Sprintf("%.2f", obsPercent),
					fmt.Sprintf("%.2f", currentTilt),
					fmt.Sprintf("%.2f", currentRot),
				})
				csvWriter.Flush()
			}

			time.Sleep(5 * time.Second)
		}
	}()

	if err := app.SetRoot(welcome, true).Run(); err != nil {
		panic(err)
	}
}

func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

func appendLog(view *tview.TextView, msg string) {
	timestamp := time.Now().Format("15:04:05")
	logLine := fmt.Sprintf("[gray]%s [white]%s", timestamp, msg)
	eventLog = append(eventLog, logLine)
	if len(eventLog) > 100 {
		eventLog = eventLog[1:]
	}
	view.Clear()
	for _, line := range eventLog {
		fmt.Fprintln(view, line)
	}
	view.ScrollToEnd()
}
