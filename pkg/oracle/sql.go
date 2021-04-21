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
	"database/sql"
	"fmt"
	"os"

	"github.com/WentaoJin/oradba/db"
	"github.com/WentaoJin/oradba/pkg/util"
)

/*
	查看TopSql信息
*/
func QueryOracleDBTOPSQL(cond string) error {
	columns, values, err := db.Query(fmt.Sprintf(`select *
  from (select s.sql_id,
               a.executions,
               round(a.elapsed_time / 1000000, 2) elapsed_time_ms,
               round(a.elapsed_time / a.executions / 1000) ms_by_exec,
               round(a.cpu_time / 1000000, 2) cpu_time,
               round(a.cpu_time / a.executions / 1000, 2) ms_cpu_time,
               a.buffer_gets buffer_gets,
               round(a.buffer_gets / a.executions, 2) avg_buffer_gets,
               a.disk_reads disk_reads,
               round(a.disk_reads / a.executions, 2) avg_disk_reads,
               a.DIRECT_WRITES direct_writes,
               a.rows_processed rows_processed,
               round(a.rows_processed / a.executions, 2) avg_rows_processed
          from v$sql a, v$sqlarea s
         where a.address = s.address
           and a.hash_value = s.hash_value
           and a.executions <> 0
         order by %s desc)
 where rownum <= 10`, cond))
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}

/*
	根据sql_text匹配查找sql_id,plan_hash_value以及child number
*/
func QueryOracleDBSQLMatchText(text string) error {
	columns, values, err := db.Query(fmt.Sprintf(`select sql_id,
       cast(address as varchar2(200)),
       to_char(hash_value),
       to_char(plan_hash_value),
       to_char(child_number),
       sql_text
  from v$sql
 where sql_text like '%s'
   and sql_text not like '%%v$sql%%'`, text))
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}

/*
	根据sql_id查看具体sql资源使用详情
*/
func QueryOracleDBSQLResourceUsedDetailBySQLID(sqlID string) error {
	columns, values, err := db.Query(fmt.Sprintf(`select s.sql_id,
       a.executions,
       round(a.elapsed_time / 1000000, 2) elapsed_time_ms,
       round(a.elapsed_time / a.executions / 1000) ms_by_exec,
       round(a.cpu_time / 1000000, 2) cpu_time,
       round(a.cpu_time / a.executions / 1000, 2) ms_cpu_time,
       a.buffer_gets buffer_gets,
       round(a.buffer_gets / a.executions, 2) avg_buffer_gets,
       a.disk_reads disk_reads,
       round(a.disk_reads / a.executions, 2) avg_disk_reads,
       a.DIRECT_WRITES,
       a.rows_processed rows_processed,
       round(a.rows_processed / a.executions, 2) avg_rows_processed,
       s.sql_fulltext
  from v$sql a, v$sqlarea s
 where a.address = s.address
   and a.hash_value = s.hash_value
   and a.executions <> 0
   and s.sql_id = '%s'`, sqlID))
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}

/*
	根据sql_id查看具体sql详情
*/
func QueryOracleDBSQLSummaryBySQLID(sqlID, startTime, endTime string) error {
	columns, values, err := db.Query(fmt.Sprintf(`SELECT /*+ LEADING(a) USE_HASH(u) */
 LPAD(ROUND(RATIO_TO_REPORT(COUNT(*)) OVER() * 100) || '%%', 5, ' ') "%%This",
 ROUND(COUNT(*) / (((sysdate) - (sysdate - 1 / 24)) * 86400), 1) AAS,
 username,
 session_id || ',' || SESSION_SERIAL# || ',@' || inst_id AS session_id,
 sql_id,
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
 PROGRAM,
 MODULE,
 MACHINE,
 to_char(SQL_PLAN_HASH_VALUE),
 child_number
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
          FROM gv$active_session_history a) a,
       dba_users u
 WHERE a.user_id = u.user_id(+)
   AND session_type = 'FOREGROUND'
   AND a.sql_id = '%s'
   AND sample_time BETWEEN '%s' AND '%s'
 GROUP BY username,
          session_id,
          SESSION_SERIAL#,
          inst_id,
          sql_id,
          SQL_OPCODE,
          PROGRAM,
          MODULE,
          child_number,
          MACHINE,
          SQL_PLAN_HASH_VALUE`, sqlID, startTime, endTime))
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}

/*
	根据sql_id和hash_value查看具体sql执行计划
*/
func QueryOracleDBSQLPlanByAWR(sqlid, planhashvalue string) error {
	columns, values, err := db.Query(fmt.Sprintf(`SELECT *
  FROM table(dbms_xplan.display_awr('%s',
                                    '%s',
                                    null,
                                    'typical +peeked_binds'))`, sqlid, planhashvalue))
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}

/*
	根据sql_id查看具体sql执行计划
*/
func QueryOracleDBSQLPlanByCursor(sqlid, childnumber, format string) error {
	columns, values, err := db.Query(fmt.Sprintf(`SELECT * FROM table(dbms_xplan.display_cursor('%s', %s, '%s'))`, sqlid, childnumber, format))
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}

/*
	根据sql_id查看具体sql执行计划,并且带有执行顺序的执行计划 order,需要编译xplan包体
*/
func QueryOracleDBSQLPlanByXplan(sqlid, childnumber, format string) error {
	var state string
	err := db.ORA.QueryRow(`select status from dba_objects where object_type = 'PACKAGE' and object_name ='XPLAN' and owner=(select user from dual)`).Scan(&state)
	if err == sql.ErrNoRows {
		return fmt.Errorf("package xplan not existed,Will make package xplan")
	}
	if state != "VALID" {
		if err = makeOracleDBXplanPackage(); err != nil {
			return err
		}
	}

	columns, values, err := db.Query(fmt.Sprintf(`SELECT * FROM table(xplan.display_cursor('%s',%s,'%s'))`, sqlid, childnumber, format))
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}

func makeOracleDBXplanPackage() error {
	sql1 := `
CREATE OR REPLACE TYPE xplan_ot force AS OBJECT( plan_table_output VARCHAR2(300) );
`
	sql2 := `CREATE OR REPLACE TYPE xplan_ntt AS  TABLE OF xplan_ot;`
	sql3 := `
CREATE OR REPLACE PACKAGE xplan AS

   FUNCTION display_cursor( p_sql_id          IN VARCHAR2 DEFAULT NULL,
                            p_cursor_child_no IN INTEGER  DEFAULT 0,
                            p_format          IN VARCHAR2 DEFAULT 'TYPICAL' )
      RETURN xplan_ntt PIPELINED;
END xplan;`
	sql4 := `
CREATE OR REPLACE PACKAGE BODY xplan AS

   TYPE ntt_order_map_binds IS TABLE OF VARCHAR2(100);

   TYPE aat_order_map IS TABLE OF PLS_INTEGER
      INDEX BY PLS_INTEGER;

   g_map  aat_order_map;
   g_hdrs PLS_INTEGER;
   g_len  PLS_INTEGER;
   g_pad  VARCHAR2(300);

   ----------------------------------------------------------------------------
   PROCEDURE reset_state IS
   BEGIN
      g_hdrs := 0;
      g_len  := 0;
      g_pad  := NULL;
      g_map.DELETE;
   END reset_state;

   ----------------------------------------------------------------------------
   PROCEDURE build_order_map( p_sql   IN VARCHAR2,
                              p_binds IN ntt_order_map_binds ) IS

      TYPE rt_id_data IS RECORD
      ( id  PLS_INTEGER
      , ord PLS_INTEGER );

      TYPE aat_id_data IS TABLE OF rt_id_data
         INDEX BY PLS_INTEGER;

      aa_ids   aat_id_data;
      v_cursor SYS_REFCURSOR;
      v_sql    VARCHAR2(32767);

   BEGIN

      -- Build SQL template...
      -- ---------------------
      v_sql := 'WITH sql_plan_data AS ( ' || 
                        p_sql || '
                        )
                ,    hierarchical_sql_plan_data AS (
                        SELECT id
                        FROM   sql_plan_data
                        START WITH id = 0
                        CONNECT BY PRIOR id = parent_id
                        ORDER SIBLINGS BY id DESC
                        )
                SELECT id
                ,      ROW_NUMBER() OVER (ORDER BY ROWNUM DESC) AS ord
                FROM   hierarchical_sql_plan_data';

      -- Binds will differ according to plan type...
      -- -------------------------------------------
      CASE p_binds.COUNT
         WHEN 0
         THEN
            OPEN v_cursor FOR v_sql;
         WHEN 1
         THEN
            OPEN v_cursor FOR v_sql USING p_binds(1);
         WHEN 2
         THEN
            OPEN v_cursor FOR v_sql USING p_binds(1),
                                          TO_NUMBER(p_binds(2));
         WHEN 3
         THEN
            OPEN v_cursor FOR v_sql USING p_binds(1), 
                                          TO_NUMBER(p_binds(2)),
                                          TO_NUMBER(p_binds(3));            
      END CASE;

      -- Fetch the ID and order data...
      -- ------------------------------
      FETCH v_cursor BULK COLLECT INTO aa_ids;
      CLOSE v_cursor;

      -- Populate the order map...
      -- -------------------------
      FOR i IN 1 .. aa_ids.COUNT LOOP      
         g_map(aa_ids(i).id) := aa_ids(i).ord;
      END LOOP;

      -- Use the map to determine padding needed to slot in our order column...
      -- ----------------------------------------------------------------------
      IF g_map.COUNT > 0 THEN
         g_len := LEAST(LENGTH(g_map.LAST) + 7, 8);
         g_pad := LPAD('-', g_len, '-');
      END IF;

   END build_order_map;

   ----------------------------------------------------------------------------
   FUNCTION prepare_row( p_curr IN VARCHAR2,
                         p_next IN VARCHAR2 ) RETURN xplan_ot IS

      v_id  PLS_INTEGER;
      v_row VARCHAR2(4000);
      v_hdr VARCHAR2(64) := '%|%Id%|%Operation%|%';

   BEGIN

      -- Intercept the plan section to include a new column for the
      -- the operation order that we mapped earlier. The plan output
      -- itself will be bound by the 2nd, 3rd and 4th dashed lines.
      -- We need to add in additional dashes, the order column heading
      -- and the order value itself...
      -- -------------------------------------------------------------

      IF p_curr LIKE '---%' THEN
  
         IF p_next LIKE v_hdr THEN
            g_hdrs := 1;
            v_row := g_pad || p_curr;
         ELSIF g_hdrs BETWEEN 1 AND 3 THEN
            g_hdrs := g_hdrs + 1;
            v_row := g_pad || p_curr;
         ELSE
            v_row := p_curr;
         END IF;

      ELSIF p_curr LIKE v_hdr THEN

         v_row := REGEXP_REPLACE(
                     p_curr, '\|',
                     RPAD('|', GREATEST(g_len-7, 2)) || 'Order |',
                     1, 2
                     ); 

      ELSIF REGEXP_LIKE(p_curr, '^\|[\* 0-9]+\|') THEN

         v_id := REGEXP_SUBSTR(p_curr, '[0-9]+');
         v_row := REGEXP_REPLACE(
                     p_curr, '\|', 
                     '|' || LPAD(g_map(v_id), GREATEST(g_len-8, 6)) || ' |',
                     1, 2
                     ); 
      ELSE
         v_row := p_curr;
      END IF;

      RETURN xplan_ot(v_row);

   END prepare_row;

   ----------------------------------------------------------------------------
   FUNCTION display_cursor( p_sql_id          IN VARCHAR2 DEFAULT NULL,
                            p_cursor_child_no IN INTEGER  DEFAULT 0,
                            p_format          IN VARCHAR2 DEFAULT 'TYPICAL' )
      RETURN xplan_ntt PIPELINED IS

      v_sql_id   v$sql_plan.sql_id%TYPE;
      v_child_no v$sql_plan.child_number%TYPE;
      v_sql      VARCHAR2(256);
      v_binds    ntt_order_map_binds := ntt_order_map_binds();

   BEGIN

      reset_state();

      -- Set a SQL_ID if default parameters passed...
      -- --------------------------------------------
      IF p_sql_id IS NULL THEN
         SELECT prev_sql_id, prev_child_number
         INTO   v_sql_id, v_child_no
         FROM   v$session
         WHERE  sid = (SELECT m.sid FROM v$mystat m WHERE ROWNUM = 1)
         AND    username IS NOT NULL 
         AND    prev_hash_value <> 0;
      ELSE
         v_sql_id := p_sql_id;
         v_child_no := p_cursor_child_no;
      END IF;

      -- Prepare the inputs for the order mapping...
      -- -------------------------------------------
      v_sql := 'SELECT id, parent_id
                FROM   v$sql_plan
                WHERE  sql_id = :bv_sql_id
                AND    child_number = :bv_child_no';

      v_binds := ntt_order_map_binds(v_sql_id, v_child_no);
      
      -- Build the plan order map from the SQL...
      -- ----------------------------------------
      build_order_map(v_sql, v_binds);

      -- Now we can call DBMS_XPLAN to output the plan...
      -- ------------------------------------------------
      FOR r_plan IN ( SELECT plan_table_output AS p
                      ,      LEAD(plan_table_output) OVER (ORDER BY ROWNUM) AS np
                      FROM   TABLE(
                                DBMS_XPLAN.DISPLAY_CURSOR(
                                   v_sql_id, v_child_no, p_format 
                                   ))
                      ORDER  BY
                             ROWNUM)
      LOOP
         IF g_map.COUNT > 0 THEN
            PIPE ROW (prepare_row(r_plan.p, r_plan.np));
         ELSE
            PIPE ROW (xplan_ot(r_plan.p));
         END IF;
      END LOOP;

      reset_state();
      RETURN;

   END display_cursor;
END xplan;
`
	sql5 := `GRANT SELECT ANY DICTIONARY TO SYSTEM`

	_, err := db.ORA.Exec(sql5)
	if err != nil {
		return fmt.Errorf("grant privilege system failed: %v", err)
	}
	_, err = db.ORA.Exec(sql1)
	if err != nil {
		return fmt.Errorf("make package type xplan_ott failed: %v", err)
	}
	_, err = db.ORA.Exec(sql2)
	if err != nil {
		return fmt.Errorf("make package type xplan_ntt failed: %v", err)
	}
	_, err = db.ORA.Exec(sql3)
	if err != nil {
		return fmt.Errorf("make package xplan failed: %v", err)
	}
	_, err = db.ORA.Exec(sql4)
	if err != nil {
		return fmt.Errorf("make package body xplan failed: %v", err)
	}
	return nil
}
