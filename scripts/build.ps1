$env:PATH = "C:\Users\HOTPOTFISH\go\go\bin;C:\Users\HOTPOTFISH\go\bin;$env:PATH"
$env:GOPROXY = "https://goproxy.cn,direct"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host " Building L4D2 VPK Manager (Wails Production)" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

wails build -clean

if ($LASTEXITCODE -eq 0) {
    $exePath = "build\bin\L4D2-VPK-Manager.exe"
    if (Test-Path $exePath) {
        $fileSize = [math]::Round((Get-Item $exePath).Length / 1MB, 2)
        Write-Host ""
        Write-Host "Build Success!" -ForegroundColor Green
        Write-Host "Output: $exePath" -ForegroundColor Green
        Write-Host "Size: $fileSize MB" -ForegroundColor Green
    }
} else {
    Write-Host "Build Failed!" -ForegroundColor Red
    exit 1
}
