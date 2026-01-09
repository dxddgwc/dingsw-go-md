package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dxddgwc/dingsw-go-md/internal/cmd"
	"github.com/dxddgwc/dingsw-go-md/internal/config"
	"github.com/dxddgwc/dingsw-go-md/internal/handler"
)

// # 建议先编译，再运行（效率更高，且容器/后台管理更方便）
// go build -o main main.go
// # 后台运行并将日志输出到 output.log
// nohup ./main server s0 > output.log 2>&1 &

func main() {

	conf := config.New("./etc/conf.yaml")
	args := os.Args
	task := args[1]
	s_tag := args[2]

	if task == "scanner" { //go run main.go scanner all
		cmd.Scanner(conf)
	} else {
		file_conf := conf.Files[s_tag]
		webPort := fmt.Sprintf(":%s", file_conf.WebPort)
		handler.MdFilePath = file_conf.MdPath
		handler.JsonOutputPath = file_conf.JsonPath
		http.HandleFunc("/", handler.MdHandler)
		fmt.Println("启动服务：http://localhost" + webPort)
		log.Fatal(http.ListenAndServe(webPort, nil))
	}
}
