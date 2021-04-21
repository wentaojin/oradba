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
	"os"

	"github.com/WentaoJin/oradba/db"
	"github.com/WentaoJin/oradba/pkg/util"
)

func QueryOracleDBInfo() error {
	columns, values, err := db.Query(`select distinct gd.name,
                gi.instance_number,
                gi.version,
                gd.log_mode,
                gd.open_mode,
                gd.database_role,
                gd.switchover_status,
                gd.force_logging,
                gd.flashback_on
  from gv$database gd, gv$instance gi`)
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}

func QueryOracleInstanceListInfo() error {
	columns, values, err := db.Query(`select inst_id,
       instance_name,
       host_name,
       startup_time,
       status,
       instance_role
  from gv$instance`)
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}

func QueryOracleComponentMemoryInfo() error {
	columns, values, err := db.Query(`select name,
       total,
       round(total - free, 2) used,
       round(free, 2) free,
       round((total - free) / total * 100, 2) pctused
  from (select 'SGA' name,
               (select sum(round(value / 1024 / 1024 / 1024, 2)) from v$sga) total,
               (select sum(bytes / 1024 / 1024 / 1024)
                  from v$sgastat
                 where name = 'free memory') free
          from dual)
union
select name,
       total,
       round(used, 2) used,
       round(total - used, 2) free,
       round(used / total * 100, 2) pctused
  from (select 'PGA' name,
               (select round(value / 1024 / 1024 / 1024, 2) total
                  from v$pgastat
                 where name = 'aggregate PGA target parameter') total,
               (select value / 1024 / 1024 / 1024 used
                  from v$pgastat
                 where name = 'total PGA allocated') used
          from dual)
union
select name,
       round(total, 2) total,
       round((total - free), 2) used,
       round(free, 2) free,
       round((total - free) / total * 100, 2) pctused
  from (select 'Shared pool' name,
               (select sum(bytes / 1024 / 1024 / 1024)
                  from v$sgastat
                 where pool = 'shared pool') total,
               (select bytes / 1024 / 1024 / 1024
                  from v$sgastat
                 where name = 'free memory'
                   and pool = 'shared pool') free
          from dual)`)
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}

func QueryOracleTablespaceSummaryInfo() error {
	columns, values, err := db.Query(`select b.tbs_totalsize, a.tbs_freesize, c.tbs_nums
  from (select round(sum(bytes) / 1024 / 1024 / 1024, 2) tbs_totalsize
          from dba_data_files) b,
       (select round(nvl(sum(bytes), 0) / 1024 / 1024 / 1024) tbs_freesize
          from dba_free_space) a,
       (select count(1) tbs_nums from v$tablespace) c`)
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}

func QueryOracleLastBackupInfo() error {
	columns, values, err := db.Query(`select b.SESSION_KEY,
       b.INPUT_TYPE,
       b.STATUS,
       b.start_time,
       b.end_time,
       b.hrs
  from (SELECT SESSION_KEY,
               INPUT_TYPE,
               STATUS,
               TO_CHAR(START_TIME, 'mm/dd/yy hh24:mi') start_time,
               TO_CHAR(END_TIME, 'mm/dd/yy hh24:mi') end_time,
               round(ELAPSED_SECONDS / 3600, 2) hrs
          FROM V$RMAN_BACKUP_JOB_DETAILS
         ORDER BY END_TIME desc) b
 where rownum <= 1`)
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}

func QueryOracleParamsListInfo() error {
	columns, values, err := db.Query(`select distinct *
  from (select p.name, sp.sid, NVL(p.value, 'UNKNOWN') value
          from gv$parameter2 p, gv$spparameter sp
         where p.name = sp.name) b
 where b.value <> 'UNKNOWN'`)
	if err != nil {
		return err
	}
	util.NewTableStyle(os.Stdout, columns, values)
	return nil
}
