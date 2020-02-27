package dialects

import (
	"testing"

	"xorm.io/xorm/schemas"

	"github.com/stretchr/testify/assert"
)

func TestQuoteFilter_Do(t *testing.T) {
	f := QuoteFilter{schemas.Quoter{"[", "]"}}
	sql := "SELECT `COLUMN_NAME` FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ? AND `COLUMN_NAME` = ?"
	res := f.Do(sql)
	assert.EqualValues(t,
		"SELECT [COLUMN_NAME] FROM [INFORMATION_SCHEMA].[COLUMNS] WHERE [TABLE_SCHEMA] = ? AND [TABLE_NAME] = ? AND [COLUMN_NAME] = ?",
		res,
	)
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
