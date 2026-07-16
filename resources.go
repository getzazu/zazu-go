package zazu

import (
	"context"
	"net/url"
)

// Attributes is a request body for create/update calls — snake_case keys,
// exactly what the API accepts (see the per-endpoint docs).
type Attributes map[string]any

func setIfPresent(v url.Values, key, value string) {
	if value != "" {
		v.Set(key, value)
	}
}

// AccountsService — accounts and their transactions.
type AccountsService struct{ client *Client }

// AccountListParams filters GET /api/accounts.
type AccountListParams struct {
	ListParams
	Status       string
	CurrencyCode string
}

// List calls GET /api/accounts.
func (s *AccountsService) List(ctx context.Context, params AccountListParams) (*Page, error) {
	base := url.Values{}
	setIfPresent(base, "status", params.Status)
	setIfPresent(base, "currency_code", params.CurrencyCode)
	return s.client.listPage(ctx, "api/accounts", base, params.ListParams)
}

// Get calls GET /api/accounts/:id.
func (s *AccountsService) Get(ctx context.Context, id string) (*Response, error) {
	return s.client.get(ctx, encodePath("api/accounts", id), nil)
}

// TransactionListParams filters GET /api/accounts/:id/transactions.
type TransactionListParams struct {
	ListParams
	Operation    string
	PostedAfter  string // ISO-8601
	PostedBefore string // ISO-8601
}

// ListTransactions calls GET /api/accounts/:account_id/transactions.
func (s *AccountsService) ListTransactions(ctx context.Context, accountID string, params TransactionListParams) (*Page, error) {
	base := url.Values{}
	setIfPresent(base, "operation", params.Operation)
	setIfPresent(base, "posted_after", params.PostedAfter)
	setIfPresent(base, "posted_before", params.PostedBefore)
	return s.client.listPage(ctx, encodePath("api/accounts", accountID, "transactions"), base, params.ListParams)
}

// GetTransaction calls GET /api/accounts/:account_id/transactions/:id.
func (s *AccountsService) GetTransaction(ctx context.Context, accountID, transactionID string) (*Response, error) {
	return s.client.get(ctx, encodePath("api/accounts", accountID, "transactions", transactionID), nil)
}

// BeneficiariesService — read-only directory of saved transfer recipients.
// Each beneficiary embeds its bank accounts; the one flagged `default` is
// used when a transfer names only the beneficiary_id. Beneficiaries are
// created and managed in the Zazu dashboard.
type BeneficiariesService struct{ client *Client }

// List calls GET /api/beneficiaries.
func (s *BeneficiariesService) List(ctx context.Context, params ListParams) (*Page, error) {
	return s.client.listPage(ctx, "api/beneficiaries", nil, params)
}

// Get calls GET /api/beneficiaries/:id.
func (s *BeneficiariesService) Get(ctx context.Context, id string) (*Response, error) {
	return s.client.get(ctx, encodePath("api/beneficiaries", id), nil)
}

// CheckoutSessionsService — one-off hosted checkout sessions. No list,
// update, or delete; sessions are created and inspected by id.
type CheckoutSessionsService struct{ client *Client }

// Create calls POST /api/checkout_sessions.
// Required attributes: account_id, amount, success_url.
func (s *CheckoutSessionsService) Create(ctx context.Context, attributes Attributes) (*Response, error) {
	return s.client.post(ctx, "api/checkout_sessions", attributes)
}

// Get calls GET /api/checkout_sessions/:id.
func (s *CheckoutSessionsService) Get(ctx context.Context, id string) (*Response, error) {
	return s.client.get(ctx, encodePath("api/checkout_sessions", id), nil)
}

// CustomersService — individuals or businesses the entity invoices.
type CustomersService struct{ client *Client }

// CustomerListParams filters GET /api/customers.
type CustomerListParams struct {
	ListParams
	Query string // matches company name, person name, email
}

// List calls GET /api/customers.
func (s *CustomersService) List(ctx context.Context, params CustomerListParams) (*Page, error) {
	base := url.Values{}
	setIfPresent(base, "q", params.Query)
	return s.client.listPage(ctx, "api/customers", base, params.ListParams)
}

// Get calls GET /api/customers/:id.
func (s *CustomersService) Get(ctx context.Context, id string) (*Response, error) {
	return s.client.get(ctx, encodePath("api/customers", id), nil)
}

// Create calls POST /api/customers.
func (s *CustomersService) Create(ctx context.Context, attributes Attributes) (*Response, error) {
	return s.client.post(ctx, "api/customers", attributes)
}

// Update calls PATCH /api/customers/:id.
func (s *CustomersService) Update(ctx context.Context, id string, attributes Attributes) (*Response, error) {
	return s.client.patch(ctx, encodePath("api/customers", id), attributes)
}

// Delete calls DELETE /api/customers/:id.
func (s *CustomersService) Delete(ctx context.Context, id string) (*Response, error) {
	return s.client.delete(ctx, encodePath("api/customers", id))
}

// EntityService — the current entity (the tenant the API key belongs to).
type EntityService struct{ client *Client }

// Get calls GET /api/entity.
func (s *EntityService) Get(ctx context.Context) (*Response, error) {
	return s.client.get(ctx, "api/entity", nil)
}

// InvoicesService — invoices and their lifecycle actions.
type InvoicesService struct{ client *Client }

// InvoiceListParams filters GET /api/invoices.
type InvoiceListParams struct {
	ListParams
	Status     string
	CustomerID string
}

// List calls GET /api/invoices.
func (s *InvoicesService) List(ctx context.Context, params InvoiceListParams) (*Page, error) {
	base := url.Values{}
	setIfPresent(base, "status", params.Status)
	setIfPresent(base, "customer_id", params.CustomerID)
	return s.client.listPage(ctx, "api/invoices", base, params.ListParams)
}

// Get calls GET /api/invoices/:id.
func (s *InvoicesService) Get(ctx context.Context, id string) (*Response, error) {
	return s.client.get(ctx, encodePath("api/invoices", id), nil)
}

// Create calls POST /api/invoices.
func (s *InvoicesService) Create(ctx context.Context, attributes Attributes) (*Response, error) {
	return s.client.post(ctx, "api/invoices", attributes)
}

// Update calls PATCH /api/invoices/:id.
func (s *InvoicesService) Update(ctx context.Context, id string, attributes Attributes) (*Response, error) {
	return s.client.patch(ctx, encodePath("api/invoices", id), attributes)
}

// Send calls POST /api/invoices/:id/send.
func (s *InvoicesService) Send(ctx context.Context, id string) (*Response, error) {
	return s.client.post(ctx, encodePath("api/invoices", id, "send"), nil)
}

// MarkAsPaid calls POST /api/invoices/:id/mark_as_paid.
func (s *InvoicesService) MarkAsPaid(ctx context.Context, id string) (*Response, error) {
	return s.client.post(ctx, encodePath("api/invoices", id, "mark_as_paid"), nil)
}

// Cancel calls POST /api/invoices/:id/cancel.
func (s *InvoicesService) Cancel(ctx context.Context, id string) (*Response, error) {
	return s.client.post(ctx, encodePath("api/invoices", id, "cancel"), nil)
}

// CreditNote calls POST /api/invoices/:id/credit_note.
func (s *InvoicesService) CreditNote(ctx context.Context, id string) (*Response, error) {
	return s.client.post(ctx, encodePath("api/invoices", id, "credit_note"), nil)
}

// Delete calls DELETE /api/invoices/:id.
func (s *InvoicesService) Delete(ctx context.Context, id string) (*Response, error) {
	return s.client.delete(ctx, encodePath("api/invoices", id))
}

// CreatePaymentLink calls POST /api/invoices/:invoice_id/payment_link.
func (s *InvoicesService) CreatePaymentLink(ctx context.Context, invoiceID, accountID string) (*Response, error) {
	return s.client.post(ctx, encodePath("api/invoices", invoiceID, "payment_link"), Attributes{"account_id": accountID})
}

// PaymentLinksService — standalone payment links (not attached to an invoice).
type PaymentLinksService struct{ client *Client }

// PaymentLinkListParams filters GET /api/payment_links.
type PaymentLinkListParams struct {
	ListParams
	Status   string
	LinkType string
}

// List calls GET /api/payment_links.
func (s *PaymentLinksService) List(ctx context.Context, params PaymentLinkListParams) (*Page, error) {
	base := url.Values{}
	setIfPresent(base, "status", params.Status)
	setIfPresent(base, "link_type", params.LinkType)
	return s.client.listPage(ctx, "api/payment_links", base, params.ListParams)
}

// Get calls GET /api/payment_links/:id.
func (s *PaymentLinksService) Get(ctx context.Context, id string) (*Response, error) {
	return s.client.get(ctx, encodePath("api/payment_links", id), nil)
}

// Create calls POST /api/payment_links.
func (s *PaymentLinksService) Create(ctx context.Context, attributes Attributes) (*Response, error) {
	return s.client.post(ctx, "api/payment_links", attributes)
}

// Cancel calls POST /api/payment_links/:id/cancel.
func (s *PaymentLinksService) Cancel(ctx context.Context, id string) (*Response, error) {
	return s.client.post(ctx, encodePath("api/payment_links", id, "cancel"), nil)
}

// TransferDraftsService — API-initiated transfers. Creating a draft routes
// it into the workspace's in-app approval flow; the API never executes a
// transfer itself. Poll Get (status: requested → processing → completed /
// failed) or subscribe to the `transfer.executed` webhook.
type TransferDraftsService struct{ client *Client }

// Create calls POST /api/transfer_drafts.
// Required: account_id, amount, and exactly one of beneficiary_id
// (external transfer) or destination_account_id (own-account move).
func (s *TransferDraftsService) Create(ctx context.Context, attributes Attributes) (*Response, error) {
	return s.client.post(ctx, "api/transfer_drafts", attributes)
}

// Get calls GET /api/transfer_drafts/:id.
func (s *TransferDraftsService) Get(ctx context.Context, id string) (*Response, error) {
	return s.client.get(ctx, encodePath("api/transfer_drafts", id), nil)
}

// WebhookEndpointsService — webhook endpoint management.
type WebhookEndpointsService struct{ client *Client }

// List calls GET /api/webhook_endpoints.
func (s *WebhookEndpointsService) List(ctx context.Context, params ListParams) (*Page, error) {
	return s.client.listPage(ctx, "api/webhook_endpoints", nil, params)
}

// Get calls GET /api/webhook_endpoints/:id.
func (s *WebhookEndpointsService) Get(ctx context.Context, id string) (*Response, error) {
	return s.client.get(ctx, encodePath("api/webhook_endpoints", id), nil)
}

// Create calls POST /api/webhook_endpoints.
func (s *WebhookEndpointsService) Create(ctx context.Context, attributes Attributes) (*Response, error) {
	return s.client.post(ctx, "api/webhook_endpoints", attributes)
}

// Update calls PATCH /api/webhook_endpoints/:id.
func (s *WebhookEndpointsService) Update(ctx context.Context, id string, attributes Attributes) (*Response, error) {
	return s.client.patch(ctx, encodePath("api/webhook_endpoints", id), attributes)
}

// Delete calls DELETE /api/webhook_endpoints/:id.
func (s *WebhookEndpointsService) Delete(ctx context.Context, id string) (*Response, error) {
	return s.client.delete(ctx, encodePath("api/webhook_endpoints", id))
}

// Test calls POST /api/webhook_endpoints/:id/test.
func (s *WebhookEndpointsService) Test(ctx context.Context, id string) (*Response, error) {
	return s.client.post(ctx, encodePath("api/webhook_endpoints", id, "test"), nil)
}

// RegenerateSecret calls POST /api/webhook_endpoints/:id/regenerate_secret.
func (s *WebhookEndpointsService) RegenerateSecret(ctx context.Context, id string) (*Response, error) {
	return s.client.post(ctx, encodePath("api/webhook_endpoints", id, "regenerate_secret"), nil)
}

// Enable calls POST /api/webhook_endpoints/:id/enable.
func (s *WebhookEndpointsService) Enable(ctx context.Context, id string) (*Response, error) {
	return s.client.post(ctx, encodePath("api/webhook_endpoints", id, "enable"), nil)
}

// Disable calls POST /api/webhook_endpoints/:id/disable.
func (s *WebhookEndpointsService) Disable(ctx context.Context, id string) (*Response, error) {
	return s.client.post(ctx, encodePath("api/webhook_endpoints", id, "disable"), nil)
}
