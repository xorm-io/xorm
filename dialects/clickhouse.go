// Copyright 2020 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package dialects

import (
	"context"
	"net/url"

	"xorm.io/xorm/core"
	"xorm.io/xorm/schemas"
)

type clickhouse struct {
	Base
}

func (db *clickhouse) Init(uri *URI) error {
	return db.Base.Init(db, uri)
}

func (db *clickhouse) IsReserved(name string) bool {
	return false
}

func (db *clickhouse) SQLType(c *schemas.Column) string {
	return ""
}

func (db *clickhouse) SetQuotePolicy(quotePolicy QuotePolicy) {
}

func (*clickhouse) AutoIncrStr() string {
	return ""
}

func (*clickhouse) CreateTableSQL(t *schemas.Table, tableName string) ([]string, bool) {
	return nil, false
}

func (*clickhouse) IsTableExist(queryer core.Queryer, ctx context.Context, tableName string) (bool, error) {
	return false, nil
}

func (*clickhouse) Filters() []Filter {
	return []Filter{}
}

func (*clickhouse) GetColumns(core.Queryer, context.Context, string) ([]string, map[string]*schemas.Column, error) {
	return nil, nil, nil
}

func (db *clickhouse) GetIndexes(queryer core.Queryer, ctx context.Context, tableName string) (map[string]*schemas.Index, error) {
	return nil, nil
}

func (db *clickhouse) IndexCheckSQL(tableName, idxName string) (string, []interface{}) {
	return "", nil
}

func (db *clickhouse) GetTables(queryer core.Queryer, ctx context.Context) ([]*schemas.Table, error) {
	return nil, nil
}

// ParseClickHouse parsed clickhouse connection string
// tcp://host1:9000?username=user&password=qwerty&database=clicks&read_timeout=10&write_timeout=20&alt_hosts=host2:9000,host3:9000
func ParseClickHouse(connStr string) (*URI, error) {
	u, err := url.Parse(connStr)
	if err != nil {
		return nil, err
	}
	forms := u.Query()
	return &URI{
		DBType: schemas.CLICKHOUSE,
		Proto:  u.Scheme,
		Host:   u.Hostname(),
		Port:   u.Port(),
		DBName: forms.Get("database"),
		User:   forms.Get("username"),
		Passwd: forms.Get("password"),
	}, nil
}