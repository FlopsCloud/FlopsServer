package logic

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"path"
	"strconv"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetObjLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	writer io.Writer
}

func NewGetObjLogic(ctx context.Context, svcCtx *svc.ServiceContext, writer io.Writer) *GetObjLogic {
	return &GetObjLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		writer: writer,
	}

}

func (l *GetObjLogic) GetObj(req *types.ObjGetRequest) error {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	logx.Infof("GetObj %s", req.File)
	body, err := ioutil.ReadFile(path.Join(l.svcCtx.Config.BucketPath, strconv.FormatUint(uint64(uid), 10), strconv.FormatUint(req.BucketId, 10), req.File))
	if err != nil {
		return err
	}

	n, err := l.writer.Write(body)
	if err != nil {
		return err
	}

	if n < len(body) {
		return io.ErrClosedPipe
	}

	return nil
}
