package page

import (
	"context"
	"fmt"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTocLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get TOC
func NewGetTocLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTocLogic {
	return &GetTocLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTocLogic) GetToc() (resp *types.HtmlResponse, err error) {
	body := fmt.Sprintf(`
	<html>	<head>		<title>Privacy policy</title>	</head>	<body>	<pre>
	Terms of Service

1. Acceptance of Terms
By accessing or using our GPU Cloud Services, you agree to be bound by these Terms of Service.

2. Service Description
We provide on-demand access to high-performance GPU computing resources. Our services are designed for various applications, including machine learning, AI, data science, and gaming.

3. User Obligations
You agree to:

Use our services responsibly and ethically.
Not engage in any illegal or harmful activities.
Protect your account credentials.
Comply with all applicable laws and regulations.
4. Service Availability
We strive to provide reliable and uninterrupted service, but we cannot guarantee 100% uptime. We may perform maintenance or make changes to the service without prior notice.

5. Intellectual Property
All intellectual property rights, including software, documentation, and trademarks, are owned by us or our licensors.

6. Data Privacy and Security
We are committed to protecting your data privacy. Please refer to our Privacy Policy for more information.

7. Payment and Billing
You agree to pay for our services according to our pricing plans. We may modify our pricing plans at any time.

8. Termination
We may terminate your access to our services at any time, with or without notice, if you violate these Terms of Service.

9. Limitation of Liability
To the maximum extent permitted by law, we are not liable for any indirect, incidental, or consequential damages arising from the use of our services.

10. Dispute Resolution
Any disputes arising from these Terms of Service shall be governed by the laws of [Jurisdiction].
	</pre>	</body>	</html>
		`, API_FILE)

	return &types.HtmlResponse{
		ContentType: "text/html",
		Body:        body,
	}, nil
}
