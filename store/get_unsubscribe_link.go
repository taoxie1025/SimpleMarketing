package store

import (
	"github.com/astaxie/beego"
)

func (s *Store) GetUnsubscribeLink(projectId, projectName string) string {
	httpAddr := beego.AppConfig.String("httpaddr")
	unsubscribeLink := "https://" + httpAddr + "/unsubscribe?projectId=" + projectId + "&projectName=" + projectName
	text := "<hr><p><small><i>To unsubscribe, click <a href=\"" + unsubscribeLink + "\" rel=\"noopener noreferrer nofollow\">here</a></i>.</small></p>"
	return text
}
