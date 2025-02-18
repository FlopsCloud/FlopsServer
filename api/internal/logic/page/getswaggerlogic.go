package page

import (
	"context"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSwaggerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get swagger json
func NewGetSwaggerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSwaggerLogic {
	return &GetSwaggerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSwaggerLogic) GetSwagger() (resp *types.HtmlResponse, err error) {
	// body := fmt.Sprintf("<html><head><title>Page Title</title>	</head>	<body>	<pre>%s</pre></body></html>", API_FILE)
	// body := fmt.Sprintf("<html><head><title>Page Title</title>	</head>	<body>	<pre>%s</pre></body></html>", API_FILE)
	body := API_FILE

	return &types.HtmlResponse{
		ContentType: "application/json",
		Body:        body,
	}, nil
}
