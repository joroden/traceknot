.PHONY: ui dev run build dist dev-install clean

ui:
	cd ui && pnpm build

dev:
	air

run: ui
	go run ./cmd/traceknot

build: ui
	go build -o bin/traceknot ./cmd/traceknot

dist: ui
	rm -rf dist
	mkdir -p dist/downloads
	cp install/install.sh install/install.ps1 install/bootstrap.sh install/bootstrap.ps1 dist/
	@set -e; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/downloads/traceknot-linux-x64 ./cmd/traceknot & \
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o dist/downloads/traceknot-linux-arm64 ./cmd/traceknot & \
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o dist/downloads/traceknot-darwin-x64 ./cmd/traceknot & \
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o dist/downloads/traceknot-darwin-arm64 ./cmd/traceknot & \
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o dist/downloads/traceknot-windows-x64.exe ./cmd/traceknot & \
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -o dist/downloads/traceknot-windows-arm64.exe ./cmd/traceknot & \
	wait
	@cd dist/downloads && for f in traceknot-*; do \
		if command -v sha256sum >/dev/null 2>&1; then \
			sha256sum "$$f" > "$$f.sha256"; \
		else \
			shasum -a 256 "$$f" > "$$f.sha256"; \
		fi; \
	done

dev-install: dist
	@echo "traceknot uninstall"
	@-"$${PREFIX:-$$HOME/.local}/bin/traceknot" uninstall
	bash ./dist/bootstrap.sh

clean:
	rm -rf ui/dist bin dist
