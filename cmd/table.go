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
	tblCmd := &grumble.Command{
		Name:     "table",
		Aliases:  []string{"tbl"},
		Help:     "table tools",
		LongHelp: "Options for query oracle db table about table, partition table, index and etc info",
	}
	App.AddCommand(tblCmd)

	tblCmd.AddCommand(&grumble.Command{
		Name: "info",
		Help: "query oracle db one username table num_rows, blocks, avg_rows_len and etc info",
		Flags: func(f *grumble.Flags) {
			f.String("u", "username", "marvin", "specify oracle one username")
			f.String("t", "table", "marvin", "specify oracle one tablename")
		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBTableDetailInfo(c.Flags.String("username"), c.Flags.String("table")); err != nil {
				return err
			}
			return nil
		},
	})
	tblCmd.AddCommand(&grumble.Command{
		Name: "column",
		Help: "query oracle db one username and one table columns about column_name, data_type, num_distinct and etc info",
		Flags: func(f *grumble.Flags) {
			f.String("u", "username", "marvin", "specify oracle one username")
			f.String("t", "table", "marvin", "specify oracle one tablename")
		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBTableColumnDetailInfo(c.Flags.String("username"), c.Flags.String("table")); err != nil {
				return err
			}
			return nil
		},
	})
	tblCmd.AddCommand(&grumble.Command{
		Name: "index",
		Help: "query oracle db one username and one table index about index_name, blevel, num_rows, distinct_keys and etc info",
		Flags: func(f *grumble.Flags) {
			f.String("u", "username", "marvin", "specify oracle one username")
			f.String("t", "table", "marvin", "specify oracle one tablename")
		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBTableIndexDetailInfo(c.Flags.String("username"), c.Flags.String("table")); err != nil {
				return err
			}
			return nil
		},
	})
	tblCmd.AddCommand(&grumble.Command{
		Name: "indexCol",
		Help: "query oracle db one username and one table name about index_name, index_type, column_list and etc info",
		Flags: func(f *grumble.Flags) {
			f.String("u", "username", "marvin", "specify oracle one username")
			f.String("t", "table", "marvin", "specify oracle one tablename")
		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBTableIndexColumnDetailInfo(c.Flags.String("username"), c.Flags.String("table")); err != nil {
				return err
			}
			return nil
		},
	})
	tblCmd.AddCommand(&grumble.Command{
		Name: "primary",
		Help: "query oracle db one username and one table name about primary key info",
		Flags: func(f *grumble.Flags) {
			f.String("u", "username", "marvin", "specify oracle one username")
			f.String("t", "table", "marvin", "specify oracle one tablename")
		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBTablePrimaryKeyDetailInfo(c.Flags.String("username"), c.Flags.String("table")); err != nil {
				return err
			}
			return nil
		},
	})
	tblCmd.AddCommand(&grumble.Command{
		Name: "unique",
		Help: "query oracle db one username and one table name about unique key info",
		Flags: func(f *grumble.Flags) {
			f.String("u", "username", "marvin", "specify oracle one username")
			f.String("t", "table", "marvin", "specify oracle one tablename")
		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBTableUniqueKeyDetailInfo(c.Flags.String("username"), c.Flags.String("table")); err != nil {
				return err
			}
			return nil
		},
	})
	tblCmd.AddCommand(&grumble.Command{
		Name: "foreign",
		Help: "query oracle db one username and one table name about foreign key info",
		Flags: func(f *grumble.Flags) {
			f.String("u", "username", "marvin", "specify oracle one username")
			f.String("t", "table", "marvin", "specify oracle one tablename")
		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBTableForeignKeyDetailInfo(c.Flags.String("username"), c.Flags.String("table")); err != nil {
				return err
			}
			return nil
		},
	})
	tblCmd.AddCommand(&grumble.Command{
		Name: "check",
		Help: "query oracle db one username and one table name about check key info",
		Flags: func(f *grumble.Flags) {
			f.String("u", "username", "marvin", "specify oracle one username")
			f.String("t", "table", "marvin", "specify oracle one tablename")
		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBTableCheckKeyDetailInfo(c.Flags.String("username"), c.Flags.String("table")); err != nil {
				return err
			}
			return nil
		},
	})
}
