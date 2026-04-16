.PHONY: run sync-contracts verify-contracts dev-clean proto clean

proto:
	mkdir -p gen/liurenpb
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/liuren.proto
	mv proto/*.pb.go gen/liurenpb/ 2>/dev/null || true
	mv proto/*_grpc.pb.go gen/liurenpb/ 2>/dev/null || true

clean:
	rm -rf gen/liurenpb

run:
	@chmod +x scripts/sync-contracts.sh
	bash scripts/sync-contracts.sh
	bash -c 'set -a; . ./.env.ports; set +a; go run ./cmd/server/main.go'

sync-contracts:
	@chmod +x scripts/sync-contracts.sh
	bash scripts/sync-contracts.sh

verify-contracts:
	@chmod +x scripts/sync-contracts.sh
	bash scripts/sync-contracts.sh --check

dev-clean:
	@chmod +x scripts/dev-clean.sh
	bash scripts/dev-clean.sh
