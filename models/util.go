package models

import (
	guuid "github.com/google/uuid"
)

func GenerateProjectUUID() string {
	return "project-" + guuid.New().String()[:8]
}

func GenerateArticleUUID() string {
	return "article-" + guuid.New().String()[:8]
}

func GenerateTicketUUID() string {
	return "ticket-" + guuid.New().String()[:8]
}

func GenerateCommentUUID() string {
	return "comment-" + guuid.New().String()[:8]
}
