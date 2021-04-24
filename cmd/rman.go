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
	"github.com/desertbit/grumble"
	"github.com/wentaojin/oradba/pkg/oracle"
)

func init() {
	rmanCmd := &grumble.Command{
		Name:     "rman",
		Help:     "rman tools",
		LongHelp: "Options for query oracle db rman info about rman process and rman lasted status",
	}
	App.AddCommand(rmanCmd)

	rmanCmd.AddCommand(&grumble.Command{
		Name: "process",
		Help: "query oracle db rman process about backup or restore",
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBRMANProcess(); err != nil {
				return err
			}
			return nil
		},
	})
	rmanCmd.AddCommand(&grumble.Command{
		Name: "status",
		Help: "query oracle db backup info in some time about success or failed or running",
		Args: func(a *grumble.Args) {
			a.String("startTime", "specify query start time")
			a.String("endTime", "specify query end time")
		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBRMANStatus(
				c.Args.String("startTime"),
				c.Args.String("endTime"),
			); err != nil {
				return err
			}
			return nil
		},
	})
}
