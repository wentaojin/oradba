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

	"github.com/WentaoJin/oradba/db"
	"github.com/WentaoJin/oradba/pkg/util"
)

/*
	根据sid查spid信息
*/
func QueryOracleDBSessionSIDBySPID(spid string) error {
	columns, values, err := db.Query(fmt.Sprintf(`select s.sid,
       s.serial#,
       s.LOGON_TIME,
       s.machine,
       s.program,
       s.module,
       p.spid,
       p.terminal
  from v$session s, v$process p
 where s.paddr = p.addr
   and s.spid = '%s'`, spid))
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}

/*
	根据spid查sid信息
*/
func QueryOracleDBSessionSPIDBySID(sid string) error {
	columns, values, err := db.Query(fmt.Sprintf(`select s.sid,
       s.serial#,
       s.LOGON_TIME,
       s.machine,
       s.program,
       s.module,
       p.spid,
       p.terminal
  from v$session s, v$process p
 where s.paddr = p.addr
   and s.sid = '%s'`, sid))
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}

/*
	根据sid查询其对应等待事件以及等待时间
*/
func QueryOracleDBWaitEventBySid(sid string) error {
	columns, values, err := db.Query(fmt.Sprintf(`SELECT event, time_waited AS time_spent
  FROM v$session_event
 WHERE sid = '%s'
   AND wait_class <> 'Idle'
   AND event NOT IN ('Null event',
                     'client message',
                     'KXFX: Execution Message Dequeue - Slave',
                     'PX Deq: Execution Msg',
                     'KXFQ: kxfqdeq - normal deqeue',
                     'PX Deq: Table Q Normal',
                     'Wait for credit - send blocked',
                     'PX Deq Credit: send blkd',
                     'Wait for credit - need buffer to send',
                     'PX Deq Credit: need buffer',
                     'Wait for credit - free buffer',
                     'PX Deq Credit: free buffer',
                     'parallel query dequeue wait',
                     'PX Deque wait',
                     'Parallel Query Idle Wait - Slaves',
                     'PX Idle Wait',
                     'slave wait',
                     'dispatcher timer',
                     'virtual circuit status',
                     'pipe get',
                     'rdbms ipc message',
                     'rdbms ipc reply',
                     'pmon timer',
                     'smon timer',
                     'PL/SQL lock timer',
                     'SQL*Net message from client',
                     'SQL*Net message to client',
                     'SQL*Net break/reset to client',
                     'SQL*Net more data to client',
                     'rdbms ipc message',
                     'WMON goes to sleep')
UNION ALL
SELECT b.name, a.VALUE
  FROM v$sesstat a, v$statname b
 WHERE a.statistic# = b.statistic#
   AND b.name = 'CPU used when call started'
   AND a.sid = '%s'`, sid, sid))
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}

/*
	查看处于等待会话的sql、等待事件等相关信息
*/
func QueryOracleDBWaitEvent() error {
	columns, values, err := db.Query(`SELECT S.USERNAME,
       S.SID,
       S.SERIAL#,
       P.SPID,
       Q.SQL_ID,
       Q.SQL_FULLTEXT,
       W.EVENT,
       W.WAIT_TIME,
       W.STATE,
       CASE
         WHEN W.STATE = 'WAITING' THEN
          W.SECONDS_IN_WAIT
         WHEN W.STATE = 'WAITING KNOWN TIME' THEN
          W.WAIT_TIME
       END AS SEC_IN_WAIT
  FROM V$SESSION S, V$SESSION_WAIT W, V$SQLAREA Q, V$PROCESS P
 WHERE S.SID = W.SID
   AND S.SQL_ID = Q.SQL_ID
   AND P.ADDR = S.PADDR
   AND W.EVENT NOT LIKE 'SQL*Net%'
   AND S.USERNAME IS NOT NULL
   AND W.WAIT_TIME >= 0
 ORDER BY W.SECONDS_IN_WAIT DESC`)
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}

/*
	查看会话使用UNDO大小 TOP 10
*/
func QueryOracleDBUsedUNDOTOP() error {
	columns, values, err := db.Query(`select *
  from (select s.sid,
               s.serial#,
               s.sql_id,
               v.usn,
               segment_name,
               r.status,
               v.rssize / 1024 / 1024 mb
          From dba_rollback_segs r,
               v$rollstat        v,
               v$transaction     t,
               v$session         s
         Where r.segment_id = v.usn
           and v.usn = t.xidusn
           and t.addr = s.taddr
         order by mb desc)
 where rownum <= 10`)
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}

/*
	查看会话使用TEMP大小TOP 10
*/
func QueryOracleDBUsedTEMPTOP() error {
	columns, values, err := db.Query(`select *
  from (select a.username,
               a.sid,
               a.serial#,
               b.tablespace,
               b.segfile#,
               b.segblk#,
               b.blocks,
               b.blocks * 32 / 1024 / 1024 as usedtempsize,
               a.osuser,
               a.status,
               c.sql_text,
               b.contents
          from v$session a, v$sort_usage b, v$sql c
         where a.saddr = b.session_addr
           and a.sql_address = c.address(+)
         order by b.blocks desc)
 where rownum <= 10`)
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}

/*
	查看会话使用PGA大小 TOP 10
*/
func QueryOracleDBUsedPGATOP() error {
	columns, values, err := db.Query(`select *
  from (Select s.sid,
               s.serial#,
               Osuser,
               Name,
               CASE
                 WHEN (Value / 1024 / 1024) LIKE '.%' THEN
                  '0' || to_char(round(Value / 1024 / 1024, 6))
                 else
                  TO_CHAR(round(Value / 1024 / 1024, 6))
               end Mb,
               NVL(s.Sql_Id, 'NULL'),
               Spid,
               s.status,
               s.machine
          From V$session s, V$sesstat St, V$statname Sn, V$process p
         Where St.Sid = s.Sid
           And St.Statistic# = Sn.Statistic#
           And Sn.Name Like 'session pga memory'
           And p.Addr = s.Paddr
         Order By Value Desc)
 where rownum <= 10`)
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}

/*
	查询会话阻塞情况(包括死锁)
*/
func QueryOracleDBSessionBlocking() error {
	columns, values, err := db.Query(`SELECT S1.USERNAME "WAITING USER",
       S1.OSUSER "OS User",
       S1.LOGON_TIME "logon time",
       W.SESSION_ID "被阻塞Sid",
       S1.SERIAL# "被阻塞Serial#",
       P1.SPID "PID",
       Q1.SQL_TEXT "SQLTEXT",
       S2.USERNAME "HOLDING User",
       S2.OSUSER "OS User",
       S2.LOGON_TIME "logon time",
       H.SESSION_ID "阻塞源头Sid",
       S2.SERIAL# "阻塞源头Serial",
       'ALTER SYSTEM KILL SESSION ' || '''' || H.SESSION_ID || ',' ||
       S2.SERIAL# || '''' || ' IMMEDIATE;' "解锁建议",
       P2.SPID "阻塞源头系统PID",
       Q2.SQL_TEXT "SQLTEXT"
  FROM SYS.V_$PROCESS P1,
       SYS.V_$PROCESS P2,
       SYS.V_$SESSION S1,
       SYS.V_$SESSION S2,
       DBA_LOCKS      W,
       DBA_LOCKS      H,
       V$SQL          Q1,
       V$SQL          Q2
 WHERE H.MODE_HELD != 'None'
   AND H.MODE_HELD != 'Null'
   AND W.MODE_REQUESTED != 'None'
   AND W.LOCK_TYPE(+) = H.LOCK_TYPE
   AND W.LOCK_ID1(+) = H.LOCK_ID1
   AND W.LOCK_ID2(+) = H.LOCK_ID2
   AND W.SESSION_ID = S1.SID(+)
   AND H.SESSION_ID = S2.SID(+)
   AND S1.PADDR = P1.ADDR(+)
   AND S2.PADDR = P2.ADDR(+)
   AND S1.SQL_ID = Q1.SQL_ID(+)
   AND S2.SQL_ID = Q2.SQL_ID(+)
 ORDER BY H.SESSION_ID`)
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}

/*
	查询会话阻塞详情包括锁模式(包括死锁)
*/
func QueryOracleDBSessionBlock() error {
	columns, values, err := db.Query(`SELECT b.type || '-' || b.id1 || '-' || b.id2 AS b_res,
       s1.sid || ',' || s1.serial# || ',@' || s1.inst_id AS blocker,
       (SELECT COUNT(*)
          FROM gv$lock t
         WHERE t.type = b.type
           AND t.id1 = b.id1
           AND t.id2 = b.id2
           AND request > 0) blocked_cnt,
       decode(b.request,
              1,
              'NULL',
              2,
              'SS',
              3,
              'SX',
              4,
              'S',
              5,
              'SSX',
              6,
              'X',
              'NONE') AS b_req,
       decode(b.lmode,
              1,
              'NULL',
              2,
              'SS',
              3,
              'SX',
              4,
              'S',
              5,
              'SSX',
              6,
              'X',
              'NONE') AS b_mode,
       s1.username,
       nvl(s1.sql_id, 'Null'),
       b.ctime AS b_ctime,
       CASE b.type
         WHEN 'TM' THEN
          (SELECT object_name FROM dba_objects WHERE object_id = b.id1)
         WHEN 'TX' THEN
          (SELECT object_name
             FROM dba_objects
            WHERE object_id = s2.row_wait_obj#)
         ELSE
          ''
       END AS tabwait,
       dbms_rowid.rowid_create(1,
                               s2.row_wait_obj#,
                               s2.row_wait_file#,
                               s2.row_wait_block#,
                               s2.row_wait_row#) AS rowait,
       s2.sid || ',' || s2.serial# || ',@' || s2.inst_id AS waiter,
       decode(w.request,
              1,
              'NULL',
              2,
              'SS',
              3,
              'SX',
              4,
              'S',
              5,
              'SSX',
              6,
              'X',
              'NONE') AS w_req,
       decode(w.lmode,
              1,
              'NULL',
              2,
              'SS',
              3,
              'SX',
              4,
              'S',
              5,
              'SSX',
              6,
              'X',
              'NONE') AS w_mode,
       s2.username,
       nvl(s2.sql_id, 'Null'),
       w.ctime AS w_ctime
  FROM gv$lock b, gv$lock w, gv$session s1, gv$session s2
 WHERE b.block > 0
   AND w.request > 0
   AND b.id1 = w.id1
   AND b.id2 = w.id2
   AND b.type = w.type
   AND b.inst_id = s1.inst_id
   AND b.sid = s1.sid
   AND w.inst_id = s2.inst_id
   AND w.sid = s2.sid
 ORDER BY b_res, w_ctime DESC`)
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}
