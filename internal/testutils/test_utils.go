package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-jet/jet/execution"
	"github.com/go-jet/jet/internal/jet"
	"gotest.tools/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func AssertExec(t *testing.T, stmt jet.Statement, db execution.DB, rowsAffected ...int64) {
	res, err := stmt.Exec(db)

	assert.NilError(t, err)
	rows, err := res.RowsAffected()
	assert.NilError(t, err)

	if len(rowsAffected) > 0 {
		assert.Equal(t, rows, rowsAffected[0])
	}
}

func AssertExecErr(t *testing.T, stmt jet.Statement, db execution.DB, errorStr string) {
	_, err := stmt.Exec(db)

	assert.Error(t, err, errorStr)
}

func getFullPath(relativePath string) string {
	goPath := os.Getenv("GOPATH")
	return filepath.Join(goPath, "src/github.com/go-jet/jet/tests", relativePath)
}

func PrintJson(v interface{}) {
	jsonText, _ := json.MarshalIndent(v, "", "\t")
	fmt.Println(string(jsonText))
}

func AssertJSON(t *testing.T, data interface{}, expectedJSON string) {
	jsonData, err := json.MarshalIndent(data, "", "\t")
	assert.NilError(t, err)

	assert.Equal(t, "\n"+string(jsonData)+"\n", expectedJSON)
}

func SaveJsonFile(v interface{}, testRelativePath string) {
	jsonText, _ := json.MarshalIndent(v, "", "\t")

	filePath := getFullPath(testRelativePath)
	err := ioutil.WriteFile(filePath, jsonText, 0644)

	if err != nil {
		panic(err)
	}
}

func AssertJSONFile(t *testing.T, data interface{}, testRelativePath string) {

	filePath := getFullPath(testRelativePath)
	fileJSONData, err := ioutil.ReadFile(filePath)
	assert.NilError(t, err)

	if runtime.GOOS == "windows" {
		fileJSONData = bytes.Replace(fileJSONData, []byte("\r\n"), []byte("\n"), -1)
	}

	jsonData, err := json.MarshalIndent(data, "", "\t")
	assert.NilError(t, err)

	assert.Assert(t, string(fileJSONData) == string(jsonData))
	//assert.DeepEqual(t, string(fileJSONData), string(jsonData))
}

func AssertStatementSql(t *testing.T, query jet.Statement, expectedQuery string, expectedArgs ...interface{}) {
	queryStr, args, err := query.Sql()
	assert.NilError(t, err)
	assert.Equal(t, queryStr, expectedQuery)

	if len(expectedArgs) == 0 {
		return
	}
	assert.DeepEqual(t, args, expectedArgs)
}

func AssertStatementSqlErr(t *testing.T, stmt jet.Statement, errorStr string) {
	_, _, err := stmt.Sql()

	assert.Assert(t, err != nil)
	assert.Error(t, err, errorStr)
}

func AssertDebugStatementSql(t *testing.T, query jet.Statement, expectedQuery string, expectedArgs ...interface{}) {
	_, args, err := query.Sql()
	assert.NilError(t, err)
	//assert.Equal(t, queryStr, expectedQuery)
	if len(expectedArgs) > 0 {
		assert.DeepEqual(t, args, expectedArgs)
	}

	debuqSql, err := query.DebugSql()

	assert.NilError(t, err)

	assert.Equal(t, debuqSql, expectedQuery)
}

func AssertClauseSerialize(t *testing.T, dialect jet.Dialect, clause jet.Serializer, query string, args ...interface{}) {
	out := jet.SqlBuilder{Dialect: dialect}
	err := jet.Serialize(clause, jet.SelectStatementType, &out)

	assert.NilError(t, err)

	//fmt.Println(out.Buff.String())

	assert.DeepEqual(t, out.Buff.String(), query)
	assert.DeepEqual(t, out.Args, args)
}

func AssertClauseSerializeErr(t *testing.T, dialect jet.Dialect, clause jet.Serializer, errString string) {
	out := jet.SqlBuilder{Dialect: dialect}
	err := jet.Serialize(clause, jet.SelectStatementType, &out)

	//fmt.Println(out.buff.String())
	assert.Assert(t, err != nil)
	assert.Error(t, err, errString)
}

func AssertProjectionSerialize(t *testing.T, dialect jet.Dialect, projection jet.Projection, query string, args ...interface{}) {
	out := jet.SqlBuilder{Dialect: dialect}
	err := jet.SerializeForProjection(projection, jet.SelectStatementType, &out)

	assert.NilError(t, err)

	assert.DeepEqual(t, out.Buff.String(), query)
	assert.DeepEqual(t, out.Args, args)
}
