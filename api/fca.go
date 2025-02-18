package main

import (
	"embed"
	"encoding/json"
	"fca/api/internal/config"
	"fca/api/internal/handler"
	"fca/api/internal/logic"
	"fca/api/internal/logic/page"
	"fca/api/internal/middleware"
	"fca/api/internal/svc"
	"fca/common/response"
	"flag"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	// "github.com/zs5460/art"
)

//go:embed etc/*
var etcFs embed.FS

var TerminalMode bool

var c config.Config

func main() {
	flag.Parse()

	TerminalMode = *flag.Bool("t", false, "Run in terminal logging mode")

	apifile, _ := etcFs.ReadFile("etc/fca.api")
	page.API_FILE = string(apifile)

	file, _ := etcFs.ReadFile("etc/fca-api.yaml")
	println(string(file))

	conf.LoadFromYamlBytes(file, &c)

	fmt.Println(c) //TODO: not working

	// Configure logx based on the -t flag

	logic.TerminalMode = TerminalMode
	if TerminalMode {
		// logic.TerminalMode = TerminalMode

		logx.DisableStat()
		logx.SetUp(logx.LogConf{
			Mode:     "console",
			Level:    "info",
			Encoding: "plain",
			Path:     "", // No file path for console logging
		})

		//fmt.Println(art.String("Running in terminal logging mode"))
		fmt.Println("########### Running in terminal logging mode ###########")

	}

	server := rest.MustNewServer(c.RestConf,
		rest.WithUnauthorizedCallback(func(w http.ResponseWriter, r *http.Request, err error) {
			res := response.Fail(response.UnauthorizedCode, "未登录")
			bts, _ := json.Marshal(res)
			w.Write(bts)
		}),
		rest.WithCors(),
	)
	defer server.Stop()

	server.Use(DemoInterceptor(c.Auth.AccessSecret))

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 或使用内置限流器

	periodLimit := middleware.NewPeriodLimit(ctx.RedisClient, 100, 60)
	server.Use(periodLimit.Handle)

	// Add CORS handler
	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		return http.StatusOK, err.Error()
	})

	fmt.Printf("Starting fca api server at %s:%d...\n", c.Host, c.Port)

	logic.InitCaptchaStore(c.CacheRedis[0].Host, c.CacheRedis[0].Pass)

	server.Start()
}

func DemoInterceptor(accessSecret string) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// authHeader := r.Header.Get("Authorization")
			// if authHeader == "" {
			// 	logx.Error("Missing Authorization header")
			// 	http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			// 	return
			// }

			// parts := strings.SplitN(authHeader, " ", 2)
			// if !(len(parts) == 2 && parts[0] == "Bearer") {
			// 	logx.Error("Invalid Authorization header format")
			// 	http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			// 	return
			// }

			// token := parts[1]
			// claims, err := jwtx.ParseToken(token, accessSecret)
			// if err != nil {
			// 	logx.Errorf("Invalid token: %v", err)
			// 	http.Error(w, "Invalid token", http.StatusUnauthorized)
			// 	return
			// }

			// ctx := context.WithValue(r.Context(), "claims", claims)
			// r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}
	}
}
