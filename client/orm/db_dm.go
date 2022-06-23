package orm

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// dm operators.
var dmOperators = map[string]string{
	"exact":       "= ?",
	"iexact":      "LIKE ? ESCAPE '\\'",
	"contains":    "LIKE ? ESCAPE '\\'",
	"icontains":   "LIKE ? ESCAPE '\\'",
	"gt":          "> ?",
	"gte":         ">= ?",
	"lt":          "< ?",
	"lte":         "<= ?",
	"eq":          "= ?",
	"ne":          "!= ?",
	"startswith":  "LIKE ? ESCAPE '\\'",
	"endswith":    "LIKE ? ESCAPE '\\'",
	"istartswith": "LIKE ? ESCAPE '\\'",
	"iendswith":   "LIKE ? ESCAPE '\\'",
}

// dm column types.
var dmTypes = map[string]string{
	"auto":                "IDENTITY(1,1)",
	"pk":                  "NOT NULL PRIMARY KEY",
	"bool":                "BIT",
	"string":              "VARCHAR(%d)",
	"string-char":         "character(%d)",
	"string-text":         "TEXT",
	"time.Time-date":      "DATE",
	"time.Time":           "TIME",
	"time.Time-precision": "TIMESTAMP",
	"int8":                "TINYINT",
	"int16":               "SMALLINT",
	"int32":               "INTEGER",
	"int64":               "BIGINT",
	"uint8":               "TINYINT unsigned",
	"uint16":              "SMALLINT unsigned",
	"uint32":              "INTEGER unsigned",
	"uint64":              "BIGINT unsigned",
	"float64":             "REAL",
	"float64-decimal":     "DECIMAL",
}

type dbBaseDM struct {
	dbBase
}

var _ dbBaser = new(dbBaseDM)

// override base db read for update behavior as DM does not support syntax
func (d *dbBaseDM) Read(ctx context.Context, q dbQuerier, mi *modelInfo, ind reflect.Value, tz *time.Location, cols []string, isForUpdate bool) error {
	if isForUpdate {
		DebugLog.Println("[WARN] DM does not support SELECT FOR UPDATE query, isForUpdate param is ignored and always as false to do the work")
	}
	return d.dbBase.Read(ctx, q, mi, ind, tz, cols, false)
}

// OperatorSQL get dm operator.
func (d *dbBaseDM) OperatorSQL(operator string) string {
	return dmOperators[operator]
}

// DbTypes get mysql table field types.
func (d *dbBaseDM) DbTypes() map[string]string {
	return dmTypes
}

// MaxLimit max int in dm.
func (d *dbBaseDM) MaxLimit() uint64 {
	return 9223372036854775807
}

// ShowTablesQuery get show tables sql in dm.
func (d *dbBaseDM) ShowTablesQuery() string {
	return "SELECT TABLE_NAME FROM USER_TABLES WHERE TABLESPACE_NAME != 'TEMP'"
}

// ShowColumnsQuery show columns sql of table for dm.
func (d *dbBaseDM) ShowColumnsQuery(table string) string {
	return fmt.Sprintf("SELECT COLUMN_NAME,DATA_TYPE,NULLABLE FROM USER_TAB_COLUMNS WHERE TABLE_NAME = '%s'", strings.ToUpper(table))
}

// IndexExists execute sql to check index exist.
func (d *dbBaseDM) IndexExists(ctx context.Context, db dbQuerier, table string, name string) bool {
	row := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM USER_IND_COLUMNS WHERE TABLE_NAME = ? AND INDEX_NAME = ?", table, name)
	var cnt int
	row.Scan(&cnt)
	return cnt > 0
}

func (d *dbBaseDM) TableQuote() string {
	return "\""
}

func (d *dbBaseDM) InsertOrUpdate(ctx context.Context, q dbQuerier, mi *modelInfo, ind reflect.Value, a *alias, args ...string) (int64, error) {
	argsMap := map[string]string{}

	// Get on the key-value pairs
	for _, v := range args {
		kv := strings.Split(v, "=")
		if len(kv) == 2 {
			argsMap[strings.ToLower(kv[0])] = kv[1]
		}
	}

	isMulti := false
	names := make([]string, 0, len(mi.fields.dbcols)-1)
	Q := d.ins.TableQuote()
	values, _, err := d.collectValues(mi, ind, mi.fields.dbcols, false, true, &names, a.TZ)
	if err != nil {
		return 0, err
	}

	marks := make([]string, 0, len(names))
	duals := make([]string, len(names))
	updateValues := make([]interface{}, 0)
	updates := make([]string, 0, len(names))
	unnames := make([]string, 0, len(names))
	pk := mi.fields.pk.column

	for i, v := range names {

		if v != pk {
			marks = append(marks, "T2."+v)
			updates = append(updates, "T1."+v+" = T2."+v)
			unnames = append(unnames, v)
		}

		valueStr := argsMap[strings.ToLower(v)]
		if valueStr != "" {
			duals[i] = fmt.Sprintf("'%s' %s", valueStr, v)
		} else {
			duals[i] = "? " + v
			updateValues = append(updateValues, values[i])
		}
	}

	sep := fmt.Sprintf("%s, %s", Q, Q)
	qmarks := strings.Join(marks, ", ")
	qduals := strings.Join(duals, ",")
	qupdates := strings.Join(updates, ", ")
	columns := strings.Join(unnames, sep)

	multi := len(values) / len(names)

	if isMulti {
		qmarks = strings.Repeat(qmarks+"), (", multi-1) + qmarks
	}

	query := fmt.Sprintf("MERGE INTO %s%s%s T1 USING (SELECT %s FROM dual) T2 ON(T1.%s = T2.%s)"+
		" WHEN NOT MATCHED THEN INSERT (%s%s%s) VALUES (%s) "+
		" WHEN MATCHED THEN UPDATE SET "+qupdates, Q, mi.table, Q, qduals, pk, pk, Q, columns, Q, qmarks)

	d.ins.ReplaceMarks(&query)

	if isMulti || !d.ins.HasReturningID(mi, &query) {
		res, err := q.ExecContext(ctx, query, updateValues...)
		if err == nil {
			if isMulti {
				return res.RowsAffected()
			}

			lastInsertId, err := res.LastInsertId()
			if err != nil {
				DebugLog.Println(ErrLastInsertIdUnavailable, ':', err)
				return lastInsertId, ErrLastInsertIdUnavailable
			} else {
				return lastInsertId, nil
			}
		}
		return 0, err
	}

	row := q.QueryRowContext(ctx, query, updateValues...)
	var id int64
	err = row.Scan(&id)
	return id, err
}

// create new DM dbBaser.
func newdbBaseDM() dbBaser {
	b := new(dbBaseDM)
	b.ins = b
	return b
}
