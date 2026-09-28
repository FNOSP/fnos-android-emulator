# androidemu — 安卓模拟器 / 云手机

**中文** | [English](README.en.md)

![version](https://img.shields.io/badge/version-v3.6.7-blue) ![arch](https://img.shields.io/badge/arch-x86__64%20%7C%20arm64-orange) ![image](https://img.shields.io/badge/image-~2GB-green) ![stars](https://img.shields.io/github/stars/lin1740/fnos-android-emulator) ![last-commit](https://img.shields.io/github/last-commit/lin1740/fnos-android-emulator) ![license](https://img.shields.io/github/license/lin1740/fnos-android-emulator)

📚 **使用手册与常见问题**：见本文档下方各章节

在飞牛 fnOS 上一键运行 Android 12 虚拟机，通过浏览器远程操控，支持 WebRTC / WebSocket 双投屏模式、ADB 连接、APK 安装、文件管理等。

基于 Android 容器+穿云投屏 scrcpy-over-webrtc（画面服务）双容器架构，适配飞牛统一网关。
<img width="1288" height="900" alt="firefox exe_20260927_095032" src="https://github.com/user-attachments/assets/da7718fc-c490-4f35-be32-572f0f1a3849" />

---

## 目录

- [功能特性](#功能特性)
- [安装要求](#安装要求)
- [安装方法](#安装方法)
- [访问方式](#访问方式)
- [默认账号](#默认账号)
- [ADB 连接](#adb-连接)
- [APK 安装](#apk-安装)
- [性能优化](#性能优化)
- [串口控制台（开发者调试）](#串口控制台开发者调试)
- [容器架构](#容器架构)
- [审核合规性说明](#审核合规性说明)
- [常见问题](#常见问题)
- [已知限制](#已知限制)
- [问题、建议反馈链接和渠道](#问题建议反馈链接和渠道)
- [对发布者和其他贡献者的支持](#对发布者和其他贡献者的支持)
- [开源许可和免责声明](#开源许可和免责声明)
- [致谢和导向链接](#致谢和导向链接)

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
- **应用本体非 root 运行**：gateway.py 等后台进程以 `docker-androidemu` 用户运行（v3.6.5+），符合上架审核要求
- **X86 / ARM 双平台**：自动检测架构和 GPU 能力，X86 用硬件加速，ARM 自动切软件渲染
- **性能优化**：webrtc/turn 进程高优先级调度、精简后台服务、CPU 动态调频（见下文）
- **安装/更新中断安全**：安装或更新中途取消会自动清理临时数据，避免占位导致下次无法安装（v3.6.0+）
- **自动容器检测**：网关自动检测安卓容器状态，容器启动后自动上线，无需手动操作

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

## 审核合规性说明

本应用已通过飞牛官方 7 条自查（基本信息、权限声明、网络端口、数据存储、启动停止、卸载清理、兼容性），以下为审核关注要点的详细说明：

### 1. 特权容器（privileged）

- **现状**：仅 `androidemu-android`（redroid 安卓主容器）使用 `privileged: true`，`androidemu-webrtc`（画面服务）为非特权运行
- **必要性**：redroid 上游（remote-android/redroid-doc）官方部署方式即要求 `--privileged`，Android 依赖内核 binder 通信，容器需要挂载/访问 binder 设备并创建设备节点
- **非特权方案已实测不可行**：非特权 + device_cgroup_rule + 挂载 binderfs + cap-add=ALL + seccomp=unconfined，容器以 ExitCode 0 静默退出、无法开机（redroid-doc issue #591 至今为开放议题）
- **影响面控制**：特权仅作用于容器内部（容器内 root = Android 系统自身初始化所需），不等于宿主 root；应用本体以 `docker-androidemu` 用户运行，不申请宿主 root

### 2. 应用本体非 root 运行（v3.6.5+）

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

### Q: 首次安装时弹出「无法安装 androidemu - 执行脚本出错且原因未知」

A: 优先检查 **binder 驱动**是否已安装：
- **x86 设备**：需先在飞牛应用中心安装「binder_linux 驱动」依赖应用，产生 `/dev/binder` 设备节点后再安装本应用
- **ARM 设备**：多数设备内核已内置 binder（如 RK3588 等），但部分精简内核可能未启用，需确认内核支持 `CONFIG_ANDROID_BINDER_IPC` / binderfs

若已安装驱动（或 ARM 设备本身支持）但仍弹出此错误，请务必通过下方「问题反馈」章节中的任一链接反馈，以便排查具体原因。

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

A: 常见原因：
- ARM 设备未启用软件渲染：检查 compose 中是否有 `androidboot.redroid_gpu_mode=guest`
- 内存不足：建议至少 2GB 可用内存
- 查看日志：`docker logs androidemu-android`
- webrtc 容器重启：检查是否有 `nice: setpriority(-10): Permission denied`，确认 compose 中包含 `cap_add: SYS_NICE`

### Q: 容器卡顿、无法点击或移动

A: 按以下步骤排查：
1. 先关闭软件页面（或网页），重新打开后再尝试点击/移动
2. 若仍无效，在飞牛应用中心的软件详情页点击「停用」，停用后再「启用」
3. 如果只有单个容器出现卡顿，可在 Docker 中找到对应容器点击「重启」即可
4. 若以上方法均无效，请将软件卡顿的截图或录屏，以及 Docker 容器中复制的日志打包为文本文档，通过下方反馈渠道任选一项进行反馈

### Q: 安卓容器内存很低（<300MB）、设备一直不在线

A: 正常 Android 12 启动后内存应在 300MB 以上。若容器在运行但内存只有 100-200MB，说明安卓系统未完成启动（`boot_completed != 1`），agent 无法部署，设备永远不在线。按以下步骤排查：

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

**3. GPU 修复命令（画面异常或 SurfaceFlinger 崩溃时用）：**
```bash
docker exec -u 0 androidemu-android chmod 666 /dev/dri/card0 /dev/dri/renderD128
docker exec -u 0 androidemu-android sh -c 'setprop ctl.restart surfaceflinger'
sleep 60
docker exec androidemu-android getprop sys.boot_completed
```

**4. 手动注入 agent（boot_completed=1 但设备仍不在线时用）：**
```bash
# 从穿云投屏镜像取出 agent（首次需要）
docker run --rm -v /var/apps/androidemu/var/agent:/out --entrypoint /bin/sh docker.fnnas.com/buutuu/scrcpy-over-webrtc:latest -c "cp /app/agent_binaries/cloudphone-agent-amd64 /app/agent_binaries/libsys_core.so /out/ && chmod 755 /out/cloudphone-agent-amd64"
# 注入并启动（将 <NAS_IP> 替换为你的 NAS 局域网 IP；x86 用 amd64，ARM 用 arm64）
docker cp /var/apps/androidemu/var/agent/cloudphone-agent-amd64 androidemu-android:/data/local/tmp/cloudphone-agent
docker cp /var/apps/androidemu/var/agent/libsys_core.so androidemu-android:/data/local/tmp/libsys_core.so
docker exec -u 0 androidemu-android sh -c "chmod 755 /data/local/tmp/cloudphone-agent && export CP_AGENT_JAR=/data/local/tmp/libsys_core.so && nohup /data/local/tmp/cloudphone-agent -signaling wss://<NAS_IP>:8443/register_agent -id androidemu -ice-servers 'turn:cloudphone_user:cloudphone_secure_password@<NAS_IP>:3478?transport=udp,turn:cloudphone_user:cloudphone_secure_password@<NAS_IP>:3478?transport=tcp,stun:<NAS_IP>:3478' -jar /data/local/tmp/libsys_core.so > /data/local/tmp/agent.log 2>&1 &"
sleep 3
docker exec androidemu-android pidof cloudphone-agent
```

### Q: WebRTC 连接失败、黑屏或一直转圈

A: 最常见原因是 `PUBLIC_IP` 被重置为 `127.0.0.1`（重建容器后 compose 默认值生效），导致 TURN 分发给客户端的中继地址是 `127.0.0.1`，客户端连不上。

**诊断命令：**
```bash
echo "===== PUBLIC_IP ====="
docker inspect androidemu-webrtc --format '{{range .Config.Env}}{{println .}}{{end}}' | grep PUBLIC_IP
echo "===== agent 是否运行 ====="
docker exec androidemu-android pidof cloudphone-agent || echo "agent 未运行"
echo "===== agent 日志最后20行 ====="
docker exec androidemu-android tail -20 /data/local/tmp/agent.log 2>/dev/null
echo "===== webrtc 容器日志 ====="
docker logs androidemu-webrtc --tail 40 2>&1 | grep -iE 'relay addr|allocation|error|fail|Unauthorized' | tail -15
```

**修复命令（确认 PUBLIC_IP=127.0.0.1 时用，将 <NAS_IP> 替换为你的 NAS 局域网 IP）：**
```bash
cd /var/apps/androidemu/target/docker
sed -i 's/PUBLIC_IP=127\.0\.0\.1/PUBLIC_IP=<NAS_IP>/g' docker-compose.yaml
docker compose -p androidemu up -d --force-recreate webrtc
sleep 10
docker inspect androidemu-webrtc --format '{{range .Config.Env}}{{println .}}{{end}}' | grep PUBLIC_IP
```
确认 `PUBLIC_IP` 为你的局域网 IP 后，重新打开云手机画面页面连接。

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

### Q: 安装/更新中途取消后无法重新安装

A: v3.6.0+ 已修复，安装/更新中途取消会自动清理临时数据。若使用旧版本遇到此问题，手动清理：
```bash
rm -rf /tmp/androidemu_*
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

1. **音频**：当前版本默认禁用音频，开启会导致断流
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

- redroid：[Apache 2.0](http://www.apache.org/licenses)
- scrcpy-over-webrtc（穿云投屏）：见上游项目
- 本项目打包脚本和配置：MIT

上游组件出处与许可状态详见包内 `LICENSE` 文件。

### 免责声明

1. 本项目为非官方第三方应用，按"现状"提供，使用风险自负。
2. 本项目仅用于学习和研究目的，不得用于任何违法用途。
3. 使用本应用产生的任何数据丢失、系统故障、服务中断等问题，开发者不承担任何责任。
4. 应用内集成的第三方组件（redroid、穿云投屏等）由各自作者维护，其功能和稳定性不受本项目控制。
5. 用户应自行备份重要数据，本应用不对容器内数据的安全性和完整性做出保证。
6. 本应用不收集任何用户数据，所有数据均存储在用户本地设备中。

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
