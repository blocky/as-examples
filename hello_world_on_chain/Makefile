.PHONY: test-local
test-local:
	npx hardhat test --grep "Local"

.PHONY: compile
compile:
	npx hardhat compile

.PHONY: deploy-base-sepolia
deploy-base-sepolia: compile
	npx hardhat run scripts/deployUser.ts --network base-sepolia
	npx hardhat verify --network base-sepolia $$(cat .user_deployed_address)

.PHONY: test-base-sepolia
test-base-sepolia:
	npx hardhat test --grep "Base Sepolia"

.PHONY: very-clean
very-clean:
	@rm -rf artifacts cache node_modules typechain-types
