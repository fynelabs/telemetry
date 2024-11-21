package telemetry

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"github.com/google/uuid"
)

type Telemetry struct {
	AccessCode string
	AppID      string

	client        *http.Client
	server        string
	session, user string
}

const prefUserKey = "fynelabs.telemetry.user"

// Init opens a new telemetry instance and logs the start of a new session.
// It uses the Fyne app to get the App ID and handle user uniqueness.
// The `accessCode“ is the developer code for accessing Fyne Labs telemetry service.
func Init(a fyne.App, accessCode string) *Telemetry {
	id := a.UniqueID()

	rand.Seed(time.Now().Unix())
	session := uuid.New()
	user := a.Preferences().String(prefUserKey)
	if user == "" {
		user = uuid.New().String()
		a.Preferences().SetString(prefUserKey, user)
	}

	return InitWithID(id, user, session.String(), accessCode)
}

// InitWithID opens a new telemetry instance and logs the start of a new session.
// The user of this must pass in a unique ID for the app, session and user.
// The app ID refers to any instance of the app, the user ID should be consistent across launches
// and the session should be unique for every invocation.
// The `accessCode“ is the developer code for accessing Fyne Labs telemetry service.
func InitWithID(appID, user, session, accessCode string) *Telemetry {
	t := initTelemetry(appID, user, session, accessCode, true)
	t.user = user

	return t
}

func initTelemetry(appID, user, session, accessCode string, native bool) *Telemetry {
	t := &Telemetry{AccessCode: accessCode, AppID: appID, user: user,
		server: "https://xavier.fynelabs.com", client: &http.Client{}}

	if env := os.Getenv("TELEMETRY_SERVER"); env != "" {
		t.server = env
	}

	t.sessionStart(session, native)
	return t
}

// Close is used to shut down the telemetry instance, it will log that the session is ended.
// This should be called at the end of an app's `main()` function.
func (t *Telemetry) Close() {
	t.sessionEnd()
}

// Error reports an error to the telemetry server.
// It will generate a stack trace starting at the function that called this method.
func (t *Telemetry) Error(err error) {
	t.sendError(err, t.session)
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

	encoded := url.QueryEscape(stack)
	t.send("error?detail=%s&stack=%s&session=%s", log, encoded, t.session)
}

// Event logs a named event to the telemetry server.
// Event names should be unique to your application for correct counting.
func (t *Telemetry) Event(name string) {
	t.send("event?name=%s&session=%s", name, t.session)
}

// UserInfo allows an app to provide a username and/or email to associate with a user.
// This data will be connected to all sessions for the current user.
func (t *Telemetry) UserInfo(username, email string) {
	t.send("user?uuid=%s&username=%s&email=%s", t.user, username, email)
}

func (t *Telemetry) send(path string, params ...any) {
	go t.sendWait(path, params...)
}

func (t *Telemetry) sendWait(path string, params ...any) {
	url := fmt.Sprintf(t.server+"/api/v1/"+path, params...)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("AccessCode", t.AccessCode)
	r, err := t.client.Do(req)

	if err != nil {
		log.Println("Failed to send telemetry", err)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Body read error", err)
		return
	}
	r.Body.Close()
	if len(data) > 0 {
		log.Println("Body returned:", string(data))
	}
}

func (t *Telemetry) sessionEnd() {
	t.sendWait("sessionend?uuid=%s", t.session)
}

func (t *Telemetry) sessionStart(id string, native bool) {
	t.session = id

	device := ""
	if native {
		device = fmt.Sprintf("os=%s&arch=%s", runtime.GOOS, runtime.GOARCH)
	} else {
		device = "device=server"
	}

	t.send("session?uuid=%s&appID=%s&user=%s&%s", id, t.AppID, t.user, device)
}
