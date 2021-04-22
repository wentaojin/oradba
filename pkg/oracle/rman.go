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

func QueryOracleDBRMANProcess() error {
	columns, values, err := db.Query(`select s.inst_id,
       o.sid,
       o.serial#,
       CLIENT_INFO ch,
       context,
       sofar,
       totalwork,
       round(sofar / totalwork * 100, 2) "% Complete"
  FROM gv$session_longops o, gv$session s
 WHERE opname LIKE 'RMAN%'
   AND opname NOT LIKE '%aggregate%'
   AND o.sid = s.sid
   AND totalwork != 0
   AND sofar <> totalwork`)
	if err != nil {
		return err
	}
	util.NewMarkdownTableStyle(os.Stdout, columns, values)
	return nil
}

func QueryOracleDBRMANStatus(startTime, endTime string) error {
	columns, values, err := db.Query(fmt.Sprintf(`select s.status as "备份状态",
       trunc((s.END_TIME - s.START_TIME) * 24 * 60, 0) "备份用时(分钟)",
       to_char(s.START_TIME, 'yyyy-mm-dd hh24:mi:ss') as "开始备份时间",
       to_char(s.END_TIME, 'yyyy-mm-dd hh24:mi:ss') as "结束备份时间",
       s.OPERATION as "命令",
       trunc(s.INPUT_BYTES / 1024 / 1024 / 1204, 2) as "INPUT/G",
       trunc(s.OUTPUT_BYTES / 1024 / 1024 / 1024, 2) as "OUTPUT/G",
       s.OBJECT_TYPE as "对象类型",
       s.MBYTES_PROCESSED as "百分比",
       s.OUTPUT_DEVICE_TYPE as "设备类型"
  from v$rman_status s
 where to_char(s.START_TIME, 'yyyy-mm-dd hh24:mi:ss') < '%s'
   and to_char(s.END_TIME, 'yyyy-mm-dd hh24:mi:ss') > '%s'
 order by s.START_TIME desc`, endTime, startTime))
	if err != nil {
		return err
	}
	util.NewMarkdownTableStyle(os.Stdout, columns, values)
	return nil
}
