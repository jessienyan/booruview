package api

import (
	"context"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/valkey-io/valkey-go"

	"golang.org/x/time/rate"
)

const (
	tokenRefillRate = 3 // per second
	maxTokens       = 20
	baseBanTime     = 5 * time.Minute
)

var (
	clientLimits      = make(map[string]*rate.Limiter)
	clientLimitsMutex sync.Mutex
)

func clientBanKey(ip string) string {
	return "clientban:" + ip
}

type clientBan struct {
	banCount    int
	bannedUntil time.Time
}

func (cb clientBan) banned() bool {
	return cb.bannedUntil.Before(time.Now())
}

func (cb clientBan) banTime() time.Duration {
	exp := int(math.Pow(2, float64(cb.banCount)))
	return baseBanTime * time.Duration(exp)
}

func getClientBan(ip string) (*clientBan, error) {
	vc := Valkey()
	resp, err := vc.Do(context.Background(),
		vc.B().
			Hmget().
			Key(clientBanKey(ip)).
			Field().
			Build()).ToMap()

	if err != nil {
		if valkey.IsValkeyNil(err) {
			return nil, nil
		}

		return nil, err
	}

	var m valkey.ValkeyMessage
	var banCount int
	var bannedUntil time.Time

	m = resp["banCount"]
	if val, err := m.AsInt64(); err != nil {
		return nil, err
	} else {
		banCount = int(val)
	}

	m = resp["bannedUntil"]
	if val, err := m.ToString(); err != nil {
		return nil, err
	} else {
		bannedUntil, err = time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, err
		}
	}

	cb := &clientBan{
		banCount:    banCount,
		bannedUntil: bannedUntil,
	}

	return cb, nil
}

func ban(ip string, cb *clientBan) error {
	if cb == nil {
		cb = &clientBan{
			banCount: 1,
		}
	} else {
		cb.banCount++
	}

	cb.bannedUntil = time.Now().Add(cb.banTime())

	vc := Valkey()
	err := vc.Do(context.Background(),
		vc.B().
			Hmset().
			Key(clientBanKey(ip)).
			FieldValue().
			FieldValue("banCount", strconv.Itoa(cb.banCount)).
			FieldValue("bannedUntil", cb.bannedUntil.Format(time.RFC3339)).
			Build()).Error()
	// TODO: set expiration

	return err
}

func getLimiter(ip string) *rate.Limiter {
	clientLimitsMutex.Lock()
	defer clientLimitsMutex.Unlock()

	lim, exists := clientLimits[ip]
	if !exists {
		lim = rate.NewLimiter(tokenRefillRate, maxTokens)
		clientLimits[ip] = lim
	}

	return lim
}

// IsRateLimited checks if a client is currently banned, and if not, whether
// they are able to make a request worth `cost` tokens.
//
// Rate limiting is handled using the token bucket algorithm. Clients have a bucket
// of tokens that slowly refill. Requests remove tokens from the bucket equal to their cost.
// If a request costs more than what's available, the client is rate limited.
func IsRateLimited(ip string, cost int) (banned bool, err error) {
	cb, err := getClientBan(ip)
	if err != nil {
		return
	}

	if banned = cb.banned(); banned {
		return
	}

	lim := getLimiter(ip)
	if !lim.AllowN(time.Now(), cost) {
		banned = true
		err = ban(ip, cb)
		return
	}

	return
}
