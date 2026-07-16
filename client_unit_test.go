package zazu_test

import (
	"context"
	"testing"

	zazu "github.com/getzazu/zazu-go"
)

func TestNewRequiresAPIKey(t *testing.T) {
	t.Setenv("ZAZU_API_KEY", "")

	_, err := zazu.New()
	if err == nil {
		t.Fatal("expected configuration error without an API key")
	}
	if _, ok := err.(*zazu.ConfigurationError); !ok {
		t.Fatalf("expected *zazu.ConfigurationError, got %T", err)
	}
}

func TestListLimitValidation(t *testing.T) {
	client, err := zazu.New(zazu.WithAPIKey("test"), zazu.WithBaseURL("http://127.0.0.1:1"))
	if err != nil {
		t.Fatalf("build client: %v", err)
	}

	_, err = client.Beneficiaries.List(context.Background(), zazu.ListParams{Limit: zazu.MaxPerPage + 1})
	if err == nil {
		t.Fatal("expected limit validation error")
	}
}
