PROTO_DIR = third_party/proto
COSMOS_VER_SHORT = 0.42.9
COSMOS_VER = v$(COSMOS_VER_SHORT)

proto-update-dep:
	@mkdir -p $(PROTO_DIR)
	@curl https://codeload.github.com/cosmos/cosmos-sdk/tar.gz/$(COSMOS_VER) | tar -xz -C $(PROTO_DIR) --strip=3 cosmos-sdk-$(COSMOS_VER_SHORT)/third_party/proto
	@curl https://codeload.github.com/cosmos/cosmos-sdk/tar.gz/$(COSMOS_VER) | tar -xz -C $(PROTO_DIR) --strip=2 cosmos-sdk-$(COSMOS_VER_SHORT)/proto