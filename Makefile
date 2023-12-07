COMPOSEFLAGS=-d
ITESTS_L2_HOST=http://localhost:9545
BEDROCK_TAGS_REMOTE?=origin
DEVNET_NODEKIT_FLAGS=--nodekit --deploy-l2 --deploy-config="devnetL1-nodekit.json" --deploy-config-template="devnetL1-nodekit-template.json" --deployment="devnetL1-nodekit" --devnet-dir=".devnet-nodekit" --l2-provider-url="http://localhost:9090"
url="http://localhost:17111"
monorepo-base := $(realpath .)
OP_STACK_GO_BUILDER?=us-docker.pkg.dev/oplabs-tools-artifacts/images/op-stack-go:latest

# Requires at least Python v3.9; specify a minor version below if needed
PYTHON?=python3

build: build-go build-ts
.PHONY: build

build-go: submodules op-node op-proposer op-batcher
.PHONY: build-go

lint-go:
	golangci-lint run -E goimports,sqlclosecheck,bodyclose,asciicheck,misspell,errorlint --timeout 5m -e "errors.As" -e "errors.Is" ./...
.PHONY: lint-go

build-ts: submodules
	if [ -n "$$NVM_DIR" ]; then \
		. $$NVM_DIR/nvm.sh && nvm use; \
	fi
	pnpm install
	pnpm build
.PHONY: build-ts

ci-builder:
	docker build -t ci-builder -f ops/docker/ci-builder/Dockerfile .

golang-docker:
	# We don't use a buildx builder here, and just load directly into regular docker, for convenience.
	GIT_COMMIT=$$(git rev-parse HEAD) \
	GIT_DATE=$$(git show -s --format='%ct') \
	IMAGE_TAGS=$$GIT_COMMIT,latest \
	docker buildx bake \
			--progress plain \
			--load \
			-f docker-bake.hcl \
			op-node op-batcher op-proposer op-challenger
.PHONY: golang-docker

submodules:
	# CI will checkout submodules on its own (and fails on these commands)
	if [ -z "$$GITHUB_ENV" ]; then \
		git submodule init; \
		git submodule update --recursive; \
	fi
.PHONY: submodules

op-bindings:
	make -C ./op-bindings
.PHONY: op-bindings

make op-bindings-docker:
	docker run -v $(monorepo-base):/work -it us-docker.pkg.dev/oplabs-tools-artifacts/images/ci-builder make -C /work op-bindings
	echo "Asking for root permissions to set owner of files to ${USER} after docker run"
	sudo chown -R ${USER} $(monorepo-base)

.PHONY: op-bindings-docker

op-node:
	make -C ./op-node op-node
.PHONY: op-node

generate-mocks-op-node:
	make -C ./op-node generate-mocks
.PHONY: generate-mocks-op-node

generate-mocks-op-service:
	make -C ./op-service generate-mocks
.PHONY: generate-mocks-op-service

op-batcher:
	make -C ./op-batcher op-batcher
.PHONY: op-batcher

op-proposer:
	make -C ./op-proposer op-proposer
.PHONY: op-proposer

op-challenger:
	make -C ./op-challenger op-challenger
.PHONY: op-challenger

op-program:
	make -C ./op-program op-program
.PHONY: op-program

op-geth-proxy:
	make -C ./op-geth-proxy op-geth-proxy
.PHONY: op-geth-proxy

cannon:
	make -C ./cannon cannon
.PHONY: cannon

cannon-prestate: op-program cannon
	./cannon/bin/cannon load-elf --path op-program/bin/op-program-client.elf --out op-program/bin/prestate.json --meta op-program/bin/meta.json
	./cannon/bin/cannon run --proof-at '=0' --stop-at '=1' --input op-program/bin/prestate.json --meta op-program/bin/meta.json --proof-fmt 'op-program/bin/%d.json' --output ""
	mv op-program/bin/0.json op-program/bin/prestate-proof.json

mod-tidy:
	# Below GOPRIVATE line allows mod-tidy to be run immediately after
	# releasing new versions. This bypasses the Go modules proxy, which
	# can take a while to index new versions.
	#
	# See https://proxy.golang.org/ for more info.
	export GOPRIVATE="github.com/ethereum-optimism" && go mod tidy
.PHONY: mod-tidy

clean:
	rm -rf ./bin
.PHONY: clean

nuke: clean devnet-clean
	git clean -Xdf
.PHONY: nuke

pre-devnet:
	@if ! [ -x "$(command -v geth)" ]; then \
		make install-geth; \
	fi
	@if [ ! -e op-program/bin ]; then \
		make cannon-prestate; \
	fi
.PHONY: pre-devnet

devnet-up: pre-devnet
	./ops/scripts/newer-file.sh .devnet/allocs-l1.json ./packages/contracts-bedrock \
		|| make devnet-allocs
	PYTHONPATH=./bedrock-devnet $(PYTHON) ./bedrock-devnet/main.py --monorepo-dir=.
.PHONY: devnet-up

#TODO the below was in devnet-up-nodekit but I removed it
#$(shell ./ops/scripts/newer-file.sh .devnet-nodekit/allocs-l1.json ./packages/contracts-bedrock)
# if [ $(.SHELLSTATUS) -ne 0 ]; then \
# 	make devnet-allocs-nodekit; \
# fi
devnet-up-nodekit: pre-devnet
	PYTHONPATH=./bedrock-devnet python3 ./bedrock-devnet/main.py --monorepo-dir=. $(DEVNET_NODEKIT_FLAGS)
.PHONY: devnet-up-nodekit

# alias for devnet-up
devnet-up-deploy: devnet-up

devnet-test: pre-devnet
	PYTHONPATH=./bedrock-devnet $(PYTHON) ./bedrock-devnet/main.py --monorepo-dir=. --test
.PHONY: devnet-test

devnet-down:
	@(cd ./ops-bedrock && GENESIS_TIMESTAMP=$(shell date +%s) docker compose stop)
.PHONY: devnet-down

devnet-clean:
	rm -rf ./packages/contracts-bedrock/deployments/devnetL1
	rm -rf ./packages/contracts-bedrock/deployments/devnetL1-nodekit
	rm -rf ./.devnet
	rm -rf ./.devnet-nodekit
	cd ./ops-bedrock && docker compose down -v
	docker image ls 'ops-bedrock*' --format='{{.Repository}}' | xargs -r docker rmi
	docker volume ls --filter name=ops-bedrock --format='{{.Name}}' | xargs -r docker volume rm
.PHONY: devnet-clean

devnet-allocs: pre-devnet
	PYTHONPATH=./bedrock-devnet $(PYTHON) ./bedrock-devnet/main.py --monorepo-dir=. --allocs

devnet-allocs-nodekit:
	PYTHONPATH=./bedrock-devnet python3 ./bedrock-devnet/main.py --monorepo-dir=. $(DEVNET_NODEKIT_FLAGS) --allocs

devnet-logs:
	@(cd ./ops-bedrock && docker compose logs -f)
.PHONY: devnet-logs

devnet-build:
	@(cd ./ops-bedrock && docker compose build)
.PHONY: devnet-build

devnet-pull:
	@(cd ./ops-bedrock && docker compose pull)
.PHONY: devnet-pull

e2e-pull:
	@(cd ./op-e2e && docker-compose pull)
.PHONY: e2e-pull

test-unit:
	make -C ./op-node test
	make -C ./op-proposer test
	make -C ./op-batcher test
	make -C ./op-e2e test
	pnpm test
.PHONY: test-unit

test-integration:
	bash ./ops-bedrock/test-integration.sh \
		./packages/contracts-bedrock/deployments/devnetL1
.PHONY: test-integration

# Remove the baseline-commit to generate a base reading & show all issues
semgrep:
	$(eval DEV_REF := $(shell git rev-parse develop))
	SEMGREP_REPO_NAME=ethereum-optimism/optimism semgrep ci --baseline-commit=$(DEV_REF)
.PHONY: semgrep

clean-node-modules:
	rm -rf node_modules
	rm -rf packages/**/node_modules

tag-bedrock-go-modules:
	./ops/scripts/tag-bedrock-go-modules.sh $(BEDROCK_TAGS_REMOTE) $(VERSION)
.PHONY: tag-bedrock-go-modules

update-op-geth:
	./ops/scripts/update-op-geth.py
.PHONY: update-op-geth

bedrock-markdown-links:
	docker run --init -it -v `pwd`:/input lycheeverse/lychee --verbose --no-progress --exclude-loopback \
		--exclude twitter.com --exclude explorer.optimism.io --exclude linux-mips.org \
		--exclude-mail /input/README.md "/input/specs/**/*.md"

install-geth:
	./ops/scripts/geth-version-checker.sh && \
	 	(echo "Geth versions match, not installing geth..."; true) || \
 		(echo "Versions do not match, installing geth!"; \
 			go install -v github.com/ethereum/go-ethereum/cmd/geth@$(shell cat .gethrc); \
 			echo "Installed geth!"; true)
.PHONY: install-geth

