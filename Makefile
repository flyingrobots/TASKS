.PHONY: setup setup-dry-run setup-win example-dot

setup:
	@bash scripts/setup-deps.sh --install --yes

setup-dry-run:
	@bash scripts/setup-deps.sh

setup-win:
	@pwsh -File scripts/setup-deps.ps1 -Install -Yes

example-dot:
	@./examples/dot-export/render.sh

