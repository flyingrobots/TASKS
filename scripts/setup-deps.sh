#!/usr/bin/env bash
set -euo pipefail

# setup-deps.sh — Detect and (optionally) install external CLI deps used by this repo.
# Safe by default: prints a plan. Pass --install to actually install. Use --yes to skip prompts.
# Currently managed deps: graphviz (dot). Also checks for Go (min 1.25) but does not auto-install.

GREEN='\033[0;32m'; YELLOW='\033[1;33m'; RED='\033[0;31m'; NC='\033[0m'

confirm() {
  local msg=${1:-Proceed?}
  if [[ "${ASSUME_YES:-false}" == "true" ]]; then return 0; fi
  read -r -p "$msg [y/N] " ans || true
  [[ "$ans" =~ ^[Yy]$ ]]
}

have() { command -v "$1" >/dev/null 2>&1; }

OS="$(uname -s || echo unknown)"
INSTALL=false
ASSUME_YES=false

while (( "$#" )); do
  case "$1" in
    --install) INSTALL=true ; shift ;;
    --yes|-y) ASSUME_YES=true ; shift ;;
    --help|-h)
      cat <<USAGE
Usage: scripts/setup-deps.sh [--install] [--yes]

Checks for external tools and offers to install them.

Manages:
  - Graphviz (dot): required to render DAG DOT -> SVG/PNG

Checks only (no auto-install):
  - Go >= 1.25 (planner builds)

Examples:
  scripts/setup-deps.sh                    # dry-run checks only
  scripts/setup-deps.sh --install --yes    # install without prompts
USAGE
      exit 0 ;;
    *) echo -e "${YELLOW}Unknown flag:${NC} $1" >&2; exit 2 ;;
  esac
done

echo "Detecting platform: $OS"

pkg_install() {
  local pkg="$1"
  case "$OS" in
    Darwin)
      if have brew; then echo "brew install $pkg"; ${INSTALL} && { confirm "Install $pkg via Homebrew?" && brew install "$pkg"; } ;
      elif have port; then echo "sudo port install $pkg"; ${INSTALL} && { confirm "Install $pkg via MacPorts?" && sudo port install "$pkg"; } ;
      else echo -e "${YELLOW}No Homebrew/MacPorts detected. Install Homebrew first: https://brew.sh${NC}"; return 1; fi ;;
    Linux)
      if have apt; then echo "sudo apt update && sudo apt install -y $pkg"; ${INSTALL} && { confirm "Install $pkg via apt?" && sudo apt update && sudo apt install -y "$pkg"; } ;
      elif have dnf; then echo "sudo dnf install -y $pkg"; ${INSTALL} && { confirm "Install $pkg via dnf?" && sudo dnf install -y "$pkg"; } ;
      elif have pacman; then echo "sudo pacman -S --noconfirm $pkg"; ${INSTALL} && { confirm "Install $pkg via pacman?" && sudo pacman -S --noconfirm "$pkg"; } ;
      elif have zypper; then echo "sudo zypper install -y $pkg"; ${INSTALL} && { confirm "Install $pkg via zypper?" && sudo zypper install -y "$pkg"; } ;
      elif have nix; then echo "nix profile install nixpkgs#$pkg"; ${INSTALL} && { confirm "Install $pkg via Nix?" && nix profile install "nixpkgs#$pkg"; } ;
      else echo -e "${YELLOW}No known package manager detected.${NC}"; return 1; fi ;;
    *) echo -e "${YELLOW}Automatic install not supported on $OS. Use winget/choco/scoop (Windows) or install manually.${NC}"; return 1 ;;
  esac
}

missing=()

# Check Graphviz (dot)
if have dot; then echo -e "${GREEN}OK${NC} graphviz (dot) is installed: $(dot -V 2>&1)"; else
  echo -e "${RED}MISSING${NC} graphviz (dot) — required for DAG rendering"
  missing+=(graphviz)
fi

# Check Go (planner builds)
if have go; then
  gov="$(go version 2>/dev/null || true)"
  echo -e "${GREEN}OK${NC} $gov"
else
  echo -e "${YELLOW}Go not found${NC} — required to build and run the planner (go 1.25+). See https://go.dev/doc/install"
fi

if ((${#missing[@]}==0)); then
  echo -e "${GREEN}All managed dependencies present.${NC}"
  exit 0
fi

echo
echo "Install plan: ${missing[*]}"
for pkg in "${missing[@]}"; do
  pkg_install "$pkg" || true
done

echo -e "${GREEN}Done.${NC}"

