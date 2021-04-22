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

func QueryOracleDBTablespaceDetailInfo() error {
	columns, values, err := db.Query(`SELECT a.tablespace_name,
       Round((total - free) / maxsize * 100, 1) USED_PCT,
       b.autoextensible,
       total total_mb,
       (total - free) USED,
       free,
       b.cnt DATAFILE_COUNT,
       c.status,
       c.CONTENTS,
       c.extent_management,
       c.allocation_type,
       b.maxsize
  FROM (SELECT tablespace_name, Round(SUM(bytes) / (1024 * 1024), 1) free
          FROM dba_free_space
         GROUP BY tablespace_name) a,
       (SELECT tablespace_name,
               Round(SUM(bytes) / (1024 * 1024), 1) total,
               Count(*) cnt,
               Max(autoextensible) autoextensible,
               sum(decode(autoextensible,
                          'YES',
                          floor(maxbytes / 1048576),
                          floor(bytes / 1048576))) maxsize
          FROM dba_data_files
         GROUP BY tablespace_name) b,
       dba_tablespaces c
 WHERE a.tablespace_name = b.tablespace_name
   AND a.tablespace_name = c.tablespace_name
UNION ALL
SELECT /*+ NO_MERGE */
 a.tablespace_name,
 Round(100 * (B.tot_gbbytes_used / A.maxsize), 1) PERC_USED,
 a.aet,
 Round(A.avail_size_gb, 1),
 Round(B.tot_gbbytes_used, 1),
 (Round(A.avail_size_gb, 1) - Round(B.tot_gbbytes_used, 1)),
 a.cnt DATAFILE_COUNT,
 c.status,
 c.CONTENTS,
 c.extent_management,
 c.allocation_type,
 a.maxsize
  FROM (SELECT tablespace_name,
               SUM(bytes) / Power(2, 20) AVAIL_SIZE_GB,
               Max(autoextensible) aet,
               Count(*) cnt,
               sum(decode(autoextensible,
                          'YES',
                          floor(maxbytes / 1048576),
                          floor(bytes / 1048576))) maxsize
          FROM dba_temp_files
         GROUP BY tablespace_name) A,
       (SELECT tablespace_name,
               SUM(bytes_used) / Power(2, 20) TOT_GBBYTES_USED
          FROM gv$temp_extent_pool
         GROUP BY tablespace_name) B,
       dba_tablespaces c
 WHERE a.tablespace_name = b.tablespace_name
   AND a.tablespace_name = c.tablespace_name`)
	if err != nil {
		return err
	}
	util.NewMarkdownTableStyle(os.Stdout, columns, values)
	return nil
}

func QueryOracleDBTablespaceIOStatInfo() error {
	columns, values, err := db.Query(`select df.tablespace_name name,
       df.file_name       "file",
       f.phyrds           phyrds,
       f.phyblkrd         phyblkrd,
       f.phywrts          phywrts,
       f.phyblkwrt        phyblkwrt
  from v$filestat f, dba_data_files df
 where f.file# = df.file_id
 order by df.tablespace_name`)
	if err != nil {
		return err
	}
	util.NewMarkdownTableStyle(os.Stdout, columns, values)
	return nil
}

func QueryOracleDBTablespaceDetailInfoByUser(username string) error {
	columns, values, err := db.Query(fmt.Sprintf(`select tablespacename,
       to_char(sum(totalContent)) total_Size,
       to_char(sum(usecontent)) used_Size,
       to_char(sum(sparecontent)) free_Size,
       Round(100 - avg(sparepercent),2) || '%%' free_Percent
  from (SELECT b.file_id as id,
               b.tablespace_name as tablespacename,
               Round(b.bytes / 1024 / 1024 / 1024, 2) as totalContent,
               (Round(b.bytes / 1024 / 1024 / 1024, 2)) -
               sum(nvl(Round(a.bytes / 1024 / 1024 / 1024, 2), 0)) as usecontent,
               sum(nvl(Round(a.bytes / 1024 / 1024 / 1024, 2), 0)) as sparecontent,
               sum(nvl(Round(a.bytes / 1024 / 1024 / 1024, 2), 0)) /
               (Round(b.bytes / 1024 / 1024 / 1024, 2)) * 100 as sparepercent
          FROM dba_free_space a, dba_data_files b
         WHERE a.file_id = b.file_id
           and b.tablespace_name in
               (select default_tablespace
                  from dba_users
                 where lower(username) = lower('%s'))
         group by b.tablespace_name,
                  b.file_name,
                  b.file_id,
                  Round(b.bytes / 1024 / 1024 / 1024, 2))
 GROUP BY tablespacename`, username))
	if err != nil {
		return err
	}
	util.NewMarkdownTableStyle(os.Stdout, columns, values)
	return nil
}

func QueryOracleDBSysauxTablespaceInfo() error {
	columns, values, err := db.Query(`SELECT occupant_name "Item",
       case
         when space_usage_kbytes / 1048576 LIKE '.%' then
          '0' || to_char(round(space_usage_kbytes / 1048576, 6))
         else
          to_char(space_usage_kbytes / 1048576)
       end "Space Used (GB)",
       schema_name "Schema",
       NVL(move_procedure, 'UNKNOWN') "Move Procedure"
  FROM v$sysaux_occupants
 ORDER BY "Space Used (GB)" DESC`)
	if err != nil {
		return err
	}
	util.NewMarkdownTableStyle(os.Stdout, columns, values)
	return nil
}
