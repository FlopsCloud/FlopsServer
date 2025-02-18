package page

import (
	"context"
	"fmt"

	"fca/api/internal/svc"
	"fca/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPPLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get Privacy policy
func NewGetPPLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPPLogic {
	return &GetPPLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPPLogic) GetPP() (resp *types.HtmlResponse, err error) {
	body := fmt.Sprintf(`
<html>	<head>		<title>Privacy policy</title>	</head>	<body>	<pre>
FlopsCloud Privacy Policy

Effective Date: 2024-11-1

FlopsCloud is committed to protecting your privacy. This Privacy Policy outlines how we collect, use, disclose, and protect your personal information when you use our services.   

Information We Collect

We may collect the following types of personal information:   

Personal Information: Name, email address, phone number, and other contact information.
Usage Data: Information about how you use our services, including your IP address, browser type, device information, and usage patterns.
Payment Information: Credit card or other payment information if you make purchases through our services.
How We Use Your Information

We may use your personal information for the following purposes:

Providing Services: To deliver our services and fulfill your requests.
Improving Services: To analyze usage patterns and improve our services.
Customer Support: To respond to your inquiries and provide customer support.
Marketing and Promotions: To send you marketing communications and promotional offers.
Legal Compliance: To comply with legal obligations.
Sharing Your Information

We may share your personal information with:

Service Providers: Third-party service providers who assist us in delivering our services.   
Legal Authorities: To comply with legal requests or requirements.
Data Security

We implement reasonable security measures to protect your personal information from unauthorized access, use, disclosure, alteration, or destruction. However, no method of transmission over the Internet or electronic storage is completely secure.   

Your Rights

You may have the right to:

Access your personal information.
Correct any inaccuracies in your personal information.
Erase your personal information.
Restrict the processing of your personal information.
Object to the processing of your personal information.   
Data portability.   
Changes to This Privacy Policy

We may update this Privacy Policy from time to time. We will notify you of any material changes.   

Contact Us

If you have any questions or concerns about this Privacy Policy or our privacy practices, please contact us at [insert contact information].   

Please note: This is a basic template and may not cover all specific privacy requirements. It's essential to consult with legal counsel to ensure compliance with local data protection laws, such as GDPR, CCPA, or other relevant regulations.

Additional Considerations:

Cookie Policy: Consider adding a separate Cookie Policy to detail how cookies are used on your website.
Data Retention: Specify how long you retain personal information.
International Data Transfers: If you process data outside the EU or other relevant jurisdictions, include information about data transfer mechanisms.
Security Measures: Provide more details about the security measures you implement.
Third-Party Services: Clearly disclose any third-party services integrated into your platform.
By carefully tailoring this template to your specific business practices and legal requirements, you can create a robust Privacy Policy that protects your users' privacy and complies with relevant laws.

</pre>	</body>	</html>
	`, API_FILE)

	return &types.HtmlResponse{
		ContentType: "text/html",
		Body:        body,
	}, nil
}
