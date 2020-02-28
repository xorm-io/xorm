// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistinctAndCols(t *testing.T) {
	type DistinctAndCols struct {
		Id   int64
		Name string
	}

	assert.NoError(t, prepareEngine())
	assertSync(t, new(DistinctAndCols))

	cnt, err := testEngine.Insert(&DistinctAndCols{
		Name: "test",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	var names []string
	err = testEngine.Table("distinct_and_cols").Cols("name").Distinct("name").Find(&names)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(names))
	assert.EqualValues(t, "test", names[0])
}

func TestUpdateIgnoreOnlyFromDBFields(t *testing.T) {
	type TestOnlyFromDBField struct {
		Id              int64  `xorm:"PK"`
		OnlyFromDBField string `xorm:"<-"`
		OnlyToDBField   string `xorm:"->"`
		IngoreField     string `xorm:"-"`
	}

	assertGetRecord := func() *TestOnlyFromDBField {
		var record TestOnlyFromDBField
		has, err := testEngine.Where("id = ?", 1).Get(&record)
		assert.NoError(t, err)
		assert.EqualValues(t, true, has)
		assert.EqualValues(t, "", record.OnlyFromDBField)
		return &record

	}
	assert.NoError(t, prepareEngine())
	assertSync(t, new(TestOnlyFromDBField))

	_, err := testEngine.Insert(&TestOnlyFromDBField{
		Id:              1,
		OnlyFromDBField: "a",
		OnlyToDBField:   "b",
		IngoreField:     "c",
	})
	assert.NoError(t, err)

	assertGetRecord()

	_, err = testEngine.ID(1).Update(&TestOnlyFromDBField{
		OnlyToDBField:   "b",
		OnlyFromDBField: "test",
	})
	assert.NoError(t, err)
	assertGetRecord()
}
