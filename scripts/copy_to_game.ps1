# ========================================================
# 同步/更新构建产物到游戏目录脚本
# ========================================================
param(
    [string]$TargetDir = "E:\Rusted Warfare\求生之路2",
    [string]$ExeName = "",
    [switch]$KillRunning
)

$ErrorActionPreference = "Stop"

# 项目根目录
$ProjectRoot = Split-Path -Parent $PSScriptRoot

if (-not $ExeName) {
    $found = Get-ChildItem (Join-Path $ProjectRoot "build\bin\*.exe") -ErrorAction SilentlyContinue | Select-Object -First 1
    if ($found) {
        $ExeName = $found.Name
    } else {
        $ExeName = "L4D2-Mod-Manager.exe"
    }
}
$SourceExe = Join-Path $ProjectRoot "build\bin\$ExeName"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host " 同步构建产物到游戏目录" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

# 1. 检查源文件
if (-not (Test-Path $SourceExe)) {
    Write-Host "[!] 错误: 未找到构建产物: $SourceExe" -ForegroundColor Red
    Write-Host "[*] 提示: 请先运行构建脚本 .\scripts\build.ps1" -ForegroundColor Yellow
    exit 1
}

$sourceItem = Get-Item $SourceExe
$sourceSizeMB = [math]::Round($sourceItem.Length / 1MB, 2)
Write-Host "[*] 源文件: $SourceExe ($sourceSizeMB MB)" -ForegroundColor Gray

# 2. 检查并确保目标目录存在
if (-not (Test-Path $TargetDir)) {
    Write-Host "[*] 目标目录不存在，正在创建: $TargetDir" -ForegroundColor Yellow
    New-Item -ItemType Directory -Path $TargetDir -Force | Out-Null
}

$TargetExe = Join-Path $TargetDir $ExeName
$isUpdate = Test-Path $TargetExe

if ($isUpdate) {
    Write-Host "[*] 检测到目标已存在该文件，执行覆盖更新..." -ForegroundColor Yellow
} else {
    Write-Host "[*] 目标不存在该文件，执行首次放置..." -ForegroundColor Green
}

# 3. 检查是否有运行中的目标进程占用
$runningProc = Get-Process -Name ([System.IO.Path]::GetFileNameWithoutExtension($ExeName)) -ErrorAction SilentlyContinue | Where-Object {
    try {
        $_.Path -and ($_.Path.TrimEnd('\') -eq $TargetExe.TrimEnd('\'))
    } catch {
        $false
    }
}

if ($runningProc) {
    if ($KillRunning) {
        Write-Host "[!] 检测到目标正在运行 (PID: $($runningProc.Id -join ', '))，正在关闭..." -ForegroundColor Yellow
        $runningProc | Stop-Process -Force
        Start-Sleep -Milliseconds 500
    } else {
        Write-Host "[!] 警告: 目标程序当前正在运行中 (PID: $($runningProc.Id -join ', '))！" -ForegroundColor Yellow
        Write-Host "    如果覆盖失败，请先关闭该程序或添加 -KillRunning 参数。" -ForegroundColor Yellow
    }
}

# 4. 执行复制覆盖
try {
    Copy-Item -Path $SourceExe -Destination $TargetExe -Force
    $targetItem = Get-Item $TargetExe
    $targetSizeMB = [math]::Round($targetItem.Length / 1MB, 2)

    Write-Host ""
    Write-Host "==========================================" -ForegroundColor Green
    if ($isUpdate) {
        Write-Host "[+] 更新成功!" -ForegroundColor Green
    } else {
        Write-Host "[+] 复制成功!" -ForegroundColor Green
    }
    Write-Host "目标路径: $TargetExe" -ForegroundColor Green
    Write-Host "文件大小: $targetSizeMB MB" -ForegroundColor Green
    Write-Host "更新时间: $($targetItem.LastWriteTime.ToString('yyyy-MM-dd HH:mm:ss'))" -ForegroundColor Green
    Write-Host "==========================================" -ForegroundColor Green
} catch {
    Write-Host ""
    Write-Host "[-] 复制失败: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "[*] 提示: 很可能是目标程序正在运行被占用，请退出程序后再试。" -ForegroundColor Yellow
    exit 1
}
