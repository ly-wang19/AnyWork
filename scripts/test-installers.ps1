$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
$temporary = Join-Path ([System.IO.Path]::GetTempPath()) ("anywork-installer-test-" + [guid]::NewGuid().ToString("N"))
$payload = Join-Path $temporary "payload"
$release = Join-Path $temporary "release"
$installDir = Join-Path $temporary "install"
$isolatedHome = Join-Path $temporary "home"
New-Item -ItemType Directory -Force -Path $payload, $release, $isolatedHome | Out-Null

try {
    $architecture = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString().ToLowerInvariant()
    switch ($architecture) {
        "x64" { $targetArchitecture = "amd64" }
        "arm64" { $targetArchitecture = "arm64" }
        default { throw "Unsupported test architecture: $architecture" }
    }
    $binary = Join-Path $payload "anywork.exe"
    & go build -C $root -trimpath -ldflags "-X main.nativeVersion=test" -o $binary .
    if ($LASTEXITCODE -ne 0) { throw "go build failed" }
    $asset = "anywork_test_windows_${targetArchitecture}.zip"
    $archive = Join-Path $release $asset
    Compress-Archive -LiteralPath $binary -DestinationPath $archive
    $hash = (Get-FileHash -LiteralPath $archive -Algorithm SHA256).Hash.ToLowerInvariant()
    Set-Content -LiteralPath (Join-Path $release "SHA256SUMS") -Value "$hash  $asset" -Encoding ascii

    $env:HOME = $isolatedHome
    $env:USERPROFILE = $isolatedHome
    & (Join-Path $root "install.ps1") -Version test -InstallDir $installDir -DownloadBase $release -Setup -Agent codex -Language en
    & (Join-Path $installDir "anywork.exe") doctor --lang en
    if (-not (Test-Path -LiteralPath (Join-Path $isolatedHome ".agents\skills\orchestrate-work\SKILL.md"))) {
        throw "setup did not install orchestrate-work for Codex"
    }
    Write-Host "Windows installer smoke test passed."
} finally {
    Remove-Item -LiteralPath $temporary -Recurse -Force -ErrorAction SilentlyContinue
}
