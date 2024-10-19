package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/google/uuid"
	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

func main() {
	// 生成新的UUID
	newUUID := uuid.NewString()

	// 1. 打开SQLite数据库
	db, err := sql.Open("sqlite3", "./your_database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 2. 开启事务
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// 3. 更新SQLite数据库
	_, err = tx.Exec("UPDATE your_table SET con_pass = ? WHERE username = ?", newUUID, "zhangsan")
	if err != nil {
		tx.Rollback() // 更新失败，回滚事务
		log.Fatal("SQLite更新失败:", err)
	}

	// 4. 更新GitHub配置文件
	// 4.1 创建GitHub客户端
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "your_github_token"}, // 使用你的GitHub Token
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// 4.2 获取文件内容
	owner := "your_github_owner"
	repo := "your_private_repo"
	filePath := "path/to/config/file"
	branch := "main" // 分支名称

	// 获取文件信息
	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, filePath, &github.RepositoryContentGetOptions{Ref: branch})
	if err != nil {
		tx.Rollback() // 获取GitHub文件失败，回滚事务
		log.Fatal("GitHub文件获取失败:", err)
	}

	// 4.3 修改文件内容
	content, err := fileContent.GetContent()
	if err != nil {
		tx.Rollback() // 文件内容解析失败，回滚事务
		log.Fatal("GitHub文件内容解析失败:", err)
	}
	newContent := updateUUIDInConfig(content, newUUID) // 自定义函数，替换配置文件中的UUID

	// 4.4 更新GitHub文件
	opts := &github.RepositoryContentFileOptions{
		Message: github.String("Update UUID in config file"),
		Content: []byte(newContent),
		SHA:     fileContent.SHA,
		Branch:  github.String(branch),
	}

	_, _, err = client.Repositories.UpdateFile(ctx, owner, repo, filePath, opts)
	if err != nil {
		tx.Rollback() // GitHub更新失败，回滚事务
		log.Fatal("GitHub文件更新失败:", err)
	}

	// 5. 提交事务
	err = tx.Commit()
	if err != nil {
		log.Fatal("事务提交失败:", err)
	}

	fmt.Println("SQLite和GitHub文件更新成功！")
}

// updateUUIDInConfig: 自定义函数，替换配置文件中的UUID
func updateUUIDInConfig(content string, newUUID string) string {
	// 假设UUID在配置文件中以某种方式存在，编写替换逻辑
	// 示例：content = strings.Replace(content, "old_uuid", newUUID, 1)
	return content
}
