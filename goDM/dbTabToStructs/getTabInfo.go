package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

/**
 自动将表生成 model结构，
通过创建数据库连接，读取数据库的所有表并对所有的表元数据封装转化实体结构体
*/

type SchemaMeta struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default interface{}
	Extra   string
}

var tabToStru = make(map[string]interface{}, 20)

func main() {
	dbString := "root:Mysql123..@tcp(192.168.1.198:3306)/ipa_manage_qd"
	db, _ := sql.Open("mysql", dbString)

	initTabFromDB("am_rule", db)
	for _, v := range tabToStru {
		fmt.Println(v)
	}
}

func initTabFromDB(tableName string, db *sql.DB) {
	tables := getTables(db)
	for _, table := range tables {
		if table == tableName {
			metas := getTableInfo(table, db)
			result := changeMetas(table, metas)
			tabToStru[tableName] = result
			// fmt.Println(result)
		}

	}

}

func getTables(db *sql.DB) []string {
	var tables []string
	res, _ := db.Query("SHOW TABLES")
	for res.Next() {
		var table string
		res.Scan(&table)
		tables = append(tables, table)
	}
	return tables
}

func getTableInfo(tableName string, db *sql.DB) (metas []SchemaMeta) {
	list, _ := db.Query(fmt.Sprintf("show columns from %s", tableName))

	for list.Next() {
		var data SchemaMeta
		err := list.Scan(&data.Field, &data.Type, &data.Null, &data.Key, &data.Default, &data.Extra)
		if err != nil {
			fmt.Println(err.Error())
		}
		metas = append(metas, data)
	}
	return metas
}

func changeMetas(tableName string, metas []SchemaMeta) string {
	var modelStr string
	var modelStr_cp string
	for _, val := range metas {
		dataType := "interface{}"
		dataType_cp := dataType
		// fmt.Println(val.Type)
		if val.Type == "int" {
			dataType = "int"
			dataType_cp = dataType
		} else if strings.Contains(val.Type, "varchar") {
			dataType = "string"
			dataType_cp = "sql.NullString"
		} else if strings.Contains(val.Type, "tinyint") {
			dataType = "int8"
			dataType_cp = dataType
		} else if strings.Contains(val.Type, "datetime") || strings.Contains(val.Type, "timestamp") {
			dataType = "time.Time"
			dataType_cp = "sql.NullTime"
		} else if val.Type == "text" || val.Type == "longtext" {
			dataType = "*dm.DmClob"
			dataType_cp = dataType
		} else if val.Type == "mediumtext" || val.Type == "blob" || val.Type == "longblob" {
			dataType = "*dm.DmBlob"
		} else if strings.Contains(val.Type, "decimal") || strings.Contains(val.Type, "double") || strings.Contains(val.Type, "float") {
			dataType = "float64"
			dataType_cp = dataType
		} else if val.Type == "bigint" {
			dataType = "int64"
			dataType_cp = dataType
		}
		// if dataType == "interface{}" {
		// 	fmt.Printf("未匹配数据类型table: %v , cloumnType: %v ,\n", tableName, val.Type)
		// }
		field := val.Field
		field = strings.ToUpper(field[:]) //+ field[:]
		modelStr += fmt.Sprintf("%s %s\n", field, dataType)
		modelStr_cp += fmt.Sprintf("%s %s\n", field, dataType_cp)
	}
	tableName = strings.ToUpper(tableName[:1]) + tableName[1:]
	return fmt.Sprintf("type %s struct {\n%s}\n", "out"+tableName, modelStr) + fmt.Sprintf("type %s struct {\n%s}\n", tableName, modelStr_cp)
}
