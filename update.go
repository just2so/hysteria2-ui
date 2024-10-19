package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/google/uuid"
	"github.com/google/go-github/v45/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"github.com/spf13/viper"  // 用于读取配置文件
)

var (
	logger *logrus.Logger
)

func init() {
	// 初始化日志
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
}

func main() {
	// 加载配置
	if err := loadConfig(); err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}

	// 从配置中读取变量
	sqliteDBPath := viper.GetString("sqlite.path")
	githubToken := viper.GetString("github.token")
	githubOwner := viper.GetString("github.owner")
	githubRepo := viper.GetString("github.repo")
	configFilePath := viper.GetString("github.configFilePath")
	githubBranch := viper.GetString("github.branch")
	username := viper.GetString("user.username")

	// 生成新的UUID
	newUUID := uuid.NewString()

	// 1. 打开SQLite数据库
	db, err := sql.Open("sqlite3", sqliteDBPath)
	if err != nil {
		logger.Fatalf("Failed to open SQLite database: %v", err)
	}
	defer db.Close()

	// 2. 开启事务
	tx, err := db.Begin()
	if err != nil {
		logger.Fatalf("Failed to begin transaction: %v", err)
	}

	// 3. 并发处理数据库和GitHub的更新
	var wg sync.WaitGroup
	errors := make(chan error, 2) // 用于捕获并发任务中的错误

	wg.Add(2)

	go func() {
		defer wg.Done()
		err := updateDatabase(tx, username, newUUID)
		if err != nil {
			errors <- wrapError("Updating database", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := updateGitHubConfig(context.Background(), githubToken, githubOwner, githubRepo, configFilePath, githubBranch, newUUID)
		if err != nil {
			errors <- wrapError("Updating GitHub config", err)
		}
	}()

	wg.Wait()
	close(errors)

	// 4. 检查是否有错误，若有则回滚事务
	if len(errors) > 0 {
		for err := range errors {
			logger.Error(err)
		}
		tx.Rollback()
		logger.Fatal("Transaction rolled back due to errors")
	}

	// 5. 提交事务
	err = tx.Commit()
	if err != nil {
		logger.Fatalf("Failed to commit transaction: %v", err)
	}

	logger.Info("Successfully updated SQLite and GitHub config!")
}

// loadConfig 加载配置文件
func loadConfig() error {
	viper.SetConfigName("config") // 配置文件名
	viper.AddConfigPath(".")       // 当前路径
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Error reading config file: %v", err)
	}
	return nil
}

// updateDatabase 更新SQLite数据库中的con_pass
func updateDatabase(tx *sql.Tx, username, newUUID string) error {
	_, err := tx.Exec("UPDATE your_table SET con_pass = ? WHERE username = ?", newUUID, username)
	if err != nil {
		return wrapError("Updating database", err)
	}
	return nil
}

// updateGitHubConfig 更新GitHub仓库中的配置文件
func updateGitHubConfig(ctx context.Context, githubToken, owner, repo, filePath, branch, newUUID string) error {
	// 创建GitHub客户端
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// 获取GitHub文件内容
	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, filePath, &github.RepositoryContentGetOptions{Ref: branch})
	if err != nil {
		return wrapError("Fetching GitHub file", err)
	}

	// 修改文件内容中的UUID
	content, err := fileContent.GetContent()
	if err != nil {
		return wrapError("Reading GitHub file content", err)
	}
	newContent := updateUUIDInConfig(content, newUUID)

	// 提交修改后的文件
	opts := &github.RepositoryContentFileOptions{
		Message: github.String("Update UUID in config file"),
		Content: []byte(newContent),
		SHA:     fileContent.SHA,
		Branch:  github.String(branch),
	}

	_, _, err = client.Repositories.UpdateFile(ctx, owner, repo, filePath, opts)
	if err != nil {
		return wrapError("Updating GitHub file", err)
	}

	return nil
}

// updateUUIDInConfig 修改配置文件中的UUID
func updateUUIDInConfig(content, newUUID string) string {
	// 假设UUID在文件内容中通过某种标记存在，使用字符串替换
	// 示例：content = strings.Replace(content, "old_uuid", newUUID, 1)
	return content // 此处需根据实际需求编写替换逻辑
}

// wrapError 包装错误信息
func wrapError(step string, err error) error {
	return fmt.Errorf("%s: %w", step, err)
}
