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
	ashCmd := &grumble.Command{
		Name:     "ash",
		Help:     "ash tools",
		LongHelp: "Options for query oracle db ash info about sql plan, top sql, db healthy and etc info",
	}
	App.AddCommand(ashCmd)

	ashCmd.AddCommand(&grumble.Command{
		Name: "top",
		Help: "query oracle db ash top 10 info in the recent an hour",
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBASHTOPRecentAnHour(); err != nil {
				return err
			}
			return nil
		},
	})
	ashCmd.AddCommand(&grumble.Command{
		Name: "history",
		Help: "query oracle db ash top 10 history info within a certain period of time",
		Args: func(a *grumble.Args) {
			a.String("startTime", "specify ash start time")
			a.String("endTime", "specify ash end time")
		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBASHTOPHistoryByTimeLimit(c.Args.String("startTime"), c.Args.String("endTime")); err != nil {
				return err
			}
			return nil
		},
	})
}
