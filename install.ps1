param(
    [string]$Version = $(if ($env:ANYWORK_VERSION) { $env:ANYWORK_VERSION } else { "v0.3.0-alpha.1" }),
    [string]$InstallDir = $(if ($env:ANYWORK_INSTALL_DIR) { $env:ANYWORK_INSTALL_DIR } else { Join-Path $HOME ".local\bin" }),
    [switch]$Setup,
    [string]$Agent = "all",
    [ValidateSet("", "en", "zh-CN", "ja")]
    [string]$Language = "",
    [string]$DownloadBase = $env:ANYWORK_DOWNLOAD_BASE
)

$ErrorActionPreference = "Stop"
$repository = "ly-wang19/AnyWork"

if (-not $IsWindows -and $PSVersionTable.PSEdition -eq "Core") {
    throw "install.ps1 supports Windows; use install.sh on macOS or Linux."
}

$architecture = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString().ToLowerInvariant()
switch ($architecture) {
    "x64" { $targetArchitecture = "amd64" }
    "arm64" { $targetArchitecture = "arm64" }
    default { throw "Unsupported Windows architecture: $architecture" }
}

$asset = "anywork_${Version}_windows_${targetArchitecture}.zip"
if (-not $DownloadBase) {
    $DownloadBase = "https://github.com/$repository/releases/download/$Version"
}
$DownloadBase = $DownloadBase.TrimEnd("/", "\")
$temporary = Join-Path ([System.IO.Path]::GetTempPath()) ("anywork-install-" + [guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $temporary | Out-Null

function Copy-OrDownload([string]$Source, [string]$Destination) {
    if (Test-Path -LiteralPath $Source) {
        Copy-Item -LiteralPath $Source -Destination $Destination
        return
    }
    Invoke-WebRequest -Uri $Source -OutFile $Destination -UseBasicParsing
}

function Release-Source([string]$Name) {
    if ($DownloadBase -match "^https?://") {
        return "$DownloadBase/$Name"
    }
    return Join-Path $DownloadBase $Name
}

try {
    $archive = Join-Path $temporary $asset
    $checksums = Join-Path $temporary "SHA256SUMS"
    Copy-OrDownload (Release-Source $asset) $archive
    Copy-OrDownload (Release-Source "SHA256SUMS") $checksums

    $escapedAsset = [regex]::Escape($asset)
    $checksumLine = Get-Content -LiteralPath $checksums | Where-Object { $_ -match "\s(?:\./)?$escapedAsset$" } | Select-Object -First 1
    if (-not $checksumLine) {
        throw "Checksum entry missing for $asset"
    }
    $expected = ($checksumLine -split "\s+")[0].ToLowerInvariant()
    $actual = (Get-FileHash -LiteralPath $archive -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($expected -ne $actual) {
        throw "Checksum mismatch for $asset"
    }

    $expanded = Join-Path $temporary "expanded"
    Expand-Archive -LiteralPath $archive -DestinationPath $expanded -Force
    $sourceBinary = Join-Path $expanded "anywork.exe"
    if (-not (Test-Path -LiteralPath $sourceBinary)) {
        throw "Archive does not contain anywork.exe"
    }
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    $destination = Join-Path $InstallDir "anywork.exe"
    Copy-Item -LiteralPath $sourceBinary -Destination $destination -Force

    & $destination --version
    Write-Host "Installed AnyWork to $destination"
    if ($Setup) {
        $setupArguments = @("setup", "--agent", $Agent, "--scope", "user")
        if ($Language) {
            $setupArguments += @("--lang", $Language)
        }
        & $destination @setupArguments
    } else {
        Write-Host "Activate all capabilities: & '$destination' setup --agent all"
    }
    if (($env:Path -split ";") -notcontains $InstallDir) {
        Write-Host "Add $InstallDir to PATH to run 'anywork' from any PowerShell session."
    }
} finally {
    Remove-Item -LiteralPath $temporary -Recurse -Force -ErrorAction SilentlyContinue
}
