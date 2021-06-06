// Copyright 2021 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integrations

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"xorm.io/builder"
)

func TestCount(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type UserinfoCount struct {
		Departname string
	}
	assert.NoError(t, testEngine.Sync2(new(UserinfoCount)))

	colName := testEngine.GetColumnMapper().Obj2Table("Departname")
	var cond builder.Cond = builder.Eq{
		"`" + colName + "`": "dev",
	}

	total, err := testEngine.Where(cond).Count(new(UserinfoCount))
	assert.NoError(t, err)
	assert.EqualValues(t, 0, total)

	cnt, err := testEngine.Insert(&UserinfoCount{
		Departname: "dev",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	total, err = testEngine.Where(cond).Count(new(UserinfoCount))
	assert.NoError(t, err)
	assert.EqualValues(t, 1, total)

	total, err = testEngine.Where(cond).Table("userinfo_count").Count()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, total)

	total, err = testEngine.Table("userinfo_count").Count()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, total)
}

func TestSQLCount(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type UserinfoCount2 struct {
		Id         int64
		Departname string
	}

	type UserinfoBooks struct {
		Id     int64
		Pid    int64
		IsOpen bool
	}

	assertSync(t, new(UserinfoCount2), new(UserinfoBooks))

	total, err := testEngine.SQL("SELECT count(id) FROM " + testEngine.TableName("userinfo_count2", true)).
		Count()
	assert.NoError(t, err)
	assert.EqualValues(t, 0, total)
}

func TestCountWithOthers(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type CountWithOthers struct {
		Id   int64
		Name string
	}

	assertSync(t, new(CountWithOthers))

	_, err := testEngine.Insert(&CountWithOthers{
		Name: "orderby",
	})
	assert.NoError(t, err)

	_, err = testEngine.Insert(&CountWithOthers{
		Name: "limit",
	})
	assert.NoError(t, err)

	total, err := testEngine.OrderBy("id desc").Limit(1).Count(new(CountWithOthers))
	assert.NoError(t, err)
	assert.EqualValues(t, 2, total)
}

type CountWithTableName struct {
	Id   int64
	Name string
}

func (CountWithTableName) TableName() string {
	return "count_with_table_name1"
}

func TestWithTableName(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	assertSync(t, new(CountWithTableName))

	_, err := testEngine.Insert(&CountWithTableName{
		Name: "orderby",
	})
	assert.NoError(t, err)

	_, err = testEngine.Insert(CountWithTableName{
		Name: "limit",
	})
	assert.NoError(t, err)

	total, err := testEngine.OrderBy("id desc").Count(new(CountWithTableName))
	assert.NoError(t, err)
	assert.EqualValues(t, 2, total)

	total, err = testEngine.OrderBy("id desc").Count(CountWithTableName{})
	assert.NoError(t, err)
	assert.EqualValues(t, 2, total)
}

func TestCountWithSelectCols(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	assertSync(t, new(CountWithTableName))

	_, err := testEngine.Insert(&CountWithTableName{
		Name: "orderby",
	})
	assert.NoError(t, err)

	_, err = testEngine.Insert(CountWithTableName{
		Name: "limit",
	})
	assert.NoError(t, err)

	total, err := testEngine.Cols("id").Count(new(CountWithTableName))
	assert.NoError(t, err)
	assert.EqualValues(t, 2, total)

	total, err = testEngine.Select("count(id)").Count(CountWithTableName{})
	assert.NoError(t, err)
	assert.EqualValues(t, 2, total)
}

func TestCountWithGroupBy(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	assertSync(t, new(CountWithTableName))

	_, err := testEngine.Insert(&CountWithTableName{
		Name: "1",
	})
	assert.NoError(t, err)

	_, err = testEngine.Insert(CountWithTableName{
		Name: "2",
	})
	assert.NoError(t, err)

	cnt, err := testEngine.GroupBy("name").Count(new(CountWithTableName))
	assert.NoError(t, err)
	assert.EqualValues(t, 2, cnt)
}
