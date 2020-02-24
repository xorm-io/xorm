// Copyright 2020 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package names

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Userinfo struct {
	Uid        int64  `xorm:"id pk not null autoincr"`
	Username   string `xorm:"unique"`
	Departname string
	Alias      string `xorm:"-"`
	Created    time.Time
	Detail     Userdetail `xorm:"detail_id int(11)"`
	Height     float64
	Avatar     []byte
	IsMan      bool
}

type Userdetail struct {
	Id      int64
	Intro   string `xorm:"text"`
	Profile string `xorm:"varchar(2000)"`
}

type MyGetCustomTableImpletation struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

const getCustomTableName = "GetCustomTableInterface"

func (MyGetCustomTableImpletation) TableName() string {
	return getCustomTableName
}

type TestTableNameStruct struct{}

func (t *TestTableNameStruct) TableName() string {
	return "my_test_table_name_struct"
}

func TestGetTableName(t *testing.T) {
	var kases = []struct {
		mapper            Mapper
		v                 reflect.Value
		expectedTableName string
	}{
		{
			SnakeMapper{},
			reflect.ValueOf(new(Userinfo)),
			"userinfo",
		},
		{
			SnakeMapper{},
			reflect.ValueOf(Userinfo{}),
			"userinfo",
		},
		{
			SameMapper{},
			reflect.ValueOf(new(Userinfo)),
			"Userinfo",
		},
		{
			SameMapper{},
			reflect.ValueOf(Userinfo{}),
			"Userinfo",
		},
		{
			SnakeMapper{},
			reflect.ValueOf(new(MyGetCustomTableImpletation)),
			getCustomTableName,
		},
		{
			SnakeMapper{},
			reflect.ValueOf(MyGetCustomTableImpletation{}),
			getCustomTableName,
		},
		{
			SnakeMapper{},
			reflect.ValueOf(MyGetCustomTableImpletation{}),
			getCustomTableName,
		},
		{
			SnakeMapper{},
			reflect.ValueOf(new(TestTableNameStruct)),
			new(TestTableNameStruct).TableName(),
		},
	}

	for _, kase := range kases {
		assert.EqualValues(t, kase.expectedTableName, GetTableName(kase.mapper, kase.v))
	}
}
