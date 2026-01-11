# dingsw-go-md

一个基于 Go 语言开发的 Markdown 目录树扫描与处理工具。该项目能够扫描指定目录，生成层级化的 JSON 结构，并提供 Web 服务支持。

## 📁 目录结构说明

```text
.
├── main.go            # 项目入口，负责启动服务
├── internal/          # 内部逻辑包（外部无法直接引用）
│   ├── cmd/           # 核心指令逻辑（如 scanner.go 负责文件扫描）
│   ├── config/        # 配置加载逻辑
│   ├── handler/       # HTTP 请求处理器（路由响应）
│   └── service/       # 业务逻辑层
├── etc/               # 配置文件存放目录 (conf.yaml)
├── json/              # 生成的中间数据或测试 JSON
└── go.mod             # Go 模块依赖管理

```

---

## 🚀 快速开始

### 1. 环境准备

确保已安装 Go 1.24 或更高版本。

### 2. 安装依赖

```bash
go mod tidy

```

### 3. 配置文件

编辑 `etc/conf.yaml`，设置你的 Markdown 根目录和监听端口：

```yaml
Files:
  s0:
    MdPath: ./mystudy
    JsonPath: ./json/tree0.json
    WebPort: 8101
  s1:
    MdPath: ./note
    JsonPath: ./json/tree1.json
    WebPort: 8102
Cache:
    s0:
      addr: 127.0.0.1:6381
      passwd: 
      timeout: 3
RateLimit:
  Retry: 10
  Limit: 1
  Burst: 3

```

### 4. 运行

```bash
go run main.go scanner all // 扫描所有文件
go run main.go server s0 // 启动s0服务
```

---

## 🛠️ 核心模块

* **Scanner (`internal/cmd`)**: 负责递归遍历文件系统，识别 `.md` 文件并构建树状结构。
* **Service (`internal/service`)**: 处理数据转换，将扫描到的结构转换为 `tree.json` 格式。
* **Handler (`internal/handler`)**: 提供 API 接口，供前端调用以展示目录树。

---

## 📝 输出示例

程序运行后会生成或更新 `tree.json`，结构如下：

```json
{
  "name": "root",
  "children": [
    { "name": "网络协议", "type": "dir", "children": [...] }
  ]
}

```

# 构建镜像（如果还没构建）
 ### 1. 配置文件 MdPathin 修改为容器内的目录
 ### 2. mkdir json
 ### 3. docker compose build

# 运行一次性扫描任务
### 方法一. docker run --rm \
  -v $(pwd)/etc/conf.yaml:/root/etc/conf.yaml \
  -v /var/www/mystudy:/root/mystudy \
  dingsw-go-md:latest ./main scanner all

### 方法二. docker compose run --rm scanner-job