package feedback

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/fynelabs/telemetry"
)

type data struct {
	Feeling telemetry.Feeling
	Detail  string
}

func ShowFeedback(t *telemetry.Telemetry, w fyne.Window) {
	var send data
	dialog.ShowCustomConfirm("Send Feedback", "Send", "Cancel", makeUI(
		func(d data) {
			send = d
		},
	), func(ok bool) {
		if !ok {
			return
		}

		t.Feedback(send.Feeling, send.Detail)
	}, w)
}

func ShowCustomFeedback(t *telemetry.Telemetry, cb func(d *data), w fyne.Window) {
}

func makeUI(cb func(data)) fyne.CanvasObject {
	var send data
	comments := widget.NewMultiLineEntry()
	comments.OnChanged = func(s string) {
		send.Detail = s
		cb(send)
	}

	var emote func(f telemetry.Feeling)
	great := widget.NewButton("😁", func() {
		emote(telemetry.Excited)
	})
	happy := widget.NewButton("🙂", func() {
		emote(telemetry.Happy)
	})
	sad := widget.NewButton("😟", func() {
		emote(telemetry.Sad)
	})
	upset := widget.NewButton("😡", func() {
		emote(telemetry.Frustrated)
	})

	emote = func(f telemetry.Feeling) {
		send.Feeling = f
		for i, b := range []*widget.Button{great, happy, sad, upset} {
			on := (i == int(f)+1) || i == 0 && f == telemetry.Excited

			if on {
				b.Importance = widget.HighImportance
			} else {
				b.Importance = widget.MediumImportance
			}

			b.Refresh()

			cb(send)
		}
	}

	buttons := container.NewThemeOverride(
		container.NewGridWithColumns(4, great, happy, sad, upset), &emojiButton{theme.Current()})

	return container.NewVBox(
		buttons,
		widget.NewLabel("Optional comments:"),
		comments)
}
