AGENTS.md
Product
Build a B2B multi-tenant payment orchestration platform with a standardized REST API and MCP server.
Initial providers: Xendit and Midtrans. This is not a new payment processor.
Core Principles

- Do not create, modify, or write code to project files unless explicitly instructed by the user. You may provide code snippets, examples, and technical explanations directly in the chat.
- Prefer a simple modular monolith; do not introduce microservices without a demonstrated need.
- REST API and MCP must call the same application services. Never duplicate payment business logic.
- Keep provider-specific behavior inside provider adapters.
- Normalize provider data while preserving raw provider statuses and external identifiers.
- Implement only the MVP scope. Avoid speculative abstractions and overengineering.
  MVP Scope
- Tenant and merchant-account isolation.
- Provider connections for Xendit and Midtrans.
- Create payment, retrieve payment, list/filter payments, and sync payment status.
- Verified webhooks with durable event processing and deduplication.
- MCP read tools for payment lookup, failed payments, summaries, and reconciliation reports.
- Basic settlement reconciliation where provider data is available.
- Refunds, payouts, split payments, and automatic provider routing are out of scope unless explicitly requested.
  Architecture
- Language: Go. Database: PostgreSQL.
- Organize code by domain/module, e.g. tenant, merchant, credentials, payment, provider/xendit, provider/midtrans, webhook, reconciliation, reporting, audit, and mcp.
- Keep handlers thin: validate input, call application services, map responses.
- Use a PaymentProvider interface for common capabilities; expose provider-specific capabilities explicitly rather than pretending all providers behave identically.
- Use database migrations for schema changes. Never make undocumented manual schema changes.
- Use context-aware I/O, bounded timeouts, structured logs, and graceful error handling.
  Payment Correctness
- Use integer minor units for money; never use floating point for amounts.
- Keep payment status, provider status, refund status, and settlement status separate.
- Define and test explicit status mapping and transition rules.
- Persist internal idempotency keys and request hashes; reject key reuse with a different payload.
- On ambiguous provider timeouts, check/reconcile the existing operation before creating another payment.
- Treat webhooks as at-least-once and potentially out of order. Verify authenticity, persist before processing, deduplicate, and process safely.
- Never assume provider features, status values, expiry rules, or idempotency guarantees are identical. Confirm them against current official documentation.
  Multi-Tenancy and Security
- Derive tenant identity from authenticated credentials; never trust a client- or model-supplied tenant ID for authorization.
- Scope every query and mutation to the authorized tenant and merchant account.
- Use PostgreSQL constraints and consider RLS as defense in depth, not as a replacement for application authorization.
- Store provider secrets in a secret manager or encrypted secret store; keep only references in ordinary database fields where practical.
- Never log secrets, authorization headers, or unnecessary sensitive customer/payment data.
- Authenticate and authorize every REST and MCP operation. Apply least privilege, rate limits, input validation, and audit logging.
- Verify each provider webhook using its documented mechanism before changing payment state.
- Do not expose raw credentials or unrestricted provider responses through MCP.
  MCP
- Tool handlers must call existing application services; never access provider APIs or SQL directly from tool definitions.
- Give tools clear descriptions, strict input schemas, bounded result sizes, and pagination for lists.
- Enforce tenant scope and per-tool permissions server-side on every call.
- Start with read-only tools. Any financial write action requires explicit authorization, audit logging, idempotency, and a confirmation policy.
- Do not claim a payment succeeded based only on a tool request or redirect response; rely on verified provider state/events.
  Engineering Workflow
- Before coding, inspect existing structure and follow established project conventions.
- Make the smallest coherent change; do not refactor unrelated code.
- Add tests for business rules, tenant isolation, idempotency, provider mapping, webhook duplicates/order, and failure paths.
- Run relevant tests, formatting, linting, and migrations checks; report commands and results accurately.
- Never claim tests passed if they were not run. Document assumptions or unsupported provider capabilities.
