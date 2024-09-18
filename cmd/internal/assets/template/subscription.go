package template

import "fmt"

var (
	confirmSubHTML    = `<h3>Dear user, you're now subscribed to the newsletter: %s! To unsubscribe, click on a link here <a href=\"%s/\">Unsubscribe</a>!</h3><br />May the delivery force be with you!`
	confirmSubText    = `Dear user, you're now subscribed to the newsletter: %s! To unsubscribe, click on link here: %s`
	confirmSubSubject = "STRV Newsletter subscription %d"
)

func GetConfirmSubHTML(newsletterName, unsubLink string) string {
	return fmt.Sprintf(confirmSubHTML, newsletterName, unsubLink)
}

func GetConfirmSubTxt(newsletterName, unsubLink string) string {
	return fmt.Sprintf(confirmSubText, newsletterName, unsubLink)
}

func GetConfirmSubSubject(newsletterName int64) string {
	return fmt.Sprintf(confirmSubSubject, newsletterName)
}
