OUTPUT_DIR := bin
EXECUTABLE_NAME := otp_server

all: linux macos_arm64

linux:
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(EXECUTABLE_NAME)

macos_arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(OUTPUT_DIR)/$(EXECUTABLE_NAME)_macos_arm64

.PHONY: all windows linux macos_arm64 clean