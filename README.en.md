# androidemu — Android Emulator / Cloud Phone

[中文](README.md) | **English**

![version](https://img.shields.io/badge/version-v3.6.7-blue) ![arch](https://img.shields.io/badge/arch-x86__64%20%7C%20arm64-orange) ![image](https://img.shields.io/badge/image-~2GB-green) ![stars](https://img.shields.io/github/stars/lin1740/fnos-android-emulator) ![last-commit](https://img.shields.io/github/last-commit/lin1740/fnos-android-emulator) ![license](https://img.shields.io/github/license/lin1740/fnos-android-emulator)

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
- [Audit Compliance Notes](#audit-compliance-notes)
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
- **Non-root App Process**: gateway.py and other background processes run as `docker-androidemu` user (v3.6.5+), compliant with app store audit requirements
- **X86 / ARM Dual Platform**: Auto-detect architecture and GPU capabilities, X86 uses hardware acceleration, ARM auto-switches to software rendering
- **Performance Optimization**: webrtc/turn process high-priority scheduling, streamlined background services, CPU dynamic frequency scaling (see below)
- **Safe Installation/Update Interruption**: Auto-cleanup of temporary data if installation or update is cancelled midway, preventing placeholder issues that block future installations (v3.6.0+)
- **Auto Container Detection**: Gateway auto-detects Android container status, container automatically comes online after startup, no manual operation needed

---


## Installation Requirements

| Item | Minimum | Recommended |
|------|---------|-------------|
| System | fnOS 1.1.8+ | Latest fnOS |
| Architecture | x86_64 / aarch64 | x86_64 |
| CPU | 2 cores | 4 cores+ |
| Memory | 2 GB | 4 GB+ |
| Storage | 4 GB free | 8 GB+ (including ~2GB image) |
| Network | LAN | — |

> **X86 devices**: Docker needs `/dev/dri` passthrough for hardware acceleration; auto-fallback to software rendering if no GPU; must install binder_linux driver before downloading (available in App Center, just search), otherwise app won't work or installation will be rejected.
> **ARM devices**: Auto uses software rendering (gpu_mode=guest), no extra driver needed.

---

## Installation Methods

### Method 1: fnOS App Center (Recommended)

1. Open fnOS App Center
2. Search for "Android Emulator" or "androidemu"
3. Click install, wait for image pull (~2GB, priority DaoCloud/fnOS mirror)
4. Click "Open" after installation to enter Scrcpy interface

### Method 2: Manual fpk Installation

1. Download latest `androidemu_all_x.x.x.fpk`
2. Select "Manual Install" in fnOS App Center, upload fpk file
3. Wait for installation and image pull to complete

### Image Acceleration

Installation automatically tries image sources in this order:
1. fnOS mirror `docker.fnnas.com`
2. DaoCloud free mirror `docker.m.daocloud.io`
3. Docker Hub official repository

If all fail, configure image accelerator (registry-mirrors) in fnOS Docker settings and retry.

---

## Access Methods

### LAN (Recommended, Full Features)

Click "Open" in fnOS App Center, or directly access:
```
https://<NAS_IP>:<NAS_PORT>/app/androidemu/
```
Can also directly access screen service native port:
```
https://<NAS_IP>:8443
```

WebRTC casting works normally in LAN, low latency, high framerate, all features (file management, terminal, casting settings etc.) available.

### fnOS Official Remote Domain (xxx.fnos.net)

When accessing via fnOS official remote domain:
- ✅ WebSocket casting auto-enabled, can view and control screen
- ⚠️ File management etc. relying on WebRTC DataChannel (UDP) may not work
- ⚠️ Framerate and latency affected by upstream bandwidth

> Reason: fnOS reverse proxy only passes TCP (HTTP/HTTPS/WebSocket), not UDP. WebRTC media stream uses UDP, so auto-downgrades to WebSocket casting.

### fnOS APP

When accessing via fnOS mobile APP built-in browser:
- ⚠️ Some versions of fnOS APP WebView have limited WebSocket proxy support, may show "waiting for device push stream timeout"
- ✅ Recommended to use mobile browser (Chrome / Edge / system browser) to open fnOS remote domain directly

### Public IP / Intranet Penetration / Reverse Proxy

If you have public IP or use intranet penetration (frp, ZeroTier, Tailscale etc.):
- ✅ Map port 8443 to external network, can use WebRTC casting directly
- ✅ Need to set "External access address (WebRTC media stream)" in "Casting Settings" to your public domain or IP
- ✅ TURN relay server built-in, ensure UDP 3478 port reachable

---

## Default Account

Scrcpy screen service default login account:

| Item | Value |
|------|-------|
| Username | `admin` |
| Password | `admin123` |

> Login state auto-injected when accessing via fnOS unified gateway, no manual input needed. Manual login required when directly accessing port 8443, please change password after login.

---

## ADB Connection

```bash
adb connect <NAS_IP>:5556
adb shell
```

Can also use Android shell directly in the "Terminal" of Scrcpy interface.

---

## APK Installation

1. Click "File" or "Upload" in Scrcpy interface
2. Select APK file to upload
3. Click APK in file manager within Android container to install

> **Note**: Built-in emulator is x86_64 architecture, no Google services. Image has built-in ARM translation layer, most ARM apps can run; but ARM64 apps strongly dependent on Google services or with anti-emulator detection/complex JIT may crash on startup (translation layer capability boundary). Such apps recommended to use Scrcpy Agent to connect real phone, or configure according to upstream redroid author's Google service recommendations.

---

## Performance Optimization

This app has multiple built-in performance optimizations to reduce access and usage latency:

### 1. Process Priority Boost

webrtc signaling service and TURN relay service both start with `nice=-10` (higher than default priority), reducing CPU scheduling latency, avoiding screen stutter.

> Technical detail: Docker containers default drop `CAP_SYS_NICE`, this app explicitly adds `cap_add: SYS_NICE` in compose, this capability only affects process scheduling priority, not device/network/filesystem access.

### 2. Background Service Streamlining

After container startup, automatically disable unnecessary Android system services, releasing CPU and memory:
- Bluetooth services (`com.android.bluetooth`, `com.android.bluetoothmidiservice`) — container has no Bluetooth hardware
- NFC service (`com.android.nfc`) — container has no NFC hardware
- Print service (`com.android.printspooler`) — container doesn't need printing
- Backup services (`com.android.backupconfirm`, `com.android.sharedstoragebackup`) — container doesn't need backup

> **Note**: `com.android.location.fused` (location fused service) cannot be disabled, it's a dependency of LocationManagerService, disabling will cause system_server crash and Android won't boot.

### 3. CPU Dynamic Frequency Scaling

Installation auto-configures CPU governor to dynamic frequency mode (schedutil/ondemand), sets minimum frequency to 40% of maximum:
- Light load: CPU won't drop too low, ensuring operation response speed
- Heavy load: auto-rise to maximum frequency, full performance
- Only works on devices supporting cpufreq, silently skips if not supported
- Won't go below quad-core CPU limit (dynamic adjustment, won't over-restrict performance)

### 4. Serial Console Disabled by Default

Android serial console (`androidboot.console=0`) disabled by default, reducing performance overhead from kernel log output. Developers can manually enable if needed (see below).

### 5. GPU Hardware Acceleration

- **X86 devices**: Auto-detect `/dev/dri`, uses `gpu_mode=host` hardware acceleration when GPU present, 60fps
- **ARM devices**: Auto uses `gpu_mode=guest` software rendering, 30fps (ARM usually no GPU passthrough)

---

## Serial Console (Developer Debugging)

Serial console disabled by default to reduce performance overhead. To enable for debugging:

```bash
# Temporarily enable (lost after container restart)
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

## Audit Compliance Notes

This app has passed fnOS official 7-point self-check (basic info, permission declaration, network ports, data storage, start/stop, uninstall cleanup, compatibility). Below are detailed explanations of audit focus points:

### 1. Privileged Container

- **Status**: Only `androidemu-android` (redroid Android main container) uses `privileged: true`, `androidemu-webrtc` (screen service) runs non-privileged
- **Necessity**: redroid upstream (remote-android/redroid-doc) official deployment requires `--privileged`, Android depends on kernel binder communication, container needs to mount/access binder device and create device nodes
- **Non-privileged alternative tested infeasible**: Non-privileged + device_cgroup_rule + mount binderfs + cap-add=ALL + seccomp=unconfined, container exits with ExitCode 0 silently, can't boot (redroid-doc issue #591 still open)
- **Impact control**: Privilege only applies inside container (container root = Android system initialization needs), not equal to host root; app itself runs as `docker-androidemu` user, doesn't request host root

### 2. Non-root App Process (v3.6.5+)

- gateway.py, audio_fix.py and other background processes all run as `docker-androidemu` user
- Auto-downgrade via `_drop_privileges()` function (auto-switch when uid=0)
- config/privilege declares `run-as=package`

### 3. Host Network Mode

- `androidemu-webrtc` container uses `network_mode: host`
- **Reason**: fnOS blocks "bridge container → host LAN IP" access, bridge mode prevents screen service backend from connecting to local ADB port
- **Impact**: Listening ports exactly same as original port mapping (8443 TCP, 3478 TCP/UDP, 50000-50100 UDP), no new ports; container still doesn't use privileged

### 4. Docker Group Permission

- App user `docker-androidemu` belongs to `docker` group
- This is standard configuration for fnOS Docker apps, required to run containers
- Only used for app scripts to orchestrate and clean up this app's containers, doesn't affect other apps

### 5. Multi-port Listening

| Port | Protocol | Purpose | Auth |
|------|----------|---------|------|
| 8443 | TCP | WebRTC signaling + screen | Account login |
| 3478 | TCP/UDP | TURN relay | TURN credentials |
| 5556 | TCP | ADB debugging | Default localhost only, manual open needed |
| 50000-50100 | UDP | WebRTC media stream | Session-level auth |

- App doesn't do any automatic port mapping (no UPnP/hole punching), public exposure decided by user
- Each port has independent auth mechanism

### 6. Paid Features Note

- This app itself is completely free
- Built-in Scrcpy screen service uses upstream third-party authorization: free version can use all basic features (screen casting, ADB debugging etc.), only "addable device count" is limited
- Paid only increases device count limit, doesn't affect any functionality
- Fees collected by upstream authorization service provider, unrelated to fnOS official
<img width="1288" height="900" alt="firefox exe_20260927_092044" src="https://github.com/user-attachments/assets/44afdf11-2112-4ae7-8544-89e6bfa0238b" />

---

## FAQ

### Q: Shows "Unauthorized" or "0 devices online" after opening

A: This is Scrcpy License authorization prompt. Newly installed devices currently have 20 devices 3-month free trial (from Nov 1, 2026, expires to 10 devices, other basic features all free). Wait for container to fully start (~1-2 minutes) then refresh page. If still unauthorized, check if container is running:
```bash
docker ps --filter name=androidemu
```

### Q: "Unable to update - script execution error with unknown reason" during upgrade

A: This is caused by `uninstall_init` script exception during upgrade in old versions (v3.6.5 and before). Fixed in v3.6.6:
- Rewrote uninstall script with more robust upgrade guard
- Added `set +e` to ensure no command failure causes abnormal script exit
- All output redirected to log, stdout stays clean

**Solution**: Download v3.6.6+ and install via "Manual Install" to overwrite. If still fails, check log on NAS:
```bash
cat /var/apps/androidemu/var/uninstall_init.log
```

### Q: WebRTC casting connection failed

A: Troubleshoot in this order:
1. Confirm using in LAN (external network auto-downgrades to WebSocket casting)
2. Confirm TURN server running: `docker exec androidemu-webrtc netstat -ulnp | grep 3478`
3. Confirm PUBLIC_IP is LAN IP not 127.0.0.1: `docker exec androidemu-webrtc env | grep PUBLIC_IP`
4. Confirm ICE_SERVERS contains correct LAN IP: `docker exec androidemu-webrtc env | grep ICE_SERVERS`
5. Try clicking "Switch to WebSocket casting"

### Q: Container restarts repeatedly

A: Common causes:
- ARM device not enabled software rendering: check if compose has `androidboot.redroid_gpu_mode=guest`
- Insufficient memory: recommend at least 2GB free memory
- Check logs: `docker logs androidemu-android`
- webrtc container restart: check for `nice: setpriority(-10): Permission denied`, confirm compose includes `cap_add: SYS_NICE`

### Q: No sound

A: Current version disables audio by default (opus encoder in redroid container is Codec2 version, scrcpy-server only recognizes OMX version, enabling audio causes `createEncoder` failure and stream disconnect). Dual protection via RUNTIME_SHIM hijacking WebSocket.send and gateway intercepting `/api/default_settings`. Future versions will attempt fix.

### Q: External network screen laggy

A: External network via fnOS reverse proxy auto-uses WebSocket casting, framerate limited by upstream bandwidth. If you have public IP, recommend directly mapping port 8443 for WebRTC casting.

### Q: Clicking "Terminal" pops up print dialog

A: This is Scrcpy frontend shortcut conflict bug. Fixed via RUNTIME_SHIM blocking `window.print()` and Ctrl+P shortcut, clicking terminal no longer triggers printing.

### Q: Audit log loading failed

A: Audit log function depends on Scrcpy backend `/api/audit` endpoint, some versions may not support. This is upstream image feature, doesn't affect core casting function.

### Q: How to uninstall

A: Click "Uninstall" in fnOS App Center. Container data (Android /data partition) kept in Docker volume, to completely remove:
```bash
docker volume rm androidemu_data androidemu-webrtc-data
```

### Q: Can't reinstall after cancelling installation/update midway

A: Fixed in v3.6.0+, auto-cleans temporary data on interruption. If using old version and encountering this, manually clean:
```bash
rm -rf /tmp/androidemu_*
docker rm -f androidemu-android androidemu-webrtc 2>/dev/null
```

### Technical Notes & Development Pitfalls

Key technical conclusions verified during development, for secondary development reference:

- **com.android.location.fused absolutely cannot be disabled**: It's a dependency of LocationManagerService, disabling causes system_server crash, Android won't boot
- **androidboot.use_memfd=1 unsafe on redroid 12**: Causes LocationManagerService crash, init kills zygote, memory drops from 600-700MB to 400MB, container can't start
- **Detect dependency apps by device node not directory name**: App Center display name and actual directory name may differ, detecting `/dev/binder` more reliable than `/var/apps/binder_linux_driver`
- **WebSocket proxy must cover all ws/wss URLs**: Can't only cover specific paths, otherwise other paths will directly connect to port 8443 causing CORS failure when accessed via gateway
- **BusyBox su resets environment variables**: Non-root transformation can't use su -c internal export to pass env vars, should set directly in docker-compose environment
- **Disable Bluetooth with pm disable not svc disable**: svc disable gets auto-restarted by system, pm disable can completely disable
- **fnOS Docker no host hairpin capability**: Bridge containers can't access host LAN IP, ADB port needs socat forwarding
- **redroid Android main container must be privileged**: Upstream official requirement, non-privileged tested can't boot, but only applies inside container, app itself doesn't request host root
- **Upgrade scripts must be robust**: fnOS upgrade process first executes old version uninstall_init, any command failure or stdout output causes "script execution error with unknown reason", must set +e + output redirect + all commands || true
- **CRLF line endings cause Linux script errors**: All shell/Python scripts must use LF line endings, otherwise `$'\r': command not found`
- **App process downgrade needs attention to IPC**: After gateway.py downgrade, `is_mine(pid)` can only check current user processes, audio watchdog must also downgrade to same user, otherwise gets repeatedly started

---

## Known Limitations

1. **Audio**: Current version disables audio by default, enabling causes stream disconnect
2. **External network access**: fnOS reverse proxy only passes TCP, WebRTC media stream (UDP) can't pass, auto-downgrades to WebSocket casting
3. **fnOS APP**: Some versions WebView has limited WebSocket proxy support, recommend mobile browser
4. **ARM app compatibility**: ARM64 apps strongly dependent on Google services or with anti-emulator detection may crash
5. **redroid privileged mode**: Android main container needs privileged (redroid upstream official requirement, non-privileged tested can't boot), but only applies inside container, app itself doesn't request host root
6. **Scrcpy device count authorization**: Free version has device count limit, paid only increases device count, doesn't affect functionality

---

## Feedback Links & Channels

1. **Communication, feedback & beta testing QQ group**: https://qm.qq.com/q/DF7nsBatFu
2. **Suggestion & issue feedback survey**: https://wj.qq.com/s2/28029808/2aab/
3. **Publisher email**: andforlin@foxmail.com
4. **Special feedback links for redroid container and Scrcpy container authors**: See "Acknowledgements & Links" section

> Feedback links and channels except #4 will be replied to, because #4 is special feedback channel for container issues and suggestions, unrelated to the app itself; if not satisfied with feedback results or want to contribute to the app, you can modify using source code, package according to fnOS official tutorial, then upload via first three feedback links and channels. Publisher will audit uploaded code and invite you to become a contributor; also thanks to those who provide feedback, suggestions or substantial help for the app.

---

## Support the Publisher & Contributors

<img width="4096" height="2926" alt="a61d0506ea65c87e4dd005f21325eda6" src="https://github.com/user-attachments/assets/1ad0e1e8-e03e-4966-a94d-24dff71981ba" />

If you find this app useful or are satisfied with the publisher, please consider supporting via donation to help the publisher and other contributors continue maintaining the app. Any amount or no donation is appreciated; please note "donation" in payment remarks, thanks.
Or give a star to support the project.

> Note: Above donation codes only for this project's maintenance and development support, please do not盗用 or use for other purposes, thanks for understanding.

---

## Open Source License & Disclaimer

### Open Source License

- redroid: [Apache 2.0](http://www.apache.org/licenses)
- scrcpy-over-webrtc (Scrcpy): See upstream project
- This project's packaging scripts and configs: MIT

Upstream component sources and license status see `LICENSE` file in package.

### Disclaimer

1. This project is unofficial third-party app, provided "as is", use at your own risk.
2. This project is for learning and research purposes only, must not be used for any illegal purpose.
3. Developers not responsible for any data loss, system failure, service interruption etc. caused by using this app.
4. Third-party components integrated in app (redroid, Scrcpy etc.) maintained by respective authors, their functionality and stability not controlled by this project.
5. Users should backup important data themselves, this app doesn't guarantee security and integrity of data inside container.
6. This app doesn't collect any user data, all data stored on user's local device.

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
