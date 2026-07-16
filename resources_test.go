package zazu_test

// Mirror of zazu-ruby's spec/zazu/resources/*_spec.rb — same cassettes,
// same assertions, per the cross-language SDK contract.

import (
	"context"
	"net/http/httptest"
	"testing"

	zazu "github.com/getzazu/zazu-go"
)

func replayClient(t *testing.T, server *httptest.Server) *zazu.Client {
	t.Helper()
	client, err := zazu.New(
		zazu.WithAPIKey("test-api-key-for-replay"),
		zazu.WithBaseURL(server.URL),
	)
	if err != nil {
		t.Fatalf("build client: %v", err)
	}
	return client
}

func TestEntityGet(t *testing.T) {
	server := startReplayServer(t, "entity/get")
	client := replayClient(t, server)

	resp, err := client.Entity.Get(context.Background())
	if err != nil {
		t.Fatalf("entity get: %v", err)
	}
	if _, ok := resp.Body["id"].(string); !ok {
		t.Fatalf("expected string id, got %#v", resp.Body["id"])
	}
}

func TestAccounts(t *testing.T) {
	server := startReplayServer(t, "accounts/list", "accounts/get", "accounts/list_transactions", "accounts/get_transaction")
	client := replayClient(t, server)
	ctx := context.Background()

	page, err := client.Accounts.List(ctx, zazu.AccountListParams{})
	if err != nil {
		t.Fatalf("accounts list: %v", err)
	}
	if page.Data == nil {
		t.Fatal("expected data rows")
	}

	accountID := fixtureID(t, "ZAZU_FIXTURE_ACCOUNT_ID")
	if _, err := client.Accounts.Get(ctx, accountID); err != nil {
		t.Fatalf("accounts get: %v", err)
	}

	if _, err := client.Accounts.ListTransactions(ctx, accountID, zazu.TransactionListParams{}); err != nil {
		t.Fatalf("list transactions: %v", err)
	}

	txID := fixtureID(t, "ZAZU_FIXTURE_TRANSACTION_ID")
	if _, err := client.Accounts.GetTransaction(ctx, accountID, txID); err != nil {
		t.Fatalf("get transaction: %v", err)
	}
}

func TestCustomers(t *testing.T) {
	server := startReplayServer(t, "customers/list", "customers/get", "customers/create", "customers/update", "customers/delete")
	client := replayClient(t, server)
	ctx := context.Background()

	if _, err := client.Customers.List(ctx, zazu.CustomerListParams{}); err != nil {
		t.Fatalf("customers list: %v", err)
	}

	customerID := fixtureID(t, "ZAZU_FIXTURE_CUSTOMER_ID")
	resp, err := client.Customers.Get(ctx, customerID)
	if err != nil {
		t.Fatalf("customers get: %v", err)
	}
	if _, ok := resp.Body["id"].(string); !ok {
		t.Fatalf("expected string id, got %#v", resp.Body["id"])
	}
}

func TestInvoices(t *testing.T) {
	server := startReplayServer(t, "invoices/list", "invoices/get")
	client := replayClient(t, server)
	ctx := context.Background()

	page, err := client.Invoices.List(ctx, zazu.InvoiceListParams{})
	if err != nil {
		t.Fatalf("invoices list: %v", err)
	}
	if page.Data == nil {
		t.Fatal("expected data rows")
	}

	invoiceID := fixtureID(t, "ZAZU_FIXTURE_INVOICE_ID")
	if _, err := client.Invoices.Get(ctx, invoiceID); err != nil {
		t.Fatalf("invoices get: %v", err)
	}
}

func TestPaymentLinks(t *testing.T) {
	server := startReplayServer(t, "payment_links/list", "payment_links/get", "payment_links/create", "payment_links/cancel")
	client := replayClient(t, server)
	ctx := context.Background()

	if _, err := client.PaymentLinks.List(ctx, zazu.PaymentLinkListParams{}); err != nil {
		t.Fatalf("payment links list: %v", err)
	}

	resp, err := client.PaymentLinks.Create(ctx, zazu.Attributes{
		"account_id":  fixtureID(t, "ZAZU_FIXTURE_ACCOUNT_ID"),
		"amount":      "100.00",
		"title":       "SDK fixture",
		"description": "Created by zazu-ruby fixture spec",
		"link_type":   "single",
	})
	if err != nil {
		t.Fatalf("payment links create: %v", err)
	}
	if resp.Status != 201 {
		t.Fatalf("expected 201, got %d", resp.Status)
	}

	if _, err := client.PaymentLinks.Cancel(ctx, fixtureID(t, "ZAZU_FIXTURE_CANCELLABLE_PAYMENT_LINK_ID")); err != nil {
		t.Fatalf("payment links cancel: %v", err)
	}
}

func TestCheckoutSessions(t *testing.T) {
	server := startReplayServer(t, "checkout_sessions/get")
	client := replayClient(t, server)

	resp, err := client.CheckoutSessions.Get(context.Background(), fixtureID(t, "ZAZU_FIXTURE_CHECKOUT_SESSION_ID"))
	if err != nil {
		t.Fatalf("checkout sessions get: %v", err)
	}
	if _, ok := resp.Body["id"].(string); !ok {
		t.Fatalf("expected string id, got %#v", resp.Body["id"])
	}
}

func TestWebhookEndpoints(t *testing.T) {
	server := startReplayServer(t, "webhook_endpoints/list", "webhook_endpoints/get")
	client := replayClient(t, server)
	ctx := context.Background()

	if _, err := client.WebhookEndpoints.List(ctx, zazu.ListParams{}); err != nil {
		t.Fatalf("webhook endpoints list: %v", err)
	}
	if _, err := client.WebhookEndpoints.Get(ctx, fixtureID(t, "ZAZU_FIXTURE_WEBHOOK_ID")); err != nil {
		t.Fatalf("webhook endpoints get: %v", err)
	}
}

func TestTransferDrafts(t *testing.T) {
	server := startReplayServer(t, "transfer_drafts/create", "transfer_drafts/get")
	client := replayClient(t, server)
	ctx := context.Background()

	resp, err := client.TransferDrafts.Create(ctx, zazu.Attributes{
		"account_id":        fixtureID(t, "ZAZU_FIXTURE_ACCOUNT_ID"),
		"beneficiary_id":    fixtureID(t, "ZAZU_FIXTURE_BENEFICIARY_ID"),
		"amount":            "150.00",
		"payment_reference": "SDK fixture",
	})
	if err != nil {
		t.Fatalf("transfer drafts create: %v", err)
	}
	if resp.Status != 201 {
		t.Fatalf("expected 201, got %d", resp.Status)
	}
	if status, _ := resp.Body["status"].(string); status != "requested" {
		t.Fatalf("expected requested status (awaiting in-app approval), got %q", status)
	}
	if resp.Body["transfer"] != nil {
		t.Fatalf("expected nil transfer before approval, got %#v", resp.Body["transfer"])
	}

	got, err := client.TransferDrafts.Get(ctx, fixtureID(t, "ZAZU_FIXTURE_TRANSFER_DRAFT_ID"))
	if err != nil {
		t.Fatalf("transfer drafts get: %v", err)
	}
	if _, ok := got.Body["status"].(string); !ok {
		t.Fatalf("expected string status, got %#v", got.Body["status"])
	}
}

func TestBeneficiaries(t *testing.T) {
	server := startReplayServer(t, "beneficiaries/list", "beneficiaries/get")
	client := replayClient(t, server)
	ctx := context.Background()

	page, err := client.Beneficiaries.List(ctx, zazu.ListParams{})
	if err != nil {
		t.Fatalf("beneficiaries list: %v", err)
	}
	if len(page.Data) == 0 {
		t.Fatal("expected at least one beneficiary")
	}
	if _, ok := page.Data[0]["external_accounts"].([]any); !ok {
		t.Fatalf("expected embedded external_accounts, got %#v", page.Data[0]["external_accounts"])
	}

	resp, err := client.Beneficiaries.Get(ctx, fixtureID(t, "ZAZU_FIXTURE_BENEFICIARY_ID"))
	if err != nil {
		t.Fatalf("beneficiaries get: %v", err)
	}
	if _, ok := resp.Body["id"].(string); !ok {
		t.Fatalf("expected string id, got %#v", resp.Body["id"])
	}
}
