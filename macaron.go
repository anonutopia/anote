package main

import (
	"github.com/cryptopragmatic/certmagic"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/i18n"
	_ "github.com/go-macaron/session/redis"
	macaron "gopkg.in/macaron.v1"
)

func initMacaron() *macaron.Macaron {
	m := macaron.Classic()

	m.Use(cache.Cacher())
	m.Use(macaron.Renderer())

	m.Use(i18n.I18n(i18n.Options{
		Langs: []string{"hr", "sr", "en-US"},
		Names: []string{"Hrvatski", "Srpski", "English"},
	}))

	if !conf.Dev {
		// certmagic.Default.Agreed = true
		certmagic.DefaultACME.Email = conf.EmailAddress
		go certmagic.HTTPS([]string{conf.Hostname}, m)
	} else {
		go m.Run("0.0.0.0", 5000)
	}

	return m
}
