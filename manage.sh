#!/bin/bash

# --- 配置区 ---
#/manage.sh start s1
# 基础命令，后面会动态拼接 s0 或 s1
BASE_CMD="go run main.go ser"
# --------------

ACTION=$1    # start, stop, restart, status
INSTANCE=$2  # s0, s1 ...

if [ -z "$INSTANCE" ]; then
    echo "错误: 请指定实例名称 (例如: ./manage.sh start s0)"
    exit 1
fi

FULL_CMD="$BASE_CMD $INSTANCE"
LOG_FILE="app_${INSTANCE}.log"

# 获取指定实例的 PID
get_pid() {
    pgrep -f "$FULL_CMD"
}

case "$ACTION" in
    start)
        PID=$(get_pid)
        if [ -n "$PID" ]; then
            echo "实例 $INSTANCE 已在运行 (PID: $PID)"
        else
            nohup $FULL_CMD > $LOG_FILE 2>&1 &
            echo "实例 $INSTANCE 已启动，日志记录在 $LOG_FILE"
        fi
        ;;
    stop)
        PID=$(get_pid)
        if [ -z "$PID" ]; then
            echo "未发现正在运行的实例 $INSTANCE"
        else
            kill $PID
            echo "实例 $INSTANCE 已停止 (PID: $PID)"
        fi
        ;;
    status)
        PID=$(get_pid)
        if [ -n "$PID" ]; then
            echo "实例 $INSTANCE 状态：运行中 (PID: $PID)"
            echo "最后 3 行日志 ($LOG_FILE)："
            tail -n 3 $LOG_FILE
        else
            echo "实例 $INSTANCE 状态：未运行"
        fi
        ;;
    restart)
        $0 stop $INSTANCE
        sleep 1
        $0 start $INSTANCE
        ;;
    *)
        echo "用法: $0 {start|stop|restart|status} {s0|s1|...}"
        exit 1
esac