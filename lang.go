package main

import ui18n "github.com/unknwon/i18n"

func trUser(user *User, message string) string {
	if len(user.Language) > 0 {
		return ui18n.Tr(user.Language, message)
	} else {
		return ui18n.Tr(lang, message)
	}
}

func trGroup(groupID int, message string) string {
	var lng string
	if groupID == tAnonBalkan {
		lng = langHr
	} else {
		lng = lang
	}
	return ui18n.Tr(lng, message)
}
