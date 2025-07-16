package limiter_test

import (
	"github.com/alicebob/miniredis/v2"
	"testing"
	"time"

	"github.com/dlcdev1/pos1/rate_limite/limiter"
)

// MockStorage simples para contar acessos
type MockStorage struct {
	counts map[string]int64
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		counts: make(map[string]int64),
	}
}

func (m *MockStorage) Increment(key string, expire time.Duration) (int64, error) {
	m.counts[key]++
	return m.counts[key], nil
}

func (m *MockStorage) Expire(key string, expire time.Duration) error {
	return nil
}

func TestLimiterAllow(t *testing.T) {
	storage := NewMockStorage()
	lim := limiter.NewLimiter(storage, 3, 5, time.Minute, "")

	// Limite IP 3 req/s
	ipKey := "ip1"

	// 3 primeiras requisições liberadas
	for i := 0; i < 3; i++ {
		allowed, err := lim.Allow("", ipKey)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !allowed {
			t.Errorf("Request %d for IP should be allowed but blocked", i+1)
		}
	}

	// 4a requisição deve ser bloqueada
	allowed, _ := lim.Allow("", ipKey)
	if allowed {
		t.Errorf("Request 4 should be blocked but was allowed")
	}

	// Teste Token com limite customizado
	tokenLimits := "mytok:2"
	lim = limiter.NewLimiter(storage, 3, 5, time.Minute, tokenLimits)

	// 2 requisições pelo token devem passar
	for i := 0; i < 2; i++ {
		allowed, _ = lim.Allow("mytok", "")
		if !allowed {
			t.Errorf("Request %d for token should be allowed", i+1)
		}
	}
	// 3a requisição token bloqueada
	allowed, _ = lim.Allow("mytok", "")
	if allowed {
		t.Errorf("Request 3 for token should be blocked")
	}
}

func TestLimiterIntegrationRedis(t *testing.T) {
	// Inicia um servidor Redis fake em memória
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("erro ao iniciar miniredis: %v", err)
	}
	defer s.Close()

	// Usa o endereço do miniredis
	storage := limiter.NewRedisStorage(s.Addr(), "", 0)
	defer storage.Close()

	lim := limiter.NewLimiter(storage, 2, 3, time.Minute, "")

	ip := "1.2.3.4"

	for i := 0; i < 2; i++ {
		allowed, err := lim.Allow("", ip)
		if err != nil {
			t.Fatal(err)
		}
		if !allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	allowed, err := lim.Allow("", ip)
	if err != nil {
		t.Fatal(err)
	}
	if allowed {
		t.Error("3rd request should be blocked")
	}
}
