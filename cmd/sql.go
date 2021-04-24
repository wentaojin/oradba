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
	sqlCmd := &grumble.Command{
		Name:     "sql",
		Help:     "sql tools",
		LongHelp: "Options for query oracle db sql info about top, activity summary, plan and etc",
	}
	App.AddCommand(sqlCmd)

	sqlCmd.AddCommand(&grumble.Command{
		Name: "top",
		Help: "query oracle db sql top 10",
		Args: func(a *grumble.Args) {
			a.String("orderCol", "specify sql order desc column, options: [executions,elapsed_time_ms,cpu_time,buffer_gets,disk_reads,direct_writes,rows_processed]", grumble.Default("elapsed_time_ms"))
		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBTOPSQL(c.Args.String("orderCol")); err != nil {
				return err
			}
			return nil
		},
	})
	sqlCmd.AddCommand(&grumble.Command{
		Name: "like",
		Help: "query oracle db sql info about sql_id, plan_hash_value and child number, according matching sql text like query",
		Args: func(a *grumble.Args) {
			a.String("text", "specify sql text used match query")
		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBSQLMatchText(c.Args.String("text")); err != nil {
				return err
			}
			return nil
		},
	})
	sqlCmd.AddCommand(&grumble.Command{
		Name: "resource",
		Help: "query oracle db sql resource usage details according to sql_id",
		Args: func(a *grumble.Args) {
			a.String("sqlid", "specify sql id")
		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBSQLResourceUsedDetailBySQLID(c.Args.String("sqlid")); err != nil {
				return err
			}
			return nil
		},
	})
	sqlCmd.AddCommand(&grumble.Command{
		Name: "summary",
		Help: "query oracle db sql a certain period of time summary info according to sql_id, startTime, endTime",
		Args: func(a *grumble.Args) {
			a.String("sqlid", "specify sql id")
			a.String("startTime", "specify query start time")
			a.String("endTime", "specify query end time")

		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBSQLSummaryBySQLID(
				c.Args.String("sqlid"),
				c.Args.String("startTime"),
				c.Args.String("endTime")); err != nil {
				return err
			}
			return nil
		},
	})

	planCmd := &grumble.Command{
		Name: "plan",
		Help: "query oracle db sql plan by awr, cursor, xplan",
	}
	sqlCmd.AddCommand(planCmd)

	planCmd.AddCommand(&grumble.Command{
		Name: "awr",
		Help: "query oracle db sql by awr",
		Args: func(a *grumble.Args) {
			a.String("sqlid", "specify sql id")
			a.String("planhashvalue", "specify sql plan hash value")
		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBSQLPlanByAWR(
				c.Args.String("sqlid"),
				c.Args.String("planhashvalue")); err != nil {
				return err
			}
			return nil
		},
	})
	planCmd.AddCommand(&grumble.Command{
		Name: "cursor",
		Help: "query oracle db sql by cursor",
		Args: func(a *grumble.Args) {
			a.String("sqlid", "specify sql id")
			a.String("childnumber", "specify sql plan child number")
			a.String("format", "specify sql plan output format, options: [basic,typical,serial,all,advanced]", grumble.Default("typical"))
		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBSQLPlanByCursor(
				c.Args.String("sqlid"),
				c.Args.String("childnumber"),
				c.Args.String("format")); err != nil {
				return err
			}
			return nil
		},
	})
	planCmd.AddCommand(&grumble.Command{
		Name: "xplan",
		Help: "query oracle db sql by xplan",
		Args: func(a *grumble.Args) {
			a.String("sqlid", "specify sql id")
			a.String("childnumber", "specify sql plan child number")
			a.String("format", "specify sql plan output format, options: [basic,typical,serial,all,advanced]", grumble.Default("typical"))
		},
		Run: func(c *grumble.Context) error {
			if err := oracle.QueryOracleDBSQLPlanByXplan(
				c.Args.String("sqlid"),
				c.Args.String("childnumber"),
				c.Args.String("format")); err != nil {
				return err
			}
			return nil
		},
	})
}
