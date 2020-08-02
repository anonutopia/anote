package main

import ui18n "github.com/unknwon/i18n"

func trUser(user *User, message string) string {
	return ui18n.Tr(user.Language, message)
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
