package telemetry

import (
	"fmt"
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

	server        string
	session, user string
}

const prefUserKey = "fynelabs.telemetry.user"

func Init(a fyne.App, accessCode string) *Telemetry {
	id := a.UniqueID()

	rand.Seed(time.Now().Unix())
	session := uuid.New()
	user := a.Preferences().String(prefUserKey)
	if user == "" {
		user = uuid.New().String()
		a.Preferences().SetString(prefUserKey, user)
	}

	return InitWithID(id, session.String(), user, accessCode)
}

func InitWithID(appID, session, user, accessCode string) *Telemetry {
	t := &Telemetry{AccessCode: accessCode, AppID: appID,
		session: session, user: user, server: "https://xavier.fynelabs.com"}

	if env := os.Getenv("TELEMETRY_SERVER"); env != "" {
		t.server = env
	}
	http.Get(fmt.Sprintf(t.server+"/api/v1/session?uuid=%s&appID=%s&user=%s",
		session, appID, user))

	return t
}

func (t *Telemetry) Error(err error) {
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
	http.Get(fmt.Sprintf(t.server+"/api/v1/error?detail=%s&stack=%s&session=%s",
		log, encoded, t.session))
}

func (t *Telemetry) Event(id string) {
	http.Get(fmt.Sprintf(t.server+"/api/v1/event?name=%s&session=%s",
		id, t.session))
}

func (t *Telemetry) UserInfo(username, email string) {
	http.Get(fmt.Sprintf(t.server+"/api/v1/user?uuid=%s&username=%s&email=%s",
		t.user, username, email))
}
