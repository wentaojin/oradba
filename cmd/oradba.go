/*
Copyright © 2020 Marvin

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/WentaoJin/oradba/db"

	"github.com/desertbit/grumble"
	"github.com/fatih/color"
)

var App = grumble.New(&grumble.Config{
	Name:                  "oradba",
	Description:           "CLI Tool For Oracle DB",
	HistoryFile:           "/tmp/oradba.hist",
	Prompt:                "oradba » ",
	PromptColor:           color.New(color.FgGreen, color.Bold),
	HelpHeadlineColor:     color.New(color.FgGreen),
	HelpHeadlineUnderline: true,
	HelpSubCommands:       true,

	Flags: func(f *grumble.Flags) {
		f.String("c", "config", "config.toml", "oracle db config info")
	},
})

func init() {
	App.OnInit(func(a *grumble.App, flags grumble.FlagMap) error {
		cfg, err := db.ReadConfigFile(flags.String("config"))
		if err != nil {
			return err
		}
		db.ORA, err = db.NewOracleDBEngine(cfg)
		if err != nil {
			return err
		}
		return nil
	})

	App.SetPrintASCIILogo(func(a *grumble.App) {
		a.Println("Welcome to")
		a.Println("                     _ _           ")
		a.Println("  ___  _ __ __ _  __| | |__   __ _ ")
		a.Println(" / _ \\| '__/ _` |/ _` | '_ \\ / _` |")
		a.Println("| (_) | | | (_| | (_| | |_) | (_| |")
		a.Println(" \\___/|_|  \\__,_|\\__,_|_.__/ \\__,_|")
	})
}
