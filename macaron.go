package main

import (
	"github.com/go-macaron/cache"
	"github.com/go-macaron/captcha"
	macaron "gopkg.in/macaron.v1"
)

func initMacaron() *macaron.Macaron {
	m := macaron.Classic()

	m.Use(cache.Cacher())
	m.Use(macaron.Renderer())
	m.Use(captcha.Captchaer())

	m.Get("/mine/:code", mineView)
	m.Post("/mine/:code", mineViewPost)
	m.Get("/withdraw/:code", withdrawView)
	m.Post("/withdraw/:code", withdrawViewPost)

	go m.Run("0.0.0.0", Port)

	return m
}
