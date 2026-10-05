package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	resConf    = "/var/apps/androidemu/var/resolution.conf"
	logFile    = "/var/apps/androidemu/var/resolution_autoswitch.log"
	listenAddr = "127.0.0.1:18443"
	cooldown   = 60 // 冷却时间（秒）
)

var (
	lastSwitch   time.Time
	currentDev   string
	mu           sync.Mutex
	logger       *log.Logger
)

func initLogger() {
	os.MkdirAll("/var/apps/androidemu/var", 0755)
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger = log.New(os.Stderr, "", log.LstdFlags)
		return
	}
	logger = log.New(f, "", log.LstdFlags)
}

func logMsg(msg string) {
	if logger != nil {
		logger.Println(msg)
	}
}

func getResolution(device string) string {
	key := strings.ToUpper(device) + "="
	data, err := os.ReadFile(resConf)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, key) {
			return strings.TrimSpace(strings.TrimPrefix(line, key))
		}
	}
	return ""
}

func switchResolution(device string) {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	if now.Sub(lastSwitch).Seconds() < cooldown {
		return
	}
	if device == currentDev {
		return
	}

	target := getResolution(device)
	if target == "" {
		return
	}
	// 校验分辨率格式
	if !regexp.MustCompile(`^\d+x\d+$`).MatchString(target) {
		return
	}

	// 执行 adb 切换
	exec.Command("adb", "connect", "127.0.0.1:5556").Run()
	time.Sleep(1 * time.Second)

	cmd := exec.Command("adb", "-s", "127.0.0.1:5556", "shell", "wm", "size", target)
	output, err := cmd.CombinedOutput()
	exec.Command("adb", "disconnect", "127.0.0.1:5556").Run()

	if err != nil {
		logMsg(fmt.Sprintf("切换失败 %s: %s, output=%s", device, target, strings.TrimSpace(string(output))))
		return
	}

	currentDev = device
	lastSwitch = now
	logMsg(fmt.Sprintf("自动切换分辨率为 %s: %s", device, target))
}

func handler(w http.ResponseWriter, r *http.Request) {
	device := r.URL.Query().Get("device")
	if device == "pc" || device == "phone" || device == "tablet" {
		go switchResolution(device)
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", "0")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
}

func main() {
	initLogger()
	logMsg("分辨率自动切换守护脚本启动 (Go版本)")

	http.HandleFunc("/", handler)
	server := &http.Server{
		Addr:         listenAddr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		logMsg(fmt.Sprintf("HTTP服务器错误: %v", err))
		os.Exit(1)
	}
}
