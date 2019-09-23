package main

import (
	"runtime"

	"github.com/gopherjs/gopherjs/js"
)

var (
	npm           = js.Global.Get("npm")
	Tray          = npm.Get("Tray")
	app           = npm.Get("app")
	BrowserWindow = npm.Get("BrowserWindow")

	Alert   = npm.Get("Alert")
	process = npm.Get("process")
)

func main() {

	// Don't show Dock icon for mac
	if runtime.GOOS == "darwin" {
		app.Get("dock").Call("hide")
	}

	// Prevent second instance of application
	if !app.Call("requestSingleInstanceLock").Bool() {
		// Second instance
		app.Call("exit", 0)
	}

	app.Call("on", "ready", func() {

		// Provide prettier error box
		process.Call("on", "uncaughtException", Alert.Call("uncaughtException", false, func(err *js.Object) {
			app.Call("exit", 1)
		}))

		// Create tray icon
		createTray()
	})

	// Prevent closing main window from quitting applicatin
	app.Call("on", "window-all-closed", func() {})
}

var tray *js.Object

func createTray() {
	tray = Tray.New("./assets/tray.png")

	tray.Call("on", "click", func() {
		openWindow()
	})
}

var win *js.Object

func openWindow() {

	if win == nil {
		opts := js.M{
			"width":       800,
			"height":      600,
			"title":       "Visual Markdown",
			"show":        false,
			"skipTaskbar": true,
			"minimizable": false,
			"webPreferences": js.M{
				"nodeIntegration": true,
				"devTools":        false,
			},
		}
		win = BrowserWindow.New(opts)
		win.Call("removeMenu")
		win.Call("loadFile", "renderer.html")
		win.Call("once", "ready-to-show", func() {
			win.Call("show")
		})
		win.Call("on", "close", func(e *js.Object) {
			// Hide window instead of closing
			e.Call("preventDefault") // Cancel closing of window
			win.Call("hide")
		})
	} else {
		if !win.Call("isVisible").Bool() {
			win.Call("show")
		} else {
			if win.Call("isFocused").Bool() {
				win.Call("hide")
			} else {
				win.Call("show")
			}
		}
	}
}
