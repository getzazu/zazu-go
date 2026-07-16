# zazu-go

Go SDK for the [Zazu](https://zazu.ma) API.

```bash
go get github.com/getzazu/zazu-go
```

```go
import zazu "github.com/getzazu/zazu-go"

client, err := zazu.New(zazu.WithAPIKey(os.Getenv("ZAZU_API_KEY")))
if err != nil { ... }

entity, err := client.Entity.Get(ctx)

page, err := client.Accounts.List(ctx, zazu.AccountListParams{})
for _, account := range page.Data {
    fmt.Println(account["id"], account["name"])
}

// Initiate a transfer — it lands in your workspace's in-app approval
// queue; the API never executes a transfer itself.
draft, err := client.TransferDrafts.Create(ctx, zazu.Attributes{
    "account_id":        accountID,
    "beneficiary_id":    beneficiaryID,
    "amount":            "150.00",
    "payment_reference": "INV-000042",
})
```

## Response shape

Response bodies are returned as-is from the API — `snake_case` keys in a
`map[string]any`, no struct mapping. The same shape ships across every Zazu
SDK (Ruby, TypeScript, Python, Go, ...) so the cassette contract is
one-to-one.

## Errors

Non-2xx responses come back as `*zazu.Error` with `Status`, `Kind`
(`authentication`, `forbidden`, `not_found`, `validation`, `rate_limit`,
`server`, `api`), the API's `Type`/`Message`/`Param`, and the `RequestID`.
Transport failures are `*zazu.ConnectionError`.

## Tests

Tests replay the canonical cassettes recorded by
[zazu-ruby](https://github.com/getzazu/zazu-ruby). The cassettes are
downloaded from the Ruby SDK's release tarball and served from an
`httptest.Server`. Same interactions, same assertions, every language.

```bash
scripts/fetch-cassettes.sh
go test ./...
```

## The SDK family

- [zazu-ruby](https://github.com/getzazu/zazu-ruby) — reference implementation (records the cassettes)
- [zazu-ts](https://github.com/getzazu/zazu-ts)
- [zazu-python](https://github.com/getzazu/zazu-python)
- [cli](https://github.com/getzazu/cli)
