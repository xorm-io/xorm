// Copyright 2019 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dialects

import (
	"fmt"
	"strings"
	"time"

	"xorm.io/xorm/core"
	"xorm.io/xorm/log"
	"xorm.io/xorm/schemas"
)

type DBType string

type URI struct {
	DBType  DBType
	Proto   string
	Host    string
	Port    string
	DBName  string
	User    string
	Passwd  string
	Charset string
	Laddr   string
	Raddr   string
	Timeout time.Duration
	Schema  string
}

// a dialect is a driver's wrapper
type Dialect interface {
	SetLogger(logger log.Logger)
	Init(*core.DB, *URI, string, string) error
	URI() *URI
	DB() *core.DB
	DBType() DBType
	SQLType(*schemas.Column) string
	FormatBytes(b []byte) string

	DriverName() string
	DataSourceName() string

	IsReserved(string) bool
	Quote(string) string

	AndStr() string
	OrStr() string
	EqStr() string
	RollBackStr() string
	AutoIncrStr() string

	SupportInsertMany() bool
	SupportEngine() bool
	SupportCharset() bool
	SupportDropIfExists() bool
	IndexOnTable() bool
	ShowCreateNull() bool

	IndexCheckSQL(tableName, idxName string) (string, []interface{})
	TableCheckSQL(tableName string) (string, []interface{})

	IsColumnExist(tableName string, colName string) (bool, error)

	CreateTableSQL(table *schemas.Table, tableName, storeEngine, charset string) string
	DropTableSQL(tableName string) string
	CreateIndexSQL(tableName string, index *schemas.Index) string
	DropIndexSQL(tableName string, index *schemas.Index) string

	ModifyColumnSQL(tableName string, col *schemas.Column) string

	ForUpdateSQL(query string) string

	// CreateTableIfNotExists(table *Table, tableName, storeEngine, charset string) error
	// MustDropTable(tableName string) error

	GetColumns(tableName string) ([]string, map[string]*schemas.Column, error)
	GetTables() ([]*schemas.Table, error)
	GetIndexes(tableName string) (map[string]*schemas.Index, error)

	Filters() []Filter
	SetParams(params map[string]string)
}

func OpenDialect(dialect Dialect) (*core.DB, error) {
	return core.Open(dialect.DriverName(), dialect.DataSourceName())
}

// Base represents a basic dialect and all real dialects could embed this struct
type Base struct {
	db             *core.DB
	dialect        Dialect
	driverName     string
	dataSourceName string
	logger         log.Logger
	uri            *URI
}

// String generate column description string according dialect
func String(d Dialect, col *schemas.Column) string {
	sql := d.Quote(col.Name) + " "

	sql += d.SQLType(col) + " "

	if col.IsPrimaryKey {
		sql += "PRIMARY KEY "
		if col.IsAutoIncrement {
			sql += d.AutoIncrStr() + " "
		}
	}

	if col.Default != "" {
		sql += "DEFAULT " + col.Default + " "
	}

	if d.ShowCreateNull() {
		if col.Nullable {
			sql += "NULL "
		} else {
			sql += "NOT NULL "
		}
	}

	return sql
}

// StringNoPk generate column description string according dialect without primary keys
func StringNoPk(d Dialect, col *schemas.Column) string {
	sql := d.Quote(col.Name) + " "

	sql += d.SQLType(col) + " "

	if col.Default != "" {
		sql += "DEFAULT " + col.Default + " "
	}

	if d.ShowCreateNull() {
		if col.Nullable {
			sql += "NULL "
		} else {
			sql += "NOT NULL "
		}
	}

	return sql
}

func (b *Base) DB() *core.DB {
	return b.db
}

func (b *Base) SetLogger(logger log.Logger) {
	b.logger = logger
}

func (b *Base) Init(db *core.DB, dialect Dialect, uri *URI, drivername, dataSourceName string) error {
	b.db, b.dialect, b.uri = db, dialect, uri
	b.driverName, b.dataSourceName = drivername, dataSourceName
	return nil
}

func (b *Base) URI() *URI {
	return b.uri
}

func (b *Base) DBType() DBType {
	return b.uri.DBType
}

func (b *Base) FormatBytes(bs []byte) string {
	return fmt.Sprintf("0x%x", bs)
}

func (b *Base) DriverName() string {
	return b.driverName
}

func (b *Base) ShowCreateNull() bool {
	return true
}

func (b *Base) DataSourceName() string {
	return b.dataSourceName
}

func (b *Base) AndStr() string {
	return "AND"
}

func (b *Base) OrStr() string {
	return "OR"
}

func (b *Base) EqStr() string {
	return "="
}

func (db *Base) RollBackStr() string {
	return "ROLL BACK"
}

func (db *Base) SupportDropIfExists() bool {
	return true
}

func (db *Base) DropTableSQL(tableName string) string {
	quote := db.dialect.Quote
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", quote(tableName))
}

func (db *Base) HasRecords(query string, args ...interface{}) (bool, error) {
	db.LogSQL(query, args)
	rows, err := db.DB().Query(query, args...)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func (db *Base) IsColumnExist(tableName, colName string) (bool, error) {
	query := fmt.Sprintf(
		"SELECT %v FROM %v.%v WHERE %v = ? AND %v = ? AND %v = ?",
		db.dialect.Quote("COLUMN_NAME"),
		db.dialect.Quote("INFORMATION_SCHEMA"),
		db.dialect.Quote("COLUMNS"),
		db.dialect.Quote("TABLE_SCHEMA"),
		db.dialect.Quote("TABLE_NAME"),
		db.dialect.Quote("COLUMN_NAME"),
	)
	return db.HasRecords(query, db.uri.DBName, tableName, colName)
}

/*
func (db *Base) CreateTableIfNotExists(table *Table, tableName, storeEngine, charset string) error {
	sql, args := db.dialect.TableCheckSQL(tableName)
	rows, err := db.DB().Query(sql, args...)
	if db.Logger != nil {
		db.Logger.Info("[sql]", sql, args)
	}
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return nil
	}

	sql = db.dialect.CreateTableSQL(table, tableName, storeEngine, charset)
	_, err = db.DB().Exec(sql)
	if db.Logger != nil {
		db.Logger.Info("[sql]", sql)
	}
	return err
}*/

func (db *Base) CreateIndexSQL(tableName string, index *schemas.Index) string {
	quotes := db.dialect.Quote("")
	quote := db.dialect.Quote
	var unique string
	var idxName string
	if index.Type == schemas.UniqueType {
		unique = " UNIQUE"
	}
	idxName = index.XName(tableName)
	return fmt.Sprintf("CREATE%s INDEX %v ON %v (%v)", unique,
		quote(idxName), quote(tableName),
		quote(strings.Join(index.Cols, fmt.Sprintf("%c,%c", quotes[1], quotes[0]))))
}

func (db *Base) DropIndexSQL(tableName string, index *schemas.Index) string {
	quote := db.dialect.Quote
	var name string
	if index.IsRegular {
		name = index.XName(tableName)
	} else {
		name = index.Name
	}
	return fmt.Sprintf("DROP INDEX %v ON %s", quote(name), quote(tableName))
}

func (db *Base) ModifyColumnSQL(tableName string, col *schemas.Column) string {
	return fmt.Sprintf("alter table %s MODIFY COLUMN %s", tableName, StringNoPk(db.dialect, col))
}

func (b *Base) CreateTableSQL(table *schemas.Table, tableName, storeEngine, charset string) string {
	var sql string
	sql = "CREATE TABLE IF NOT EXISTS "
	if tableName == "" {
		tableName = table.Name
	}

	sql += b.dialect.Quote(tableName)
	sql += " ("

	quotes := b.dialect.Quote("")

	if len(table.ColumnsSeq()) > 0 {
		pkList := table.PrimaryKeys

		for _, colName := range table.ColumnsSeq() {
			col := table.GetColumn(colName)
			if col.IsPrimaryKey && len(pkList) == 1 {
				sql += String(b.dialect, col)
			} else {
				sql += StringNoPk(b.dialect, col)
			}
			sql = strings.TrimSpace(sql)
			if b.DriverName() == schemas.MYSQL && len(col.Comment) > 0 {
				sql += " COMMENT '" + col.Comment + "'"
			}
			sql += ", "
		}

		if len(pkList) > 1 {
			sql += "PRIMARY KEY ( "
			sql += b.dialect.Quote(strings.Join(pkList, fmt.Sprintf("%c,%c", quotes[1], quotes[0])))
			sql += " ), "
		}

		sql = sql[:len(sql)-2]
	}
	sql += ")"

	if b.dialect.SupportEngine() && storeEngine != "" {
		sql += " ENGINE=" + storeEngine
	}
	if b.dialect.SupportCharset() {
		if len(charset) == 0 {
			charset = b.dialect.URI().Charset
		}
		if len(charset) > 0 {
			sql += " DEFAULT CHARSET " + charset
		}
	}

	return sql
}

func (b *Base) ForUpdateSQL(query string) string {
	return query + " FOR UPDATE"
}

func (b *Base) LogSQL(sql string, args []interface{}) {
	if b.logger != nil && b.logger.IsShowSQL() {
		if len(args) > 0 {
			b.logger.Infof("[SQL] %v %v", sql, args)
		} else {
			b.logger.Infof("[SQL] %v", sql)
		}
	}
}

func (b *Base) SetParams(params map[string]string) {
}

var (
	dialects = map[string]func() Dialect{}
)

// RegisterDialect register database dialect
func RegisterDialect(dbName DBType, dialectFunc func() Dialect) {
	if dialectFunc == nil {
		panic("core: Register dialect is nil")
	}
	dialects[strings.ToLower(string(dbName))] = dialectFunc // !nashtsai! allow override dialect
}

// QueryDialect query if registered database dialect
func QueryDialect(dbName DBType) Dialect {
	if d, ok := dialects[strings.ToLower(string(dbName))]; ok {
		return d()
	}
	return nil
}

func regDrvsNDialects() bool {
	providedDrvsNDialects := map[string]struct {
		dbType     DBType
		getDriver  func() Driver
		getDialect func() Dialect
	}{
		"mssql":    {"mssql", func() Driver { return &odbcDriver{} }, func() Dialect { return &mssql{} }},
		"odbc":     {"mssql", func() Driver { return &odbcDriver{} }, func() Dialect { return &mssql{} }}, // !nashtsai! TODO change this when supporting MS Access
		"mysql":    {"mysql", func() Driver { return &mysqlDriver{} }, func() Dialect { return &mysql{} }},
		"mymysql":  {"mysql", func() Driver { return &mymysqlDriver{} }, func() Dialect { return &mysql{} }},
		"postgres": {"postgres", func() Driver { return &pqDriver{} }, func() Dialect { return &postgres{} }},
		"pgx":      {"postgres", func() Driver { return &pqDriverPgx{} }, func() Dialect { return &postgres{} }},
		"sqlite3":  {"sqlite3", func() Driver { return &sqlite3Driver{} }, func() Dialect { return &sqlite3{} }},
		"oci8":     {"oracle", func() Driver { return &oci8Driver{} }, func() Dialect { return &oracle{} }},
		"goracle":  {"oracle", func() Driver { return &goracleDriver{} }, func() Dialect { return &oracle{} }},
	}

	for driverName, v := range providedDrvsNDialects {
		if driver := QueryDriver(driverName); driver == nil {
			RegisterDriver(driverName, v.getDriver())
			RegisterDialect(v.dbType, v.getDialect)
		}
	}
	return true
}

func init() {
	regDrvsNDialects()
}
