package main

import (
	"github.com/caddyserver/certmagic"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/captcha"
	macaron "gopkg.in/macaron.v1"
)

func initMacaron() *macaron.Macaron {
	m := macaron.Classic()

	m.Use(cache.Cacher())
	m.Use(macaron.Renderer())
	m.Use(captcha.Captchaer())

	if !conf.Dev {
		certmagic.DefaultACME.Email = conf.EmailAddress
		certmagic.DefaultACME.Agreed = true
		go certmagic.HTTPS([]string{conf.Hostname}, m)
	} else {
		go m.Run("0.0.0.0", 5000)
	}

	m.Get("/mine/:code", mineView)
	m.Post("/mine/:code", mineViewPost)
	m.Get("/withdraw/:code", withdrawView)
	m.Post("/withdraw/:code", withdrawViewPost)

	return m
}
