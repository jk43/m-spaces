package main

import (
	"hook/hook"
)

func (app *application) setHookHandlers() {
	hook := hook.HookHandler{}

	app.HookHandlers.AddHook("LocalhostPostPutUser", hook.LocalhostPostPutUser)
	app.HookHandlers.AddHook("LocalhostPrePutUser", hook.LocalhostPrePutUser)
}
