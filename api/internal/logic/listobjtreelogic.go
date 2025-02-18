package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"fca/api/internal/svc"
	"fca/api/internal/types"
	"fca/common/response"
	"fca/model"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"gopkg.in/yaml.v2"
)

type ListObjTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListObjTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListObjTreeLogic {
	return &ListObjTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// buildTree constructs a tree structure from a list of tree nodes
func (l *ListObjTreeLogic) buildTree(nodes *[]model.ObjectsTree) *types.TreeNode {
	// Create a map of nodes by ID for easy lookup
	nodeMap := make(map[uint64]*types.TreeNode)

	// First pass: create all nodes
	for _, node := range *nodes {
		nodeMap[node.TreeId] = &types.TreeNode{
			TreeId:       node.TreeId,
			TreeName:     node.TreeName,
			FullTreeName: node.FullTreeName,
			CreatedAt:    uint64(node.CreatedAt.Unix()),
			ParantTreeId: node.ParantTreeId,
			Children:     make([]*types.TreeNode, 0),
		}
	}

	// Second pass: build the tree structure
	var root *types.TreeNode

	root = &types.TreeNode{
		TreeName: "/",
		Children: make([]*types.TreeNode, 0),
	}
	nodeMap[0] = root

	for _, node := range *nodes {
		treeNode := nodeMap[node.TreeId]
		if node.TreeId == 0 {
			root = treeNode
		} else if parent := nodeMap[node.ParantTreeId]; parent != nil {
			parent.Children = append(parent.Children, treeNode)
		}
	}

	return root
}

func (l *ListObjTreeLogic) ListObjTree(req *types.ListObjTreeRequest) (resp *types.ListObjTreeResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	// Check bucket ownership
	bucket, err := l.svcCtx.BucketsModel.FindOne(l.ctx, req.BucketId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.ListObjTreeResponse{
				Response: types.Response{
					Code:    response.ServerErrorCode,
					Message: "Bucket not found",
				},
			}, nil
		}
		return &types.ListObjTreeResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to check bucket",
				Info:    err.Error(),
			},
		}, nil
	}
	if bucket.UserId != uint64(uid) {
		return &types.ListObjTreeResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Permission denied",
				Info:    fmt.Sprintf("bucket.UserId: %d, uid: %d", bucket.UserId, uid),
			},
		}, nil
	}

	// Get all objects for the bucket
	objectsTrees, err := l.svcCtx.ObjectsTreeModel.FindAllByBucketId(l.ctx, req.BucketId)
	if err != nil {
		return &types.ListObjTreeResponse{
			Response: types.Response{
				Code:    response.ServerErrorCode,
				Message: "Failed to get objects tree",
				Info:    err.Error(),
			},
		}, nil
	}
	yaml, _ := yaml.Marshal(objectsTrees)
	logc.Infof(l.ctx, "objectsTrees: %+v", string(yaml))
	// Build the tree structure
	tree := l.buildTree(objectsTrees)

	return &types.ListObjTreeResponse{
		Response: types.Response{
			Code:    response.SuccessCode,
			Message: "Success",
		},
		Data: types.ListObjTreeResponseData{
			Tree: *tree,
		},
	}, nil
}
