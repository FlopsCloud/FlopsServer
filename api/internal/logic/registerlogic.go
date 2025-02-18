package logic

import (
	"context"
	"fca/common/cryptx"
	"fca/common/response"
	"regexp"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp response.Response) {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	if !emailRegex.MatchString(req.Email) {
		return response.Fail(response.InvalidRequestParamCode, "邮箱格式不正确")
	}

	user, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, req.Email)
	if err != nil && err != model.ErrNotFound {
		return response.Error(err.Error())
	}
	if user != nil {
		return response.Fail(response.InvalidRequestParamCode, "邮箱已注册")
	}

	var userId int64
	var respData *types.RegisterResponseData

	err = sqlx.NewMysql(l.svcCtx.Config.MySQL.DataSource).Transact(func(session sqlx.Session) error {
		// Create user with transaction
		userModel := l.svcCtx.UserModel.WithSession(session)
		res, err := userModel.Insert(l.ctx, &model.Users{
			Username:     req.Username,
			Email:        req.Email,
			Phone:        req.Phone,
			Nickname:     req.Nickname,
			IsMaster:     1,
			ShareBalance: 1,
			PasswordHash: cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, req.Password),
			Balance:      0,
		})

		if err != nil {
			return err
		}

		userId, err = res.LastInsertId()
		if err != nil {
			return err
		}

		// Link user to org with transaction
		orgsUsersModel := l.svcCtx.OrgsUsersModel.WithSession(session)
		_, err = orgsUsersModel.Insert(l.ctx, &model.OrgsUsers{
			OrgId:  1,
			UserId: uint64(userId),
			Role:   "member", // Default role
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return response.Error(err.Error())
	}

	// Create initial balance with transaction
	balanceLogic := NewUserBalancesLogic(l.ctx, l.svcCtx)
	err, _ = balanceLogic.ManualAdjust(uint64(userId), 9900, "CNY", "注册用户，赠送价值99元体验额度")

	respData = &types.RegisterResponseData{
		UserId:   uint64(userId),
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Nickname: req.Nickname,
	}

	logx.Info(respData, req)

	// Get organization info

	org, err := l.svcCtx.OrgsUsersModel.FindAllByUserId(l.ctx, uint64(userId))
	if err != nil {
		return response.Error(err.Error())
	}

	var orgUser []types.Org
	for _, aorg := range *org {
		orgUser = append(orgUser, types.Org{
			OrgId:   aorg.OrgId,
			OrgName: aorg.OrgName,
			Role:    aorg.Role,
		})
	}

	respData.Orgs = orgUser
	logx.Info(respData, req)
	return response.OK(respData)
}
