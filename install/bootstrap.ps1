$ErrorActionPreference = "Stop"

$DaemonUrl = "http://127.0.0.1:4318"
$WindowsArch = if ($env:PROCESSOR_ARCHITECTURE -eq 'ARM64') { 'arm64' } else { 'x64' }
$BinaryName = "traceknot-windows-$WindowsArch.exe"
$InstallRoot = if ($env:PREFIX) { $env:PREFIX } else { Join-Path $env:LOCALAPPDATA "traceknot" }
$BinDir = Join-Path $InstallRoot "bin"
$Binary = Join-Path $BinDir "traceknot.exe"
$InstallLog = Join-Path $env:LOCALAPPDATA "traceknot\install.log"
$SourceBinary = Join-Path $PSScriptRoot "downloads\$BinaryName"

function Log-TkInstall {
	param([string]$Message)
	New-Item -ItemType Directory -Force -Path (Split-Path -Parent $InstallLog) | Out-Null
	Add-Content -Path $InstallLog -Value ("{0:HH:mm:ss}  {1}" -f (Get-Date), $Message)
}

if (-not (Test-Path $SourceBinary)) {
	Write-Error "$SourceBinary not found next to bootstrap.ps1."
	exit 1
}

Write-Host "Installing traceknot..."
Log-TkInstall "install started"

New-Item -ItemType Directory -Force -Path $BinDir | Out-Null
$Temp = Join-Path $BinDir ".traceknot-install-$PID"
try {
	Copy-Item -Force $SourceBinary $Temp
	if (Test-Path $Binary) {
		& $Binary stop *> $null
		Start-Sleep -Milliseconds 300
	}
	Move-Item -Force $Temp $Binary
} finally {
	Remove-Item $Temp -ErrorAction SilentlyContinue
}
Log-TkInstall "binary installed at $Binary"

& $Binary
Log-TkInstall "install: interactive menu completed"
Log-TkInstall "install complete"

Write-Host ""
Write-Host "Installation complete!"
Write-Host "Dashboard available at: $DaemonUrl/"
Write-Host "Reconfigure anytime with: traceknot"
Write-Host "Installation log: $InstallLog"
