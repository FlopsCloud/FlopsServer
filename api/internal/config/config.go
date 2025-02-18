package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	MySQL struct {
		DataSource string
	}

	CacheRedis cache.CacheConf

	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	Salt string

	WeixinPay struct {
		AppId          string
		MerchantId     string
		MerchantKey    string
		SerialNumber   string
		NotifyUrl      string
		PrivateKeyPath string
	}

	Mail struct {
		MailAccount string
		MailHost    string
		MailPass    string
		MailPort    int
	}
	Path string

	Google struct {
		Key    string
		Client string
	}
}
