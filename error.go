package telemetry

import (
	"errors"
	"fmt"
	"net/url"
	"runtime"
	"runtime/debug"

	"fyne.io/fyne/v2"
)

// Error reports an error to the telemetry server.
// It will generate a stack trace starting at the function that called this method.
func (t *Telemetry) Error(err error) {
	t.sendError(err, t.session)
}

// Run wraps the standard fyne App.Run() with a crash reporting wrapper.
// If a panic occurs it will be reported as an error event to the telemetry server.
func (t *Telemetry) Run(a fyne.App) {
	defer func() {
		if r := recover(); r != nil {
			err := errors.New(fmt.Sprintf("%v", r))
			fyne.LogError("Handling panic", err)

			t.Error(err)
			debug.PrintStack()
		}
	}()

	a.Run()
}

// ShowAndRun wraps the standard fyne Window.ShowAndRun() with a crash reporting wrapper.
// If a panic occurs it will be reported as an error event to the telemetry server.
func (t *Telemetry) ShowAndRun(w fyne.Window, a fyne.App) {
	w.Show()
	t.Run(a)
}

func (t *Telemetry) sendError(err error, session string) {
	log := err.Error()

	stack := ""
	for i := 0; ; i++ {
		_, file, line, ok := runtime.Caller(i + 1)
		if !ok {
			break
		}

		stack += fmt.Sprintf("  %s:%d\n", file, line)
	}

	t.send("error?detail=%s&stack=%s&session=%s", url.QueryEscape(log), url.QueryEscape(stack), t.session)
}
