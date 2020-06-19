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

// DefaultMaxlag is a Maxlag representing sensible defaults for any non-urgent
// task being run by a Yapperbot.
var DefaultMaxlag mwclient.Maxlag = mwclient.Maxlag{
	On:      true,
	Timeout: "5",
	Retries: 3,
}

// NoMaxlag is a Maxlag for urgent tasks, where maxlag must be disabled.
var NoMaxlag mwclient.Maxlag = mwclient.Maxlag{
	On:      true,
	Timeout: "5",
	Retries: 3,
}

// CreateAndAuthenticateClient uses the details already passed into ybtools
// in setup.go to return a fully-authenticated mwclient
func CreateAndAuthenticateClient(maxlag mwclient.Maxlag) *mwclient.Client {
	if settings.TaskName == "" || settings.BotUser == "" {
		PanicErr("Call ybtools.SetupBot first!")
	}

	var err error

	w, err = mwclient.New(config.APIEndpoint, "Yapperbot-"+settings.TaskName+" on User:"+settings.BotUser+" - Golang, licensed GNU GPL")
	if err != nil {
		PanicErr("Failed to create MediaWiki client with error ", err)
	}

	// This is necessary because maxlag.sleep is unexported,
	// and is only configured correctly for production within
	// mwclient.New. Annoying, but this is a bit of an edge case /shrug
	w.Maxlag.On = maxlag.On
	w.Maxlag.Retries = maxlag.Retries
	w.Maxlag.Timeout = maxlag.Timeout

	err = w.Login(config.BotUsername, botPassword)
	if err != nil {
		PanicErr("Failed to authenticate with MediaWiki with username ", config.BotUsername, " - error was ", err)
	}

	// runs here to make sure we have a client authenticated when we run it
	killTaskIfNeeded()

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
