# 一键启动 Wails 开发模式 (热重载)
$env:PATH = "C:\Users\HOTPOTFISH\go\go\bin;C:\Users\HOTPOTFISH\go\bin;$env:PATH"
$env:GOPROXY = "https://goproxy.cn,direct"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host " 启动 L4D2 VPK 管理器 (Go + Wails 开发模式)" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

wails dev
