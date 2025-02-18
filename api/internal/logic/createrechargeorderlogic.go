package logic

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fca/api/internal/config"
	"fca/api/internal/utils/random"
	"fca/common/response"
	"fca/model"
	"time"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRechargeOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRechargeOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRechargeOrderLogic {
	return &CreateRechargeOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRechargeOrderLogic) CreateRechargeOrder(req *types.CreateRechargeOrderRequest) response.Response {

	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	user, err := l.svcCtx.UserModel.FindOne(l.ctx, uint64(uid))
	if err != nil {
		return response.Fail(response.UserNotExistCode, err.Error())
	}
	orderNo := time.Now().Format("P20060120") + random.GetRandomNumberString(8)
	order := &model.RechargeOrders{
		Id:            0,
		UserId:        user.UserId,
		OrgId:         req.OrgId,
		TransactionId: "",
		Remark:        "",
		Status:        1, // pending
		PaidAt:        0,
		PayCodeUrl:    "",
		OrderTitle:    "FlopsCloud平台充值",
		Amount:        req.PayMoney,
		OrderNo:       orderNo,
		PayMethod:     req.PayMethod,
	}

	resp := response.Fail(response.ServerErrorCode, "该支付方式暂不支持，请选择其他支付方式")
	if req.PayMethod == 1 {
		resp = l.CreateWeixinRechargeOrder(order)
	}

	return resp
}

func (l *CreateRechargeOrderLogic) CreateWeixinRechargeOrder(order *model.RechargeOrders) response.Response {

	mchPrivateKey, err := GetPrivateKey(l.svcCtx.Config)
	if err != nil {
		// new wechat pay client err
		return response.Fail(response.ServerErrorCode, err.Error())
	}
	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchId, serialNumber, mchPrivateKey, mchKey),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		// new wechat pay client err
		return response.Fail(response.ServerErrorCode, err.Error())
	}
	ip := "127.0.0.1"
	svc := native.NativeApiService{Client: client}
	resp, _, err := svc.Prepay(ctx,
		native.PrepayRequest{
			Appid:         core.String(appId),
			Mchid:         core.String(mchId),
			Description:   core.String(order.OrderTitle),
			OutTradeNo:    core.String(order.OrderNo),
			TimeExpire:    core.Time(time.Now()),
			Attach:        core.String(""),
			NotifyUrl:     core.String(notifyUrl),
			GoodsTag:      core.String("FlopsCloudAI"),
			LimitPay:      []string{},
			SupportFapiao: core.Bool(false),
			Amount: &native.Amount{
				Currency: core.String("CNY"),
				Total:    core.Int64(order.Amount),
			},
			Detail: &native.Detail{
				CostPrice: core.Int64(order.Amount),
				GoodsDetail: []native.GoodsDetail{native.GoodsDetail{
					GoodsName:        core.String(order.OrderTitle),
					MerchantGoodsId:  core.String("AI-2024"),
					Quantity:         core.Int64(1),
					UnitPrice:        core.Int64(order.Amount),
					WechatpayGoodsId: core.String("AI202409"),
				}},
				InvoiceId: core.String("ai123"),
			},
			SettleInfo: &native.SettleInfo{
				ProfitSharing: core.Bool(false),
			},
			SceneInfo: &native.SceneInfo{
				PayerClientIp: core.String(ip),
			},
		},
	)

	if err != nil {
		// 请求微信支付失败
		return response.Fail(response.ServerErrorCode, err.Error())
	}

	order.PayCodeUrl = *resp.CodeUrl
	_, err = l.svcCtx.RechargeOrderModel.Insert(l.ctx, order)
	if err != nil {
		// 创建订单失败
		return response.Fail(response.ServerErrorCode, err.Error())
	}

	respData := types.CreateRechargeOrderData{
		OrgId:     order.OrgId,
		OrderNo:   order.OrderNo,
		PayMethod: order.PayMethod,
		Amount:    order.Amount,
		Title:     order.OrderTitle,
		Url:       order.PayCodeUrl,
	}

	return response.OK(respData)
}

var (
	privateKey   *rsa.PrivateKey
	mchId        string
	mchKey       string
	serialNumber string
	appId        string
	notifyUrl    string
)

func GetPrivateKey(conf config.Config) (*rsa.PrivateKey, error) {
	var err error
	if privateKey == nil {
		mchId = conf.WeixinPay.MerchantId
		mchKey = conf.WeixinPay.MerchantKey
		serialNumber = conf.WeixinPay.SerialNumber
		appId = conf.WeixinPay.AppId
		notifyUrl = conf.WeixinPay.NotifyUrl
		privateKey, err = utils.LoadPrivateKeyWithPath(conf.WeixinPay.PrivateKeyPath)
	}
	return privateKey, err
}
