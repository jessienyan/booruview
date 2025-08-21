package api

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/valkey-io/valkey-go"

	"golang.org/x/time/rate"
)

const (
	// These limits are nearly impossible for a human to hit even when trying to
	// do it intentionally. However it should block most scripts (just scrape gelbooru directly??)
	tokenRefillRate = 3 // per second
	maxTokens       = 25
	banResetsAfter  = 24 * time.Hour
)

var (
	banTimes = []time.Duration{
		// Very short initially to prevent punishing legit users
		30 * time.Second,
		60 * time.Second,
		30 * time.Minute,
		4 * time.Hour,
		24 * time.Hour,
	}

	clientLimits      = make(map[string]*rate.Limiter)
	clientLimitsMutex sync.Mutex
)

func clientBanKey(ip string) string {
	return "clientban:" + ip
}

type ClientBan struct {
	BanCount    int
	BannedUntil time.Time
}

func (cb ClientBan) Banned() bool {
	return cb.BannedUntil.After(time.Now())
}

func (cb ClientBan) BanTime() time.Duration {
	return banTimes[min(cb.BanCount-1, len(banTimes)-1)]
}

func getClientBan(ip string) (*ClientBan, error) {
	vc := Valkey()
	resp, err := vc.Do(context.Background(),
		vc.B().
			Hgetall().
			Key(clientBanKey(ip)).
			Build()).ToMap()

	if err != nil {
		return nil, err
	} else if len(resp) == 0 {
		return &ClientBan{}, nil
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

	cb := &ClientBan{
		BanCount:    banCount,
		BannedUntil: bannedUntil,
	}

	return cb, nil
}

func ban(ip string, cb *ClientBan) (*ClientBan, error) {
	if cb == nil {
		cb = &ClientBan{
			BanCount: 1,
		}
	} else {
		cb.BanCount++
	}

	cb.BannedUntil = time.Now().Add(cb.BanTime())

	vc := Valkey()
	key := clientBanKey(ip)
	err := vc.Do(
		context.Background(),
		vc.B().
			Hmset().
			Key(key).
			FieldValue().
			FieldValue("banCount", strconv.Itoa(cb.BanCount)).
			FieldValue("bannedUntil", cb.BannedUntil.Format(time.RFC3339)).
			Build(),
	).Error()
	if err != nil {
		return nil, err
	}

	// Reset the ban count after a waiting period (w.p. is relative to unban time)
	expires := cb.BannedUntil.Add(banResetsAfter).Unix()
	err = vc.Do(
		context.Background(),
		vc.B().
			Expireat().
			Key(key).
			Timestamp(expires).
			Build(),
	).Error()
	if err != nil {
		return nil, err
	}

	log.Warn().
		Str("ip", ip).
		Int("banCount", cb.BanCount).
		Time("bannedUntil", cb.BannedUntil).
		Msg("client banned (rate limited)")

	return cb, nil
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
func IsRateLimited(ip string, cost int) (cb *ClientBan, err error) {
	cb, err = getClientBan(ip)
	if err != nil {
		return
	}

	if banned := cb.Banned(); banned {
		log.Info().
			Str("ip", ip).
			Int("banCount", cb.BanCount).
			Time("bannedUntil", cb.BannedUntil).
			Msg("request blocked (rate limited)")
		return
	}

	lim := getLimiter(ip)
	if !lim.AllowN(time.Now(), cost) {
		cb, err = ban(ip, cb)
	}

	return
}
