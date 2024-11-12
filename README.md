# Telemetry

A simple library for reporting data to the Fyne labs telemetry server.

Telemetry will track user sessions and aim to determine if it is the same
user across new invocations of your app.

## Usage

This is designed to be dropped into a [Fyne](https://fyne.io) app simply:

```go
	a := app.NewWithID("com.example.myapp")
	t := telemetry.Init(a, "ACCESSCODE")
```

And then you can report telemetry events simply as:

```go
    t.Event("eventname")
```

Or report errors

```go
    t.Error(err)
```

The library will work out the stack trace for your error and upload that too!
