# `Agent Cat Agent`

某种神秘又高效的**自动化系统**, “协作式 AI 工具链”

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
