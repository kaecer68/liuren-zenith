---
description: Sync contract ports into local env
---
1. Confirm contract runtime ports file exists at `../destiny-contracts/runtime/ports.env` (or set `CONTRACT_PORTS_FILE`).
2. Run `make sync-contracts` to regenerate `.env.ports` from contract values.
3. Run `make verify-contracts` to ensure local generated ports stay aligned with contract.
4. Start service with contract-driven ports using `make run`.
5. If verification fails, update contract-side ports file first, then re-run steps 2-4.
