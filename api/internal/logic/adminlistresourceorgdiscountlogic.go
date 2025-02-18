package logic

import (
	"context"
	"fmt"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/api/internal/utils"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListResourceOrgDiscountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取组织资源折扣
func NewAdminListResourceOrgDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListResourceOrgDiscountLogic {
	return &AdminListResourceOrgDiscountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListResourceOrgDiscountLogic) AdminListResourceOrgDiscount(req *types.AdminListResourceOrgDiscountRequest) (resp *types.AdminListResourceOrgDiscountResponse, err error) {

	sql := `SELECT 
				ro.id,
				ro.org_id,
				ro.resource_id,
				ro.discount_id,
				UNIX_TIMESTAMP(ro.created_at) as created_at,  
				o.org_name,
				r.resource_name, 
				r.resource_type, 
				r.unit_hourly_price,
				r.unit_daily_price,
				r.unit_monthly_price,
				r.unit_yearly_price,
				r.total,
				r.remains,
				IFNULL( d.memo,"") as memo,
				IFNULL(d.min_discount,0) min_discount,
				IFNULL(d.hourly_discount,0) hourly_discount,
				IFNULL(d.daily_discount,0) daily_discount,
				IFNULL(d.monthly_discount,0) monthly_discount,
				IFNULL(d.yearly_discount,0) yearly_discount 
			
			FROM 
				resource_orgs ro
			LEFT JOIN resources r ON ro.resource_id = r.resource_id
			LEFT JOIN organizations o ON ro.org_id = o.org_id
			LEFT JOIN resource_discounts d ON ro.discount_id = d.discount_id and ro.org_id= d.org_id`

	var resourceOrgDiscounts []types.ResourceOrgDiscount
	var resourceOrgDiscount []model.ResourceOrgDiscount
	err = l.svcCtx.Mysql.QueryRowsCtx(l.ctx, &resourceOrgDiscount, sql)
	if err != nil {
		return nil, err
	}
	fmt.Println(resourceOrgDiscount)

	for _, v := range resourceOrgDiscount {

		var des types.ResourceOrgDiscount
		utils.CopyStruct(&v, &des, false)
		resourceOrgDiscounts = append(resourceOrgDiscounts, des)
	}

	return &types.AdminListResourceOrgDiscountResponse{
		Response: types.Response{Code: response.SuccessCode},
		Data: types.ListResourceOrgDiscountResponseData{
			ResourceOrgDiscounts: resourceOrgDiscounts,
			Total:                uint64(len(resourceOrgDiscounts)),
		},
	}, nil

}
func (l *AdminListResourceOrgDiscountLogic) OldAdminListResourceOrgDiscount(req *types.AdminListResourceOrgDiscountRequest) (resp *types.AdminListResourceOrgDiscountResponse, err error) {
	// todo: add your logic here and delete this line

	sql := `
SELECT 
    ro.id,
    ro.org_id,
    ro.resource_id,
    ro.discount_id,
    UNIX_TIMESTAMP(ro.created_at) as created_at,
    o.org_name,
    r.resource_name, 
    r.resource_type, 
    r.unit_hourly_price,
    r.unit_daily_price,
    r.unit_monthly_price,
    r.unit_yearly_price,
    r.total,
    r.remains,
    ifnull(d.memo,'') as memo,
    UNIX_TIMESTAMP(ifnull(d.startdate,current_timestamp())) as startdate,
    UNIX_TIMESTAMP(ifnull(d.enddate,current_timestamp())) as enddate,
    ifnull(d.discount,0) as discount  
FROM 
    resource_orgs ro
LEFT JOIN resources r ON ro.resource_id = r.resource_id
LEFT JOIN organizations o ON ro.org_id = o.org_id
LEFT JOIN discounts d ON ro.discount_id = d.discount_id and d.is_deleted=0
`
	var resourceOrgDiscounts []types.ResourceOrgDiscount
	var resourceOrgDiscount []model.ResourceOrgDiscount
	err = l.svcCtx.Mysql.QueryRowsCtx(l.ctx, &resourceOrgDiscount, sql)
	if err != nil {
		return nil, err
	}
	fmt.Println(resourceOrgDiscount)

	for _, v := range resourceOrgDiscount {
		var des types.ResourceOrgDiscount
		utils.CopyStruct(&v, &des, false)
		resourceOrgDiscounts = append(resourceOrgDiscounts, des)
	}

	return &types.AdminListResourceOrgDiscountResponse{
		Response: types.Response{Code: response.SuccessCode},
		Data: types.ListResourceOrgDiscountResponseData{
			ResourceOrgDiscounts: resourceOrgDiscounts,
			Total:                uint64(len(resourceOrgDiscounts)),
		},
	}, nil
}
