package wx

import (
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
)

var (
	wechatClient *wechat.Wechat
)

func init() {

	wechatClient = wechat.NewWechat()

	memCache := cache.NewMemory()
	wechatClient.SetCache(memCache)

}
