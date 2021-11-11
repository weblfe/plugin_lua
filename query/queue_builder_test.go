package query

import (
	"fmt"
	"testing"
)

func TestCommandQueryBuilder_CreateIndex(t *testing.T) {
	var (
		command = NewCommandBuilder()
		columns = []string{
			"name", "state",
		}
		table = "{{%test}}"
		name  = "uq_test"
	)
	var sql = command.CreateIndex(name, table, columns, false)
	if sql == "" {
		t.Error("创建 index SQL 失败")
		return
	}
	fmt.Println("sql:", sql)
}

func TestCommandQueryBuilder_DropIndex(t *testing.T) {
	var (
		command = NewCommandBuilder()
		table   = "{{%test}}"
		name    = "uq_test"
	)
	var sql = command.DropIndex(name, table)
	if sql == "" {
		t.Error("Drop index SQL 失败")
		return
	}
	fmt.Println("sql:", sql)
}

func TestCommandQueryBuilder_AlterColumn(t *testing.T) {
	var (
		command = NewCommandBuilder()
		table   = "{{%test}}"
		column  = "state"
		tyStr   = `varchar(100)`
	)
	var sql = command.AlterColumn(table, column, NewString(tyStr))
	if sql == "" {
		t.Error("Drop index SQL 失败")
		return
	}
	fmt.Println("sql:", sql)
}

func TestRegexpMatches(t *testing.T) {
	var (
		pattern  = `/^(\w+)\((.+?)\)(.*)$/`
		pattern2 = `/^(\w+)\s+/`
		subject  = `varchar(100)`
		matches  = RegexpMatches(pattern, subject)
		matches2 = RegexpMatches(pattern2, subject)
	)
	if len(matches) <= 0 {
		t.Error("matches error")
	}
	if len(matches2) != 0 {
		t.Error("matches2 error")
	}
	for i, match := range matches {
		fmt.Println("match item: ", i, "=>", match)
	}
}
