package ybtools

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/metal3d/go-slugify"
	"gopkg.in/yaml.v2"
)

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

type configObject struct {
	APIEndpoint string
	BotUsername string
	EditLimits  map[string]int64
}

// acts like an interface for config files
// edit limits from tool configs are unloaded into here
type toolConfigWithEditLimit struct {
	EditLimit int64
}

const thisDirConfigPath string = "config.yml"
const subDirConfigPath string = "../config-global.yml"

var botPassword string
var config configObject
var taskConfigFile []byte

// ParseTaskConfig takes in a pointer to a config object
// and populates the config object with the task configuration
func ParseTaskConfig(cobj *interface{}) {
	yaml.Unmarshal(taskConfigFile, cobj)
}

func init() {
	botConfigFile, err := ioutil.ReadFile(findConfigFile())
	if err != nil {
		log.Fatal("Bot config file could not be read at detected path!")
	}
	err = yaml.UnmarshalStrict(botConfigFile, &config)
	if err != nil {
		log.Fatal("Bot config file was invalid!")
	}
}

func setupTaskConfigFile() {
	var err error
	var taskConfigForEditLimit toolConfigWithEditLimit

	taskConfigFile, err = ioutil.ReadFile("config-" + slugify.Marshal(taskName) + ".yml")
	if err != nil {
		log.Println("No task-specific config file found, ignoring")
	}

	// Immediately parse the file for an edit limit and only an edit limit
	yaml.Unmarshal(taskConfigFile, &taskConfigForEditLimit)
	if taskConfigForEditLimit.EditLimit > 0 {
		setupEditLimit(taskConfigForEditLimit.EditLimit)
	}
}

func findConfigFile() string {
	if _, err := os.Stat(thisDirConfigPath); os.IsNotExist(err) {
		if _, err := os.Stat(subDirConfigPath); os.IsNotExist(err) {
			return subDirConfigPath
		}
		log.Fatal("Couldn't find a config file either in this directory or in the one above it!")
	}
	return thisDirConfigPath
}
