package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dxddgwc/dingsw-go-md/internal/config"
)

var MdFilePath = ""
var JsonOutputPath = "" // 输出路径
// FileNode 结构
type FileNode struct {
	Name     string
	Path     string
	IsDir    bool
	Children []*FileNode
}

func buildTree(rootPath string) (*FileNode, error) {
	node := &FileNode{
		Name:  filepath.Base(rootPath),
		Path:  filepath.ToSlash(strings.TrimPrefix(rootPath, MdFilePath)),
		IsDir: true,
	}

	files, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fullPath := filepath.Join(rootPath, file.Name())
		if strings.HasPrefix(file.Name(), ".") {
			fmt.Println("正在处理:", fullPath, file.Name(), rootPath)
			// fmt.Println("正在处理111:", strings.HasPrefix(file.Name(), "."))
			// os.Exit(1)
			continue
		} else if file.IsDir() {
			child, err := buildTree(fullPath)
			if err == nil {
				node.Children = append(node.Children, child)
			}
		} else if strings.HasSuffix(file.Name(), ".md") {
			name := strings.TrimSuffix(file.Name(), ".md")
			relPath := filepath.ToSlash(strings.TrimPrefix(fullPath, MdFilePath))
			node.Children = append(node.Children, &FileNode{
				Name:  name,
				Path:  strings.TrimSuffix(relPath, ".md"),
				IsDir: false,
			})
		}
	}

	sort.Slice(node.Children, func(i, j int) bool {
		if node.Children[i].IsDir != node.Children[j].IsDir {
			return node.Children[i].IsDir
		}
		return node.Children[i].Name < node.Children[j].Name
	})
	return node, nil
}

// go run main.go scanner s0
func Scanner(conf *config.Config) {
	for _, file := range conf.Files {
		MdFilePath = file.MdPath
		JsonOutputPath = file.JsonPath
		batch(&file)
	}
}
func batch(conf *config.File) {
	exists, err := config.FileExists(MdFilePath)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic(fmt.Errorf("配置文件 : %s 不存在", MdFilePath))
	}
	tree, err := buildTree(MdFilePath)
	if err != nil {
		fmt.Printf("扫描失败: %v\n", err)
		return
	}

	data, _ := json.MarshalIndent(tree, "", "  ")
	dir := filepath.Dir(JsonOutputPath)
	EnsureDirectoryExists(dir)
	err = os.WriteFile(JsonOutputPath, data, 0644)
	if err != nil {
		fmt.Printf("保存失败: %v\n", err)
	} else {
		fmt.Println("✅ 目录索引已更新: " + JsonOutputPath)
	}
}

// 判断目录是否存在，不存在则创建
func EnsureDirectoryExists(directoryPath string) error {
	// 判断目录是否存在
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		// 目录不存在，创建目录
		err := os.MkdirAll(directoryPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory: %v", err)
		}
	}
	return nil
}
