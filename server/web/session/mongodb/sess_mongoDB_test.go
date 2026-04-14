package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/beego/beego/v2/server/web/session"
)

func TestMongoProvider(t *testing.T) {
	// Note: This test requires a local MongoDB instance running on port 27017
	config := `mongodb://localhost:27017`
	provider := &MongoProvider{}
	
	err := provider.SessionInit(context.Background(), 3600, config)
	if err != nil {
		t.Skipf("Skipping test, MongoDB not running: %v", err)
	}

	// Test Session Creation
	store, err := provider.SessionRead(context.Background(), "test-session-id")
	if err != nil {
		t.Fatalf("Failed to read session: %v", err)
	}

	// Test Set and Get
	err = store.Set(context.Background(), "username", "beego_test")
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	val := store.Get(context.Background(), "username")
	if val != "beego_test" {
		t.Fatalf("Expected 'beego_test', got %v", val)
	}

	// Test Destroy
	err = provider.SessionDestroy(context.Background(), "test-session-id")
	if err != nil {
		t.Fatalf("Failed to destroy session: %v", err)
	}
}