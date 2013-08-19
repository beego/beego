package orm

import (
	"fmt"
	"os"
	"strings"
)

func getDbAlias(name string) *alias {
	if al, ok := dataBaseCache.get(name); ok {
		return al
	} else {
		fmt.Println(fmt.Sprintf("unknown DataBase alias name %s", name))
		os.Exit(2)
	}

	return nil
}

func getDbDropSql(al *alias) (sqls []string) {
	if len(modelCache.cache) == 0 {
		fmt.Println("no Model found, need register your model")
		os.Exit(2)
	}

	Q := al.DbBaser.TableQuote()

	for _, mi := range modelCache.allOrdered() {
		sqls = append(sqls, fmt.Sprintf(`DROP TABLE IF EXISTS %s%s%s`, Q, mi.table, Q))
	}
	return sqls
}

func getDbCreateSql(al *alias) (sqls []string) {
	if len(modelCache.cache) == 0 {
		fmt.Println("no Model found, need register your model")
		os.Exit(2)
	}

	Q := al.DbBaser.TableQuote()
	T := al.DbBaser.DbTypes()

	for _, mi := range modelCache.allOrdered() {
		sql := fmt.Sprintf("-- %s\n", strings.Repeat("-", 50))
		sql += fmt.Sprintf("--  Table Structure for `%s`\n", mi.fullName)
		sql += fmt.Sprintf("-- %s\n", strings.Repeat("-", 50))

		sql += fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s%s%s (\n", Q, mi.table, Q)

		columns := make([]string, 0, len(mi.fields.fieldsDB))

		for _, fi := range mi.fields.fieldsDB {

			fieldType := fi.fieldType
			column := fmt.Sprintf("    %s%s%s ", Q, fi.column, Q)
			col := ""

		checkColumn:
			switch fieldType {
			case TypeBooleanField:
				col = T["bool"]
			case TypeCharField:
				col = fmt.Sprintf(T["string"], fi.size)
			case TypeTextField:
				col = T["string-text"]
			case TypeDateField:
				col = T["time.Time-date"]
			case TypeDateTimeField:
				col = T["time.Time"]
			case TypeBitField:
				col = T["int8"]
			case TypeSmallIntegerField:
				col = T["int16"]
			case TypeIntegerField:
				col = T["int32"]
			case TypeBigIntegerField:
				if al.Driver == DR_Sqlite {
					fieldType = TypeIntegerField
					goto checkColumn
				}
				col = T["int64"]
			case TypePositiveBitField:
				col = T["uint8"]
			case TypePositiveSmallIntegerField:
				col = T["uint16"]
			case TypePositiveIntegerField:
				col = T["uint32"]
			case TypePositiveBigIntegerField:
				col = T["uint64"]
			case TypeFloatField:
				col = T["float64"]
			case TypeDecimalField:
				s := T["float64-decimal"]
				if strings.Index(s, "%d") == -1 {
					col = s
				} else {
					col = fmt.Sprintf(s, fi.digits, fi.decimals)
				}
			case RelForeignKey, RelOneToOne:
				fieldType = fi.relModelInfo.fields.pk.fieldType
				goto checkColumn
			}

			if fi.auto {
				if al.Driver == DR_Postgres {
					column += T["auto"]
				} else {
					column += col + " " + T["auto"]
				}
			} else if fi.pk {
				column += col + " " + T["pk"]
			} else {
				column += col

				if fi.null == false {
					column += " " + "NOT NULL"
				}

				if fi.unique {
					column += " " + "UNIQUE"
				}
			}

			if strings.Index(column, "%COL%") != -1 {
				column = strings.Replace(column, "%COL%", fi.column, -1)
			}

			columns = append(columns, column)
		}

		sql += strings.Join(columns, ",\n")
		sql += "\n)"

		if al.Driver == DR_MySQL {
			sql += " ENGINE=INNODB"
		}

		sqls = append(sqls, sql)
	}

	return sqls
}
