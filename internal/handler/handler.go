package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// FileNode 结构
type FileNode struct {
	Name     string
	Path     string
	IsDir    bool
	Children []*FileNode
}

// 模板数据
type PageData struct {
	Content     template.HTML
	Tree        *FileNode
	CurrentPath string
}

var MdFilePath = ""
var JsonOutputPath = "" // 输出路径

func loadTreeFromJSON() (*FileNode, error) {
	data, err := os.ReadFile(JsonOutputPath)
	if err != nil {
		return nil, err
	}
	var tree FileNode
	err = json.Unmarshal(data, &tree)
	return &tree, err
}

// 递归构建树并排序
const pageTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>我的文档库</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/5.2.0/github-markdown.min.css">
    <style>
        body { display: flex; margin: 0; background: #f6f8fa; font-family: -apple-system, Segoe UI, Helvetica, Arial, sans-serif; height: 100vh; }
        
        /* 侧边栏优化 */
        .sidebar { 
            width: 300px; background: #ffffff; border-right: 1px solid #d0d7de; 
            padding: 20px 10px; overflow-y: auto; flex-shrink: 0;
        }
        .sidebar h3 { padding-left: 15px; color: #1f2328; font-size: 16px; border-bottom: 1px solid #d0d7de; padding-bottom: 10px; }
        
        /* 内容区优化 */
        .main-content { flex: 1; overflow-y: auto; padding: 40px 20px; scroll-behavior: smooth; }
        .markdown-body { 
            max-width: 880px; margin: 0 auto; background: #fff; padding: 45px; 
            border: 1px solid #d0d7de; border-radius: 8px; box-shadow: 0 1px 2px rgba(0,0,0,0.05); 
        }

        /* 树状结构 CSS */
        details { margin: 2px 0; }
        summary { 
            cursor: pointer; padding: 4px 8px; border-radius: 4px; font-size: 14px;
            color: #1f2328; display: flex; align-items: center; 
        }
        summary::before { content: "▸"; display: inline-block; width: 15px; transition: transform 0.2s; }
        details[open] > summary::before { transform: rotate(90deg); }
		/* 鼠标悬停在 summary 上时显示背景，增加交互感 */
		summary:hover { 
			background: #f0f2f5; 
			border-radius: 4px;
		}

		/* 给打开状态的文件夹名字加粗 */
		details[open] > summary {
			color: #0969da;
			font-weight: 600;
		}

		/* 调整缩进线，增强层级感 */
		ul {
			border-left: 1px solid #e1e4e8;
			margin-left: 10px;
			padding-left: 15px;
		}
        
        .file-link { 
            text-decoration: none; color: #444d56; font-size: 14px; display: block; 
            padding: 4px 8px 4px 15px; border-radius: 4px; margin: 1px 0;
        }
        .file-link:hover { background: #f3f4f6; color: #0969da; }
        
        /* 当前页面高亮样式 */
        .active { background: #ddf4ff !important; color: #0969da !important; font-weight: 600; border-left: 3px solid #0969da; }

        @media (max-width: 768px) { body { flex-direction: column; } .sidebar { width: 100%; height: 300px; } }
    </style>
</head>
<body>
    <div class="sidebar">
        <h3>📖 文档目录</h3>
        {{template "tree" .}}
    </div>
    <div class="main-content">
        <article class="markdown-body">
            {{.Content}}
        </article>
    </div>
</body>
</html>

{{define "tree"}}
<ul>
    {{$current := .CurrentPath}}
    {{range .Tree.Children}}
        <li>
            {{if .IsDir}}
                <details {{if isAncestor $current .Path}}open{{end}}>
                    <summary>📁 {{.Name}}</summary>
                    {{template "tree" dict "Tree" . "CurrentPath" $current}}
                </details>
            {{else}}
                <a class="file-link {{if eq $current .Path}}active{{end}}" href="{{.Path}}">
                    📄 {{.Name}}
                </a>
            {{end}}
        </li>
    {{end}}
</ul>
{{end}}
`

func MdHandler(w http.ResponseWriter, r *http.Request) {
	relPath := r.URL.Path
	if relPath == "/" || relPath == "" {
		relPath = "/README"
	}
	// 1. 直接读 JSON 生成侧边栏
	tree, err := loadTreeFromJSON()
	if err != nil {
		log.Println("JSON加载失败，请先运行 scanner.go")
	}
	// 2. 拼接绝对路径读取文件
	// 注意：filepath.Join 会处理多余的斜杠
	filePath := filepath.Join(MdFilePath, relPath+".md")

	input, err := os.ReadFile(filePath)
	log.Printf("md 文件目录: %s", relPath+".md")
	// log.Printf("读取成功2: %s", input)
	if err != nil {
		http.Error(w, "文件未找到: "+relPath, http.StatusNotFound)
		log.Printf("读取失败: %s", filePath)
		http.Error(w, "文件未找到", http.StatusNotFound)
		return
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(highlighting.WithStyle("github")),
		),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)
	// log.Printf("转换成功: %s", input)
	var buf strings.Builder
	md.Convert(input, &buf)

	// 准备模板函数和数据
	// log.Printf("模板成功: %s", buf)

	// 必须先 New，再注册 Funcs，最后 Parse
	tmpl, err := template.New("page").Funcs(template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			dict := make(map[string]interface{})
			for i := 0; i < len(values); i += 2 {
				dict[values[i].(string)] = values[i+1]
			}
			return dict, nil
		},
		// 新增：判断当前路径是否在某个文件夹内
		"isAncestor": func(currentPath, folderPath string) bool {
			if folderPath == "/" || folderPath == "" {
				return true
			}
			// 如果当前访问的是 /A/B/C，那么 /A 和 /A/B 都是它的祖先，应该展开
			return strings.HasPrefix(currentPath, folderPath)
		},
	}).Parse(pageTemplate)

	if err != nil {
		log.Fatalf("模板语法错误: %v", err) // 这里会告诉你具体的行号和错误原因
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// log.Printf("执行成功: %s", tmpl)
	// fmt.Println("三级/四级目录支", relPath, tree)
	// fmt.Println("三级/四级目录支111111", buf.String())
	// fmt.Println("三级/四级目录支222222", template.HTML(buf.String()))
	tmpl.Execute(w, PageData{
		Content:     template.HTML(buf.String()),
		Tree:        tree,
		CurrentPath: relPath,
	})

	// fmt.Println("三级/四级目录支DDDDDDD")
}
