# Local Service Registry

一个轻量级的本地服务注册中心，带有健康检查和内置 Web UI。它为你提供一个统一的地方，记录本机上运行的所有 Web 服务，避免忘记哪些服务正在运行、监听在哪个端口。

## 功能特性

- **服务注册** -- 通过 REST API 注册，服务可在启动时自动注册自身
- **定期健康检查** -- 每 5 分钟（可配置）对已注册的 URL 发起 GET 请求，2xx 状态码视为健康
- **注册后立即检查** -- 新注册的服务会立刻进行一次异步健康检查
- **内置 Web UI** -- 在浏览器中查看所有服务及其健康状态，支持注册和删除操作
- **SQLite 持久化** -- 注册信息在重启后不会丢失
- **AI 提示词助手** -- 一键复制提示词，让 AI 为你的其他项目添加自注册功能
- **零配置** -- 直接运行二进制即可使用

## 安装

### 从源码安装

```bash
go install github.com/greatbody/local-service-registry@latest
```

### 本地构建

```bash
git clone https://github.com/greatbody/local-service-registry.git
cd local-service-registry
go build -o local-service-registry .
```

## 使用方法

```bash
# 使用默认配置启动（端口 8500，5 分钟健康检查间隔）
./local-service-registry

# 自定义端口和间隔
./local-service-registry -addr :1234 -interval 30s -db /path/to/registry.db
```

启动后在浏览器中打开 `http://localhost:1234` 即可看到 Web UI。

### 启动参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-addr` | `:8500` | HTTP 监听地址 |
| `-db` | `registry.db` | SQLite 数据库文件路径 |
| `-interval` | `5m` | 健康检查间隔 |

## API 接口

### 注册服务

```bash
curl -X POST http://localhost:1234/services \
  -H 'Content-Type: application/json' \
  -d '{"name": "my-app", "url": "http://localhost:3000", "description": "我的 Web 应用"}'
```

### 列出所有服务

```bash
curl http://localhost:1234/services
```

### 查询单个服务

```bash
curl http://localhost:1234/services/<id>
```

### 删除服务

```bash
curl -X DELETE http://localhost:1234/services/<id>
```

## 作为 macOS 系统服务运行

创建 LaunchAgent 实现登录后自动启动：

```bash
cat > ~/Library/LaunchAgents/com.local-service-registry.plist << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.local-service-registry</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/local-service-registry</string>
        <string>-addr</string>
        <string>:1234</string>
        <string>-db</string>
        <string>/Users/YOU/Library/Application Support/local-service-registry/registry.db</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/Users/YOU/Library/Logs/local-service-registry.log</string>
    <key>StandardErrorPath</key>
    <string>/Users/YOU/Library/Logs/local-service-registry.log</string>
</dict>
</plist>
EOF

# 加载并启动
launchctl load ~/Library/LaunchAgents/com.local-service-registry.plist
```

将 `/Users/YOU` 替换为你的实际用户目录。

## 让其他服务自动注册

Web UI 中有一个 **"Copy AI Prompt"** 按钮。点击复制提示词，然后粘贴给你的 AI 助手，它会为当前项目生成自注册代码——启动时异步注册到本注册中心，静默忽略所有错误（适用于可能公开部署、目标机器没有此注册中心的场景）。

## 开源协议

[MIT](LICENSE)
