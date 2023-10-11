package mail

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/mailgun/mailgun-go"
	"github.com/ritsec/ops-bot-iii/config"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// SendVerificationEmail sends a verification email to the specified address
func SendVerificationEmail(to string, ctx ddtrace.SpanContext) (string, error) {
	span := tracer.StartSpan(
		"mail.SendVerificationEmail",
		tracer.ResourceName("mail.SendVerificationEmail"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	// Generate a 6-digit verification code
	verificationCode, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	code := fmt.Sprintf("%06d", verificationCode)

	// Build the HTML email
	subject := "Your 6-digit verification code"
	body := fmt.Sprintf(`<!DOCTYPE html><html><head> <meta charset="utf-8"> <meta http-equiv="x-ua-compatible" content="ie=edge"> <title>Email Confirmation</title> <meta name="viewport" content="width=device-width, initial-scale=1"> <style type="text/css"> /** * Google webfonts. Recommended to include the .woff version for cross-client compatibility. */ @media screen{@font-face{font-family: 'Source Sans Pro'; font-style: normal; font-weight: 400; src: local('Source Sans Pro Regular'), local('SourceSansPro-Regular'), url(https://fonts.gstatic.com/s/sourcesanspro/v10/ODelI1aHBYDBqgeIAH2zlBM0YzuT7MdOe03otPbuUS0.woff) format('woff');}@font-face{font-family: 'Source Sans Pro'; font-style: normal; font-weight: 700; src: local('Source Sans Pro Bold'), local('SourceSansPro-Bold'), url(https://fonts.gstatic.com/s/sourcesanspro/v10/toadOcfmlt9b38dHJxOBGFkQc6VGVFSmCnC_l7QZG60.woff) format('woff');}}/** * Avoid browser level font resizing. * 1. Windows Mobile * 2. iOS / OSX */ body, table, td, a{-ms-text-size-adjust: 100%%; /* 1 */ -webkit-text-size-adjust: 100%%; /* 2 */}/** * Remove extra space added to tables and cells in Outlook. */ table, td{mso-table-rspace: 0pt; mso-table-lspace: 0pt;}/** * Better fluid images in Internet Explorer. */ img{-ms-interpolation-mode: bicubic;}/** * Remove blue links for iOS devices. */ a[x-apple-data-detectors]{font-family: inherit !important; font-size: inherit !important; font-weight: inherit !important; line-height: inherit !important; color: inherit !important; text-decoration: none !important;}/** * Fix centering issues in Android 4.4. */ div[style*="margin: 16px 0;"]{margin: 0 !important;}body{width: 100%% !important; height: 100%% !important; padding: 0 !important; margin: 0 !important;}/** * Collapse table borders to avoid space between cells. */ table{border-collapse: collapse !important;}a{color: #1a82e2;}img{height: auto; line-height: 100%%; text-decoration: none; border: 0; outline: none;}</style></head><body style="background-color: #e9ecef;"> <table border="0" cellpadding="0" cellspacing="0" width="100%%"> <tr> <td align="center" bgcolor="#e9ecef"><!--[if (gte mso 9)|(IE)]> <table align="center" border="0" cellpadding="0" cellspacing="0" width="600"> <tr> <td align="center" valign="top" width="600"><![endif]--> <table border="0" cellpadding="0" cellspacing="0" width="100%%" style="max-width: 600px;"> <tr> <td align="center" valign="top" style="padding: 36px 24px;"> <a href="https://www.ritsec.club" target="_blank" style="display: inline-block;"> <img src="https://www.ritsec.club/assets/images/ritsec_logo_final-05.png" alt="Logo" border="0" width="200px"> </a> </td></tr></table><!--[if (gte mso 9)|(IE)]> </td></tr></table><![endif]--> </td></tr><tr> <td align="center" bgcolor="#e9ecef"><!--[if (gte mso 9)|(IE)]> <table align="center" border="0" cellpadding="0" cellspacing="0" width="600"> <tr> <td align="center" valign="top" width="600"><![endif]--> <table border="0" cellpadding="0" cellspacing="0" width="100%%" style="max-width: 600px;"> <tr> <td align="center" bgcolor="#ffffff" style="padding: 36px 24px 0; font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; border-top: 3px solid #d4dadf;"> <h1 style="margin: 0; font-size: 32px; font-weight: 700; letter-spacing: -1px; line-height: 48px;">OBIII Confirmation Code</h1> </td></tr></table><!--[if (gte mso 9)|(IE)]> </td></tr></table><![endif]--> </td></tr><tr> <td align="center" bgcolor="#e9ecef"><!--[if (gte mso 9)|(IE)]> <table align="center" border="0" cellpadding="0" cellspacing="0" width="600"> <tr> <td align="center" valign="top" width="600"><![endif]--> <table border="0" cellpadding="0" cellspacing="0" width="100%%" style="max-width: 600px;"> <tr> <td align="left" bgcolor="#ffffff"> <table border="0" cellpadding="0" cellspacing="0" width="100%%"> <tr> <td align="center" bgcolor="#ffffff" style="padding: 12px;"> <table border="0" cellpadding="0" cellspacing="0"> <tr> <td align="center" bgcolor="#1a82e2" style="border-radius: 6px;"> <a target="_blank" style="display: inline-block; padding: 16px 36px; font-family: 'Source Sans Pro', Helvetica, Arial, sans-serif; font-size: 16px; color: #ffffff; text-decoration: none; border-radius: 6px;">%v</a> </td></tr></table> </td></tr></table> </td></tr></table><!--[if (gte mso 9)|(IE)]> </td></tr></table><![endif]--> </td></tr><tr> <td align="center" bgcolor="#e9ecef" style="padding: 24px;"><!--[if (gte mso 9)|(IE)]> <table align="center" border="0" cellpadding="0" cellspacing="0" width="600"> <tr> <td align="center" valign="top" width="600"><![endif]--> <table border="0" cellpadding="0" cellspacing="0" width="100%%" style="max-width: 600px;"> </table><!--[if (gte mso 9)|(IE)]> </td></tr></table><![endif]--> </td></tr></table> </body></html>`, code)

	// Initialize the Mailgun client
	mg := mailgun.NewMailgun(config.MailGun.Domain, config.MailGun.APIKey)

	// Create the message
	message := mg.NewMessage(
		"noreply@"+config.MailGun.Domain, // From
		subject,                          // Subject
		"",                               // Text body (empty since we are using an HTML body)
		to,                               // To
	)

	// Set the HTML body of the email
	message.SetHtml(body)

	// Send the email
	_, _, err = mg.Send(message)
	if err != nil {
		return "", err
	}

	return code, nil
}
