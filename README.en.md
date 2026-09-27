# androidemu — Android Emulator / Cloud Phone

[中文](README.md) | **English**

![version](https://img.shields.io/badge/version-v3.6.4-blue) ![arch](https://img.shields.io/badge/arch-x86__64%20%7C%20arm64-orange) ![image](https://img.shields.io/badge/image-~2GB-green) ![stars](https://img.shields.io/github/stars/lin1740/fnos-android-emulator) ![last-commit](https://img.shields.io/github/last-commit/lin1740/fnos-android-emulator) ![license](https://img.shields.io/github/license/lin1740/fnos-android-emulator)

📚 **User Manual & FAQ**: See sections below

Run an Android 12 virtual machine on fnOS with one click, remote control via browser, supporting WebRTC / WebSocket dual screen casting modes, ADB connection, APK installation, file management, and more.

Based on Android container + Scrcpy over WebRTC (screen service) dual-container architecture, adapted for fnOS unified gateway.
<img width="1288" height="900" alt="firefox exe_20260927_095032" src="https://github.com/user-attachments/assets/da7718fc-c490-4f35-be32-572f0f1a3849" />

---

## Table of Contents

- [Features](#features)
- [Installation Requirements](#installation-requirements)
- [Installation Methods](#installation-methods)
- [Access Methods](#access-methods)
- [Default Account](#default-account)
- [ADB Connection](#adb-connection)
- [APK Installation](#apk-installation)
- [Performance Optimization](#performance-optimization)
- [Serial Console (Developer Debugging)](#serial-console-developer-debugging)
- [Container Architecture](#container-architecture)
- [FAQ](#faq)
- [Known Limitations](#known-limitations)
- [Feedback Links & Channels](#feedback-links--channels)
- [Support the Publisher & Contributors](#support-the-publisher--contributors)
- [Open Source License & Disclaimer](#open-source-license--disclaimer)
- [Acknowledgements & Links](#acknowledgements--links)

---

## Features

- **Android 12 System**: x86_64 architecture, built-in ARM translation layer (libndk_translation), most ARM apps can be installed and run directly
- **Browser Remote Control**: No client installation needed, open browser to control Android desktop
- **Dual Screen Casting Modes**:
  - WebRTC casting (low latency, high framerate, recommended for LAN)
  - WebSocket casting (good compatibility, auto-switch for external network/penetration environments)
- **fnOS Unified Gateway Adaptation**: One-click open via fnOS App Center, auto-inject login state, no secondary login required
- **ADB Debugging**: `adb connect <NAS_IP>:5556` to connect
- **File Management**: Manage files directly in the Android container within the Scrcpy interface
- **Terminal Access**: Built-in Android shell terminal
- **Serial Console**: Disabled by default to reduce performance overhead, developers can manually enable for debugging (see below)
- **Scrcpy Agent Integration**: Support connecting multiple Android devices/real phones, unified management
- **Simplified Chinese + China Timezone**: Container defaults to `zh_CN` + `Asia/Shanghai`
- **Non-root Operation**: Screen service container runs as normal user (appuser), compliant with app store review requirements
- **X86 / ARM Dual Platform**: Auto-detect architecture and GPU capabilities, X86 uses hardware acceleration, ARM auto-switches to software rendering
- **Performance Optimization**: webrtc/turn process high-priority scheduling, streamlined background services, CPU dynamic frequency scaling (see below)

---

## Installation Requirements

| Item | Minimum | Recommended |
|------|----------|-------------|
| System | fnOS 1.1.8+ | Latest fnOS |
| Architecture | x86_64 / aarch64 | x86_64 |
| CPU | 2 cores | 4 cores+ |
| Memory | 2 GB | 4 GB+ |
| Storage | 4 GB available | 8 GB+ (including ~2GB image) |
| Network | LAN | — |

> **X86 Devices**: Requires Docker `/dev/dri` passthrough for hardware acceleration; auto-fallback to software rendering if no GPU; must install binder_linux driver before downloading the app (available in App Center, just search), otherwise installation will be rejected.
> **ARM Devices**: Auto-use software rendering (gpu_mode=guest), no additional driver needed.

---

## Installation Methods

### Method 1: fnOS App Center (Recommended)

1. Open fnOS App Center
2. Search for "Android Emulator" or "androidemu"
3. Click install, wait for image pull to complete (~2GB, priority DaoCloud/fnOS mirror)
4. After installation, click "Open" to enter the Scrcpy interface

### Method 2: Manual fpk Installation

1. Download the latest `androidemu_all_x.x.x.fpk`
2. Select "Manual Install" in fnOS App Center, upload the fpk file
3. Wait for installation and image pull to complete

### Image Acceleration

Installation automatically tries image sources in the following order:
1. fnOS mirror `docker.fnnas.com`
2. DaoCloud free mirror `docker.m.daocloud.io`
3. Docker Hub official repository

If all fail, please configure image accelerator (registry-mirrors) in fnOS Docker settings and retry.

---

## Access Methods

### LAN (Recommended, most complete features)

Click "Open" in fnOS App Center, or directly access:
```
https://<NAS_IP>:<NAS_PORT>/app/androidemu/
```
You can also directly access the screen service native port:
```
https://<NAS_IP>:8443
```

WebRTC casting works normally under LAN, low latency, high framerate, all features (file management, terminal, casting settings, etc.) available.

### fnOS Official Remote Domain (xxx.fnos.net)

When accessing through fnOS official remote domain:
- ✅ WebSocket casting auto-enabled, can view and control screen normally
- ⚠️ Features relying on WebRTC DataChannel (UDP) such as file management may not be available
- ⚠️ Framerate and latency affected by upstream bandwidth

> Reason: fnOS reverse proxy only passes TCP (HTTP/HTTPS/WebSocket), not UDP. WebRTC media streams use UDP, so auto-degrade to WebSocket casting.

### fnOS APP

When accessing through fnOS mobile APP built-in browser:
- ⚠️ Some versions of fnOS APP WebView have limited WebSocket proxy support, may show "waiting for device push stream timeout"
- ✅ Recommended to use mobile browser (Chrome / Edge / system browser) to directly open fnOS remote domain

### Public IP / Intranet Penetration / Reverse Proxy

If you have a public IP or use intranet penetration (frp, ZeroTier, Tailscale, etc.):
- ✅ Map port 8443 to external network, can use WebRTC casting directly
- ✅ Need to set "External Access Address (WebRTC Media Stream)" to your public domain or IP in "Casting Settings"
- ✅ TURN relay server is built-in, ensure UDP 3478 port is reachable

---

## Default Account

Scrcpy screen service default login account:

| Item | Value |
|------|-------|
| Username | `admin` |
| Password | `admin123` |

> Auto-inject login state when accessing through fnOS unified gateway, no manual input needed. Manual login required when directly accessing port 8443, please change account and password after login.

---

## ADB Connection

```bash
adb connect <NAS_IP>:5556
adb shell
```

You can also use Android shell directly in the "Terminal" of the Scrcpy interface.

---

## APK Installation

1. Click "File" or "Upload" in the Scrcpy interface
2. Select APK file to upload
3. Click the APK in the file manager within the Android container to install

> **Note**: Built-in emulator is x86_64 architecture, no Google services. Image has built-in ARM translation layer, most ARM apps can run; but ARM64 apps strongly dependent on Google services or with anti-emulator detection may crash. Such apps are recommended to use real phone via Scrcpy Agent or configure according to upstream redroid container author's Google service recommendations.

---

## Performance Optimization

This app has multiple built-in performance optimizations to reduce access and usage latency:

### 1. Process Priority Boost

Both webrtc signaling service and TURN relay service start with `nice=-10` (higher than default priority), reducing CPU scheduling latency and avoiding screen stuttering.

> Technical detail: Docker containers default drop `CAP_SYS_NICE`, this app explicitly adds `cap_add: SYS_NICE` in compose, this capability only affects process scheduling priority, does not involve device/network/filesystem access.

### 2. Background Service Streamlining

After container startup, automatically disable unnecessary Android system services to free CPU and memory:
- Bluetooth services (`com.android.bluetooth`, `com.android.bluetoothmidiservice`) — container has no Bluetooth hardware
- NFC service (`com.android.nfc`) — container has no NFC hardware
- Print service (`com.android.printspooler`) — container does not need printing
- Backup services (`com.android.backupconfirm`, `com.android.sharedstoragebackup`) — container does not need backup
- Location fusion service (`com.android.location.fused`) — container has no GPS hardware

### 3. CPU Dynamic Frequency Scaling

Installation automatically configures CPU scheduler to dynamic frequency scaling mode (schedutil/ondemand), and sets minimum frequency to 40% of maximum frequency:
- Under light load, CPU won't drop too low, ensuring operation response speed
- Under heavy load, auto-rise to maximum frequency,发挥 full performance
- Only effective on devices supporting cpufreq, silently skip if not supported

### 4. Serial Console Disabled by Default

Android serial console (`androidboot.console=0`) is disabled by default, reducing performance overhead from kernel log output. Developers can manually enable if needed (see below).

### 5. GPU Hardware Acceleration

- **X86 Devices**: Auto-detect `/dev/dri`, use `gpu_mode=host` hardware acceleration when GPU available, 60fps
- **ARM Devices**: Auto-use `gpu_mode=guest` software rendering, 30fps (ARM usually has no GPU passthrough)

---

## Serial Console (Developer Debugging)

Serial console is disabled by default to reduce performance overhead. To enable for debugging:

```bash
# Temporarily enable (invalid after container restart)
docker exec -u 0 androidemu-android setprop persist.sys.serialconsole 1
docker exec -u 0 androidemu-android start console

# View console output
docker exec -u 0 androidemu-android dmesg -w

# Disable
docker exec -u 0 androidemu-android stop console
docker exec -u 0 androidemu-android setprop persist.sys.serialconsole 0
```

> Note: Enabling serial console increases kernel log output, may slightly affect performance, recommended to disable after debugging.

---

## Container Architecture

```
┌─────────────────────────────────────────────┐
│  fnOS Host                                  │
│                                             │
│  ┌──────────────────┐  ┌─────────────────┐  │
│  │  androidemu-     │  │  androidemu-    │  │
│  │  android         │  │  webrtc         │  │
│  │  (redroid)       │  │  (scrcpy+TURN)  │  │
│  │                  │  │                 │  │
│  │  Android 12      │  │  Signaling :8443│  │
│  │  ADB :5556       │  │  TURN  :3478    │  │
│  │  bridge network  │  │  host network   │  │
│  │  privileged      │  │  SYS_NICE       │  │
│  └──────────────────┘  └─────────────────┘  │
│          │                       │          │
│          └───── scrcpy ─────────┘           │
│                  (ADB over TCP)             │
└─────────────────────────────────────────────┘
         │
         ▼
   fnOS Unified Gateway
         │
         ▼
      Browser
```

---

## FAQ

### Q: Shows "Unauthorized" or "0 devices online" after opening

A: This is Scrcpy's License authorization prompt. Newly installed devices currently have a 3-month free trial for 20 devices (from November 1, 2026, expires to 10 devices, other basic features are free). Wait for container to fully start (~1-2 minutes) then refresh page. If持续 unauthorized, check if container is running normally:
```bash
docker ps --filter name=androidemu
```
<img width="1288" height="900" alt="firefox exe_20260927_092044" src="https://github.com/user-attachments/assets/44afdf11-2112-4ae7-8544-89e6bfa0238b" />

### Q: WebRTC casting connection failed

A: Troubleshoot in following steps:
1. Confirm using within LAN (external network auto-degrades to WebSocket casting)
2. Confirm TURN server running normally: `docker exec androidemu-webrtc netstat -ulnp | grep 3478`
3. Confirm PUBLIC_IP is LAN IP not 127.0.0.1: `docker exec androidemu-webrtc env | grep PUBLIC_IP`
4. Confirm ICE_SERVERS contains correct LAN IP: `docker exec androidemu-webrtc env | grep ICE_SERVERS`
5. Try clicking "Switch to WebSocket casting"

### Q: Container restarts repeatedly

A: Common causes:
- ARM device not enabled software rendering: check if compose has `androidboot.redroid_gpu_mode=guest`
- Insufficient memory: recommend at least 2GB available memory
- View logs: `docker logs androidemu-android`
- webrtc container restart: check if `nice: setpriority(-10): Permission denied`, confirm compose includes `cap_add: SYS_NICE`

### Q: No sound

A: Current version disables audio by default (opus encoder in redroid container is Codec2 version, scrcpy-server only recognizes OMX version, enabling audio causes `createEncoder` failure and stream disconnection). Dual protection through RUNTIME_SHIM hijacking WebSocket.send and gateway intercepting `/api/default_settings` to disable audio. Future versions will attempt to fix.

### Q: External network access screen laggy

A: External network through fnOS reverse proxy auto-uses WebSocket casting, framerate limited by upstream bandwidth. If you have public IP, recommend directly mapping port 8443 to use WebRTC casting.

### Q: Clicking "Terminal" pops up print page

A: This is a Scrcpy frontend shortcut conflict bug. Fixed through RUNTIME_SHIM blocking `window.print()` and Ctrl+P shortcut, clicking terminal no longer triggers printing.

### Q: Audit log loading failed

A: Audit log function depends on Scrcpy backend `/api/audit` interface, some versions may not support. This is upstream image function, does not affect core casting function.

### Q: How to uninstall

A: Click "Uninstall" in fnOS App Center. Container data (Android /data partition) remains in Docker volume, to completely remove:
```bash
docker volume rm androidemu_data androidemu-webrtc-data
```

### Technical Notes & Development Pitfalls

Key technical conclusions verified during development, for reference for secondary development:

- **com.android.location.fused must never be disabled**: It is a dependency of LocationManagerService, disabling causes system_server crash and Android fails to boot
- **androidboot.use_memfd=1 is unsafe on redroid 12**: Causes LocationManagerService crash, init kills zygote, memory drops from 600-700MB to 400MB, container cannot start
- **Detect device nodes rather than directory names for dependency apps**: App Center display name and actual directory name may differ; detecting `/dev/binder` is more reliable than `/var/apps/binder_linux_driver`
- **WebSocket proxy must cover all ws/wss URLs**: Cannot only cover specific paths, otherwise other paths will directly connect to port 8443 causing cross-origin failure when accessed through gateway
- **BusyBox su resets environment variables**: Cannot use su -c internal export to pass environment variables during non-root transformation, should set directly in docker-compose environment
- **Use pm disable not svc disable for Bluetooth**: svc disable gets auto-restarted by system, pm disable can completely disable
- **fnOS Docker has no host loopback capability**: Bridge containers cannot access host LAN IP, ADB port needs socat forwarding
- **redroid Android main container must be privileged**: Upstream official requirement, non-privileged solution tested unable to boot, but only acts inside container, app itself does not request host root

---

## Known Limitations

1. **Audio**: Current version disables audio by default, enabling causes stream disconnection
2. **External Network Access**: fnOS reverse proxy only passes TCP, WebRTC media streams (UDP) cannot pass, auto-degrade to WebSocket casting
3. **fnOS APP**: Some versions WebView has limited WebSocket proxy support, recommend using mobile browser
4. **ARM App Compatibility**: ARM64 apps strongly dependent on Google services or with anti-emulator detection may crash
5. **redroid Privileged Mode**: Android main container requires privileged (redroid upstream official requirement, non-privileged solution tested unable to boot), but only acts inside container, app itself does not request host root

---

## Feedback Links & Channels

1. **Communication, Feedback & Beta Testing QQ Group**: https://qm.qq.com/q/DF7nsBatFu
2. **Suggestion & Issue Feedback Survey**: https://wj.qq.com/s2/28029808/2aab/
3. **Publisher Email**: andforlin@foxmail.com
4. **redroid container and Scrcpy container author's dedicated feedback links**: See "Acknowledgements & Links" section

> All feedback links and channels except the fourth will be replied, because the fourth is dedicated feedback channel for container issues and suggestions, unrelated to the software itself; if you are unsatisfied with subsequent processing results after feedback or want to contribute to the software, you can modify using source code, follow fnOS official packaging tutorial then upload through the first three feedback links and channels, after publisher audits the uploaded code, you will be invited to become a contributor together; also thanks to those who provide feedback, suggestions or substantial help for software issues.

---

## Support the Publisher & Contributors

<img width="4096" height="2926" alt="a61d0506ea65c87e4dd005f21325eda6" src="https://github.com/user-attachments/assets/1ad0e1e8-e03e-4966-a94d-24dff71981ba" />

If you find this software helpful or are satisfied with the publisher, hope you can support through appreciation, so that the publisher and other contributors can continue maintaining the software, whether appreciating any amount or not, thank you here; please note "appreciation" in the remark when paying, thanks.
Or give a star to support the project.

> Note: The above donation codes are only used for the maintenance and development support of this project. Please do not misappropriate or use them for other purposes. Thanks for understanding.

---

## Open Source License & Disclaimer

### Open Source License

- redroid: [Apache 2.0](http://www.apache.org/licenses)
- scrcpy-over-webrtc (Scrcpy): See upstream project
- This project's packaging scripts and configuration: MIT

Upstream component sources and license status see `LICENSE` file in the package.

### Disclaimer

1. This project is an unofficial third-party application, provided "as is", use at your own risk.
2. This project is for learning and research purposes only, and must not be used for any illegal purposes.
3. The developer is not responsible for any data loss, system failure, service interruption or other issues caused by using this application.
4. Third-party components integrated in the application (redroid, Scrcpy, etc.) are maintained by their respective authors, and their functionality and stability are not controlled by this project.
5. Users should back up important data by themselves. This application does not guarantee the security and integrity of data in the container.
6. This application does not collect any user data; all data is stored on the user's local device.

---

## Acknowledgements & Links

### Links

1. redroid container project: https://github.com/remote-android/redroid-doc
2. Scrcpy container project: https://github.com/hqw700/ScrcpyOverWebRTC
3. Scrcpy official docs: https://webrtc-phone.com/docs/
4. Scrcpy official website: https://webrtc-phone.com/

### Acknowledgements

- [redroid project] — Android in Docker
- [Scrcpy scrcpy-over-webrtc] — WebRTC screen service
- fnOS development community
