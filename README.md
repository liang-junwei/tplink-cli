<h1 align="center">TPlink CLI</h1>
<p align="center">
  <img src="https://img.shields.io/badge/Go-1.25.0-cyan" alt="Go"/>
  <img src="https://img.shields.io/badge/Cobra-1.10.2-royalblue" alt="Cobra"/>
  <img src="https://img.shields.io/badge/Yaml-3.0.1-orange" alt="Yaml"/>
    <img src="https://img.shields.io/badge/TPlink-WAR1200L-deepskyblue" alt="Yaml"/>
</p>
<hr>




## 简介

tplink-cli 是一个 Go 实现的 TP-Link 路由器命令行管理工具，单二进制、零依赖、跨平台。覆盖端口映射、DHCP、Wi-Fi、行为管控、安全防护等常见运维场景。所有操作通过 HTTP API 与路由器交互，无需登录 Web 管理界面。

> **⚠️ 适用设备**：本工具仅适用于 TP-Link 企业路由器 **TL-WAR1200L**。其他型号（含 WAR 系列其他版本、家用路由器等）未做适配，不保证可用。

---

## 安装

二进制安装

```bash
# 下载二进制文件 重命名为 tplink
# Windows 放置在C:\Windows\tplink.exe
# Linux 放置在/usr/bin/tplink
```

从源码构建：

```bash
git clone https://github.com/ljw/tplink-cli.git
cd tplink-cli && go build -o tplink .
```

---

## 快速开始

### 添加路由器

```bash
tplink context add home -U http://192.168.0.1 -u admin -p your_password --default
```

配置自动持久化到 `~/.tplink.json`，支持管理多台设备。

> your_password 需要从浏览器 web F12控制台中 登录接口截获的 加密密码 而非明文密码。若使用 明文密码 参考下面的`动态认证方式`。

### 动态认证（`--dynamic-auth`）

某些 TP-Link 固件版本在登录时会对密码进行编码处理（非标准 HTTP Basic/Digest），编码密钥存储在前端 JS 文件中，且**随固件版本变化**。

为兼容此类设备，提供了动态认证机制：

**工作原理**

1. `context add/update --dynamic-auth` 时，立即从路由器 `{url}/web-static/js/su/su.js` 提取编码密钥
2. 使用密钥对密码进行编码，将**编码后的密码**持久化到 `~/.tplink.json`（字段 `encoded_password`）
3. 后续登录直接发送已编码的密码，不再每次动态获取密钥

**使用方式**

```bash
# 添加设备时开启动态认证
tplink context add tp-war1200l \
  -U http://192.168.1.1 \
  -u admin \
  -p your_password \
  --dynamic-auth

# 为已有设备开启动态认证（会立即连接路由器获取密钥）
tplink context update home --dynamic-auth

# 登录时自动使用已编码的密码，无需额外操作
tplink -S home forward list

# 关闭动态认证（回退到发送原始密码）
tplink context update home --dynamic-auth=false
```

**错误处理**

若动态获取密钥失败，会输出详细错误信息并指引回退：

```
动态认证开启失败（无法从路由器获取密钥）
  错误: 动态获取认证密钥失败（获取 JS 文件）: ...
  提示: 请确认路由器 URL 正确且可访问，或关闭动态认证：
        tplink context update home --dynamic-auth=false
```

> **⚠️ 注意**：动态认证默认为关闭状态。建议使用从TP-Link web 登录接口中截获的密钥，无需开启此选项。



### 全局参数

| 参数 | 简写 | 说明 |
|------|------|------|
| `--server` | `-S` | 指定目标设备 |
| `--output` | `-o` | 输出格式：`table` / `json` / `yaml` |
| `--dry-run` | | 模拟执行，仅打印请求不发送 |

---

## 命令概览

| 分类 | 命令 | 说明 |
|------|------|------|
| 配置 | `context` | 管理路由器连接配置 |
| 端口映射 | `forward` | NAT 规则 CRUD、启停 |
| 接口 | `ifmode` `brv6mode` `port` `wan` `lan` | 接口模式、WAN/LAN 管理 |
| DHCP | `dhcp` | DHCP 配置、客户端列表、静态绑定 |
| 无线 | `wireless` | Wi-Fi 配置、访客网络、MAC 过滤、客户端 |
| 行为管控 | `ipgroup` `timerange` `qos` `acl` | IP 组、时间段、QoS、访问控制 |
| 安全 | `arp` `macfilter` `dos` | ARP 防护、MAC 过滤、DoS 防御 |
| VPN | `vpn` | VPN配置 |
| 高级功能 | `route` `napt` `alg` `phddns` | 路由、NAPT、动态DNS |
| 系统 | `system` | 设备信息、重启 |
| 通用 | `api` | 发送原始 API 请求 |

---

## 典型用法

### 端口映射

最常用的场景——把内网服务暴露到公网。

```bash
# 查看所有规则
tplink forward list

# 按名称/协议/IP 过滤
tplink forward list -n myapp --proto TCP -d 192.168.0.100

# 添加映射
tplink forward add -n myapp -p 8080 -d 192.168.0.100

# 端口范围映射
tplink forward add -n range-demo -p 8000-8090 -d 192.168.0.100

# 外部端口与内部端口不同
tplink forward add -n web -p 8443 -P 443 -d 192.168.0.100 --proto TCP

# 启停 / 删除
tplink forward disable redirect_1779863436
tplink forward enable redirect_1779863436
tplink forward delete redirect_1779863436
```

### 查看设备状态

```bash
# 设备信息
tplink system info

# WAN 口状态
tplink wan list

# LAN 口状态
tplink lan list

# 有线端口状态
tplink port list

# DHCP 客户端列表（谁连了我的 Wi-Fi）
tplink dhcp client list
```

### DHCP 静态绑定

```bash
# 查看现有绑定
tplink dhcp static list

# 添加绑定
tplink dhcp static add -m 90-E2-BA-7F-6D-58 -i 192.168.0.250 -n my-server

# 删除绑定
tplink dhcp static del dhcp_static_5
```

### Wi-Fi 管理

```bash
# 查看 Wi-Fi 配置
tplink wireless config list

# 修改 2.4G Wi-Fi
tplink wireless config set --2g --ssid "MyWiFi" -k newpassword

# 查看连接客户端
tplink wireless client list
```

### 多设备管理

```bash
# 添加多台设备
tplink context add office -U http://10.0.0.1 -u admin -p pass123

# 查看列表
tplink context list

# 切换默认设备
tplink context use office

# 临时指定设备执行命令
tplink -S office forward list
```

### 原始 API 调用

当内置命令不满足需求时，直接调用 API：

```bash
# GET 请求
tplink api get /users -q dept=2

# POST 请求
tplink api post /users -d '{"account":"test","password":"Abc123456"}'
```

---

## 认证流程

```
CLI 命令 → stok 缓存有效? ──是──→ 执行业务请求
                │
               否
                ↓
           调用登录接口获取 stok
                │
                ↓
           缓存到 ~/.tplink.json (有效期 30min)
```

- 首次使用通过 `context add` 保存密码
- stok 自动管理，过期自动刷新，无需手动干预
- `--dry-run` 模式下仅输出请求内容，不实际发送


---

## 输出格式

支持三种格式，通过 `-o` 切换：

```bash
tplink forward list                # 默认 table
tplink forward list -o json        # JSON
tplink forward list -o yaml        # YAML
```

---

## 项目结构

```
tplink-cli/
├── main.go
├── cmd/           # CLI 命令定义
├── internal/
│   ├── api/       # HTTP 客户端、业务 API 封装
│   ├── config/    # 配置持久化管理
│   ├── model/     # 数据结构
│   └── format/    # 输出格式化 (table/json/yaml)
├── go.mod
└── go.sum
```
