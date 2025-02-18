// middleware/ratelimit.go
package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type PeriodLimit struct {
	limiters map[string]*limit.PeriodLimit
	redis    *redis.Redis
	rate     int
	period   int
	lock     sync.Mutex
}

func NewPeriodLimit(store *redis.Redis, rate, period int) *PeriodLimit {
	return &PeriodLimit{
		limiters: make(map[string]*limit.PeriodLimit),
		redis:    store,
		rate:     rate,
		period:   period,
	}
}

func getClientIP(r *http.Request) string {
	// 尝试从 X-Real-IP 获取
	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// 尝试从 X-Forwarded-For 获取
	ip = r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 使用 RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func (l *PeriodLimit) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		if ip == "" {
			http.Error(w, "无法获取IP地址", http.StatusInternalServerError)

			return
		}

		limiter := l.getLimiter(ip)
		if limiter == nil {
			http.Error(w, "限流服务异常", http.StatusInternalServerError)

			return
		}

		ok, _ := limiter.Take(ip)

		if ok != limit.Allowed {
			http.Error(w, "请求太频繁，请稍后再试", http.StatusTooManyRequests)

			return
		}

		next(w, r)
	}
}

func (l *PeriodLimit) getLimiter(ip string) *limit.PeriodLimit {
	l.lock.Lock()
	defer l.lock.Unlock()

	if limiter, ok := l.limiters[ip]; ok {
		return limiter
	}

	// 创建新的限流器
	limiter := limit.NewPeriodLimit(
		l.period,
		l.rate,
		l.redis,
		"ratelimit:"+ip, // redis key 前缀
	)

	l.limiters[ip] = limiter
	return limiter
}
