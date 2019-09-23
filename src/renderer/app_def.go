package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/rocketlaunchr/react"
	"github.com/rocketlaunchr/react/elements"
)

var (
	npm          = js.Global.Get("npm")
	Markdown     = npm.Get("Markdown")
	AppComponent *js.Object
)

type AppProps struct {
	DefaultMarkdown string `react:"default_markdown"`
}

type AppState struct {
	Markdown string `react:"markdown"`
}

func init() {

	appDef := react.NewClassDef("App")

	appDef.GetInitialState(func(this *js.Object, props react.Map) interface{} {

		var appProps AppProps
		react.UnmarshalProps(this, &appProps)

		return AppState{
			Markdown: appProps.DefaultMarkdown,
		}
	})

	appDef.SetEventHandler("change", func(this *js.Object, e *react.SyntheticEvent, props, state react.Map, setState react.SetState) {

		newValue := e.Target().Get("value").String()

		setState(AppState{
			Markdown: newValue,
		})
	})

	appDef.Render(func(this *js.Object, props, state react.Map) interface{} {

		var appState AppState
		react.UnmarshalState(this, &appState)

		return elements.Div(&elements.DivProps{Class: "app"},
			elements.TextArea(&elements.TextAreaProps{Class: "editor", DefaultValue: appState.Markdown, OnChange: this.Get("change")}),
			react.JSX(Markdown, js.M{"className": "preview", "source": appState.Markdown, "escapeHtml": false}),
		)
	})

	AppComponent = react.CreateClass(appDef)
}
