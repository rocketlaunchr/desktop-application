package main

import (
	"github.com/rocketlaunchr/react"
)

var markdown = `
# Visual Markdown
`

func main() {
	domTarget := react.GetElementByID("app")
	react.Render(react.JSX(AppComponent, &AppProps{DefaultMarkdown: markdown}), domTarget)
}
