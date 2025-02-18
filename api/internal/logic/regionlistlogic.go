package logic

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegionListLogic {
	return &RegionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegionListLogic) RegionList() (resp *types.RegionsListResp, err error) {
	// Get all regions from database
	regions, err := l.svcCtx.RegionsModel.FindAll(l.ctx)
	if err != nil {
		return &types.RegionsListResp{
			Response: types.Response{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	// Convert model.Regions to types.Region
	var regionList []types.Region
	for _, r := range regions {
		regionList = append(regionList, types.Region{
			RegionId:   r.RegionId,
			RegionName: r.RegionName,
			RegionCode: r.RegionCode,
			CreatedAt:  uint64(r.CreatedAt.Unix()),
			UpdatedAt:  uint64(r.UpdatedAt.Unix()),
		})
	}

	return &types.RegionsListResp{
		Response: types.Response{
			Code:    200,
			Message: "success",
		},
		Data: types.RegionsListRespData{
			Regions: regionList,
			Total:   uint64(len(regionList)),
		},
	}, nil
}
