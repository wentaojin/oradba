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
	"github.com/WentaoJin/oradba/pkg/oracle"
	"github.com/desertbit/grumble"
)

func init() {
	sessionCmd := &grumble.Command{
		Name:     "session",
		Help:     "session tools",
		LongHelp: "Options for query oracle db about sid, spid, sql_id and etc info",
	}
	App.AddCommand(sessionCmd)

	sessionCmd.AddCommand(&grumble.Command{
		Name: "spid",
		Help: "query oracle db session sid by use spid",
		Args: func(a *grumble.Args) {
			a.String("spid", "specify oracle db spid")
		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBSessionSIDBySPID(c.Args.String("spid")); err != nil {
				return err
			}
			return nil
		},
	})
	sessionCmd.AddCommand(&grumble.Command{
		Name: "sid",
		Help: "query oracle db session spid by use sid",
		Args: func(a *grumble.Args) {
			a.String("sid", "specify oracle db sid")
		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBSessionSPIDBySID(c.Args.String("sid")); err != nil {
				return err
			}
			return nil
		},
	})
	sessionCmd.AddCommand(&grumble.Command{
		Name: "waitBySid",
		Help: "query oracle db session wait event and wait time by use sid",
		Args: func(a *grumble.Args) {
			a.String("sid", "specify oracle db sid")
		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBWaitEventBySid(c.Args.String("sid")); err != nil {
				return err
			}
			return nil
		},
	})
	sessionCmd.AddCommand(&grumble.Command{
		Name: "wait",
		Help: "query oracle db current wait session about wait sid, event, username, sql_id, wait_time and etc info",
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBWaitEvent(); err != nil {
				return err
			}
			return nil
		},
	})
	sessionCmd.AddCommand(&grumble.Command{
		Name: "undo",
		Help: "query oracle db current session used undo size top 10",
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBUsedUNDOTOP(); err != nil {
				return err
			}
			return nil
		},
	})
	sessionCmd.AddCommand(&grumble.Command{
		Name: "temp",
		Help: "query oracle db current session used temp size top 10",
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBUsedTEMPTOP(); err != nil {
				return err
			}
			return nil
		},
	})
	sessionCmd.AddCommand(&grumble.Command{
		Name: "pga",
		Help: "query oracle db current session used pga size top 10",
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBUsedPGATOP(); err != nil {
				return err
			}
			return nil
		},
	})
	sessionCmd.AddCommand(&grumble.Command{
		Name: "blocking",
		Help: "query oracle db current session block info without lock mode include dead-lock",
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBSessionBlocking(); err != nil {
				return err
			}
			return nil
		},
	})
	sessionCmd.AddCommand(&grumble.Command{
		Name: "block",
		Help: "query oracle db current blocking details with lock mode include dead-lock",
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBSessionBlock(); err != nil {
				return err
			}
			return nil
		},
	})
}
