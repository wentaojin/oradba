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
	infoCmd := &grumble.Command{
		Name:     "info",
		Help:     "info tools",
		LongHelp: "Options for query oracle db base, table and tablespace etc info",
	}
	App.AddCommand(infoCmd)

	infoCmd.AddCommand(&grumble.Command{
		Name: "db",
		Help: "query oracle db name, instance name, version etc info",
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBInfo(); err != nil {
				return err
			}
			return nil
		},
	})
	infoCmd.AddCommand(&grumble.Command{
		Name: "instance",
		Help: "query oracle instance, host, startup time etc info",
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleInstanceListInfo(); err != nil {
				return err
			}
			return nil
		},
	})
	infoCmd.AddCommand(&grumble.Command{
		Name: "memory",
		Help: "query oracle component memory used info [size: G]",
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleComponentMemoryInfo(); err != nil {
				return err
			}
			return nil
		},
	})
	infoCmd.AddCommand(&grumble.Command{
		Name: "tablespace",
		Help: "query oracle tablespace all total size, all free size info [size: G]",
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleTablespaceSummaryInfo(); err != nil {
				return err
			}
			return nil
		},
	})
	infoCmd.AddCommand(&grumble.Command{
		Name: "backup",
		Help: "query oracle last time backup start time, end time, status etc info",
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleLastBackupInfo(); err != nil {
				return err
			}
			return nil
		},
	})
	infoCmd.AddCommand(&grumble.Command{
		Name: "params",
		Help: "query oracle param name, instance name, params value info",
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleParamsListInfo(); err != nil {
				return err
			}
			return nil
		},
	})
}
