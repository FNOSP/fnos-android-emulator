# androidemu — 飞牛 fnOS 安卓模拟器 / 云手机

在飞牛 fnOS 上一键运行 Android 12 虚拟机，通过浏览器远程操控，支持 WebRTC / WebSocket 双投屏模式、ADB 连接、APK 安装、文件管理等。

基于 Android 容器+穿云投屏 scrcpy-over-webrtc（画面服务）双容器架构，适配飞牛统一网关。

---

## 功能特性

- **Android 12 系统**：x86_64 架构，内置 ARM 翻译层（libndk_translation），大多数 ARM 应用可直接安装运行
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
- **非 root 运行**：画面服务容器内以普通用户（appuser）运行，符合上架审核要求
- **X86 / ARM 双平台**：自动检测架构和 GPU 能力，X86 用硬件加速，ARM 自动切软件渲染
- **性能优化**：webrtc/turn 进程高优先级调度、精简后台服务、CPU 动态调频（见下文）

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

> **X86 设备**：需要 Docker 支持 `/dev/dri` 直通以启用硬件加速；无 GPU 时自动回退软件渲染。
> **ARM 设备**：自动使用软件渲染（gpu_mode=guest），无需额外驱动。

---

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

### 飞牛 APP

通过飞牛手机 APP 内置浏览器访问时：
- ⚠️ 部分版本的飞牛 APP WebView 对 WebSocket 代理支持有限，可能出现「等待设备推流超时」
- ✅ 建议在手机浏览器（Chrome / Edge / 系统浏览器）中直接打开飞牛远程域名使用

### 公网 IP / 内网穿透 / 反向代理

如果你有公网 IP 或使用内网穿透（如 frp、ZeroTier、Tailscale 等）：
- ✅ 将 8443 端口映射到外网，可直接使用 WebRTC 投屏
- ✅ 需要在「投屏设置」中将「外网访问地址（WebRTC 媒体流）」设置为你的公网域名或 IP
- ✅ TURN 中继服务器已内置，确保 UDP 3478 端口可达

---

## 默认账号

穿云投屏画面服务默认登录账号：

| 项目 | 值 |
|------|-----|
| 用户名 | `admin` |
| 密码 | `admin123` |

> 通过飞牛统一网关访问时会自动注入登录态，无需手动输入。直接访问 8443 端口时需要手动登录，登录之后请更改账号和密码。

---

## ADB 连接

```bash
adb connect <NAS_IP>:5556
adb shell
```

也可以在穿云投屏界面的「终端」中直接使用安卓 shell。

---

## APK 安装

1. 在穿云投屏界面点击「文件」或「上传」
2. 选择 APK 文件上传
3. 在安卓容器中点击文件管理器中的 APK 进行安装

> **注意**：内置模拟器为 x86_64 架构，不含谷歌服务。镜像已内置 ARM 翻译层，大多数 ARM 应用可运行；但强依赖谷歌服务、或含反模拟器检测的 ARM64 应用可能闪退。此类应用建议通过穿云投屏 Agent 接入真机使用或者根据上游redroid容器作者的谷歌服务推荐配置来进行。

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
- 位置融合服务（`com.android.location.fused`）— 容器无 GPS 硬件

### 3. CPU 动态调频

安装时自动配置 CPU 调度器为动态调频模式（schedutil/ondemand），并设置最低频率为最大频率的 40%：
- 轻载时 CPU 不会降频过低，保证操作响应速度
- 重载时自动升到最高频率，发挥全部性能
- 仅在支持 cpufreq 的设备上生效，不支持则静默跳过

### 4. 串口控制台默认关闭

安卓串口控制台（`androidboot.console=0`）默认关闭，减少内核日志输出带来的性能损耗。开发者如需调试可手动开启（见下文）。

### 5. GPU 硬件加速

- **X86 设备**：自动检测 `/dev/dri`，有 GPU 时使用 `gpu_mode=host` 硬件加速，60fps
- **ARM 设备**：自动使用 `gpu_mode=guest` 软件渲染，30fps（ARM 通常无 GPU 直通）

---

## 串口控制台（开发者调试）

串口控制台默认关闭以减少性能损耗。如需开启用于调试：

```bash
# 临时开启（容器重启后失效）
docker exec -u 0 androidemu-android setprop persist.sys.serialconsole 1
docker exec -u 0 androidemu-android start console

# 查看控制台输出
docker exec -u 0 androidemu-android dmesg -w

# 关闭
docker exec -u 0 androidemu-android stop console
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
│          │                       │          │
│          └───── scrcpy ─────────┘           │
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

## 常见问题

### Q: 打开后显示「未授权」或「0 台在线」

A: 这是穿云投屏的 License 授权提示。新安装的设备目前有20台设备三个月免费试用期（自2026年11月1日，到期为10台设备，其他基础功能均为免费），等待容器完全启动（约 1-2 分钟）后刷新页面即可。若持续未授权，请检查容器是否正常运行：
```bash
docker ps --filter name=androidemu
```

### Q: WebRTC 投屏连接失败

A: 按以下步骤排查：
1. 确认在局域网内使用（外网自动降级为 WebSocket 投屏）
2. 确认 TURN 服务器正常运行：`docker exec androidemu-webrtc netstat -ulnp | grep 3478`
3. 确认 PUBLIC_IP 是局域网 IP 而非 127.0.0.1：`docker exec androidemu-webrtc env | grep PUBLIC_IP`
4. 确认 ICE_SERVERS 中包含正确的局域网 IP：`docker exec androidemu-webrtc env | grep ICE_SERVERS`
5. 尝试点击「改用 WebSocket 投屏」

### Q: 容器反复重启

A: 常见原因：
- ARM 设备未启用软件渲染：检查 compose 中是否有 `androidboot.redroid_gpu_mode=guest`
- 内存不足：建议至少 2GB 可用内存
- 查看日志：`docker logs androidemu-android`
- webrtc 容器重启：检查是否有 `nice: setpriority(-10): Permission denied`，确认 compose 中包含 `cap_add: SYS_NICE`

### Q: 没有声音

A: 当前版本默认禁用音频（redroid 容器内的 opus 编码器为 Codec2 版本，scrcpy-server 只识别 OMX 版本，开启音频会导致 `createEncoder` 失败并断流）。已通过 RUNTIME_SHIM 劫持 WebSocket.send 和 gateway 拦截 `/api/default_settings` 双重保障禁用音频。后续版本将尝试修复。

### Q: 外网访问画面卡顿

A: 外网通过飞牛反向代理时自动使用 WebSocket 投屏，帧率受上行带宽限制。如有公网 IP，建议直接映射 8443 端口使用 WebRTC 投屏。

### Q: 点击「终端」弹出打印页面

A: 这是穿云投屏前端的快捷键冲突 bug。已通过 RUNTIME_SHIM 屏蔽 `window.print()` 和 Ctrl+P 快捷键，点击终端不会再触发打印。

### Q: 审计日志加载失败

A: 审计日志功能依赖穿云投屏后端的 `/api/audit` 接口，部分版本可能不支持。此为上游镜像功能，不影响核心投屏功能。

### Q: 如何卸载

A: 在飞牛应用中心点击「卸载」即可。容器数据（安卓 /data 分区）会保留在 Docker volume 中，如需彻底清除：
```bash
docker volume rm androidemu_data androidemu-webrtc-data
```

---

## 版本历史与踩坑记录

### v3.6.x — 容器启动稳定性修复

- **v3.6.1**：加回 console=0（确认安全，之前一直没问题）；精简后台服务（NFC/打印/备份/蓝牙，共6个）；nice=-10 保持；use_memfd=1 永久移除
- **v3.6.0**：修复安卓容器无法启动 — 根因是 androidboot.use_memfd=1 导致 LocationManagerService 崩溃（"Unable to find a direct boot aware fused location provider"），init 杀掉 zygote 及所有系统服务，内存从600-700MB降到400MB。移除 use_memfd=1 和 console=0 后恢复
- **踩坑教训**：com.android.location.fused 绝对不能禁用（LocationManagerService 依赖项，禁用会导致 system_server 崩溃）；use_memfd=1 在 redroid 12 上不安全

### v3.5.x — 非 root 改造与 ICE 修复

- **v3.5.9**：性能优化（webrtc/turn 进程 nice=-10、禁用后台服务、CPU 动态调频最低40%）；动画和壁纸保持系统默认
- **v3.5.8**：首次尝试性能优化（禁用动画+黑色壁纸），后因用户反馈回退动画和壁纸
- **v3.5.7**：根本性修复 ICE_SERVERS 传递问题 — install_callback 同时替换 ICE_SERVERS 中的 IP；entrypoint 新增 IP 兜底检测（自动修正 127.0.0.1 / Docker 网桥地址 172.x）
- **v3.5.6**：ICE_SERVERS 直接在 docker-compose environment 中设置，不依赖 entrypoint 的 export/su 传递（BusyBox su 会重置环境变量）
- **v3.5.3-3.5.5**：多轮修复非 root 改造导致的 ICE_SERVERS 丢失问题 — 尝试 setpriv（BusyBox 不支持 --reuid）、su -c 内部 export（子 shell 展开为空），最终在 compose 中直接设置环境变量解决
- **v3.5.2**：修复源文件硬编码 ARM IP（192.168.9.2）导致 X86 部署后 PUBLIC_IP 错误的问题
- **v3.5.1**：修复 install_callback sed 只匹配硬编码 IP、不匹配变量引用格式的问题；修复 hostname -I 返回 Docker 网桥地址（172.x）而非局域网 IP 的问题
- **v3.5.0**：画面服务容器内非 root 运行（appuser uid=1000），turnserver 和 webrtc-signaling 均降权启动

### v3.4.x — 投屏稳定性修复

- 修复 WS 投屏后自动刷新回首页的 bug（RUNTIME_SHIM 劫持了所有 WebSocket，断开时无条件 location.reload()）
- 修复 watchdog 误判 audio:false 为"缺失默认键"并反复重建容器的问题（ensure_once 端口在听就直接 return）
- 屏蔽终端按钮触发打印页面（RUNTIME_SHIM 屏蔽 window.print）
- 外网访问自动切换 WebSocket 投屏 + 提示条

### v3.3.x — 容器优化与稳定性

- 简体中文 + 中国时区默认（zh_CN / Asia/Shanghai）
- 串口控制台默认关闭（androidboot.console=0，减少性能损耗）
- 禁用蓝牙（pm disable 彻底禁用，svc disable 会被系统自动重启）
- 自定义 entrypoint 修复 TURN 监听地址（listening-ip=0.0.0.0，原镜像只监听 Docker 网桥 172.x）
- 禁用音频（redroid 只有 Codec2 版 Opus 编码器，scrcpy 只识别 OMX 版，开启会 createEncoder 失败并断流）

### v3.0.x — Python 重写网关

- Python 重写网关服务（gateway.py），纯标准库，无外部依赖
- 飞牛统一网关完整适配（子路径 /app/androidemu/）
- 自动登录注入，无需二次输入账号密码
- RUNTIME_SHIM 运行时前缀拦截器（fetch/XHR/WebSocket/setAttribute/MutationObserver）
- WebSocket 代理（保留 Connection: Upgrade 头）

### v2.0.x — 双架构合一

- X86 / ARM 双架构合一安装包（androidemu_all_x.x.x.fpk）
- 安装时自动检测架构和 GPU 能力（tune_compose.sh）
- ADB 端口改用 socat 转发（fnOS Docker 无宿主回环能力，bridge 容器无法访问宿主局域网 IP）
- webrtc 容器改用 host 网络（确保云手机 Agent 能直连宿主信令/TURN 端口）

### v1.0.x — 初始版本

- 初始版本：redroid + scrcpy-over-webrtc 双容器
- 镜像加速源（DaoCloud 免注册 + Docker Hub 回退）
- 经历 Docker Hub 连接超时、镜像拉取失败等问题

---

## 已知限制

1. **音频**：当前版本默认禁用音频，开启会导致断流
2. **外网访问**：飞牛反向代理只透传 TCP，WebRTC 媒体流（UDP）无法通过，自动降级为 WebSocket 投屏
3. **飞牛 APP**：部分版本 WebView 对 WebSocket 代理支持有限，建议用手机浏览器
4. **ARM 应用兼容性**：强依赖谷歌服务或含反模拟器检测的 ARM64 应用可能闪退
5. **redroid 特权模式**：安卓主容器需要 privileged（redroid 上游官方要求，非特权方案实测无法开机），但仅作用于容器内部，应用本体不申请宿主 root

---

## 开源许可

本项目为非官方第三方应用，按"现状"提供，使用风险自负。

- redroid：[Apache 2.0](https://github.com/remote-android/redroid-doc)
- scrcpy-over-webrtc（穿云投屏）：见上游项目
- 本项目打包脚本和配置：MIT

上游组件出处与许可状态详见包内 `LICENSE` 文件。

---

## 致谢

- [redroid 项目](https://github.com/remote-android/redroid-doc) — Android in Docker
- [穿云投屏 scrcpy-over-webrtc](https://github.com/hqw700/ScrcpyOverWebRTC) — WebRTC 画面服务
- 飞牛 fnOS 开发社区
