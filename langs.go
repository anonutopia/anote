package main

import "github.com/bykovme/gotrans"

func initLangs() {
	gotrans.InitLocales("langs")
	gotrans.SetDefaultLocale("en")
}
