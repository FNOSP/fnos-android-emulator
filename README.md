# androidemu — 安卓模拟器 / 云手机

**中文** | [English](README.en.md)

![version](https://img.shields.io/badge/version-v3.8.6-blue) ![arch](https://img.shields.io/badge/arch-x86__64%20%7C%20arm64-orange) ![image](https://img.shields.io/badge/image-~2GB-green) ![stars](https://img.shields.io/github/stars/lin1740/fnos-android-emulator) ![last-commit](https://img.shields.io/github/last-commit/lin1740/fnos-android-emulator) ![license](https://img.shields.io/github/license/lin1740/fnos-android-emulator)

📚 **使用手册与常见问题**：见本文档下方各章节

在飞牛 fnOS 上一键运行 Android 12 虚拟机，通过浏览器远程操控，支持 WebRTC / WebSocket 双投屏模式、ADB 连接、APK 安装、文件管理等。

基于 Android 容器+穿云投屏 scrcpy-over-webrtc（画面服务）双容器架构，适配飞牛统一网关。
<img width="1288" height="900" alt="firefox exe_20260927_095032" src="https://github.com/user-attachments/assets/da7718fc-c490-4f35-be32-572f0f1a3849" />

---

## 目录

- [功能特性](#功能特性)
- [安装要求](#安装要求)
- [端口说明](#端口说明)
- [安装方法](#安装方法)
- [访问方式](#访问方式)
- [默认账号](#默认账号)
- [快速上手（穿云投屏操作指南）](#快速上手穿云投屏操作指南)
- [ADB 连接](#adb-连接)
- [APK 安装](#apk-安装)
- [性能优化](#性能优化)
- [串口控制台（开发者调试）](#串口控制台开发者调试)
- [容器架构](#容器架构)
- [技术栈与翻译层](#技术栈与翻译层)
- [审核合规性说明](#审核合规性说明)
- [常见问题](#常见问题)
- [已知限制](#已知限制)
- [问题、建议反馈链接和渠道](#问题建议反馈链接和渠道)
- [对发布者和其他贡献者的支持](#对发布者和其他贡献者的支持)
- [开源许可和免责声明](#开源许可和免责声明)
- [致谢和导向链接](#致谢和导向链接)

---

## 功能特性

- **Android 12 系统**：x86_64 架构，镜像默认内置 **libndk_translation**（Google 官方 NDK 翻译层）并已启用，支持 x86_64/arm64-v8a/x86/armeabi-v7a/armeabi 五种 ABI，大多数 ARM 应用可直接安装运行
- **浏览器远程操控**：无需安装客户端，打开浏览器即可操控安卓桌面
- **双投屏模式**：
  - WebRTC 投屏（低延迟、高帧率，局域网推荐）
  - WebSocket 投屏（兼容性好，外网/穿透环境自动切换）
- **飞牛统一网关适配**：通过飞牛应用中心一键打开，自动注入登录态，无需二次登录
- **ADB 调试**：`adb connect <NAS_IP>:5556` 即可连接
- **文件管理**：在穿云投屏界面中直接管理安卓容器内文件
- **终端访问**：内置安卓 shell 终端
- **串口控制台**：默认关闭以减少性能损耗，开发者可手动开启用于调试（见下文）
- **穿云投屏 Agent 接入**：支持接入多台安卓设备/真机，统一管理
- **简体中文 + 中国时区**：容器默认 `zh_CN` + `Asia/Shanghai`
- **应用本体非 root 运行**：gateway.py 等后台进程以 `docker-androidemu` 用户运行，符合上架审核要求
- **X86 / ARM 双平台**：自动检测架构和 GPU 能力，X86 用硬件加速，ARM 自动切软件渲染
- **性能优化**：webrtc/turn 进程高优先级调度、精简后台服务、CPU 动态调频（见下文）
- **60fps 高帧率**：redroid 容器和 webrtc 编码均提升至 60fps，编码长边 1280，码率 4-16Mbps 动态调整，画面更流畅、触摸滑动更跟手
- **安装/更新中断安全**：安装或更新中途取消会自动清理临时数据，避免占位导致下次无法安装
- **自动容器检测**：网关自动检测安卓容器状态，容器启动后自动上线，无需手动操作
- **安装预检查**：安装前自动检测 binder 驱动、内存（<1GB 阻断）、Docker 可用性、磁盘空间（<2GB 阻断）、GPU 能力，不通过时给出明确原因，不再显示"执行脚本出错且原因未知"
- **容器健康检查**：实时检测 boot 状态、运行时长、OOM、surfaceflinger/agent 进程，自动识别"启动超时""内存不足被杀死""画面服务异常"等问题
- **一键修复**：状态页提供"修复GPU/画面""重启安卓容器""重启画面服务"三个按钮，无需 SSH 命令行即可自助修复常见问题
- **友好状态页**：上游服务不可用时显示美观的状态页（容器状态表格、常见问题排查、刷新按钮），不再是纯文本 "Bad Gateway"
- **连接稳定性优化**：WebSocket 自动重连（指数退避 1s→30s）+ 25 秒心跳保活，TURN 中继配置优化（no-loopback-peers、bps-capacity、max-allocate-lifetime=3600）
- **沉浸式全屏**：电脑端 object-fit:contain、手机端 cover，多选择器兼容不同版本穿云投屏，点击全屏键自动全部全屏
- **VAAPI 硬件编解码动态检测**：自动检测 GPU 是否支持 VAAPI 编码（vainfo 含 EncSlice/EncPicture）和解码（H264 VLD），支持才启用硬件加速，不支持自动回退 Google 软件编解码器，避免黑屏/花屏
- **NVIDIA GPU 支持**：自动检测 NVIDIA GPU 并挂载设备和驱动，智能 GPU 选择优先级 Intel > AMD > NVIDIA
- **双翻译层自动管理**：内置 libndk_translation（默认）和 libhoudini（自动下载），通过 bind mount 覆盖 /system/lib*/libnb.so 切换；auto 模式自动检测 ARMv8.1 指令 SIGILL 崩溃并切换到 houdini，5 分钟防循环重启冷却
- **启动性能优化**：lmkd 阈值调高（最高 315MB→3072MB）减少启动期频繁杀进程，dex2oat 首次启动用 verify-only 模式加快启动，提升 system_server/surfaceflinger 进程优先级
- **国内 APP 性能优化**：针对抖音/快手等国内 APP 占用高的问题，系统级限制后台进程数（4个）、禁用自动同步/后台数据/网络扫描、CPU/GPU 深度优化、动画半速（0.5）、lmkd 内存管理优化
- **Go 合并守护进程**：音频修复 + 分辨率自动切换合并为单个 androidemu_daemon 进程，减少 Go runtime 内存占用（约省 5-8MB）
- **可选资源限制**：docker-compose 中预置注释好的 CPU/内存限制配置，用户可根据 NAS 性能自行启用

---


## 安装要求

| 项目 | 最低要求 | 推荐配置 |
|------|----------|----------|
| 系统 | 飞牛 fnOS 1.1.8+ | fnOS 最新版 |
| 架构 | x86_64 / aarch64 | x86_64 |
| CPU | 2 核 | 4 核+ |
| 内存 | 2 GB | 4 GB+ |
| 存储 | 4 GB 可用 | 8 GB+（含镜像约 2GB） |
| 网络 | 局域网 | — |

> **X86 设备**：需要 Docker 支持 `/dev/dri` 直通以启用硬件加速；无 GPU 时自动回退软件渲染；在下载应用之前必须安装 binder_linux 驱动（应用中心里面有，直接搜索就可以），否则无法使用或者被拒绝安装。
> **ARM 设备**：自动使用软件渲染（gpu_mode=guest），无需额外驱动，也不需要 binder_linux。
>
> **架构兼容性说明**：X86_64 已在飞牛 fnOS 上充分测试；ARM64（aarch64）代码层面已做适配（自动软件渲染、容器内 binderfs、架构判断），但因测试设备有限，建议 ARM 用户安装后留意启动状态，如有问题欢迎通过下方渠道反馈。

---

## 端口说明

本应用实际使用以下端口，其中仅 8443 由飞牛应用中心自动反代，其余端口需根据使用场景自行处理：

| 端口 | 协议 | 用途 | 局域网使用 | 外网访问 |
|------|------|------|-----------|---------|
| 8443 | TCP | 穿云投屏 Web 界面 + 信令 | 自动可用（飞牛反代） | 需自行配置反向代理/内网穿透 |
| 3478 | TCP+UDP | TURN/STUN 中继（WebRTC 投屏必需） | 自动可用（host 网络） | 需自行端口映射或内网穿透 |
| 5556 | TCP | ADB 调试（外部 adb connect） | **默认仅本机可用**（需手动开放后局域网可用） | 需自行端口映射或内网穿透 |
| 50000-50100 | UDP | WebRTC 媒体流（TURN relay 备用） | 自动可用（host 网络） | 通常无需外网开放，TURN 走 3478 即可 |

> **注意**：
> - 飞牛应用中心的 `service_port` 仅声明 8443，用于统一网关反代，不代表软件只监听这一个端口。
> - webrtc 容器使用 `host` 网络模式，3478 和 50000-50100 直接监听宿主机，无需 Docker 端口映射。
> - ADB 5556 由宿主机 socat 进程转发到安卓容器的 5555，非 Docker 映射。
> - 8443、3478、50000-50100 在局域网内自动可用，无需配置。
> - 需要外网访问时，8443 走反向代理即可；WebRTC 投屏若在外网使用，必须同时让 3478（TCP+UDP）可达，否则会黑屏或一直转圈。



## 安装方法

### 方法一：飞牛应用中心（推荐）

1. 打开飞牛 fnOS 应用中心
2. 搜索「安卓模拟器」或「androidemu」
3. 点击安装，等待镜像拉取完成（约 2GB，优先 DaoCloud/飞牛加速源）
4. 安装完成后点击「打开」即可进入穿云投屏界面

### 方法二：手动安装 fpk

1. 下载最新版 `androidemu_all_x.x.x.fpk`
2. 在飞牛应用中心选择「手动安装」，上传 fpk 文件
3. 等待安装和镜像拉取完成

### 镜像加速

安装时自动按以下顺序尝试镜像源：
1. 飞牛加速源 `docker.fnnas.com`
2. DaoCloud 免注册加速源 `docker.m.daocloud.io`
3. Docker Hub 官方仓库

若全部失败，请在飞牛 Docker 设置中配置镜像加速器（registry-mirrors）后重试。

---

## 访问方式

### 局域网（推荐，功能最完整）

在飞牛应用中心点击「打开」，或直接访问：
```
https://<NAS_IP>:<NAS_端口>/app/androidemu/
```
也可直接访问画面服务原生端口：
```
https://<NAS_IP>:8443
```

局域网下 WebRTC 投屏可正常使用，低延迟、高帧率，所有功能（文件管理、终端、投屏设置等）均可用。

### 飞牛官方远程域名（xxx.fnos.net）

通过飞牛官方提供的远程域名访问时：
- ✅ WebSocket 投屏自动启用，可正常观看和操控画面
- ⚠️ 文件管理等依赖 WebRTC DataChannel（UDP）的功能可能不可用
- ⚠️ 帧率和延迟受上行带宽影响

> 原因：飞牛反向代理只透传 TCP（HTTP/HTTPS/WebSocket），不透传 UDP。WebRTC 媒体流走 UDP，因此自动降级为 WebSocket 投屏。
>
> v3.6.7+ 增强了外网自动适配：自动检测并切换 WebSocket 投屏模式，若切换失败则提供 MJPEG 降级画面，确保外网访问始终可用。

### 飞牛 APP

通过飞牛手机 APP 内置浏览器访问时：
- ⚠️ 部分版本的飞牛 APP WebView 对 WebSocket 代理支持有限，可能出现「等待设备推流超时」
- ✅ 建议在手机浏览器（Chrome / Edge / 系统浏览器）中直接打开飞牛远程域名使用

### 公网 IP / 内网穿透 / 反向代理

如果你有公网 IP 或使用内网穿透（如 frp、ZeroTier、Tailscale 等）：
- ✅ 将 8443 端口映射到外网，可直接使用 WebRTC 投屏
- ✅ 需要在「投屏设置」中将「外网访问地址（WebRTC 媒体流）」设置为你的公网域名或 IP
- ✅ TURN 中继服务器已内置，确保 UDP 3478 端口可达

#### 外网访问专项说明（重要）

**飞牛官方远程域名（xxx.fnos.net）的限制：**
- 飞牛反向代理**只透传 TCP**，WebRTC 的 UDP 媒体流无法通过
- 因此通过飞牛远程域名访问时，**只能使用 WS 投屏**（走 TCP），WebRTC 投屏会失败
- 文件管理、终端（ADB）等功能依赖 WebRTC DataChannel（UDP），**外网不可用**
- WS 投屏连接时会先尝试 WebRTC（UDP），超时后才回退到 WS，因此**首次连接可能需要等待 10-30 秒**

**推荐的外网访问方案（功能最完整）：**

使用 frp 或其他内网穿透工具，**同时暴露 TCP 8443 和 UDP 3478**：

```ini
# frpc.ini 示例
[androidemu_web]
type = tcp
local_ip = 127.0.0.1
local_port = 8443
remote_port = 8443

[androidemu_turn]
type = udp
local_ip = 127.0.0.1
local_port = 3478
remote_port = 3478
```

配置后通过 `http://<你的域名>:8443` 访问，WebRTC 投屏、文件管理、终端等全部功能可用。

**WS 投屏连接慢的解决方法：**
1. 耐心等待 10-30 秒，WebRTC 超时后会自动回退到 WS
2. 或在连接页面手动点击「切换为 WebSocket 投屏」按钮，立即使用 WS
3. 使用 frp 暴露 UDP 3478 后，WebRTC 可直接连接，无需等待

---

## 默认账号

穿云投屏画面服务默认登录账号：

| 项目 | 值 |
|------|-----|
| 用户名 | `admin` |
| 密码 | `admin123` |

> 通过飞牛统一网关访问时会自动注入登录态，无需手动输入。直接访问 8443 端口时需要手动登录，登录之后请更改账号和密码。

---

## 快速上手（穿云投屏操作指南）

### 1. 登录

- 通过飞牛应用中心点击「打开」进入，自动登录，无需输入账号
- 直接访问 `http://<NAS_IP>:8443` 时，使用默认账号 `admin` / `admin123` 登录
- 登录后建议在「设置」中修改密码

### 2. 主界面

左侧菜单栏：

| 菜单 | 功能 |
|------|------|
| 云虚机 | 查看和管理已连接的安卓设备，点击设备进入投屏画面 |
| 大盘 | 设备状态总览（在线数、CPU、内存等） |
| 文件 | 文件中心，批量安装/传输APK和文件 |
| 部署 | 设备部署配置 |
| 终端 | 安卓 shell 命令行（ADB调试） |
| 外设 | 外设管理 |
| 用户管理 | 管理穿云投屏的登录用户 |
| 设备运营 | 设备分组、标签管理 |
| 分享 | 生成设备分享链接 |
| 审计 | 操作审计日志 |
| 设置 | 系统设置（修改密码、端口等） |

### 3. 查看设备

- 点击左侧「云虚机」，右侧会显示已注册的安卓设备列表
- 设备状态：**在线**（绿色）/ **离线**（灰色）
- 正常安装后，设备会自动注册并显示为在线（设备ID默认为 `androidemu`）
- 如果显示「0 台在线」，参考常见问题中的排查方法

### 4. 投屏操作

1. 在「云虚机」中点击在线的设备，进入投屏画面
2. 首次连接会自动选择 **WebRTC** 模式（画质好、延迟低）
3. 如果 WebRTC 连接失败，画面顶部会有模式切换按钮，可切换为 **WebSocket（WS）** 模式
4. 画面操作：
   - **单击**：鼠标左键点击 = 安卓触摸
   - **滑动**：按住鼠标左键拖动 = 安卓滑动手势
   - **缩放**：鼠标滚轮 = 双指缩放
   - **返回**：画面侧边工具栏的返回按钮，或按键盘 Esc
   - **主页**：侧边工具栏主页按钮，或按键盘 Home
   - **多任务**：侧边工具栏多任务按钮
   - **键盘输入**：直接用键盘打字，会输入到当前聚焦的输入框
   - **音量调节**：侧边工具栏音量加减按钮

### 5. 安装 APK

**方法一：批量安装（推荐）**
1. 点击左侧「文件」→ 选择「批量安装/传输」标签
2. 拖拽APK文件到上传区域，或点击上传
3. 选择目标设备（勾选要安装的设备）
4. 任务类型选「静默安装 APK」
5. 点击「下发批量任务」，等待安装完成

**方法二：单设备安装**
1. 进入设备投屏画面
2. 在侧边工具栏找到「安装APK」按钮
3. 选择本地APK文件上传安装

> 安装完成后，上传的APK文件不会自动删除，需在「文件」页面手动点「移除」，或参考常见问题中的清理方法。

### 6. 终端（ADB调试）

1. 点击左侧「终端」
2. 选择目标设备
3. 直接输入安卓 shell 命令，常用命令：

   列出已安装应用：
   ```
   pm list packages
   ```

   卸载应用（将 `<包名>` 替换为实际包名）：
   ```
   pm uninstall <包名>
   ```

   查看安卓是否启动完成：
   ```
   getprop sys.boot_completed
   ```
4. 也可以在外部用 `adb connect <NAS_IP>:5556` 连接

### 7. 文件管理

- 「文件」→「文件中心」可以浏览和管理安卓设备内的文件
- 支持上传文件到设备、从设备下载文件、删除文件
- 「批量安装/传输」用于向多台设备同时下发APK或文件

### 8. 使用手机 APP 控制（可选）

穿云投屏官方提供独立的 **Android APP 客户端**，手机端体验优于浏览器（真正全屏、无地址栏、后台保活）。

**下载安装：**
1. 手机浏览器访问 https://webrtc-phone.com/#download
2. 下载 `ScrcpyOverWebRTC-release.apk` 并安装
3. 打开 APP，在地址栏输入你的访问地址：
   - 局域网：`http://<NAS_IP>:8443`
   - 外网：你的飞牛远程域名或反向代理地址
4. 使用默认账号 `admin` / `admin123` 登录（或你修改后的账号）
5. 点击设备即可进入投屏控制

> APP 与浏览器访问的是同一个服务端，数据和配置完全同步。APP 的优势在于移动端体验优化和后台保活，基础功能与浏览器一致。

### 9. 常见操作速查

| 操作 | 方法 |
|------|------|
| 截图 | 投屏画面侧边工具栏截图按钮 |
| 旋转屏幕 | 侧边工具栏旋转按钮 |
| 锁屏/唤醒 | 单击侧边工具栏电源键按钮 |
| 关机/重启 | 长按侧边工具栏电源键按钮，弹出安卓电源菜单后选择 |
| 切换投屏模式 | 画面顶部 WebRTC/WS 切换按钮 |
| 修改密码 | 左侧「设置」→ 修改密码 |
| 分享设备 | 左侧「分享」→ 生成分享链接 |

---

## ADB 连接

> **重要说明**：ADB 是 Android 调试通道，**只用于执行命令、安装 APK、调试，不传输画面**。查看安卓画面请通过穿云投屏 Web 界面（`http://<NAS_IP>:8443`）或穿云投屏 APP。Scrcpy / QtScrcpy 等依赖 scrcpy-server 的远程控制软件**不兼容**（本容器使用穿云投屏 cloudphone-agent，无 scrcpy-server），会出现一直转圈无画面的情况。

ADB 5556 **默认仅监听 127.0.0.1**（安全考虑，ADB 无密码），NAS 本机可直接连接。如需从局域网其他设备（如电脑）连接，需先手动开放端口：

### 一、开放 ADB 5556 端口

**第 1 步：SSH 登录 NAS，编辑配置文件**

```bash
nano /var/apps/androidemu/var/ports.conf
```

添加或修改以下内容：

```
ADB_BIND=0.0.0.0
```

> nano 操作：按 `Ctrl+O` 保存，按 `回车` 确认文件名，按 `Ctrl+X` 退出。

**第 2 步：重启 ADB 转发进程**

```bash
sudo pkill -f redroid_adb_forward.sh && sudo bash /vol1/@appcenter/androidemu/scripts/redroid_adb_forward.sh install
```

**第 3 步：验证端口是否开放**

```bash
ss -tln | grep 5556
```

应显示 `0.0.0.0:5556`，表示已监听所有网卡。

### 二、连接安卓容器

**第 4 步：从电脑连接**

将 `<NAS_IP>` 替换为你的 NAS 实际局域网 IP 地址：

```bash
adb connect <NAS_IP>:5556
```

连接成功会显示 `connected to <NAS_IP>:5556`。

**第 5 步：进入安卓 shell**

```bash
adb shell
```

进入后即可执行安卓命令（如 `pm list packages` 查看已安装应用）。

> 也可以不连 ADB，直接在穿云投屏界面的「终端」中使用安卓 shell。

> ⚠️ **安全提醒**：ADB 无密码验证，用完后建议改回 `ADB_BIND=127.0.0.1` 并重启转发，避免端口长期暴露。

---

## APK 安装

1. 在穿云投屏界面点击「文件」或「上传」
2. 选择 APK 文件上传
3. 在安卓容器中点击文件管理器中的 APK 进行安装

> **注意**：内置模拟器为 x86_64 架构，不含谷歌服务。镜像已内置 ARM 翻译层，大多数 ARM 应用可运行；但强依赖谷歌服务、或含反模拟器检测/复杂 JIT 的 ARM64 应用可能启动即闪退（属翻译层能力边界）。此类应用建议通过穿云投屏 Agent 接入真机使用，或者根据上游 redroid 容器作者的谷歌服务推荐配置来进行。

---

## 性能优化

本应用内置多项性能优化，降低访问延迟和使用延迟：

### 1. 进程优先级提升

webrtc 信令服务和 TURN 中继服务均以 `nice=-10` 启动（高于默认优先级），减少 CPU 调度延迟，避免画面卡顿。

> 技术细节：Docker 容器默认 drop `CAP_SYS_NICE`，本应用在 compose 中显式添加 `cap_add: SYS_NICE`，该 capability 仅影响进程调度优先级，不涉及设备/网络/文件系统访问。

### 2. 后台服务精简

容器启动后自动禁用以下不必要的安卓系统服务，释放 CPU 和内存：
- 蓝牙服务（`com.android.bluetooth`、`com.android.bluetoothmidiservice`）— 容器无蓝牙硬件
- NFC 服务（`com.android.nfc`）— 容器无 NFC 硬件
- 打印服务（`com.android.printspooler`）— 容器无需打印
- 备份服务（`com.android.backupconfirm`、`com.android.sharedstoragebackup`）— 容器无需备份

> **注意**：`com.android.location.fused`（位置融合服务）不能禁用，它是 LocationManagerService 的依赖项，禁用会导致 system_server 崩溃，安卓无法启动。

### 3. CPU 动态调频

安装时自动配置 CPU 调度器为动态调频模式（schedutil/ondemand），并设置最低频率为最大频率的 40%：
- 轻载时 CPU 不会降频过低，保证操作响应速度
- 重载时自动升到最高频率，发挥全部性能
- 仅在支持 cpufreq 的设备上生效，不支持则静默跳过
- 不低于四核 CPU 的极限（动态调整，不会过度限制性能）

### 4. 串口控制台默认关闭

安卓串口控制台（`androidboot.console=0`）默认关闭，减少内核日志输出带来的性能损耗。开发者如需调试可手动开启（见下文）。

### 5. GPU 硬件加速

- **X86 设备**：自动检测 `/dev/dri`，有 GPU 时使用 `gpu_mode=host` 硬件加速，60fps
- **ARM 设备**：自动使用 `gpu_mode=guest` 软件渲染，60fps（统一高帧率，低配设备如感觉卡顿可在 compose 中改回 30fps）

### 6. 国内 APP 专项优化

针对抖音、快手等国内 APP 后台服务多、占用高的问题，做了以下系统级优化：

| 优化项 | 说明 |
|--------|------|
| **限制后台进程** | `background_process_limit=4`，超过后系统更积极回收 |
| **禁用自动同步** | `auto_sync=0`，减少后台同步开销 |
| **禁用后台数据** | `mobile_data_always_on=0`，防止后台偷跑流量和 CPU |
| **禁用扫描** | WiFi/BLE 扫描关闭，减少定位和后台唤醒 |
| **禁用网络优化** | 关闭网络推荐、自适应连接等后台服务 |
| **动画半速** | 窗口/过渡/动画师缩放设为 0.5，既流畅又省 CPU |
| **GPU 强制渲染** | `debug.hwui.renderer=skiagl`，用 GPU 渲染 UI |
| **图层合成优化** | `disable_backpressure=1`、`latch_unsignaled=1`，减少合成延迟 |
| **内存管理优化** | 调整 lmkd 阈值，更积极回收后台 APP 内存 |

### 7. 可选资源限制

如果 NAS 性能有限，可在 `docker-compose.yaml` 中启用资源限制（默认注释，取消注释即可）：

```yaml
deploy:
  resources:
    limits:
      cpus: '4.0'      # 最多使用4核
      memory: 4G       # 最多使用4GB内存
```

> 修改后需重启容器生效：`docker compose -f /vol1/@appcenter/androidemu/app/docker/docker-compose.yaml restart`

---

## 串口控制台（开发者调试）

串口控制台默认关闭以减少性能损耗。如需开启用于调试：

**第 1 步：开启串口控制台**（临时开启，容器重启后失效）

```bash
docker exec -u 0 androidemu-android setprop persist.sys.serialconsole 1
```

```bash
docker exec -u 0 androidemu-android start console
```

**第 2 步：查看控制台输出**

```bash
docker exec -u 0 androidemu-android dmesg -w
```

> 按 `Ctrl+C` 退出实时日志查看。

**第 3 步：调试完成后关闭**

```bash
docker exec -u 0 androidemu-android stop console
```

```bash
docker exec -u 0 androidemu-android setprop persist.sys.serialconsole 0
```

> 注意：开启串口控制台会增加内核日志输出，可能略微影响性能，调试完成后建议关闭。

---

## 容器架构

```
┌─────────────────────────────────────────────┐
│  飞牛 fnOS 主机                              │
│                                             │
│  ┌──────────────────┐  ┌─────────────────┐  │
│  │  androidemu-     │  │  androidemu-    │  │
│  │  android         │  │  webrtc         │  │
│  │  (redroid)       │  │  (scrcpy+TURN)  │  │
│  │                  │  │                 │  │
│  │  Android 12      │  │  信令服务 :8443  │  │
│  │  ADB :5556       │  │  TURN  :3478    │  │
│  │  bridge 网络     │  │  host 网络       │  │
│  │  privileged      │  │  SYS_NICE       │  │
│  └──────────────────┘  └─────────────────┘  │
│            │                    │           │
│            └────── scrcpy ──────┘           │
│                  (ADB over TCP)             │
└─────────────────────────────────────────────┘
                       │
                       ▼
                 飞牛统一网关
                       │
                       ▼
                    浏览器
```

---

## 技术栈与翻译层

androidemu 从硬件到浏览器画面共经过 **3 个核心翻译/转换层**，另有 1 个可选的指令集翻译层：

### 第 1 层：容器化层（Docker）

- 不是全虚拟机，而是**进程级容器隔离**，Android 用户态直接运行在宿主 Linux 内核上
- 与宿主共享同一个内核，**不翻译 CPU 指令**，性能接近原生
- 提供文件系统、网络、进程隔离；安卓容器以 `privileged` 模式运行（redroid 上游官方要求，用于 binder 设备访问）

### 第 2 层：GPU 渲染翻译层

| 模式 | 适用场景 | 原理 | 帧率 |
|------|----------|------|------|
| GPU 直通（guest） | X86 有核显/独显 | Android 的 OpenGL ES 指令直接发给宿主 GPU 驱动，几乎无翻译开销 | 60fps |
| 软件渲染（swiftshader） | 无 GPU / ARM 设备 | **swiftshader** 把 OpenGL ES 指令翻译成 CPU 指令执行，有翻译开销 | 60fps（低配设备可改回 30fps） |

- 安装脚本自动检测宿主 GPU 能力，有 `/dev/dri` 时用 GPU 直通，否则自动回退软件渲染
- ARM 设备默认使用软件渲染（swiftshader）

### 第 3 层：画面采集与编码层

- **scrcpy** 通过 Android 的 surfaceflinger 采集画面帧
- 编码成 **H.264** 视频流（码率 2-10Mbps，长边 960px）
- 通过 **WebRTC**（TURN/STUN 中继 + P2P）传输到浏览器
- 浏览器解码后显示，**音频已启用**（3.7.3+ 修复：通过启用 Codec2 框架加载 c2.android.opus.encoder 软件编码器）

### 第 4 层：ABI 指令集翻译层（libndk_translation，默认已启用）

- redroid 镜像默认内置 **libndk_translation**（Google 官方 NDK 翻译方案），通过 Native Bridge 机制实现 ARM→x86 二进制翻译
- 已验证配置：`ro.dalvik.vm.native.bridge=libnb.so`（符号链接指向 `libndk_translation.so`）
- 支持 ABI：`x86_64, arm64-v8a, x86, armeabi-v7a, armeabi`（五种架构，ARM 应用可直接运行）
- 镜像内**不包含** libhoudini（Intel 方案）和 QEMU translator（仅有相关属性文件，非实际翻译器）
- X86 设备：Android x86_64 原生 + libndk_translation 翻译 ARM 应用
- ARM 设备：Android arm64 原生运行，无需翻译层

### 完整数据流

```
用户点击浏览器
    │
    ▼
WebRTC 接收 H.264 流 ←── TURN/STUN 中继 ←── scrcpy 编码 ←── surfaceflinger 采集
    │                                                          │
    │                                                          ▼
    │                                                   Android 12 (redroid)
    │                                                          │
    │                                                          ▼
    │                                                   GPU 渲染翻译层
    │                                                   (GPU直通 / swiftshader)
    │                                                          │
    ▼                                                          ▼
浏览器显示 ←──── 飞牛统一网关 ←──── Docker 容器 ←──── 宿主 Linux 内核
```

---

## 审核合规性说明

本应用已通过飞牛官方 7 条自查（基本信息、权限声明、网络端口、数据存储、启动停止、卸载清理、兼容性），以下为审核关注要点的详细说明：

### 1. 特权容器（privileged）

- **现状**：仅 `androidemu-android`（redroid 安卓主容器）使用 `privileged: true`，`androidemu-webrtc`（画面服务）为非特权运行
- **必要性**：redroid 上游（remote-android/redroid-doc）官方部署方式即要求 `--privileged`，Android 依赖内核 binder 通信，容器需要挂载/访问 binder 设备并创建设备节点
- **非特权方案已实测不可行**：非特权 + device_cgroup_rule + 挂载 binderfs + cap-add=ALL + seccomp=unconfined，容器以 ExitCode 0 静默退出、无法开机（redroid-doc issue #591 至今为开放议题）
- **影响面控制**：特权仅作用于容器内部（容器内 root = Android 系统自身初始化所需），不等于宿主 root；应用本体以 `docker-androidemu` 用户运行，不申请宿主 root

### 2. 应用本体非 root 运行

- gateway.py、audio_fix.py 等后台进程均以 `docker-androidemu` 用户运行
- 通过 `_drop_privileges()` 函数实现自动降权（uid=0 时自动切换）
- config/privilege 声明 `run-as=package`

### 3. Host 网络模式

- `androidemu-webrtc` 容器使用 `network_mode: host`
- **原因**：fnOS 会拦截「bridge 容器 → 宿主局域网 IP」的访问，bridge 模式下云手机后台无法连接本机 ADB 端口
- **影响面**：监听端口与原先端口映射完全一致（8443 TCP、3478 TCP/UDP、50000-50100 UDP），不新增任何端口；容器仍不使用 privileged

### 4. Docker 组权限

- 应用用户 `docker-androidemu` 属于 `docker` 组
- 这是飞牛 Docker 应用的标准配置，运行容器必须
- 仅用于应用脚本对本应用容器的编排与清理，不涉及其他应用

### 5. 多端口监听

| 端口 | 协议 | 用途 | 鉴权 |
|------|------|------|------|
| 8443 | TCP | WebRTC 信令+画面 | 账号登录 |
| 3478 | TCP/UDP | TURN 中继 | TURN 凭据 |
| 5556 | TCP | ADB 调试 | 默认仅本机，需手动开放 |
| 50000-50100 | UDP | WebRTC 媒体流 | 会话级鉴权 |

- 应用侧不做任何自动端口映射（无 UPnP/打洞），公网暴露由用户自行决定
- 各端口均有独立鉴权机制

### 6. 付费说明

- 本应用自身完全免费
- 内置的穿云投屏画面服务采用上游第三方授权：免费版即可使用云手机画面、ADB 调试等全部基础功能，仅「可添加设备数量」受限
- 付费仅增加设备数量上限，不影响任何功能
- 费用由上游授权服务方收取，与飞牛官方无关
<img width="1288" height="900" alt="firefox exe_20260927_092044" src="https://github.com/user-attachments/assets/44afdf11-2112-4ae7-8544-89e6bfa0238b" />

---

## 常见问题

### Q: 手机飞牛 APP 应用中心安装时弹出「无法安装 安卓模拟器（国内版）- 未知错误」

A: 这是飞牛 **手机 APP 端**应用中心的通用报错，发生在应用中心层面（还未执行安装脚本），不是本应用的问题。

**可能原因与解决方法：**
1. **优先用网页版安装**：在电脑浏览器打开飞牛管理页面，通过网页版应用中心安装。手机 APP 端应用中心偶发此类「未知错误」，网页版通常正常
2. **检查磁盘空间**：安卓镜像约 2GB，加上临时文件需至少 4GB 空闲空间
3. **重启飞牛 APP**：完全关闭 APP 后重新打开再试
4. **手动安装**：下载 fpk 安装包，通过网页版「手动安装」上传安装

若以上方法均无效，请通过下方「问题反馈」章节反馈，并附上 NAS 型号、系统版本和截图。

### Q: 手机飞牛 APP 安装时弹出「无法安装 - 状态操作不支持，并返回当前应用状态和业务状态」

A: 这是飞牛 **手机 APP 端**应用中心的状态机错误，通常是因为应用处于异常状态（如上次安装/更新中途取消、重复点击安装按钮、安装进程残留）导致状态冲突。

**解决方法：**
1. **不要重复点击**：等待当前操作完成，安装过程中不要反复点击安装/更新按钮
2. **重启飞牛 APP**：完全关闭 APP 后重新打开，刷新应用状态
3. **网页版操作**：在电脑浏览器打开飞牛管理页面，通过网页版应用中心进行安装/更新/卸载操作
4. **清理残留状态**：如果应用显示「已安装」但实际无法使用，先在网页版中卸载，再重新安装

### Q: 首次安装时弹出「无法安装 androidemu - 执行脚本出错且原因未知」

A: 优先检查 **binder 驱动**是否已安装：
- **x86 设备**：需先在飞牛应用中心安装「binder_linux 驱动」依赖应用，产生 `/dev/binder` 设备节点后再安装本应用
- **ARM 设备**：多数设备内核已内置 binder（如 RK3588 等），但部分精简内核可能未启用，需确认内核支持 `CONFIG_ANDROID_BINDER_IPC` / binderfs

若已安装驱动（或 ARM 设备本身支持）但仍弹出此错误，请务必通过下方「问题反馈」章节中的任一链接反馈，以便排查具体原因。

### Q: 安装后打不开，页面显示 400 错误（可能多个应用同时出现）

A: 这是飞牛系统 nginx 网关的请求头大小限制导致的，**不是本应用的问题**。

**原因**：飞牛 nginx 默认 `large_client_header_buffers 4 8k`，当浏览器 cookie 过大（飞牛网关注入的 `fnos-token`、`osrt` 等 cookie 累计超过 8KB）时，nginx 会直接返回 400。此问题会同时影响所有走飞牛网关的应用（如 gxsales、hddlocator 等），清除浏览器缓存后短暂恢复，但重新登录后 cookie 又会变大。

**解决方案 A（推荐，立即可用）：绕开飞牛 nginx，直接访问端口**

局域网内直接访问：
```
http://<NAS_IP>:8443
```

或通过自己的反向代理（如 Lucky、Nginx Proxy Manager）访问，不经过飞牛网关。

**解决方案 B（根治）：修改飞牛 nginx 配置，增大请求头缓冲区**

> 需要 SSH 登录 NAS 并具有 root 权限，修改前请备份配置文件。

第 1 步：备份 nginx 配置
```bash
sudo cp /usr/trim/nginx/conf/nginx.conf /usr/trim/nginx/conf/nginx.conf.bak
```

第 2 步：编辑配置文件
```bash
sudo nano /usr/trim/nginx/conf/nginx.conf
```

在 `http {}` 块中添加以下两行：
```
client_header_buffer_size 16k;
large_client_header_buffers 4 32k;
```

第 3 步：测试配置并重新加载
```bash
sudo nginx -t && sudo nginx -s reload
```

> 注意：飞牛系统更新可能会覆盖 nginx 配置，若更新后问题复现，重新添加即可。

### Q: 打开后显示「Bad Gateway」（502 错误）

A: 这是飞牛 nginx 网关找不到后端服务时的默认提示，**不是本应用的报错页面**。

**最常见原因：刚安装/刚启动，后端服务还没起来**

安卓容器首次启动需要 1-2 分钟（下载镜像、初始化、启动系统服务），此时 gateway.py 或 webrtc 容器尚未就绪，飞牛网关直接返回 502。**等待 1-2 分钟后刷新页面即可**。

**按版本区分：**

- **v3.6.8 及以下版本**：显示纯文本「Bad Gateway」是正常现象，旧版本没有友好状态页，上游未就绪时就是这个样子。等容器启动完成后刷新即可正常。
- **v3.7.0 及以上版本**：正常情况下上游未就绪时会显示友好状态页（带容器状态表格和刷新按钮）。如果仍然显示纯文本「Bad Gateway」，说明 gateway.py 进程根本没启动，需要排查：

第 1 步：检查 8443 端口是否在监听
```bash
ss -tln | grep 8443
```

第 2 步：检查 gateway.py 进程是否在运行
```bash
ps aux | grep gateway.py | grep -v grep
```

第 3 步：如果端口没在监听或进程不存在，在应用中心点击「停用」再「启用」，或直接重启应用。

第 4 步：也可以绕开飞牛网关直接访问，确认是不是网关层的问题：
```
http://<NAS_IP>:8443
```

### Q: 打开后显示「未授权」或「0 台在线」

A: 这是穿云投屏的 License 授权提示。新安装的设备目前有 20 台设备三个月免费试用期（自 2026 年 11 月 1 日，到期为 10 台设备，其他基础功能均为免费），等待容器完全启动（约 1-2 分钟）后刷新页面即可。若持续未授权，请检查容器是否正常运行：
```bash
docker ps --filter name=androidemu
```

### Q: 升级时提示「无法更新 - 执行脚本出错且原因未知」

A: 这是旧版本（v3.6.5 及之前）`uninstall_init` 脚本在升级流程中执行异常导致的。v3.6.6 已修复：
- 重写卸载脚本，增加更健壮的升级守卫
- 添加 `set +e` 确保任何命令失败都不会导致脚本异常退出
- 所有输出重定向到日志，stdout 保持干净

**解决方案**：下载 v3.6.6+ 版本，通过「手动安装」覆盖安装即可。若仍失败，可在 NAS 上查看日志：
```bash
cat /var/apps/androidemu/var/uninstall_init.log
```

### Q: WebRTC 投屏连接失败

A: 按以下步骤排查：
1. 确认在局域网内使用（外网自动降级为 WebSocket 投屏）
2. 确认 TURN 服务器正常运行：`docker exec androidemu-webrtc netstat -ulnp | grep 3478`
3. 确认 PUBLIC_IP 是局域网 IP 而非 127.0.0.1：`docker exec androidemu-webrtc env | grep PUBLIC_IP`
4. 确认 ICE_SERVERS 中包含正确的局域网 IP：`docker exec androidemu-webrtc env | grep ICE_SERVERS`
5. 尝试点击「改用 WebSocket 投屏」

### Q: WebSocket（WS）投屏显示连接失败

A: 使用WebSocket投屏（局域网或外网环境）时，若出现「连接失败」页面，请耐心等待：
- **2分钟内**：系统会自动重试并连接，最终显示安卓容器主页，无需任何操作
- **超过2分钟仍显示连接失败**：请按以下步骤排查：
  1. 确认安卓容器已完全启动（首次启动约需1-2分钟），可在飞牛Docker中查看容器状态
  2. 刷新页面重新进入
  3. 局域网用户可尝试直接访问 `https://<NAS_IP>:8443`
  4. 外网用户可使用页面右上角提示中的MJPEG降级画面查看
  5. 如仍无法解决，请查看容器日志或通过反馈渠道联系我们

### Q: 容器反复重启

A: 常见原因及排查：

**1. 查看安卓容器日志**

```bash
docker logs androidemu-android
```

**2. 根据日志关键词判断原因：**

| 日志关键词 | 原因 | 解决方法 |
|-----------|------|---------|
| 与 `gpu`、`dri`、`SurfaceFlinger` 相关的报错 | GPU 直通失败 | 见下文「安卓容器内存很低」中的 GPU 修复命令 |
| `out of memory`、`lowmemory` | 内存不足 | 关闭其他应用，至少保留 2GB 可用内存 |
| `binder`、`binderfs` 相关报错 | binder 驱动缺失 | x86 设备先在应用中心安装 `binder_linux` 驱动 |

**3. ARM 设备额外检查：** 确认 compose 中包含 `androidboot.redroid_gpu_mode=guest`（软件渲染），ARM 设备通常无 GPU 直通。

**4. webrtc 容器反复重启：** 在日志中搜索是否有 `nice: setpriority(-10): Permission denied`，如果有，确认 compose 中包含 `cap_add: SYS_NICE`。

查看 webrtc 容器日志：
```bash
docker logs androidemu-webrtc --tail 50
```

### Q: 终端执行 docker 命令报 `permission denied while trying to connect to the Docker daemon socket`

A: 当前用户不在 docker 组中，没有权限直接访问 Docker daemon。

**解决方法（任选一种）：**
1. **临时方案**：在所有 docker 命令前加 `sudo`，例如 `sudo docker ps`、`sudo docker exec androidemu-android ...`
2. **永久方案**：将用户加入 docker 组，重新登录后生效：
```bash
sudo usermod -aG docker <你的用户名>
```
> 注意：加入 docker 组后需要**退出终端重新登录**才能生效。

### Q: 如何设置多端分辨率并动态切换？

A: v3.8.4 起支持多端分辨率配置。安装向导中可分别填写**电脑端、手机端、平板端**的分辨率，容器默认以手机端分辨率启动，使用过程中可通过脚本动态切换（即时生效，无需重启容器）。

**1. 安装时设置（推荐）**

安装向导的「分辨率设置」步骤中分别填写三端分辨率（格式：宽×高，如 1920×1080）：
- **电脑端**：如 1920×1080（横屏）或 1280×720
- **手机端**：如 720×1440（9:18 全面屏，默认）或 720×1280（9:16）
- **平板端**：如 1200×2000（3:5）或 800×1280

留空则该端使用默认值（电脑 1920×1080、手机 720×1440、平板 1200×2000）；三端都不填则全部使用默认。

**2. 自动切换分辨率**

安装后应用会自动启动「分辨率自动切换守护」，当你在不同设备上打开穿云投屏页面时，页面会自动检测设备类型并上报，守护脚本接收后**自动切换到对应分辨率**（即时生效）。

- 电脑端打开 → 自动切换到电脑端分辨率
- 手机端打开 → 自动切换到手机端分辨率
- 平板端打开 → 自动切换到平板端分辨率

> 冷却时间 60 秒：切换后 60 秒内不会再次切换，避免频繁切换影响体验。多个设备同时打开时，以最新打开的设备为准。

**查看自动切换日志：**
```bash
cat /var/apps/androidemu/var/resolution_autoswitch.log
```

**3. 手动切换分辨率**

如果自动切换不符合预期，可手动执行：
```bash
# 切换到电脑端分辨率
bash /vol1/@appcenter/androidemu/scripts/switch_resolution.sh pc

# 切换到手机端分辨率
bash /vol1/@appcenter/androidemu/scripts/switch_resolution.sh phone

# 切换到平板端分辨率
bash /vol1/@appcenter/androidemu/scripts/switch_resolution.sh tablet

# 直接指定自定义分辨率
bash /vol1/@appcenter/androidemu/scripts/switch_resolution.sh 1080x1920
```

切换通过 `adb shell wm size` 实现，**即时生效**，不需要重启容器或重新连接。

**3. 查看当前配置**
```bash
cat /var/apps/androidemu/var/resolution.conf
```

**4. 永久修改默认分辨率（修改 compose）**

如果需要修改容器启动时的默认分辨率（重启后生效），编辑 compose 文件：
```bash
nano /vol1/@appcenter/androidemu/docker/docker-compose.yaml
```
找到 `redroid` 服务的 `command` 部分，修改：
```yaml
- androidboot.redroid_width=720
- androidboot.redroid_height=1600
```
然后重启容器：
```bash
cd /vol1/@appcenter/androidemu/docker
docker compose up -d --force-recreate redroid
```

> **注意**：
> - 容器同一时间只能使用一个分辨率，切换后所有连接的客户端都会看到新分辨率
> - 动态切换不丢失容器内数据和已安装应用
> - 分辨率比例与设备屏幕一致时，可实现沉浸式真全屏（无黑边）
> - 1080p 及以上分辨率对 NAS 性能要求较高，低配设备建议 720p 系列

### Q: 容器卡顿、无法点击或移动、一动不动

A: **v3.7.0+ 用户**：打开应用页面，如果上游服务暂时不可用，会自动显示友好状态页，其中包含：
- **健康状态**：自动检测 boot 状态、运行时长、是否 OOM、surfaceflinger/agent 是否运行
- **一键修复按钮**：
  - 「修复GPU/画面」：自动 chmod /dev/dri + 重启 surfaceflinger（解决 GPU 权限问题导致的画面卡住）
  - 「重启安卓容器」：重启整个安卓容器
  - 「重启画面服务」：只重启 surfaceflinger，不影响容器内其他进程

**所有版本通用排查步骤：**
1. 先关闭软件页面（或网页），重新打开后再尝试点击/移动
2. 若仍无效，在飞牛应用中心的软件详情页点击「停用」，停用后再「启用」
3. 如果只有单个容器出现卡顿，可在 Docker 中找到对应容器点击「重启」即可
4. 手动检查健康状态（逐条执行）：

   查看容器是否在运行：
   ```bash
   docker ps --filter name=androidemu
   ```

   查看安卓是否启动完成：
   ```bash
   docker exec androidemu-android getprop sys.boot_completed
   ```

   查看是否被OOM杀死：
   ```bash
   docker inspect -f '{{.State.OOMKilled}}' androidemu-android
   ```
5. 若以上方法均无效，请将软件卡顿的截图或录屏，以及 Docker 容器中复制的日志打包为文本文档，通过下方反馈渠道任选一项进行反馈

### Q: 安卓容器内存很低（<300MB）、设备一直不在线

A: 正常 Android 12 启动后内存应在 300MB 以上（单个安卓容器内的系统占用，不含宿主其他服务）。若容器在运行但内存只有 100-200MB，说明安卓系统未完成启动（`boot_completed != 1`），agent 无法部署，设备永远不在线。按以下步骤排查：

**1. 确认启动状态：**
```bash
docker exec androidemu-android getprop sys.boot_completed
```
- 返回 `1` → 已启动，跳到第3步
- 返回空或 `0` → 未启动，继续第2步

**2. 查看启动日志找原因：**
```bash
docker logs androidemu-android --tail 80 2>&1 | grep -iE "error|fail|panic|binder|surface|zygote|boot"
```

| 日志关键词 | 原因 | 修复方法 |
|-----------|------|---------|
| `binder`、`binderfs` | binder 驱动缺失 | x86 设备先在应用中心安装 `binder_linux`，再 `docker restart androidemu-android` |
| `SurfaceFlinger`、`gpu`、`dri` | GPU 直通失败 | 执行下方 GPU 修复命令 |
| `out of memory`、`lowmemory` | 内存不足 | 关闭其他应用，至少保留 2GB 可用内存 |
| `zygote` 反复重启 | 系统服务崩溃 | 删除数据卷重新初始化（会清空安卓数据）：`docker compose -p androidemu down && docker volume rm androidemu-data && docker compose -p androidemu up -d` |

**3. GPU 修复命令（画面异常或 SurfaceFlinger 崩溃时用，按顺序执行）：**

第 1 步：修复 GPU 设备权限
```bash
docker exec -u 0 androidemu-android chmod 666 /dev/dri/card0 /dev/dri/renderD128
```

第 2 步：重启画面服务
```bash
docker exec -u 0 androidemu-android sh -c 'setprop ctl.restart surfaceflinger'
```

第 3 步：等待 60 秒让服务重启完成
```bash
sleep 60
```

第 4 步：确认安卓启动完成
```bash
docker exec androidemu-android getprop sys.boot_completed
```

**4. 手动注入 agent（boot_completed=1 但设备仍不在线时用，按顺序执行）：**

> 注意：以下命令中的 `<NAS_IP>` 需替换为你的 NAS 实际局域网 IP；x86 设备用 `amd64`，ARM 设备用 `arm64`。
>
> **权限提示**：如果执行 docker 命令时出现 `permission denied while trying to connect to the Docker daemon socket`，说明当前用户不在 docker 组中。请在所有 docker 命令前加 `sudo`，或将用户加入 docker 组（`sudo usermod -aG docker <用户名>`，重新登录后生效）。

第 1 步：从穿云投屏镜像取出 agent（首次需要）
```bash
docker run --rm -v /var/apps/androidemu/var/agent:/out --entrypoint /bin/sh docker.fnnas.com/buutuu/scrcpy-over-webrtc:latest -c "cp /app/agent_binaries/cloudphone-agent-amd64 /app/agent_binaries/libsys_core.so /out/ && chmod 755 /out/cloudphone-agent-amd64"
```

第 2 步：复制 agent 二进制到安卓容器
```bash
docker cp /var/apps/androidemu/var/agent/cloudphone-agent-amd64 androidemu-android:/data/local/tmp/cloudphone-agent
```

第 3 步：复制依赖库到安卓容器
```bash
docker cp /var/apps/androidemu/var/agent/libsys_core.so androidemu-android:/data/local/tmp/libsys_core.so
```

第 4 步：注入并启动 agent（将 `<NAS_IP>` 替换为你的 NAS 局域网 IP）
```bash
docker exec -u 0 androidemu-android sh -c "chmod 755 /data/local/tmp/cloudphone-agent && export CP_AGENT_JAR=/data/local/tmp/libsys_core.so && nohup /data/local/tmp/cloudphone-agent -signaling wss://<NAS_IP>:8443/register_agent -id androidemu -ice-servers 'turn:cloudphone_user:cloudphone_secure_password@<NAS_IP>:3478?transport=udp,turn:cloudphone_user:cloudphone_secure_password@<NAS_IP>:3478?transport=tcp,stun:<NAS_IP>:3478' -jar /data/local/tmp/libsys_core.so > /data/local/tmp/agent.log 2>&1 &"
```

第 5 步：等待 3 秒
```bash
sleep 3
```

第 6 步：确认 agent 已运行（返回进程号即成功）
```bash
docker exec androidemu-android pidof cloudphone-agent
```

### Q: 安卓容器启动很慢（超过5分钟）、内存波动大、系统不稳定

A: v3.8.3 之前的版本存在 lmkd（低内存杀手）阈值过低的问题：默认最高阈值仅 315MB，对大内存系统过于激进，导致启动期频繁杀空进程、系统服务反复重启，表现为启动慢、内存在 2.0-2.5GB 波动、设备长时间不在线。

**v3.8.3+ 已修复**：lmkd 阈值调高到 512/768/1024/1280/2048/3072MB，dex2oat 首次启动用 verify-only 模式加快启动，启动后自动提升关键进程优先级。升级到 v3.8.3+ 即可解决。

若仍有问题，请检查：
1. 宿主可用内存是否 ≥2GB（`free -h`）
2. 是否有其他应用占用大量内存（如飞牛照片、下载等）
3. 容器日志是否有 OOM 或崩溃（`docker logs androidemu-android --tail 50`）

### Q: WebRTC 连接失败、黑屏或一直转圈

A: 最常见原因是 `PUBLIC_IP` 被重置为 `127.0.0.1`（重建容器后 compose 默认值生效），导致 TURN 分发给客户端的中继地址是 `127.0.0.1`，客户端连不上。

**诊断命令（逐条执行查看结果）：**

查看 PUBLIC_IP 配置：
```bash
docker inspect androidemu-webrtc --format '{{range .Config.Env}}{{println .}}{{end}}' | grep PUBLIC_IP
```

查看 agent 是否运行：
```bash
docker exec androidemu-android pidof cloudphone-agent || echo "agent 未运行"
```

查看 agent 日志最后20行：
```bash
docker exec androidemu-android tail -20 /data/local/tmp/agent.log 2>/dev/null
```

查看 webrtc 容器日志中的错误：
```bash
docker logs androidemu-webrtc --tail 40 2>&1 | grep -iE 'relay addr|allocation|error|fail|Unauthorized' | tail -15
```

**修复命令（确认 PUBLIC_IP=127.0.0.1 时用，按顺序执行，将 `<NAS_IP>` 替换为你的 NAS 局域网 IP）：**

第 1 步：进入 docker 配置目录
```bash
cd /var/apps/androidemu/target/docker
```

第 2 步：替换 PUBLIC_IP（将 `<NAS_IP>` 替换为你的 NAS 局域网 IP）
```bash
sed -i 's/PUBLIC_IP=127\.0\.0\.1/PUBLIC_IP=<NAS_IP>/g' docker-compose.yaml
```

第 3 步：重建 webrtc 容器
```bash
docker compose -p androidemu up -d --force-recreate webrtc
```

第 4 步：等待 10 秒
```bash
sleep 10
```

第 5 步：确认 PUBLIC_IP 已更新
```bash
docker inspect androidemu-webrtc --format '{{range .Config.Env}}{{println .}}{{end}}' | grep PUBLIC_IP
```

确认 `PUBLIC_IP` 为你的局域网 IP 后，重新打开云手机画面页面连接。

### Q: 没有声音 / 声音很小

A: **3.7.3+ 版本已修复音频功能**，默认启用 Opus 软件编码器。如仍无声音，请检查以下两点：

1. **连接设置中需手动开启音频**：
   - 电脑端（飞牛 Web 界面）：在穿云投屏的连接设置中勾选「音频」选项
   - 手机端（穿云投屏 APP）：音频流透传，无需额外设置，确保 APP 版本支持音频
2. **音量与宿主系统音量联动**：
   - 最终音量 = 容器内音量 × 宿主系统音量
   - 例如：容器内设 100%，但宿主系统音量只设 50%，实际输出按 50% 计算
   - 请同时检查容器内媒体音量和宿主系统音量

3. **使用第三方控制/连接软件时**：
   - 如果使用穿云投屏以外的软件（如 scrcpy、QtScrcpy、ADB 远程控制等）连接，请务必在该软件的设置中**关闭或禁用音频**，或根据软件说明正确配置音频
   - 第三方软件可能不兼容穿云投屏的 Opus 音频流协议，强行开启可能导致连接失败或无声
   - 穿云投屏 APP 和飞牛 Web 界面已内置音频支持，推荐使用官方客户端

> 技术说明：redroid 镜像默认 `debug.stagefright.ccodec=0` 禁用了 Codec2 框架，导致 opus 编码器不加载。3.7.3 通过 `audio_fix.py` 守护进程自动设置 `debug.stagefright.ccodec=1` 启用 Codec2，容器重启后自动补回。

### Q: 外网访问画面卡顿

A: 外网通过飞牛反向代理时自动使用 WebSocket 投屏，帧率受上行带宽限制。如有公网 IP，建议直接映射 8443 端口使用 WebRTC 投屏。

### Q: 点击「终端」弹出打印页面

A: 这是穿云投屏前端的快捷键冲突 bug。已通过 RUNTIME_SHIM 屏蔽 `window.print()` 和 Ctrl+P 快捷键，点击终端不会再触发打印。

### Q: 审计日志加载失败

A: 审计日志功能依赖穿云投屏后端的 `/api/audit` 接口，部分版本可能不支持。此为上游镜像功能，不影响核心投屏功能。

### Q: 穿云投屏上传的APK文件怎么删除？点了「移除」还在

A: 界面上的「移除」只是从当前安装任务里去掉，**不会删除云端文件中心的文件**。文件实际存在穿云投屏容器的 `/app/data/downloads/` 目录里，需要进容器删除：

第 1 步：查看已上传的文件
```bash
docker exec androidemu-webrtc ls -la /app/data/downloads/
```

第 2 步：删除所有APK（也可以指定文件名删除单个）
```bash
sudo docker exec androidemu-webrtc sh -c 'rm -rf /app/data/downloads/*'
```

第 3 步：清空文件元数据记录（否则下拉框还会显示文件名）
```bash
docker exec androidemu-webrtc sh -c 'echo "{}" > /app/data/files_meta.json'
```

删完刷新页面，「从信令云端选择已有文件」下拉框就空了。

### Q: 穿云投屏上传APK后安装失败，怎么手动安装？

A: 穿云投屏上传的APK存放在穿云投屏容器的 `/app/data/downloads/` 目录。如果界面安装失败，可以通过命令行手动安装：

**第 1 步：查看已上传的APK文件**
```bash
sudo docker exec androidemu-webrtc ls -la /app/data/downloads/
```

**第 2 步：把APK复制到宿主机临时目录**
```bash
sudo docker cp androidemu-webrtc:/app/data/downloads/你的应用.apk /tmp/
```

**第 3 步：复制到安卓容器**
```bash
sudo docker cp /tmp/你的应用.apk androidemu-android:/data/local/tmp/
```

**第 4 步：在安卓容器内安装**
```bash
sudo docker exec androidemu-android pm install /data/local/tmp/你的应用.apk
```

**第 5 步（可选）：清理临时文件**
```bash
sudo docker exec androidemu-android rm /data/local/tmp/你的应用.apk
sudo rm /tmp/你的应用.apk
```

> 以上命令已在 x86 和 ARM 平台验证通过。将 `你的应用.apk` 替换为实际的文件名。

### Q: 如何卸载

A: 在飞牛应用中心点击「卸载」即可。容器数据（安卓 /data 分区）会保留在 Docker volume 中，如需彻底清除：
```bash
docker volume rm androidemu_data androidemu-webrtc-data
```

### Q: 安装/更新中途取消后无法重新安装

A: v3.6.0+ 已修复，安装/更新中途取消会自动清理临时数据。若使用旧版本遇到此问题，手动清理：
第 1 步：清理临时文件
```bash
rm -rf /tmp/androidemu_*
```

第 2 步：删除残留容器
```bash
docker rm -f androidemu-android androidemu-webrtc 2>/dev/null
```

### 技术说明与开发踩坑

以下是开发过程中验证过的关键技术结论，供二次开发参考：

- **com.android.location.fused 绝对不能禁用**：它是 LocationManagerService 的依赖项，禁用会导致 system_server 崩溃，安卓无法启动
- **androidboot.use_memfd=1 在 redroid 12 上不安全**：会导致 LocationManagerService 崩溃，init 杀掉 zygote，内存从 600-700MB 降到 400MB，容器无法启动
- **检测依赖应用应检测设备节点而非目录名**：应用中心显示名和实际目录名可能不一致，检测 `/dev/binder` 比检测 `/var/apps/binder_linux_driver` 更可靠
- **WebSocket 代理必须覆盖所有 ws/wss URL**：不能只覆盖特定路径，否则通过网关访问时其他路径会直连 8443 端口导致跨域失败
- **BusyBox su 会重置环境变量**：非 root 改造时不能用 su -c 内部 export 传递环境变量，应直接在 docker-compose environment 中设置
- **禁用蓝牙要用 pm disable 而非 svc disable**：svc disable 会被系统自动重启，pm disable 才能彻底禁用
- **fnOS Docker 无宿主回环能力**：bridge 容器无法访问宿主局域网 IP，ADB 端口需要用 socat 转发
- **redroid 安卓主容器必须 privileged**：上游官方要求，非特权方案实测无法开机，但仅作用于容器内部，应用本体不申请宿主 root
- **升级脚本必须健壮**：飞牛升级流程先执行旧版本 uninstall_init，任何命令失败或 stdout 输出都会导致「执行脚本出错且原因未知」，必须 set +e + 输出重定向 + 所有命令 || true
- **CRLF 行尾会导致 Linux 脚本报错**：所有 shell/Python 脚本必须使用 LF 行尾，否则会出现 `$'\r': command not found`
- **应用本体降权需注意进程间通信**：gateway.py 降权后 `is_mine(pid)` 只能检查当前用户进程，音频守护也必须降权到同一用户，否则会被反复拉起

---

## 已知限制

1. **音频**：3.7.3+ 已启用音频（Codec2 Opus 软件编码器），需在连接设置中手动开启；3.8.x 进一步优化了音频稳定性。注意：音量大小与宿主系统实际音量匹配（容器内设100但宿主系统只设50，则按50输出）
2. **外网访问**：飞牛反向代理只透传 TCP，WebRTC 媒体流（UDP）无法通过，自动降级为 WebSocket 投屏
3. **飞牛 APP**：部分版本 WebView 对 WebSocket 代理支持有限，建议用手机浏览器
4. **ARM 应用兼容性**：强依赖谷歌服务或含反模拟器检测的 ARM64 应用可能闪退
5. **redroid 特权模式**：安卓主容器需要 privileged（redroid 上游官方要求，非特权方案实测无法开机），但仅作用于容器内部，应用本体不申请宿主 root
6. **穿云投屏设备数量授权**：免费版有设备数量限制，付费仅增加设备数量，不影响功能

---

## 问题、建议反馈链接和渠道

1. **交流、反馈和内测体验 QQ 群**：https://qm.qq.com/q/DF7nsBatFu
2. **建议、问题反馈问卷**：https://wj.qq.com/s2/28029808/2aab/
3. **发布者邮箱**：andforlin@foxmail.com
4. **redroid 容器和穿云投屏容器作者的容器专门反馈链接**：见「致谢和导向链接」章节

> 提供的反馈链接和渠道除第四条以外都会回复，因为第四条是对容器的问题和建议的专门反馈渠道，与软件本身无关；如果在反馈之后后续处理结果不满意或者想为软件添砖加瓦的，你可以使用源代码进行修改，按照飞牛官方打包教程之后还是按照前三条的反馈链接和渠道进行上传，发布者对上传的代码进行审计之后，会邀请你一起成为贡献者，为软件作出奉献；也感谢对软件本身的问题提出反馈、建议或者提供实质性帮助的人员。

---

## 对发布者和其他贡献者的支持
<img width="4096" height="2926" alt="a61d0506ea65c87e4dd005f21325eda6" src="https://github.com/user-attachments/assets/1ad0e1e8-e03e-4966-a94d-24dff71981ba" />

如果你觉得这个软件很好或者对发布者本身感到满意的，希望能赞赏多多支持一下，让发布者和其他贡献者得以继续维护软件，无论赞赏多少或者是不赞赏，都在此感谢；在支付的时候请在备注上标注赞赏，感谢。
或者给一个 star 支持一下项目也行。

> 注：以上赞赏码仅用于本项目的维护和开发支持，请勿盗用或用于其他用途，感谢理解。

---

## 开源许可和免责声明

### 开源许可

本应用运行时通过 Docker 拉取以下公开镜像，不修改、不捆绑、不分发其源码与二进制：

#### 1. redroid（安卓容器）
- 项目地址：https://github.com/remote-android/redroid-doc
- 原作者：zhouziyang（remote-android 组织）
- 许可证状态：
  - redroid 本身：[Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0)（上游 README 明确声明）
  - redroid-modules 内核模块仓库：[GPL-2.0](https://github.com/remote-android/redroid-modules/blob/master/LICENSE)
  - 容器内 AOSP（Android 开源项目）：[Apache 2.0](https://source.android.com/setup/start/licenses)
  - 容器内 Linux 内核相关：GPL-2.0，项目地址：https://www.kernel.org/
- 内置翻译层：
  - libndk_translation（Google 官方 NDK 翻译层）：Google 专有组件，随 redroid 镜像内置，许可条款见 Google 相关协议
  - libhoudini（Intel 翻译层，v3.8.1+ 自动下载）：Intel 专有组件，从公开渠道下载，许可条款见 Intel 相关协议

#### 2. scrcpy-over-webrtc（穿云投屏 / 云手机画面服务）
- 项目地址：https://github.com/hqw700/ScrcpyOverWebRTC
- 原作者：hqw700（buutuu）
- 许可证状态：
  - 前端源码（web-app）：[MIT License](https://opensource.org/licenses/MIT)
  - 官方二进制核心组件（服务端、Agent 部署包、APK 运行环境）：仅供个人学习交流、技术研究与非商业测试使用
- 内部依赖：
  - scrcpy（作者：Genymobile）：[Apache 2.0](https://github.com/Genymobile/scrcpy/blob/master/LICENSE)，项目地址：https://github.com/Genymobile/scrcpy
  - ya-webadb / Tango（作者：yume-chan）：[MIT](https://github.com/yume-chan/ya-webadb/blob/master/LICENSE)，项目地址：https://github.com/yume-chan/ya-webadb
  - Pion WebRTC（作者：pion 组织）：[MIT](https://github.com/pion/webrtc/blob/master/LICENSE)，项目地址：https://github.com/pion/webrtc
  - xterm.js（作者：xtermjs 组织）：[MIT](https://github.com/xtermjs/xterm.js/blob/master/LICENSE)，项目地址：https://github.com/xtermjs/xterm.js
  - coturn TURN 服务器（作者：coturn 项目）：[BSD 3-Clause](https://github.com/coturn/coturn/blob/master/LICENSE)，项目地址：https://github.com/coturn/coturn

#### 3. 本项目打包脚本和配置
- 项目地址：https://github.com/lin1740/fnos-android-emulator
- 原作者：键盘敲粥香（lin1740）
- 许可证：[MIT License](https://opensource.org/licenses/MIT)
  > 简要说明：MIT 许可证允许任何人免费使用、复制、修改、合并、发布、分发、再许可和销售本软件的副本，前提是在所有副本或重要部分中包含上述版权声明和本许可声明。本软件按"现状"提供，不提供任何明示或默示的担保。
- 包含：docker-compose 配置、安装/升级脚本、gateway.py 网关、状态页、性能优化脚本等（均为本项目自行开发，适配飞牛 fnOS 平台）
- 项目源码链接：见飞牛应用中心本应用详情页的「发布者」蓝色链接，或应用介绍中的项目链接
- 说明：本应用未自行开发 UI 界面，画面管理页面依赖穿云投屏（scrcpy-over-webrtc）的原生 UI，该 UI 不在本项目修改范围内

> **许可提示**：穿云投屏官方二进制核心组件仅供个人学习交流、技术研究与非商业测试使用，商用前请与作者确认授权；libndk_translation 和 libhoudini 为厂商专有组件，仅随镜像使用或自动下载，不进行再分发。

### 免责声明

1. 本项目为非官方第三方应用，按"现状"（AS IS）提供，使用风险自负。
2. 本项目仅用于学习和研究目的，不得用于任何违法用途。
3. 使用本应用产生的任何数据丢失、系统故障、服务中断等问题，发布者不承担任何责任。
4. 应用内集成的第三方组件（redroid、穿云投屏等）由各自作者维护，其功能、稳定性和合规性不受本项目控制。
5. 用户应自行备份重要数据，本应用不对容器内数据的安全性和完整性做出保证。
6. 本应用不收集任何用户数据，所有数据均存储在用户本地设备中。

### 开源义务说明

1. **GPL-2.0 组件义务**：redroid-modules（内核模块）及容器内 Linux 内核相关代码遵循 GPL-2.0 许可证。本应用仅运行时从公开仓库拉取 redroid 镜像，不修改、不重分发其源码与二进制，因此不触发 GPL-2.0 的传染条款；若用户自行修改、重编译或再分发上述 GPL-2.0 组件，须遵守 GPL-2.0 的开源义务，包括但不限于公开修改后的源代码、保留版权声明、以相同许可证分发。
2. **Apache 2.0 组件义务**：AOSP、scrcpy 等遵循 Apache 2.0 许可证的组件，再分发时须保留版权声明、许可证副本和 NOTICE 文件。
3. **专有组件**：libndk_translation（Google）、libhoudini（Intel）为厂商专有组件，本应用不进行再分发，仅随上游镜像使用或运行时自动下载；用户应遵守对应厂商的使用条款。
4. **穿云投屏组件**：前端源码为 MIT 许可证，可自由二次开发；官方二进制核心组件仅供个人学习交流、技术研究与非商业测试使用，商用前请与作者确认授权。
5. **本项目代码**：打包脚本和配置以 MIT 许可证开放，可自由使用、修改和分发，须保留版权声明和许可声明。

### 其他说明

1. **关于商标与标识**：redroid、穿云投屏、scrcpy、WebRTC 等名称及相关标识的商标权、著作权均归各自原作者或组织所有，本项目仅在技术集成层面使用这些名称进行说明，不主张任何商标权利，也不暗示与上述项目存在官方合作或背书关系。
2. **关于项目背书**：本项目对上游开源组件的集成和使用，仅代表技术层面的兼容性适配，不代表上游作者对本项目的认可、推荐或背书；各上游组件的质量、安全性及更新维护由其原作者负责。
3. **关于上游变更**：上游开源项目可能随版本迭代调整功能、接口或许可证条款，本项目将尽力跟进适配，但不对上游变更导致的兼容性问题或许可状态变化承担责任；如遇重大变更，建议以各上游项目官方公告为准。
4. **关于组件完整性**：本应用运行时从公开容器仓库拉取上游官方镜像，不对镜像内代码进行修改或二次打包；如用户自行替换或修改上游镜像，由此产生的功能异常或合规问题由用户自行负责。
5. **关于专利授权**：Apache 2.0 许可证包含贡献者的专利授权条款，MIT 和 BSD 许可证不涉及明确的专利授权；用户在使用、修改或再分发相关组件时，应自行评估专利风险。
6. **关于出口管制**：本项目涉及的部分编解码、加密技术可能受某些国家或地区的出口管制法规约束，用户在跨境使用或再分发时，应确保遵守所在地的相关法律法规。
7. **关于遗漏与勘误声明**：由于上游开源项目的依赖关系较为复杂，部分传递依赖或子组件的许可证信息、项目链接可能未能在本章节中逐一完整列举或准确标注。若您发现有应列而未列的开源项目、许可证状态有误，或项目链接存在错误，我们深表歉意，欢迎通过下方「问题、建议反馈链接和渠道」中的任意渠道（第 4 条上游组件专门反馈渠道除外）告知我们，我们将在核实后及时补充、更正。
8. **关于源代码发布**：本项目的打包脚本和配置以 MIT 许可证开放，但源代码的发布可能基于实际情况酌情处理。例如，内测版本的源代码可能因稳定性、安全性或其他原因暂不公开，公测版本的源代码通常会同步发布至本仓库。具体以本仓库实际发布的内容为准。

---

## 致谢和导向链接

### 导向链接

1. redroid 容器项目链接：https://github.com/remote-android/redroid-doc
2. 穿云投屏容器项目链接：https://github.com/hqw700/ScrcpyOverWebRTC
3. 穿云投屏官方文档：https://webrtc-phone.com/docs/
4. 穿云投屏官方网站：https://webrtc-phone.com/

### 致谢

- [redroid 项目] — Android in Docker
- [穿云投屏 scrcpy-over-webrtc] — WebRTC 画面服务
- 飞牛 fnOS 开发社区
