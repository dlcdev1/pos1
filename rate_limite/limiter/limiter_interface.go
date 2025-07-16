package limiter

type LimiterInterface interface {
	Allow(token, ip string) (bool, error)
}
