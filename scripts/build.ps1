param(
    [string]$TargetDir = "E:\Rusted Warfare\求生之路2",
    [switch]$NoDeploy
)

$env:PATH = "C:\Users\HOTPOTFISH\go\go\bin;C:\Users\HOTPOTFISH\go\bin;$env:PATH"
$env:GOPROXY = "https://goproxy.cn,direct"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host " Building L4D2 Mod Manager (Wails Production)" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

wails build -clean

if ($LASTEXITCODE -eq 0) {
    $exeItem = Get-ChildItem "build\bin\*.exe" | Select-Object -First 1
    if ($exeItem -and (Test-Path $exeItem.FullName)) {
        $exePath = $exeItem.FullName
        $fileSize = [math]::Round($exeItem.Length / 1MB, 2)
        Write-Host ""
        Write-Host "Build Success!" -ForegroundColor Green
        Write-Host "Output: $exePath" -ForegroundColor Green
        Write-Host "Size: $fileSize MB" -ForegroundColor Green

        if (-not $NoDeploy) {
            Write-Host ""
            $deployScript = Join-Path $PSScriptRoot "copy_to_game.ps1"
            if (Test-Path $deployScript) {
                & $deployScript -TargetDir $TargetDir -ExeName $exeItem.Name
            }
        }
    }
} else {
    Write-Host "Build Failed!" -ForegroundColor Red
    exit 1
}
