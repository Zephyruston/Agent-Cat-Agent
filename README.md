# `Agent Cat Agent`

某种神秘又高效的**自动化系统**, "协作式 AI 工具链"

## 😺 Why `Cat`？

> "Like the `cat` command in Linux, our agents are curious observers and executors. One acts like a cat who does the work in a box (Docker), the other checks the box (like Schrödinger's cat) to see if it's alive and done its job."

### 👨‍💻 Agent 1：Coder Agent（执行型）

- **功能**：

  - 接收任务（写代码、构建 Docker 镜像、部署运行等）
  - 自动生成代码、打包到容器，执行任务

### 🕵️‍♀️ Agent 2：Watcher Agent（监督型）

- **功能**：

  - 实时监控运行状态（容器是否启动成功、任务是否完成）
  - 如果任务异常/超时/失败，发出通知或反馈，让 Agent 1 重试或改进

## 运行

```bash
git clone https://github.com/Zephyruston/Agent-Cat-Agent.git
cd Agent-Cat-Agent
go mod tidy
cp ./etc/config.example.yaml ./etc/config.yaml # 然后在config.yaml 填入你的相关内容

go build ./cmd/aca.go
# 推荐使用如下命令行格式
./aca --mode=<gen|test> --prompt "生成一个矩阵乘法" --language go --config ./etc/config.yaml --workdir ./tmp --mount /app

# 参数说明：
# --mode        生成代码(gen) 或单元测试(test)
# --prompt      代码/测试生成的自然语言描述（必填）
# --language/-l 编程语言（目前仅支持 go）
# --config/-c   配置文件路径
# --workdir     本地工作目录（生成代码/测试文件的存放目录）
# --mount       容器内挂载路径（如 /app）

# 例如：
./aca --mode=gen --prompt "生成一个矩阵乘法" --language go --config ./etc/config.yaml --workdir ./tmp --mount /app
```

## 🗺️ Roadmap

### 已完成

- [x] 实现基于 OpenAI 的代码生成功能
- [x] 支持 Go 语言的代码生成(mode=gen)和 Docker 容器内运行

### TODO

- [ ] 完善 Coder Agent 和 Watcher Agent 基础架构
- [ ] 扩展语言支持，增加 Python, Rust 等语言的代码生成和测试运行
- [ ] 实现自动依赖管理，根据生成代码自动检测并安装依赖
- [ ] 代码安全检查，防止生成包含漏洞的代码
- [ ] 优化代码生成质量，添加代码审查和修复功能
- [ ] 开发 Web 界面，提供友好操作和可视化监控
- [ ] 实现 Agent 间智能协作，自动分配任务和资源
- [ ] 引入强化学习机制，持续优化任务执行策略
