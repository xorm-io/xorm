// Copyright 2018 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/ziutek/mymysql/godrv"
	"xorm.io/xorm/caches"
	"xorm.io/xorm/log"
	"xorm.io/xorm/names"
	"xorm.io/xorm/schemas"
)

var (
	testEngine EngineInterface
	dbType     string
	connString string

	db                 = flag.String("db", "sqlite3", "the tested database")
	showSQL            = flag.Bool("show_sql", true, "show generated SQLs")
	ptrConnStr         = flag.String("conn_str", "./test.db?cache=shared&mode=rwc", "test database connection string")
	mapType            = flag.String("map_type", "snake", "indicate the name mapping")
	cacheFlag          = flag.Bool("cache", false, "if enable cache")
	cluster            = flag.Bool("cluster", false, "if this is a cluster")
	splitter           = flag.String("splitter", ";", "the splitter on connstr for cluster")
	schema             = flag.String("schema", "", "specify the schema")
	ignoreSelectUpdate = flag.Bool("ignore_select_update", false, "ignore select update if implementation difference, only for tidb")
	ingoreUpdateLimit  = flag.Bool("ignore_update_limit", false, "ignore update limit if implementation difference, only for cockroach")
	tableMapper        names.Mapper
	colMapper          names.Mapper
)

func createEngine(dbType, connStr string) error {
	if testEngine == nil {
		var err error

		if !*cluster {
			switch schemas.DBType(strings.ToLower(dbType)) {
			case schemas.MSSQL:
				db, err := sql.Open(dbType, strings.Replace(connStr, "xorm_test", "master", -1))
				if err != nil {
					return err
				}
				if _, err = db.Exec("If(db_id(N'xorm_test') IS NULL) BEGIN CREATE DATABASE xorm_test; END;"); err != nil {
					return fmt.Errorf("db.Exec: %v", err)
				}
				db.Close()
				*ignoreSelectUpdate = true
			case schemas.POSTGRES:
				db, err := sql.Open(dbType, strings.Replace(connStr, "xorm_test", "postgres", -1))
				if err != nil {
					return err
				}
				rows, err := db.Query("SELECT 1 FROM pg_database WHERE datname = 'xorm_test'")
				if err != nil {
					return fmt.Errorf("db.Query: %v", err)
				}
				defer rows.Close()

				if !rows.Next() {
					if _, err = db.Exec("CREATE DATABASE xorm_test"); err != nil {
						return fmt.Errorf("CREATE DATABASE: %v", err)
					}
				}
				if *schema != "" {
					db.Close()
					db, err = sql.Open(dbType, connStr)
					if err != nil {
						return err
					}
					defer db.Close()
					if _, err = db.Exec("CREATE SCHEMA IF NOT EXISTS " + *schema); err != nil {
						return fmt.Errorf("CREATE SCHEMA: %v", err)
					}
				}
				db.Close()
				*ignoreSelectUpdate = true
			case schemas.MYSQL:
				db, err := sql.Open(dbType, strings.Replace(connStr, "xorm_test", "mysql", -1))
				if err != nil {
					return err
				}
				if _, err = db.Exec("CREATE DATABASE IF NOT EXISTS xorm_test"); err != nil {
					return fmt.Errorf("db.Exec: %v", err)
				}
				db.Close()
			default:
				*ignoreSelectUpdate = true
			}

			testEngine, err = NewEngine(dbType, connStr)
		} else {
			testEngine, err = NewEngineGroup(dbType, strings.Split(connStr, *splitter))
			if dbType != "mysql" && dbType != "mymysql" {
				*ignoreSelectUpdate = true
			}
		}
		if err != nil {
			return err
		}

		if *schema != "" {
			testEngine.SetSchema(*schema)
		}
		testEngine.ShowSQL(*showSQL)
		testEngine.SetLogLevel(log.LOG_DEBUG)
		if *cacheFlag {
			cacher := caches.NewLRUCacher(caches.NewMemoryStore(), 100000)
			testEngine.SetDefaultCacher(cacher)
		}

		if len(*mapType) > 0 {
			switch *mapType {
			case "snake":
				testEngine.SetMapper(names.SnakeMapper{})
			case "same":
				testEngine.SetMapper(names.SameMapper{})
			case "gonic":
				testEngine.SetMapper(names.LintGonicMapper)
			}
		}
	}

	tableMapper = testEngine.GetTableMapper()
	colMapper = testEngine.GetColumnMapper()

	tables, err := testEngine.DBMetas()
	if err != nil {
		return err
	}
	var tableNames = make([]interface{}, 0, len(tables))
	for _, table := range tables {
		tableNames = append(tableNames, table.Name)
	}
	if err = testEngine.DropTables(tableNames...); err != nil {
		return err
	}
	return nil
}

func prepareEngine() error {
	return createEngine(dbType, connString)
}

func TestMain(m *testing.M) {
	flag.Parse()

	dbType = *db
	if *db == "sqlite3" {
		if ptrConnStr == nil {
			connString = "./test.db?cache=shared&mode=rwc"
		} else {
			connString = *ptrConnStr
		}
	} else {
		if ptrConnStr == nil {
			fmt.Println("you should indicate conn string")
			return
		}
		connString = *ptrConnStr
	}

	dbs := strings.Split(*db, "::")
	conns := strings.Split(connString, "::")

	var res int
	for i := 0; i < len(dbs); i++ {
		dbType = dbs[i]
		connString = conns[i]
		testEngine = nil
		fmt.Println("testing", dbType, connString)

		if err := prepareEngine(); err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}

		code := m.Run()
		if code > 0 {
			res = code
		}
	}

	os.Exit(res)
}

func TestPing(t *testing.T) {
	if err := testEngine.Ping(); err != nil {
		t.Fatal(err)
	}
}
