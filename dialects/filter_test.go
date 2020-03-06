package dialects

import (
	"testing"

	"xorm.io/xorm/schemas"

	"github.com/stretchr/testify/assert"
)

func TestQuoteFilter_Do(t *testing.T) {
	f := QuoteFilter{schemas.Quoter{'[', ']', schemas.AlwaysReserve}}
	var kases = []struct {
		source   string
		expected string
	}{
		{
			"SELECT `COLUMN_NAME` FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ? AND `COLUMN_NAME` = ?",
			"SELECT [COLUMN_NAME] FROM [INFORMATION_SCHEMA].[COLUMNS] WHERE [TABLE_SCHEMA] = ? AND [TABLE_NAME] = ? AND [COLUMN_NAME] = ?",
		},
		{
			"SELECT 'abc```test```''', `a` FROM b",
			"SELECT 'abc```test```''', [a] FROM b",
		},
		{
			"UPDATE table SET `a` = ~ `a`, `b`='abc`'",
			"UPDATE table SET [a] = ~ [a], [b]='abc`'",
		},
	}

	for _, kase := range kases {
		t.Run(kase.source, func(t *testing.T) {
			assert.EqualValues(t, kase.expected, f.Do(kase.source))
		})
	}
}

func TestSeqFilter(t *testing.T) {
	var kases = map[string]string{
		"SELECT * FROM TABLE1 WHERE a=? AND b=?":                               "SELECT * FROM TABLE1 WHERE a=$1 AND b=$2",
		"SELECT 1, '???', '2006-01-02 15:04:05' FROM TABLE1 WHERE a=? AND b=?": "SELECT 1, '???', '2006-01-02 15:04:05' FROM TABLE1 WHERE a=$1 AND b=$2",
		"select '1''?' from issue":                                             "select '1''?' from issue",
		"select '1\\??' from issue":                                            "select '1\\??' from issue",
		"select '1\\\\',? from issue":                                          "select '1\\\\',$1 from issue",
		"select '1\\''?',? from issue":                                         "select '1\\''?',$1 from issue",
	}
	for sql, result := range kases {
		assert.EqualValues(t, result, convertQuestionMark(sql, "$", 1))
	}
}
