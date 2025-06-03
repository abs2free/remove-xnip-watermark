# define var
BINARY_NAME=remove-xnip-watermark
GO_FILES=main.go
BIN_DIR=bin

# target platform
PLATFORMS= darwin/amd64 darwin/arm64

all: $(PLATFORMS)

# rule：build different platform
$(PLATFORMS):
	@IFS=/; for platform in $@; do \
		OS=$${platform%%/*}; \
		ARCH=$${platform##*/}; \
		OUTPUT=$(BIN_DIR)/$(BINARY_NAME)_$${OS}_$${ARCH}; \
		echo "Building for $$OS/$$ARCH..."; \
		GOOS=$$OS GOARCH=$$ARCH go build -o "$$OUTPUT" $(GO_FILES); \
	done

# clean all binary file
clean:
	rm -f $(BIN_DIR)/$(BINARY_NAME)*

.PHONY: all clean

