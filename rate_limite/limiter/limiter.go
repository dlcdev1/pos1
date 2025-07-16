package limiter

import (
	"strconv"
	"strings"
	"time"
)

type TokenLimits map[string]int

type Limiter struct {
	storage       Storage
	blockDuration time.Duration

	ipLimit         int
	defaultTokLimit int
	tokenLimits     TokenLimits
}

func NewLimiter(storage Storage, ipLimit, defaultTokenLimit int, blockDuration time.Duration, tokenLimitsStr string) *Limiter {
	tokenLimits := make(TokenLimits)
	if tokenLimitsStr != "" {
		for _, pair := range strings.Split(tokenLimitsStr, ",") {
			parts := strings.Split(pair, ":")
			if len(parts) == 2 {
				lim, err := strconv.Atoi(parts[1])
				if err == nil {
					tokenLimits[parts[0]] = lim
				}
			}
		}
	}

	return &Limiter{
		storage:         storage,
		blockDuration:   blockDuration,
		ipLimit:         ipLimit,
		defaultTokLimit: defaultTokenLimit,
		tokenLimits:     tokenLimits,
	}
}

func (l *Limiter) getLimit(token string, ip string) int {
	if token != "" {
		if val, ok := l.tokenLimits[token]; ok {
			return val
		}
		return l.defaultTokLimit
	}
	return l.ipLimit
}

func (l *Limiter) Key(token, ip string) string {
	if token != "" {
		return "rl:token:" + token
	}
	return "rl:ip:" + ip
}

func (l *Limiter) Allow(token, ip string) (bool, error) {
	limit := l.getLimit(token, ip)
	key := l.Key(token, ip)

	count, err := l.storage.Increment(key, l.blockDuration)
	if err != nil {
		return false, err
	}

	if count > int64(limit) {
		return false, nil
	}
	return true, nil
}
