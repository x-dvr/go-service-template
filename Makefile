BUILD_DIR := ./bin

.PHONY: debug
debug:
	@mkdir -p $(BUILD_DIR)
	go build -gcflags="all=-N -l" -o $(BUILD_DIR)/http ./cmd/http
