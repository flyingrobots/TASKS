.PHONY: all test clean setup setup-dry-run setup-win example-dot

all: test

test:
	@cd planner && go test ./...

clean:
	@cd planner && go clean

setup:
	@bash scripts/setup-deps.sh --install --yes

setup-dry-run:
	@bash scripts/setup-deps.sh

setup-win:
	@pwsh -File scripts/setup-deps.ps1 -Install -Yes

example-dot:
	@./examples/dot-export/render.sh
