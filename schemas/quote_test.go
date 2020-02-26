// Copyright 2019 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schemas

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuoteTo(t *testing.T) {
	var quoter = Quoter{"[", "]"}

	test := func(t *testing.T, expected string, value string) {
		buf := &strings.Builder{}
		quoter.QuoteTo(buf, value)
		assert.EqualValues(t, expected, buf.String())
	}

	test(t, "[mytable]", "mytable")
	test(t, "[mytable]", "`mytable`")
	test(t, "[mytable]", `[mytable]`)

	test(t, `["mytable"]`, `"mytable"`)

	test(t, "[myschema].[mytable]", "myschema.mytable")
	test(t, "[myschema].[mytable]", "`myschema`.mytable")
	test(t, "[myschema].[mytable]", "myschema.`mytable`")
	test(t, "[myschema].[mytable]", "`myschema`.`mytable`")
	test(t, "[myschema].[mytable]", `[myschema].mytable`)
	test(t, "[myschema].[mytable]", `myschema.[mytable]`)
	test(t, "[myschema].[mytable]", `[myschema].[mytable]`)

	test(t, `["myschema].[mytable"]`, `"myschema.mytable"`)

	test(t, "[message_user] AS [sender]", "`message_user` AS `sender`")

	assert.EqualValues(t, "[a],[b]", quoter.Join([]string{"a", " b"}, ","))

	buf := &strings.Builder{}
	quoter = Quoter{"", ""}
	quoter.QuoteTo(buf, "noquote")
	assert.EqualValues(t, "noquote", buf.String())
}

func TestJoin(t *testing.T) {
	cols := []string{"f1", "f2", "f3"}
	quoter := Quoter{"[", "]"}

	assert.EqualValues(t, "[f1], [f2], [f3]", quoter.Join(cols, ", "))

	quoter = Quoter{"", ""}
	assert.EqualValues(t, "f1, f2, f3", quoter.Join(cols, ", "))
}
