# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build -tags $(env)
GOTEST=$(GOCMD) test

# Set and confirm environment `BRIDGE`, it should be one of devnet/testnet/mainnet
env=$(BRIDGE)
BaseDir=build/$(env)

.PHONY: all test clean

deploy_tool:
	@mkdir -p $(BaseDir)/deploy_tool
	@$(GOBUILD) -o $(BaseDir)/deploy_tool/deploy_tool chain_tool/*.go

clean:
	@rm -rf $(BaseDir)/deploy_tool/deploy_tool