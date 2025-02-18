package svc

import (
	"fca/api/internal"
	"fca/api/internal/config"
	"fca/model"
	"net/http"

	"github.com/zeromicro/go-zero/core/stores/redis"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config                   config.Config
	UserModel                model.UsersModel
	BaseDataModel            model.BaseDataModel
	ServerModel              model.ServersModel
	InstanceModel            model.InstancesModel
	InstancePortMappingModel model.InstancePortMappingModel
	RechargeOrderModel       model.RechargeOrdersModel
	TransactionRecordsModel  model.TransactionRecordsModel
	VerificationCodesModel   model.VerificationCodesModel
	OrderRecordsModel        model.OrderRecordsModel
	OrganizationModel        model.OrganizationsModel
	BalancesModel            model.BalancesModel
	OrgsUsersModel           model.OrgsUsersModel
	InvitationModel          model.InvitationModel
	ApplyJoinModel           model.ApplyJoinModel
	ServerTagsModel          model.ServerTagsModel
	ServerOrgsModel          model.ServerOrgsModel
	TagsModel                model.TagsModel
	RolesModel               model.RolesModel
	PermissionsModel         model.PermissionsModel
	RolePermissionsModel     model.RolePermissionsModel
	UserRolesModel           model.UserRolesModel
	DailyUsageModel          model.DailyUsageModel
	MinuteUsageModel         model.MinuteUsageModel
	ResourcesModel           model.ResourcesModel
	ImagesModel              model.ImagesModel
	ResourceOrgsModel        model.ResourceOrgsModel
	DiscountsModel           model.DiscountsModel
	ServerDiscountsModel     model.ServerDiscountsModel
	BucketsModel             model.BucketsModel
	ObjectsModel             model.ObjectsModel
	AuthInterceptor          rest.Middleware
	RedisClient              *redis.Redis
	W                        *internal.MyWords
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySQL.DataSource)

	config := redis.RedisConf{
		Host: c.CacheRedis[0].Host,
		Type: c.CacheRedis[0].Type,
		Pass: c.CacheRedis[0].Pass,
	}

	redisClient := redis.MustNewRedis(config)

	return &ServiceContext{
		Config:                   c,
		UserModel:                model.NewUsersModel(conn),
		BaseDataModel:            model.NewBaseDataModel(conn),
		ServerModel:              model.NewServersModel(conn),
		InstanceModel:            model.NewInstancesModel(conn),
		InstancePortMappingModel: model.NewInstancePortMappingModel(conn),
		RechargeOrderModel:       model.NewRechargeOrdersModel(conn),
		TransactionRecordsModel:  model.NewTransactionRecordsModel(conn),
		VerificationCodesModel:   model.NewVerificationCodesModel(conn),
		OrderRecordsModel:        model.NewOrderRecordsModel(conn),
		OrganizationModel:        model.NewOrganizationsModel(conn),
		BalancesModel:            model.NewBalancesModel(conn),
		OrgsUsersModel:           model.NewOrgsUsersModel(conn),
		InvitationModel:          model.NewInvitationModel(conn),
		ApplyJoinModel:           model.NewApplyJoinModel(conn),
		ServerTagsModel:          model.NewServerTagsModel(conn),
		ServerOrgsModel:          model.NewServerOrgsModel(conn),
		TagsModel:                model.NewTagsModel(conn),
		RolesModel:               model.NewRolesModel(conn),
		PermissionsModel:         model.NewPermissionsModel(conn),
		RolePermissionsModel:     model.NewRolePermissionsModel(conn),
		UserRolesModel:           model.NewUserRolesModel(conn),
		DailyUsageModel:          model.NewDailyUsageModel(conn),
		MinuteUsageModel:         model.NewMinuteUsageModel(conn),
		ResourcesModel:           model.NewResourcesModel(conn),
		ImagesModel:              model.NewImagesModel(conn),
		ResourceOrgsModel:        model.NewResourceOrgsModel(conn),
		DiscountsModel:           model.NewDiscountsModel(conn),
		ServerDiscountsModel:     model.NewServerDiscountsModel(conn),
		BucketsModel:             model.NewBucketsModel(conn),
		ObjectsModel:             model.NewObjectsModel(conn),
		AuthInterceptor:          AuthInterceptor(c.Auth.AccessSecret),
		RedisClient:              redisClient,
		W:                        internal.NewMyWords(),
	}
}

func AuthInterceptor(accessSecret string) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Add your authentication logic here
			// For example, you can check for tokens, validate user sessions, etc.
			// If authentication fails, you can return an error response

			// For now, we'll just call the next handler
			next(w, r)
		}
	}
}
