.SILENT: check-go build build-linux-amd64 build-linux-arm64 build-macos-arm64 clean

MAIN_PATH=.
BINARY_NAME=bin/gotail

all: clean build install

check-go:
	@echo "Cheking Go..."
	@if ! command -v go >/dev/null 2>&1; then \
		echo "Error: Go is not installed. Please install Go before proceeding."; \
		exit 1; \
	fi

build:
	@if [ "$$(uname -s)" = "Linux" ]; then \
		if [ "$$(uname -m)" = "x86_64" ]; then \
			$(MAKE) build-linux-amd64; \
		elif [ "$$(uname -m)" = "aarch64" ]; then \
			$(MAKE) build-linux-arm64; \
		else \
			echo "Unsupported Linux architecture: $$(uname -m)"; \
			exit 1; \
		fi; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		if [ "$$(uname -m)" = "arm64" ]; then \
			$(MAKE) build-macos-arm64; \
		else \
			echo "Unsupported macOS architecture: $$(uname -m)"; \
			exit 1; \
		fi; \
	else \
		echo "Unsupported operating system: $$(uname -s)"; \
		exit 1; \
	fi
		
build-linux-amd64: check-go
	@echo "Building linux amd64 binary..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/borisdvlpr/gotail/cmd.version=$$(cat version.txt)" -o ${BINARY_NAME}-linux-amd64 ${MAIN_PATH}
	@echo "Build completed."

build-linux-arm64: check-go
	@echo "Building linux arm64 binary..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-X github.com/borisdvlpr/gotail/cmd.version=$$(cat version.txt)" -o ${BINARY_NAME}-linux-arm64 ${MAIN_PATH}
	@echo "Build completed."

build-macos-arm64: check-go
	@echo "Building macOS arm64 binary..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-X github.com/borisdvlpr/gotail/cmd.version=$$(cat version.txt)" -o ${BINARY_NAME}-darwin-arm64 ${MAIN_PATH}
	@echo "Build completed."

install:
	@echo "Installing gotail on $$HOME/.local/bin..."
	@if [ ! -d "$$HOME/.local/bin" ]; then mkdir -p "$$HOME/.local/bin"; fi;

	@if [ "$$(uname -s)" = "Linux" ]; then \
		if [ "$$(uname -m)" = "x86_64" ]; then \
			cp bin/gotail-linux-amd64 $$HOME/.local/bin/gotail; \
		elif [ "$$(uname -m)" = "aarch64" ]; then \
			cp bin/gotail-linux-arm64 $$HOME/.local/bin/gotail; \
		else \
			echo "Unsupported Linux architecture: $$(uname -m)"; \
			exit 1; \
		fi; \
	elif [ "$$(uname -s)" = "Darwin" ]; then \
		if [ "$$(uname -m)" = "arm64" ]; then \
			cp bin/gotail-darwin-arm64 $$HOME/.local/bin/gotail; \
		else \
			echo "Unsupported macOS architecture: $$(uname -m)"; \
			exit 1; \
		fi; \
	else \
		echo "Unsupported operating system: $$(uname -s)"; \
		exit 1; \
	fi

clean:
	@echo "Cleaning local installation..."
	@go clean && rm -r bin/	&& rm ~/.local/bin/gotail
