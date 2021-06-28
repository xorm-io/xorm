// Copyright 2020 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tags

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"xorm.io/xorm/caches"
	"xorm.io/xorm/dialects"
	"xorm.io/xorm/names"
	"xorm.io/xorm/schemas"

	"github.com/stretchr/testify/assert"
)

type ParseTableName1 struct{}

type ParseTableName2 struct{}

func (p ParseTableName2) TableName() string {
	return "p_parseTableName"
}

func TestParseTableName(t *testing.T) {
	parser := NewParser(
		"xorm",
		dialects.QueryDialect("mysql"),
		names.SnakeMapper{},
		names.SnakeMapper{},
		caches.NewManager(),
	)
	table, err := parser.Parse(reflect.ValueOf(new(ParseTableName1)))
	assert.NoError(t, err)
	assert.EqualValues(t, "parse_table_name1", table.Name)

	table, err = parser.Parse(reflect.ValueOf(new(ParseTableName2)))
	assert.NoError(t, err)
	assert.EqualValues(t, "p_parseTableName", table.Name)

	table, err = parser.Parse(reflect.ValueOf(ParseTableName2{}))
	assert.NoError(t, err)
	assert.EqualValues(t, "p_parseTableName", table.Name)
}

func TestUnexportField(t *testing.T) {
	parser := NewParser(
		"xorm",
		dialects.QueryDialect("mysql"),
		names.SnakeMapper{},
		names.SnakeMapper{},
		caches.NewManager(),
	)

	type VanilaStruct struct {
		private int // unexported fields will be ignored
		Public  int
	}
	table, err := parser.Parse(reflect.ValueOf(new(VanilaStruct)))
	assert.NoError(t, err)
	assert.EqualValues(t, "vanila_struct", table.Name)
	assert.EqualValues(t, 1, len(table.Columns()))

	for _, col := range table.Columns() {
		assert.EqualValues(t, "public", col.Name)
		assert.NotEqual(t, "private", col.Name)
	}

	type TaggedStruct struct {
		private int `xorm:"private"` // unexported fields will be ignored
		Public  int `xorm:"-"`
	}
	table, err = parser.Parse(reflect.ValueOf(new(TaggedStruct)))
	assert.NoError(t, err)
	assert.EqualValues(t, "tagged_struct", table.Name)
	assert.EqualValues(t, 0, len(table.Columns()))
}

func TestParseWithOtherIdentifier(t *testing.T) {
	parser := NewParser(
		"xorm",
		dialects.QueryDialect("mysql"),
		names.SameMapper{},
		names.SnakeMapper{},
		caches.NewManager(),
	)

	type StructWithDBTag struct {
		FieldFoo string `db:"foo"`
	}

	parser.SetIdentifier("db")
	table, err := parser.Parse(reflect.ValueOf(new(StructWithDBTag)))
	assert.NoError(t, err)
	assert.EqualValues(t, "StructWithDBTag", table.Name)
	assert.EqualValues(t, 1, len(table.Columns()))

	for _, col := range table.Columns() {
		assert.EqualValues(t, "foo", col.Name)
	}
}

func TestParseWithIgnore(t *testing.T) {
	parser := NewParser(
		"db",
		dialects.QueryDialect("mysql"),
		names.SameMapper{},
		names.SnakeMapper{},
		caches.NewManager(),
	)

	type StructWithIgnoreTag struct {
		FieldFoo string `db:"-"`
	}

	table, err := parser.Parse(reflect.ValueOf(new(StructWithIgnoreTag)))
	assert.NoError(t, err)
	assert.EqualValues(t, "StructWithIgnoreTag", table.Name)
	assert.EqualValues(t, 0, len(table.Columns()))
}

func TestParseWithAutoincrement(t *testing.T) {
	parser := NewParser(
		"db",
		dialects.QueryDialect("mysql"),
		names.SnakeMapper{},
		names.GonicMapper{},
		caches.NewManager(),
	)

	type StructWithAutoIncrement struct {
		ID int64
	}

	table, err := parser.Parse(reflect.ValueOf(new(StructWithAutoIncrement)))
	assert.NoError(t, err)
	assert.EqualValues(t, "struct_with_auto_increment", table.Name)
	assert.EqualValues(t, 1, len(table.Columns()))
	assert.EqualValues(t, "id", table.Columns()[0].Name)
	assert.True(t, table.Columns()[0].IsAutoIncrement)
	assert.True(t, table.Columns()[0].IsPrimaryKey)
}

func TestParseWithAutoincrement2(t *testing.T) {
	parser := NewParser(
		"db",
		dialects.QueryDialect("mysql"),
		names.SnakeMapper{},
		names.GonicMapper{},
		caches.NewManager(),
	)

	type StructWithAutoIncrement2 struct {
		ID int64 `db:"pk autoincr"`
	}

	table, err := parser.Parse(reflect.ValueOf(new(StructWithAutoIncrement2)))
	assert.NoError(t, err)
	assert.EqualValues(t, "struct_with_auto_increment2", table.Name)
	assert.EqualValues(t, 1, len(table.Columns()))
	assert.EqualValues(t, "id", table.Columns()[0].Name)
	assert.True(t, table.Columns()[0].IsAutoIncrement)
	assert.True(t, table.Columns()[0].IsPrimaryKey)
	assert.False(t, table.Columns()[0].Nullable)
}

func TestParseWithNullable(t *testing.T) {
	parser := NewParser(
		"db",
		dialects.QueryDialect("mysql"),
		names.SnakeMapper{},
		names.GonicMapper{},
		caches.NewManager(),
	)

	type StructWithNullable struct {
		Name     string `db:"notnull"`
		FullName string `db:"null comment('column comment,字段注释')"`
	}

	table, err := parser.Parse(reflect.ValueOf(new(StructWithNullable)))
	assert.NoError(t, err)
	assert.EqualValues(t, "struct_with_nullable", table.Name)
	assert.EqualValues(t, 2, len(table.Columns()))
	assert.EqualValues(t, "name", table.Columns()[0].Name)
	assert.EqualValues(t, "full_name", table.Columns()[1].Name)
	assert.False(t, table.Columns()[0].Nullable)
	assert.True(t, table.Columns()[1].Nullable)
	assert.EqualValues(t, "column comment,字段注释", table.Columns()[1].Comment)
}

func TestParseWithTimes(t *testing.T) {
	parser := NewParser(
		"db",
		dialects.QueryDialect("mysql"),
		names.SnakeMapper{},
		names.GonicMapper{},
		caches.NewManager(),
	)

	type StructWithTimes struct {
		Name      string    `db:"notnull"`
		CreatedAt time.Time `db:"created"`
		UpdatedAt time.Time `db:"updated"`
		DeletedAt time.Time `db:"deleted"`
	}

	table, err := parser.Parse(reflect.ValueOf(new(StructWithTimes)))
	assert.NoError(t, err)
	assert.EqualValues(t, "struct_with_times", table.Name)
	assert.EqualValues(t, 4, len(table.Columns()))
	assert.EqualValues(t, "name", table.Columns()[0].Name)
	assert.EqualValues(t, "created_at", table.Columns()[1].Name)
	assert.EqualValues(t, "updated_at", table.Columns()[2].Name)
	assert.EqualValues(t, "deleted_at", table.Columns()[3].Name)
	assert.False(t, table.Columns()[0].Nullable)
	assert.True(t, table.Columns()[1].Nullable)
	assert.True(t, table.Columns()[1].IsCreated)
	assert.True(t, table.Columns()[2].Nullable)
	assert.True(t, table.Columns()[2].IsUpdated)
	assert.True(t, table.Columns()[3].Nullable)
	assert.True(t, table.Columns()[3].IsDeleted)
}

func TestParseWithExtends(t *testing.T) {
	parser := NewParser(
		"db",
		dialects.QueryDialect("mysql"),
		names.SnakeMapper{},
		names.GonicMapper{},
		caches.NewManager(),
	)

	type StructWithEmbed struct {
		Name      string
		CreatedAt time.Time `db:"created"`
		UpdatedAt time.Time `db:"updated"`
		DeletedAt time.Time `db:"deleted"`
	}

	type StructWithExtends struct {
		SW StructWithEmbed `db:"extends"`
	}

	table, err := parser.Parse(reflect.ValueOf(new(StructWithExtends)))
	assert.NoError(t, err)
	assert.EqualValues(t, "struct_with_extends", table.Name)
	assert.EqualValues(t, 4, len(table.Columns()))
	assert.EqualValues(t, "name", table.Columns()[0].Name)
	assert.EqualValues(t, "created_at", table.Columns()[1].Name)
	assert.EqualValues(t, "updated_at", table.Columns()[2].Name)
	assert.EqualValues(t, "deleted_at", table.Columns()[3].Name)
	assert.True(t, table.Columns()[0].Nullable)
	assert.True(t, table.Columns()[1].Nullable)
	assert.True(t, table.Columns()[1].IsCreated)
	assert.True(t, table.Columns()[2].Nullable)
	assert.True(t, table.Columns()[2].IsUpdated)
	assert.True(t, table.Columns()[3].Nullable)
	assert.True(t, table.Columns()[3].IsDeleted)
}

func TestParseWithCache(t *testing.T) {
	parser := NewParser(
		"db",
		dialects.QueryDialect("mysql"),
		names.SnakeMapper{},
		names.GonicMapper{},
		caches.NewManager(),
	)

	type StructWithCache struct {
		Name string `db:"cache"`
	}

	table, err := parser.Parse(reflect.ValueOf(new(StructWithCache)))
	assert.NoError(t, err)
	assert.EqualValues(t, "struct_with_cache", table.Name)
	assert.EqualValues(t, 1, len(table.Columns()))
	assert.EqualValues(t, "name", table.Columns()[0].Name)
	assert.True(t, table.Columns()[0].Nullable)
	cacher := parser.cacherMgr.GetCacher(table.Name)
	assert.NotNil(t, cacher)
}

func TestParseWithNoCache(t *testing.T) {
	parser := NewParser(
		"db",
		dialects.QueryDialect("mysql"),
		names.SnakeMapper{},
		names.GonicMapper{},
		caches.NewManager(),
	)

	type StructWithNoCache struct {
		Name string `db:"nocache"`
	}

	table, err := parser.Parse(reflect.ValueOf(new(StructWithNoCache)))
	assert.NoError(t, err)
	assert.EqualValues(t, "struct_with_no_cache", table.Name)
	assert.EqualValues(t, 1, len(table.Columns()))
	assert.EqualValues(t, "name", table.Columns()[0].Name)
	assert.True(t, table.Columns()[0].Nullable)
	cacher := parser.cacherMgr.GetCacher(table.Name)
	assert.Nil(t, cacher)
}

func TestParseWithEnum(t *testing.T) {
	parser := NewParser(
		"db",
		dialects.QueryDialect("mysql"),
		names.SnakeMapper{},
		names.GonicMapper{},
		caches.NewManager(),
	)

	type StructWithEnum struct {
		Name string `db:"enum('alice', 'bob')"`
	}

	table, err := parser.Parse(reflect.ValueOf(new(StructWithEnum)))
	assert.NoError(t, err)
	assert.EqualValues(t, "struct_with_enum", table.Name)
	assert.EqualValues(t, 1, len(table.Columns()))
	assert.EqualValues(t, "name", table.Columns()[0].Name)
	assert.True(t, table.Columns()[0].Nullable)
	assert.EqualValues(t, schemas.Enum, strings.ToUpper(table.Columns()[0].SQLType.Name))
	assert.EqualValues(t, map[string]int{
		"alice": 0,
		"bob":   1,
	}, table.Columns()[0].EnumOptions)
}

func TestParseWithSet(t *testing.T) {
	parser := NewParser(
		"db",
		dialects.QueryDialect("mysql"),
		names.SnakeMapper{},
		names.GonicMapper{},
		caches.NewManager(),
	)

	type StructWithSet struct {
		Name string `db:"set('alice', 'bob')"`
	}

	table, err := parser.Parse(reflect.ValueOf(new(StructWithSet)))
	assert.NoError(t, err)
	assert.EqualValues(t, "struct_with_set", table.Name)
	assert.EqualValues(t, 1, len(table.Columns()))
	assert.EqualValues(t, "name", table.Columns()[0].Name)
	assert.True(t, table.Columns()[0].Nullable)
	assert.EqualValues(t, schemas.Set, strings.ToUpper(table.Columns()[0].SQLType.Name))
	assert.EqualValues(t, map[string]int{
		"alice": 0,
		"bob":   1,
	}, table.Columns()[0].SetOptions)
}

func TestParseWithIndex(t *testing.T) {
	parser := NewParser(
		"db",
		dialects.QueryDialect("mysql"),
		names.SnakeMapper{},
		names.GonicMapper{},
		caches.NewManager(),
	)

	type StructWithIndex struct {
		Name  string `db:"index"`
		Name2 string `db:"index(s)"`
		Name3 string `db:"unique"`
	}

	table, err := parser.Parse(reflect.ValueOf(new(StructWithIndex)))
	assert.NoError(t, err)
	assert.EqualValues(t, "struct_with_index", table.Name)
	assert.EqualValues(t, 3, len(table.Columns()))
	assert.EqualValues(t, "name", table.Columns()[0].Name)
	assert.EqualValues(t, "name2", table.Columns()[1].Name)
	assert.EqualValues(t, "name3", table.Columns()[2].Name)
	assert.True(t, table.Columns()[0].Nullable)
	assert.True(t, table.Columns()[1].Nullable)
	assert.True(t, table.Columns()[2].Nullable)
	assert.EqualValues(t, 1, len(table.Columns()[0].Indexes))
	assert.EqualValues(t, 1, len(table.Columns()[1].Indexes))
	assert.EqualValues(t, 1, len(table.Columns()[2].Indexes))
}

func TestParseWithVersion(t *testing.T) {
	parser := NewParser(
		"db",
		dialects.QueryDialect("mysql"),
		names.SnakeMapper{},
		names.GonicMapper{},
		caches.NewManager(),
	)

	type StructWithVersion struct {
		Name    string
		Version int `db:"version"`
	}

	table, err := parser.Parse(reflect.ValueOf(new(StructWithVersion)))
	assert.NoError(t, err)
	assert.EqualValues(t, "struct_with_version", table.Name)
	assert.EqualValues(t, 2, len(table.Columns()))
	assert.EqualValues(t, "name", table.Columns()[0].Name)
	assert.EqualValues(t, "version", table.Columns()[1].Name)
	assert.True(t, table.Columns()[0].Nullable)
	assert.True(t, table.Columns()[1].Nullable)
	assert.True(t, table.Columns()[1].IsVersion)
}

func TestParseWithLocale(t *testing.T) {
	parser := NewParser(
		"db",
		dialects.QueryDialect("mysql"),
		names.SnakeMapper{},
		names.GonicMapper{},
		caches.NewManager(),
	)

	type StructWithLocale struct {
		UTCLocale   time.Time `db:"utc"`
		LocalLocale time.Time `db:"local"`
	}

	table, err := parser.Parse(reflect.ValueOf(new(StructWithLocale)))
	assert.NoError(t, err)
	assert.EqualValues(t, "struct_with_locale", table.Name)
	assert.EqualValues(t, 2, len(table.Columns()))
	assert.EqualValues(t, "utc_locale", table.Columns()[0].Name)
	assert.EqualValues(t, "local_locale", table.Columns()[1].Name)
	assert.EqualValues(t, time.UTC, table.Columns()[0].TimeZone)
	assert.EqualValues(t, time.Local, table.Columns()[1].TimeZone)
}

func TestParseWithDefault(t *testing.T) {
	parser := NewParser(
		"db",
		dialects.QueryDialect("mysql"),
		names.SnakeMapper{},
		names.GonicMapper{},
		caches.NewManager(),
	)

	type StructWithDefault struct {
		Default1 time.Time `db:"default '1970-01-01 00:00:00'"`
		Default2 time.Time `db:"default(CURRENT_TIMESTAMP)"`
	}

	table, err := parser.Parse(reflect.ValueOf(new(StructWithDefault)))
	assert.NoError(t, err)
	assert.EqualValues(t, "struct_with_default", table.Name)
	assert.EqualValues(t, 2, len(table.Columns()))
	assert.EqualValues(t, "default1", table.Columns()[0].Name)
	assert.EqualValues(t, "default2", table.Columns()[1].Name)
	assert.EqualValues(t, "'1970-01-01 00:00:00'", table.Columns()[0].Default)
	assert.EqualValues(t, "CURRENT_TIMESTAMP", table.Columns()[1].Default)
}

func TestParseWithOnlyToDB(t *testing.T) {
	parser := NewParser(
		"db",
		dialects.QueryDialect("mysql"),
		names.GonicMapper{
			"DB": true,
		},
		names.SnakeMapper{},
		caches.NewManager(),
	)

	type StructWithOnlyToDB struct {
		Default1 time.Time `db:"->"`
		Default2 time.Time `db:"<-"`
	}

	table, err := parser.Parse(reflect.ValueOf(new(StructWithOnlyToDB)))
	assert.NoError(t, err)
	assert.EqualValues(t, "struct_with_only_to_db", table.Name)
	assert.EqualValues(t, 2, len(table.Columns()))
	assert.EqualValues(t, "default1", table.Columns()[0].Name)
	assert.EqualValues(t, "default2", table.Columns()[1].Name)
	assert.EqualValues(t, schemas.ONLYTODB, table.Columns()[0].MapType)
	assert.EqualValues(t, schemas.ONLYFROMDB, table.Columns()[1].MapType)
}

func TestParseWithJSON(t *testing.T) {
	parser := NewParser(
		"db",
		dialects.QueryDialect("mysql"),
		names.GonicMapper{
			"JSON": true,
		},
		names.SnakeMapper{},
		caches.NewManager(),
	)

	type StructWithJSON struct {
		Default1 []string `db:"json"`
	}

	table, err := parser.Parse(reflect.ValueOf(new(StructWithJSON)))
	assert.NoError(t, err)
	assert.EqualValues(t, "struct_with_json", table.Name)
	assert.EqualValues(t, 1, len(table.Columns()))
	assert.EqualValues(t, "default1", table.Columns()[0].Name)
	assert.True(t, table.Columns()[0].IsJSON)
}

func TestParseWithSQLType(t *testing.T) {
	parser := NewParser(
		"db",
		dialects.QueryDialect("mysql"),
		names.GonicMapper{
			"SQL": true,
		},
		names.GonicMapper{
			"UUID": true,
		},
		caches.NewManager(),
	)

	type StructWithSQLType struct {
		Col1     string    `db:"varchar(32)"`
		Col2     string    `db:"char(32)"`
		Int      int64     `db:"bigint"`
		DateTime time.Time `db:"datetime"`
		UUID     string    `db:"uuid"`
	}

	table, err := parser.Parse(reflect.ValueOf(new(StructWithSQLType)))
	assert.NoError(t, err)
	assert.EqualValues(t, "struct_with_sql_type", table.Name)
	assert.EqualValues(t, 5, len(table.Columns()))
	assert.EqualValues(t, "col1", table.Columns()[0].Name)
	assert.EqualValues(t, "col2", table.Columns()[1].Name)
	assert.EqualValues(t, "int", table.Columns()[2].Name)
	assert.EqualValues(t, "date_time", table.Columns()[3].Name)
	assert.EqualValues(t, "uuid", table.Columns()[4].Name)

	assert.EqualValues(t, "VARCHAR", table.Columns()[0].SQLType.Name)
	assert.EqualValues(t, "CHAR", table.Columns()[1].SQLType.Name)
	assert.EqualValues(t, "BIGINT", table.Columns()[2].SQLType.Name)
	assert.EqualValues(t, "DATETIME", table.Columns()[3].SQLType.Name)
	assert.EqualValues(t, "UUID", table.Columns()[4].SQLType.Name)
}
