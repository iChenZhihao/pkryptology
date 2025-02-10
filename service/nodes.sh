#!/bin/bash

# 默认启动的节点数
DEFAULT_NODE_COUNT=5
# 默认端口起始值
DEFAULT_START_PORT=9080

# 启动 gg20-node 的函数
start_nodes() {
    local node_count=$1
    local start_port=$DEFAULT_START_PORT

    echo "正在启动 $node_count 个节点..."

    for ((i = 0; i < node_count; i++)); do
        local port=$((start_port + i))
        if lsof -i :$port > /dev/null 2>&1; then
            echo "端口 $port 已被占用，跳过启动。"
        else
            echo "启动节点，端口: $port"
            ./gg20-node --server.port=$port --config=./config.yaml > "./logs/gg20-node_$port.log" 2>&1 &
        fi
    done

    echo "所有节点启动完成。"
}

# 停止所有 gg20-node 进程的函数
stop_nodes() {
    echo "正在停止所有 gg20-node 进程..."
    pkill -f "gg20-node"
    if [ $? -eq 0 ]; then
        echo "所有 gg20-node 进程已停止。"
    else
        echo "没有找到 gg20-node 进程。"
    fi
}

# 根据参数执行操作
case "$1" in
    start)
        # 如果提供了第二个参数，则使用它作为节点数，否则使用默认值
        node_count=${2:-$DEFAULT_NODE_COUNT}
        start_nodes $node_count
        ;;
    stop)
        stop_nodes
        ;;
    *)
        echo "用法: $0 {start [节点数]|stop}"
        echo "示例:"
        echo "  $0 start       # 启动 5 个节点（默认）"
        echo "  $0 start 3     # 启动 3 个节点"
        echo "  $0 stop        # 停止所有 gg20-node 进程"
        exit 1
        ;;
esac
