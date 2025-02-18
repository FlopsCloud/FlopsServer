package logic

import (
	"context"
	"fmt"
	"os"
	"time"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type InfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoLogic {
	return &InfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoLogic) getMySQLStatus() (string, error) {
	conn := sqlx.NewMysql(l.svcCtx.Config.MySQL.DataSource)

	var version string
	err := conn.QueryRowCtx(l.ctx, &version, "SELECT VERSION()")
	if err != nil {
		return "", fmt.Errorf("error getting MySQL version: %v", err)
	}

	var connectionCount int
	err = conn.QueryRowCtx(l.ctx, &connectionCount, "SELECT COUNT(1) FROM information_schema.processlist")
	if err != nil {
		return "", fmt.Errorf("error getting MySQL connection count: %v", err)
	}

	return fmt.Sprintf(l.svcCtx.W.T("MySQL Version: %s, Connection Count: %d"), version, connectionCount), nil
}

func (l *InfoLogic) getRedisStatus() string {
	// Create Redis client using the cache configuration
	//redisClient := redis.New(l.svcCtx.Config.CacheRedis[0].Host)
	config := redis.RedisConf{
		Host: l.svcCtx.Config.CacheRedis[0].Host,
		Type: l.svcCtx.Config.CacheRedis[0].Type,
		Pass: l.svcCtx.Config.CacheRedis[0].Pass,
	}

	redisClient := redis.MustNewRedis(config)

	// Try to ping Redis to check connection and get basic info
	pinok := redisClient.Ping()
	if !pinok {
		return "Redis: Error connecting"
	} else {
		//"Redis: Connected to %s"
		return fmt.Sprintf(l.svcCtx.W.T("Redis: Connected to %s"), l.svcCtx.Config.CacheRedis[0].Host)
		// return fmt.Sprintf("Redis: success connecting ")
	}

}

func (l *InfoLogic) Info(ip string) (resp *types.Response, err error) {
	aresp := new(types.Response)
	aresp.Code = 200

	ctx := context.Background()
	logc.Info(ctx, "info message")

	// Get current execution time
	currentTime := time.Now().Format(time.RFC3339)

	// Get the current executable path
	execPath, err := os.Executable()
	if err != nil {
		aresp.Info = fmt.Sprintf("Error getting executable path: %v", err)
		aresp.Message = "Error retrieving info"
		return aresp, nil
	}

	// Get file info for the current executable
	fileInfo, err := os.Stat(execPath)
	if err != nil {
		aresp.Info = fmt.Sprintf("Error reading executable file info: %v", err)
		aresp.Message = "Error retrieving info"
		return aresp, nil
	}

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown"
	}

	// Get MySQL status
	mysqlStatus, err := l.getMySQLStatus()
	if err != nil {
		mysqlStatus = fmt.Sprintf("Error getting MySQL status: %v", err)
	}

	// Get Redis status
	redisStatus := l.getRedisStatus()

	// Store the file info, hostname, MySQL and Redis status in the info field
	aresp.Info = fmt.Sprintf("Hostname: %s, Executable: %s, Size: %d bytes, LastModified: %s, Executed at: %s, %s, %s,%s",
		hostname, execPath, fileInfo.Size(), fileInfo.ModTime().Format(time.RFC3339), currentTime, mysqlStatus, redisStatus, ip)

	aresp.Message = "OK"

	return aresp, nil
}
