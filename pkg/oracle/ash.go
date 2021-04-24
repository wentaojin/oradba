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
package oracle

import (
	"fmt"
	"os"

	"github.com/wentaojin/oradba/db"
	"github.com/wentaojin/oradba/pkg/util"
)

/*
	查看近一小时AShTop信息
*/
func QueryOracleDBASHTOPRecentAnHour() error {
	columns, values, err := db.Query(`SELECT *
  FROM (SELECT /*+ LEADING(a) USE_HASH(u) */
         LPAD(ROUND(RATIO_TO_REPORT(COUNT(*)) OVER() * 100) || '%', 5, ' ') "%This",
         username,
         session_id,
         session_serial#,
         nvl(sql_id, 'Null') sql_id,
         case SQL_OPCODE
           when 1 then
            'create table'
           when 2 then
            'INSERT'
           when 3 then
            'SELECT'
           when 6 then
            'UPDATE'
           when 7 then
            'DELETE'
           when 9 then
            'create index'
           when 11 then
            'ALTER INDEX'
           when 26 then
            'LOCK table'
           when 42 then
            'ALTER_SESSION (NOT ddl)'
           when 44 then
            'COMMIT'
           when 45 then
            'rollback'
           when 46 then
            'savepoint'
           when 47 then
            'PL/SQL BLOCK or begin/declare'
           when 48 then
            'set transaction'
           when 50 then
            'explain'
           when 62 then
            'analyze table'
           when 90 then
            'set constraints'
           when 170 then
            'call'
           when 189 then
            'merge'
           else
            'other'
         end command_type,
         NVl(event, 'Null') event,
         ROUND(COUNT(*) / (((sysdate) - (sysdate - 1 / 24)) * 86400), 1) AAS,
         COUNT(*) "TotalSeconds"
        --, SUM(CASE WHEN wait_class IS NULL           THEN 1 ELSE 0 END) "CPU"
        --, SUM(CASE WHEN wait_class ='User I/O'       THEN 1 ELSE 0 END) "User I/O"
        --, SUM(CASE WHEN wait_class ='Application'    THEN 1 ELSE 0 END) "Application"
        --, SUM(CASE WHEN wait_class ='Concurrency'    THEN 1 ELSE 0 END) "Concurrency"
        --, SUM(CASE WHEN wait_class ='Commit'         THEN 1 ELSE 0 END) "Commit"
        --, SUM(CASE WHEN wait_class ='Configuration'  THEN 1 ELSE 0 END) "Configuration"
        --, SUM(CASE WHEN wait_class ='Cluster'        THEN 1 ELSE 0 END) "Cluster"
        --, SUM(CASE WHEN wait_class ='Idle'           THEN 1 ELSE 0 END) "Idle"
        --, SUM(CASE WHEN wait_class ='Network'        THEN 1 ELSE 0 END) "Network"
        --, SUM(CASE WHEN wait_class ='System I/O'     THEN 1 ELSE 0 END) "System I/O"
        --, SUM(CASE WHEN wait_class ='Scheduler'      THEN 1 ELSE 0 END) "Scheduler"
        --, SUM(CASE WHEN wait_class ='Administrative' THEN 1 ELSE 0 END) "Administrative"
        --, SUM(CASE WHEN wait_class ='Queueing'       THEN 1 ELSE 0 END) "Queueing"
        --, SUM(CASE WHEN wait_class ='Other'          THEN 1 ELSE 0 END) "Other"
        --, MIN(sample_time)
        --, MAX(sample_time)
        --, MAX(sql_exec_id) - MIN(sql_exec_id)
          FROM (SELECT a.*,
                       TO_CHAR(CASE
                                 WHEN session_state = 'ON CPU' THEN
                                  p1
                                 ELSE
                                  null
                               END,
                               '0XXXXXXXXXXXXXXX') p1hex,
                       TO_CHAR(CASE
                                 WHEN session_state = 'ON CPU' THEN
                                  p2
                                 ELSE
                                  null
                               END,
                               '0XXXXXXXXXXXXXXX') p2hex,
                       TO_CHAR(CASE
                                 WHEN session_state = 'ON CPU' THEN
                                  p3
                                 ELSE
                                  null
                               END,
                               '0XXXXXXXXXXXXXXX') p3hex
                  FROM v$active_session_history a) a,
               dba_users u
         WHERE a.user_id = u.user_id(+)
           AND session_type = 'FOREGROUND'
           AND sample_time BETWEEN sysdate - 1 / 24 AND sysdate
         GROUP BY username,
                  session_id,
                  session_serial#,
                  sql_id,
                  SQL_OPCODE,
                  event
         ORDER BY "TotalSeconds" DESC,
                  username,
                  session_id,
                  session_serial#,
                  sql_id,
                  SQL_OPCODE,
                  event)
 WHERE ROWNUM <= 10`)
	if err != nil {
		return err
	}
	util.NewMarkdownTableStyle(os.Stdout, columns, values)
	return nil
}

/*
	查看某段时间的AShTop信息
*/
func QueryOracleDBASHTOPHistoryByTimeLimit(startTime, endTime string) error {
	columns, values, err := db.Query(fmt.Sprintf(`SELECT *
  FROM (SELECT /*+ LEADING(a) USE_HASH(u) */
         LPAD(ROUND(RATIO_TO_REPORT(COUNT(*)) OVER() * 100) || '%%', 5, ' ') "%%This",
         username,
         session_id,
         session_serial#,
         nvl(sql_id, 'Null') sqlid,
         case SQL_OPCODE
           when 1 then
            'create table'
           when 2 then
            'INSERT'
           when 3 then
            'SELECT'
           when 6 then
            'UPDATE'
           when 7 then
            'DELETE'
           when 9 then
            'create index'
           when 11 then
            'ALTER INDEX'
           when 26 then
            'LOCK table'
           when 42 then
            'ALTER_SESSION (NOT ddl)'
           when 44 then
            'COMMIT'
           when 45 then
            'rollback'
           when 46 then
            'savepoint'
           when 47 then
            'PL/SQL BLOCK or begin/declare'
           when 48 then
            'set transaction'
           when 50 then
            'explain'
           when 62 then
            'analyze table'
           when 90 then
            'set constraints'
           when 170 then
            'call'
           when 189 then
            'merge'
           else
            'other'
         end command_type,
         NVl(event, 'Null') event,
         ROUND(COUNT(*) / (((sysdate) - (sysdate - 1 / 24)) * 86400), 1) AAS,
         COUNT(*) "TotalSeconds"
        --, SUM(CASE WHEN wait_class IS NULL           THEN 1 ELSE 0 END) "CPU"
        --, SUM(CASE WHEN wait_class ='User I/O'       THEN 1 ELSE 0 END) "User I/O"
        --, SUM(CASE WHEN wait_class ='Application'    THEN 1 ELSE 0 END) "Application"
        --, SUM(CASE WHEN wait_class ='Concurrency'    THEN 1 ELSE 0 END) "Concurrency"
        --, SUM(CASE WHEN wait_class ='Commit'         THEN 1 ELSE 0 END) "Commit"
        --, SUM(CASE WHEN wait_class ='Configuration'  THEN 1 ELSE 0 END) "Configuration"
        --, SUM(CASE WHEN wait_class ='Cluster'        THEN 1 ELSE 0 END) "Cluster"
        --, SUM(CASE WHEN wait_class ='Idle'           THEN 1 ELSE 0 END) "Idle"
        --, SUM(CASE WHEN wait_class ='Network'        THEN 1 ELSE 0 END) "Network"
        --, SUM(CASE WHEN wait_class ='System I/O'     THEN 1 ELSE 0 END) "System I/O"
        --, SUM(CASE WHEN wait_class ='Scheduler'      THEN 1 ELSE 0 END) "Scheduler"
        --, SUM(CASE WHEN wait_class ='Administrative' THEN 1 ELSE 0 END) "Administrative"
        --, SUM(CASE WHEN wait_class ='Queueing'       THEN 1 ELSE 0 END) "Queueing"
        --, SUM(CASE WHEN wait_class ='Other'          THEN 1 ELSE 0 END) "Other"
        --, MIN(sample_time)
        --, MAX(sample_time)
        --, MAX(sql_exec_id) - MIN(sql_exec_id)
          FROM (SELECT a.*,
                       TO_CHAR(CASE
                                 WHEN session_state = 'ON CPU' THEN
                                  p1
                                 ELSE
                                  null
                               END,
                               '0XXXXXXXXXXXXXXX') p1hex,
                       TO_CHAR(CASE
                                 WHEN session_state = 'ON CPU' THEN
                                  p2
                                 ELSE
                                  null
                               END,
                               '0XXXXXXXXXXXXXXX') p2hex,
                       TO_CHAR(CASE
                                 WHEN session_state = 'ON CPU' THEN
                                  p3
                                 ELSE
                                  null
                               END,
                               '0XXXXXXXXXXXXXXX') p3hex
                  FROM v$active_session_history a) a,
               dba_users u
         WHERE a.user_id = u.user_id(+)
           AND session_type = 'FOREGROUND'
           AND to_char(sample_time,'yyyy-mm-dd hh24:mi:ss') BETWEEN '%s' AND '%s'
         GROUP BY username,
                  session_id,
                  session_serial#,
                  sql_id,
                  SQL_OPCODE,
                  event
         ORDER BY "TotalSeconds" DESC,
                  username,
                  session_id,
                  session_serial#,
                  sql_id,
                  SQL_OPCODE,
                  event)
 WHERE ROWNUM <= 10`, startTime, endTime))
	if err != nil {
		return err
	}
	util.NewMarkdownTableStyle(os.Stdout, columns, values)
	return nil
}
