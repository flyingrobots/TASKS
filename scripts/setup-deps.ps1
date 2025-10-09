<#
.SYNOPSIS
  Detects and optionally installs external dependencies on Windows.

.DESCRIPTION
  - Checks for Graphviz (dot). Installs via winget/choco/scoop if available.
  - Checks for Go (1.25+) and prints guidance if missing.

.PARAMETER Install
  Switch to perform installation actions (default is check-only).

.PARAMETER Yes
  Assume yes to prompts.
#>

param(
  [switch]$Install,
  [switch]$Yes
)

function Have($cmd) {
  return $null -ne (Get-Command $cmd -ErrorAction SilentlyContinue)
}
function Confirm-Action($Message) {
  if ($Yes) { return $true }
  $r = Read-Host "$Message [y/N]"
  return $r -match '^[Yy]$'
}

Write-Host "Detecting Windows tooling..."

$missing = @()

# Graphviz
if (Have 'dot') {
  $v = (& dot -V) 2>&1
  Write-Host "OK graphviz: $v"
} else {
  Write-Host "MISSING graphviz (dot) — required for DAG rendering"
  $missing += 'graphviz'
}

# Go
if (Have 'go') {
  $gv = (& go version) 2>&1
  Write-Host "OK $gv"
} else {
  Write-Host "Go not found — required to build and run the planner (go 1.25+). See https://go.dev/doc/install"
}

if ($missing.Count -eq 0) { Write-Host "All managed dependencies present."; exit 0 }

if (-not $Install) {
  Write-Host "Install plan:" ($missing -join ', ')
  Write-Host "Run with -Install [-Yes] to perform installation."
  exit 0
}

foreach ($pkg in $missing) {
  if ($pkg -eq 'graphviz') {
    if (Have 'winget') {
      if (Confirm-Action "Install Graphviz via winget?") { winget install -e --id Graphviz.Graphviz }
      continue
    }
    if (Have 'choco') {
      if (Confirm-Action "Install Graphviz via chocolatey?") { choco install graphviz -y }
      continue
    }
    if (Have 'scoop') {
      if (Confirm-Action "Install Graphviz via scoop?") { scoop install graphviz }
      continue
    }
    Write-Host "No winget/choco/scoop detected. Install Graphviz from https://graphviz.org/download/"
  }
}

Write-Host "Done."
