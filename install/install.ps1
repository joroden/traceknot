$ErrorActionPreference = "Stop"

$BaseUrl = if ($env:BASE_URL) { $env:BASE_URL } else { "https://downloads.traceknot.com" }
$DaemonUrl = "http://127.0.0.1:4318"
$WindowsArch = if ($env:PROCESSOR_ARCHITECTURE -eq 'ARM64') { 'arm64' } else { 'x64' }
$InstallRoot = if ($env:PREFIX) { $env:PREFIX } else { Join-Path $env:LOCALAPPDATA "traceknot" }
$BinDir = Join-Path $InstallRoot "bin"
$Binary = Join-Path $BinDir "traceknot.exe"
$InstallLog = Join-Path $env:LOCALAPPDATA "traceknot\install.log"

function Log-TkInstall {
	param([string]$Message)
	New-Item -ItemType Directory -Force -Path (Split-Path -Parent $InstallLog) | Out-Null
	Add-Content -Path $InstallLog -Value ("{0:HH:mm:ss}  {1}" -f (Get-Date), $Message)
}

function Invoke-TkDownload {
	param([string]$Url, [string]$Output)
	if (Get-Command curl.exe -ErrorAction SilentlyContinue) {
		& curl.exe -fsSL $Url -o $Output
		if ($LASTEXITCODE -ne 0) { throw "download failed: $Url" }
		return
	}
	Invoke-WebRequest $Url -UseBasicParsing -OutFile $Output
}

function Get-TkSha256 {
	param([string]$Path)
	return (Get-FileHash -Path $Path -Algorithm SHA256).Hash.ToLowerInvariant()
}

Write-Host "Installing traceknot..."
Log-TkInstall "install started"

if ($env:TRACEKNOT_VERSION) {
	$Release = $env:TRACEKNOT_VERSION
} else {
	$TempLatest = New-TemporaryFile
	try {
		Invoke-TkDownload "$BaseUrl/latest.txt" $TempLatest
		$Release = (Get-Content -Path $TempLatest -Raw).Trim()
	} finally {
		Remove-Item $TempLatest -ErrorAction SilentlyContinue
	}
}

if ([string]::IsNullOrEmpty($Release)) {
	Write-Error "Could not determine which release to install from $BaseUrl/latest.txt"
	exit 1
}

Write-Host "Installing traceknot $Release..."
Log-TkInstall "resolved release $Release"

$DownloadUrl = "$BaseUrl/releases/$Release/downloads/traceknot-windows-$WindowsArch.exe"

New-Item -ItemType Directory -Force -Path $BinDir | Out-Null
$Temp = Join-Path $BinDir ".traceknot-install-$PID"
$TempSum = Join-Path $BinDir ".traceknot-install-$PID.sha256"
try {
	Invoke-TkDownload $DownloadUrl $Temp
	Invoke-TkDownload "$DownloadUrl.sha256" $TempSum
	try { Unblock-File -Path $Temp -ErrorAction SilentlyContinue } catch {
	}

	$ExpectedSum = ((Get-Content -Path $TempSum -TotalCount 1) -split '\s+')[0].ToLowerInvariant()
	$ActualSum = Get-TkSha256 $Temp
	if ([string]::IsNullOrEmpty($ExpectedSum) -or $ActualSum -ne $ExpectedSum) {
		Log-TkInstall "ERROR: checksum verification failed for $DownloadUrl"
		Write-Error "The download from $DownloadUrl failed checksum verification (possible corruption)."
		exit 1
	}
	if (Test-Path $Binary) {
		& $Binary stop *> $null
		Start-Sleep -Milliseconds 300
	}
	Move-Item -Force $Temp $Binary
} finally {
	Remove-Item $Temp -ErrorAction SilentlyContinue
	Remove-Item $TempSum -ErrorAction SilentlyContinue
}
Log-TkInstall "binary installed at $Binary"

& $Binary post-install
Log-TkInstall "post-install done"
Log-TkInstall "install complete"

Write-Host ""
Write-Host "Installation complete!"
Write-Host "Dashboard available at: $DaemonUrl/"
Write-Host "Reconfigure anytime with: traceknot"
Write-Host "Installation log: $InstallLog"
