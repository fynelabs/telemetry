package telemetry

import "github.com/google/uuid"

type ServerTelemetry struct {
	*Telemetry
}

// InitServer opens a new telemetry instance and logs the start of a new server session.
// The app ID refers to any instance of the app.
// The `accessCode“ is the developer code for accessing Fyne Labs telemetry service.
func InitServer(appID, accessCode string) *ServerTelemetry {
	session := uuid.New().String() // just a random new session for the "server" context
	return &ServerTelemetry{Telemetry: initTelemetry(appID, "", session, accessCode, false)}
}

// ClientError reports an error to the telemetry server for the specified user session.
// It will generate a stack trace starting at the function that called this method.
// The session should have been started using `ClientSessionStart`.
func (t *Telemetry) UserError(err error, session string) {
	t.sendError(err, session)
}

// ClientEvent logs a named event to the telemetry server associated with a client session.
// Event names should be unique to your application for correct counting.
// The session should have been started using `ClientSessionStart`.
func (t *ServerTelemetry) ClientEvent(name, session string) {
	t.send("event?name=%s&session=%s", name, session)
}

// ClientUserInfo allows an app to provide a username and/or email to associate with a user.
// This data will be connected to all sessions for the specified user.
// The userID should have been connected to a session using `ClientSessionStart`.
func (t *ServerTelemetry) ClientUserInfo(id, username, email string) {
	t.send("user?uuid=%s&username=%s&email=%s", id, username, email)
}

// ClientSessionEnd will mark a client session as ended, where possible.
// The id should belong to a session opened using `ClientSessionStart`
func (t *ServerTelemetry) ClientSessionEnd(id string) {
	t.sendWait("sessionend?uuid=%s", id)
}

// ClientSessionStart starts a new session for a specific client of a server.
// The id parameter is a globally unique ID for this session, and the user ID should be
// globally unique and re-used accross sessions for that user.
func (t *ServerTelemetry) ClientSessionStart(id, user string) {
	t.send("session?uuid=%s&appID=%s&user=%s&device=web", id, t.AppID, user)
}
