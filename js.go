// +build js,wasm

package main

import (
	"syscall/js"
)

func setupJavascriptCallback(game *Game) {
	toggleDebug := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 0 {
			return "Invalid no of arguments passed"
		}
		if game == nil {
			return "Game isn't running"
		}
		game.debug = !game.debug
		return game.debug
	})
	js.Global().Set("toggleDebug", toggleDebug)
}
