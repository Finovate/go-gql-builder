package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/xwb1989/sqlparser"
)

func main() {
	var (
		sqlString = flag.String("sql", "", "Direct SQL string")
		sqlFile   = flag.String("file", "", "Path to SQL file")
		sqlDir    = flag.String("dir", "", "Directory with SQL files")
	)
	flag.Parse()
	var err error

	switch {
	case *sqlString != "":
		// 处理直接输入的 SQL 字符串
		err = processSqlString(*sqlString)
	case *sqlFile != "":
		// 处理单个 SQL 文件
		err = processSqlFile(*sqlFile)
	case *sqlDir != "":
		// 处理 SQL 文件夹
		err = processSQLDir(*sqlDir)
	default:
		fmt.Println("No input provided")
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Done")
}

func processSqlString(sql string) error {
	fileName, content, err := parseSQL(sql)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	if err = saveContentToFile(fileName, content); err != nil {
		return err
	}

	return nil
}

func processSqlFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filePath, err)
	}

	return processSqlString(string(content))
}

func processSQLDir(dirPath string) error {
	fmt.Println("Processing SQL directory:", dirPath)
	// 遍历目录中的 SQL 文件并处理
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".sql" {
			err := processSqlFile(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// 其他辅助函数，如解析 SQL、生成 Go 代码等
func parseSQL(sql string) (fileName string, content []byte, err error) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		err = fmt.Errorf("error parsing SQL: %v", err)
		return
	}

	createTableStmt, ok := stmt.(*sqlparser.DDL)
	if !ok {
		err = fmt.Errorf("not a DDL statement")
		return
	}
	if createTableStmt.Action != sqlparser.CreateStr {
		err = fmt.Errorf("not a CREATE statement")
		return
	}

	primaryKeyMap := make(map[string]bool)

	for _, index := range createTableStmt.TableSpec.Indexes {
		if !index.Info.Primary {
			continue
		}
		for _, column := range index.Columns {
			primaryKeyMap[column.Column.String()] = true
		}
	}

	tableName := createTableStmt.NewName.Name.String()
	parser := NewParser(tableName)
	for _, column := range createTableStmt.TableSpec.Columns {
		c := &Column{
			Name:  column.Name.String(),
			Alias: column.Name.String(),
			Type:  column.Type.Type,
		}
		if column.Type.KeyOpt == 1 || primaryKeyMap[column.Name.String()] {
			c.IsPrimaryKey = true
		}
		parser.AddColumns(c)
	}

	result, err := parser.ParseTemplate()
	if err != nil {
		err = fmt.Errorf("error parsing template: %v", err)
		return
	}

	//result.WriteString(fmt.Sprintf("type %s struct {\n", createTableStmt.NewName.Name.String()))
	//for _, colDef := range createTableStmt.Columns {
	//	goType := sqlTypeToGoType(colDef.Type.Type)
	//	result.WriteString(fmt.Sprintf("    %s %s\n", colDef.Name.String(), goType))
	//}
	//result.WriteString("}\n")

	return tableName, result, nil
}

func saveContentToFile(fileName string, content []byte) error {
	once := sync.Once{}
	once.Do(func() {
		err := os.Mkdir("model", 0755)
		if err != nil {
			if os.IsExist(err) {
				return
			}
			fmt.Printf("Error creating directory: %v\n", err)
			return
		}
	})

	filePath := fmt.Sprintf("./model/%s.go", fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("error writing file %s: %v", filePath, err)
	}
	return nil
}
