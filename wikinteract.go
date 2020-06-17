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
	"cgt.name/pkg/go-mwclient"
	"cgt.name/pkg/go-mwclient/params"
)

var w *mwclient.Client

// CreateAndAuthenticateClient uses the details already passed into ybtools
// in setup.go to return a fully-authenticated mwclient
func CreateAndAuthenticateClient() *mwclient.Client {
	if taskName == "" || botUser == "" {
		PanicErr("Call ybtools.SetupBot first!")
	}

	var err error

	w, err = mwclient.New(config.APIEndpoint, "Yapperbot-"+taskName+" on User:"+botUser+" - Golang, licensed GNU GPL")
	if err != nil {
		PanicErr("Failed to create MediaWiki client with error ", err)
	}

	err = w.Login(config.BotUsername, botPassword)
	if err != nil {
		PanicErr("Failed to authenticate with MediaWiki with username ", config.BotUsername, " - error was ", err)
	}

	return w
}

// FetchWikitext takes a pageId and gets the wikitext of that page.
// The default functionality in the library does not work for this in
// my experience; it just returns an empty string for some reason. So we're rolling our own!
func FetchWikitext(pageID string) (content string, err error) {
	return fetchWikitextFrom("pageid", pageID)
}

// FetchWikitextFromTitle takes a title and gets the wikitext of that page.
func FetchWikitextFromTitle(pageTitle string) (content string, err error) {
	return fetchWikitextFrom("page", pageTitle)
}

func fetchWikitextFrom(identifierName string, identifier string) (string, error) {
	pageContent, err := w.Get(params.Values{
		"action":       "parse",
		identifierName: identifier,
		"prop":         "wikitext",
	})
	if err != nil {
		return "", err
	}
	text, err := pageContent.GetString("parse", "wikitext")
	if err != nil {
		return "", err
	}
	return text, nil
}
