package ybtools

//
// Yapperbot Tools, the internal system bits for Yapperbot and co.
// Copyright (C) 2020 Naypta

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//

import (
	"fmt"
	"regexp"
)

// botName is used to match Template:Bots and the like.
var botName string

var botBanRegex *regexp.Regexp

const botBanRegexTemplate string = `{{nobots}}|{{bots\|deny=(?:[^,|}]*,)*%[1]s|{{bots\|allow=(?![^}|]*%[1]s[},|])`

// BotAllowed take a page content and determines if the botName given is allowed
// to edit the page per the applicable templates.
func BotAllowed(pageContent string) bool {
	if botName == "" {
		panic("BotAllowed called with no botName set!")
	}
	return !botBanRegex.MatchString(pageContent)
}

// SetupBot sets the bot name, ready for future calls to BotAllowed.
func SetupBot(bn string) {
	botName = bn
	botBanRegex = regexp.MustCompile(fmt.Sprintf(botRegexTemplate, bn))
}
