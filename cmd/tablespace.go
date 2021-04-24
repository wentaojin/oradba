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
	tblsCmd := &grumble.Command{
		Name:     "tablespace",
		Aliases:  []string{"tbls"},
		Help:     "tablespace tools",
		LongHelp: "Options for query oracle db tablespace about path, size, used and etc info",
	}
	App.AddCommand(tblsCmd)

	tblsCmd.AddCommand(&grumble.Command{
		Name: "detail",
		Help: "query oracle db all tableSpaces info about total size,use size,autoextend,datafile numbers and etc info [size: M]",
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBTablespaceDetailInfo(); err != nil {
				return err
			}
			return nil
		},
	})
	tblsCmd.AddCommand(&grumble.Command{
		Name: "io",
		Help: "query oracle db tablespace io info about physical read adn write info",
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBTablespaceIOStatInfo(); err != nil {
				return err
			}
			return nil
		},
	})
	tblsCmd.AddCommand(&grumble.Command{
		Name:    "username",
		Aliases: []string{"u"},
		Help:    "query oracle db one username default tablespace info",
		Args: func(a *grumble.Args) {
			a.String("username", "specify oracle one username")
		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBTablespaceDetailInfoByUser(c.Args.String("username")); err != nil {
				return err
			}
			return nil
		},
	})
	tblsCmd.AddCommand(&grumble.Command{
		Name: "sysaux",
		Help: "query oracle db sysaux tablespace occupancy info",
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBSysauxTablespaceInfo(); err != nil {
				return err
			}
			return nil
		},
	})
}
