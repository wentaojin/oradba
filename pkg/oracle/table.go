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

	"github.com/WentaoJin/oradba/pkg/util"

	"github.com/WentaoJin/oradba/db"
)

func QueryOracleDBTableDetailInfo(username, tablename string) error {
	switch {
	case username != "" && tablename == "":
		if err := queryOracleDBTableDetailInfo(username); err != nil {
			return err
		}
		return nil
	case username != "" && tablename != "":
		if err := queryOracleDBTableDetailInfoByUser(username, tablename); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("flag username and table specify exist error")
	}
}

func QueryOracleDBTableColumnDetailInfo(username, tablename string) error {
	isExist, err := isExistUserAndTable(username, tablename)
	if !isExist {
		return fmt.Errorf("oracle db user [%s] or table [%s] not exists", username, tablename)
	}
	if err != nil {
		return err
	}
	columns, values, err := db.Query(fmt.Sprintf(`select COLUMN_NAME,
       decode(t.DATA_TYPE,
              'NUMBER',
              t.DATA_TYPE || '(' ||
              decode(t.DATA_PRECISION,
                     null,
                     t.DATA_LENGTH || ')',
                     t.DATA_PRECISION || ',' || t.DATA_SCALE || ')'),
              'DATE',
              t.DATA_TYPE,
              'LONG',
              t.DATA_TYPE,
              'LONG RAW',
              t.DATA_TYPE,
              'ROWID',
              t.DATA_TYPE,
              'MLSLABEL',
              t.DATA_TYPE,
              t.DATA_TYPE || '(' || t.DATA_LENGTH || ')') || ' ' ||
       decode(t.nullable, 'N', 'NOT NULL', 'n', 'NOT NULL', NULL) COL_TYPE,
       NUM_DISTINCT,
       DENSITY,
       NUM_BUCKETS,
       NUM_NULLS,
       histogram,
       GLOBAL_STATS,
       USER_STATS,
       SAMPLE_SIZE,
       to_char(t.last_analyzed, 'MM-DD-YYYY') last_analyzed
  from dba_tab_columns t
 where upper(owner) = upper('%s')
   and upper(table_name) = upper('%s')`, username, tablename))
	if err != nil {
		return err
	}
	util.NewMarkdownTableStyle(os.Stdout, columns, values)
	return nil
}

func QueryOracleDBTableIndexDetailInfo(username, tablename string) error {
	isExist, err := isExistUserAndTable(username, tablename)
	if !isExist {
		return fmt.Errorf("oracle db user [%s] or table [%s] not exists", username, tablename)
	}
	if err != nil {
		return err
	}
	columns, values, err := db.Query(fmt.Sprintf(`select INDEX_NAME,
       UNIQUENESS,
       BLEVEL BLev,
       LEAF_BLOCKS,
       DISTINCT_KEYS,
       NUM_ROWS,
       CLUSTERING_FACTOR,
       GLOBAL_STATS,
       USER_STATS,
       SAMPLE_SIZE,
       to_char(t.last_analyzed, 'YYYY-MM-DD HH24:MM:SS') last_analyzed
  from dba_indexes t
 where upper(table_owner) = upper('%s')
   and upper(table_name) = upper('%s')`, username, tablename))
	if err != nil {
		return err
	}
	util.NewMarkdownTableStyle(os.Stdout, columns, values)
	return nil
}

func QueryOracleDBTableIndexColumnDetailInfo(username, tablename string) error {
	isExist, err := isExistUserAndTable(username, tablename)
	if !isExist {
		return fmt.Errorf("oracle db user [%s] or table [%s] not exists", username, tablename)
	}
	if err != nil {
		return err
	}
	columns, values, err := db.Query(fmt.Sprintf(`select T.TABLE_NAME,
       T.INDEX_NAME,
       I.INDEX_TYPE,
       I.UNIQUENESS, --是否唯一索引
       --T.COLUMN_POSITION,
       LISTAGG(T.COLUMN_NAME, ',') WITHIN GROUP(ORDER BY T.COLUMN_POSITION) AS column_list
  FROM ALL_IND_COLUMNS T, ALL_INDEXES I, ALL_CONSTRAINTS C
 WHERE T.INDEX_NAME = I.INDEX_NAME
   AND T.INDEX_NAME = C.CONSTRAINT_NAME(+)
   -- AND C.CONSTRAINT_TYPE is Null --排除主键、唯一约束索引
   AND T.TABLE_OWNER = upper('%s')
   AND T.TABLE_NAME = upper('%s')
 group by T.TABLE_NAME,
          I.UNIQUENESS, --是否唯一索引
          T.INDEX_NAME,
          I.INDEX_TYPE`, username, tablename))
	if err != nil {
		return err
	}
	util.NewMarkdownTableStyle(os.Stdout, columns, values)
	return nil
}

func QueryOracleDBTableForeignKeyDetailInfo(username, tablename string) error {
	isExist, err := isExistUserAndTable(username, tablename)
	if !isExist {
		return fmt.Errorf("oracle db user [%s] or table [%s] not exists", username, tablename)
	}
	if err != nil {
		return err
	}
	columns, values, err := db.Query(fmt.Sprintf(`with temp1 as
 (select t1.r_owner,
         t1.constraint_name,
         t1.r_constraint_name,
         LISTAGG(a1.column_name, ',') WITHIN GROUP(ORDER BY a1.POSITION) AS COLUMN_LIST
    from all_constraints t1, all_cons_columns a1
   where t1.constraint_name = a1.constraint_name
     AND upper(t1.owner) = upper('%s')
     AND upper(t1.table_name) = upper('%s')
     AND t1.STATUS = 'ENABLED'
     AND t1.Constraint_Type = 'R'
   group by t1.r_owner, t1.constraint_name, t1.r_constraint_name),
temp2 as
 (select t1.owner,
         t1.constraint_name,
         LISTAGG(a1.column_name, ',') WITHIN GROUP(ORDER BY a1.POSITION) AS COLUMN_LIST
    from all_constraints t1, all_cons_columns a1
   where t1.constraint_name = a1.constraint_name
     AND upper(t1.owner) = upper('%s')
     AND t1.STATUS = 'ENABLED'
     AND t1.Constraint_Type = 'P'
   group by t1.owner, t1.r_owner, t1.constraint_name)
select x.constraint_name,
       x.COLUMN_LIST,
       x.r_owner,
       x.r_constraint_name as RCONSTRAINT_NAME,
       y.COLUMN_LIST       as RCOLUMN_LIST
  from temp1 x, temp2 y
 where x.r_owner = y.owner
   and x.r_constraint_name = y.constraint_name`, username, tablename, username))
	if err != nil {
		return err
	}
	util.NewMarkdownTableStyle(os.Stdout, columns, values)
	return nil
}

func QueryOracleDBTableCheckKeyDetailInfo(username, tablename string) error {
	isExist, err := isExistUserAndTable(username, tablename)
	if !isExist {
		return fmt.Errorf("oracle db user [%s] or table [%s] not exists", username, tablename)
	}
	if err != nil {
		return err
	}
	columns, values, err := db.Query(fmt.Sprintf(`select cu.constraint_name, SEARCH_CONDITION
  from all_cons_columns cu, all_constraints au
 where cu.constraint_name = au.constraint_name
   and au.constraint_type = 'C'
   and au.STATUS = 'ENABLED'
   and upper(au.table_name) = upper('%s')
   and upper(cu.owner) = upper('%s')`, tablename, username))
	if err != nil {
		return err
	}
	util.NewMarkdownTableStyle(os.Stdout, columns, values)
	return nil
}

func QueryOracleDBTableUniqueKeyDetailInfo(username, tablename string) error {
	isExist, err := isExistUserAndTable(username, tablename)
	if !isExist {
		return fmt.Errorf("oracle db user [%s] or table [%s] not exists", username, tablename)
	}
	if err != nil {
		return err
	}
	columns, values, err := db.Query(fmt.Sprintf(`select cu.constraint_name,
       LISTAGG(cu.column_name, ',') WITHIN GROUP(ORDER BY cu.POSITION) AS column_list
  from all_cons_columns cu, all_constraints au
 where cu.constraint_name = au.constraint_name
   and au.constraint_type = 'U'
   and au.STATUS = 'ENABLED'
   and upper(au.table_name) = upper('%s')
   and upper(cu.owner) = upper('%s')
 group by cu.constraint_name`, tablename, username))
	if err != nil {
		return err
	}
	util.NewMarkdownTableStyle(os.Stdout, columns, values)
	return nil
}

func QueryOracleDBTablePrimaryKeyDetailInfo(username, tablename string) error {
	isExist, err := isExistUserAndTable(username, tablename)
	if !isExist {
		return fmt.Errorf("oracle db user [%s] or table [%s] not exists", username, tablename)
	}
	if err != nil {
		return err
	}
	// for the primary key of an Engine table, you can use the following command to set whether the primary key takes effect.
	// disable the primary key: alter table tableName disable primary key;
	// enable the primary key: alter table tableName enable primary key;
	// primary key status Disabled will not do primary key processing
	columns, values, err := db.Query(fmt.Sprintf(`select cu.constraint_name,
       LISTAGG(cu.column_name, ',') WITHIN GROUP(ORDER BY cu.POSITION) AS COLUMN_LIST
  from all_cons_columns cu, all_constraints au
 where cu.constraint_name = au.constraint_name
   and au.constraint_type = 'P'
   and au.STATUS = 'ENABLED'
   and upper(au.table_name) = upper('%s')
   and upper(cu.owner) = upper('%s')
 group by cu.constraint_name`, tablename, username))
	if err != nil {
		return err
	}
	util.NewMarkdownTableStyle(os.Stdout, columns, values)
	return nil
}

func queryOracleDBTableDetailInfo(username string) error {
	isExist, err := isExistUser(username)
	if !isExist {
		return fmt.Errorf("oracle db user [%s] not exists", username)
	}
	if err != nil {
		return err
	}
	columns, values, err := db.Query(fmt.Sprintf(`select owner,
       table_name,
       num_rows,
       blocks,
       avg_row_len,
       TEMPORARY,
       GLOBAL_STATS,
       SAMPLE_SIZE,
       degree,
       to_char(LAST_ANALYZED, 'YYYY-MM-DD HH24:MI:SS') LAST_ANALYZED
  from dba_tables
 where owner = upper('%s')`, username))
	if err != nil {
		return err
	}
	util.NewMarkdownTableStyle(os.Stdout, columns, values)
	return nil
}

func queryOracleDBTableDetailInfoByUser(username, tablename string) error {
	isExist, err := isExistUserAndTable(username, tablename)
	if !isExist {
		return fmt.Errorf("oracle db user [%s] or table [%s] not exists", username, tablename)
	}
	if err != nil {
		return err
	}
	columns, values, err := db.Query(fmt.Sprintf(`select owner,
       table_name,
       num_rows,
       blocks,
       avg_row_len,
       TEMPORARY,
       GLOBAL_STATS,
       SAMPLE_SIZE,
       degree,
       to_char(LAST_ANALYZED, 'YYYY-MM-DD HH24:MI:SS') LAST_ANALYZED
  from dba_tables
 where owner = upper('%s')
   and table_name = upper('%s')`, username, tablename))
	if err != nil {
		return err
	}
	util.NewMarkdownTableStyle(os.Stdout, columns, values)
	return nil
}

func isExistUser(username string) (bool, error) {
	_, vals, err := db.Query(fmt.Sprintf(`select count(*) AS COUNT
  from dba_tables
 where owner = upper('%s')`, username))
	if err != nil {
		return false, err
	}
	if vals[0][0] == "0" {
		return false, nil
	}
	return true, nil
}

func isExistUserAndTable(username, tablename string) (bool, error) {
	_, vals, err := db.Query(fmt.Sprintf(`select count(*) AS COUNT
  from dba_tables
 where owner = upper('%s')
   and table_name = upper('%s')`, username, tablename))
	if err != nil {
		return false, err
	}
	if vals[0][0] == "0" {
		return false, nil
	}
	return true, nil
}
