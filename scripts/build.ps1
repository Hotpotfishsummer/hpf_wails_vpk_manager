# 构建生产版本单文件 .exe
$env:PATH = "C:\Users\HOTPOTFISH\go\go\bin;C:\Users\HOTPOTFISH\go\bin;$env:PATH"
$env:GOPROXY = "https://goproxy.cn,direct"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host " 构建 L4D2 VPK 管理器 (Go + Wails 生产版本)" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

wails build -clean

if ($LASTEXITCODE -eq 0) {
    $exePath = "build\bin\L4D2-VPK-Manager.exe"
    if (Test-Path $exePath) {
        $fileSize = [math]::Round((Get-Item $exePath).Length / 1MB, 2)
        Write-Host ""
        Write-Host "构建成功！" -ForegroundColor Green
        Write-Host "输出路径: $exePath" -ForegroundColor Green
        Write-Host "文件大小: $fileSize MB" -ForegroundColor Green
    }
} else {
    Write-Host "构建失败！" -ForegroundColor Red
    exit 1
}
