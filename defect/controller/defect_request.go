package controller

type bulletinRequest struct {
	IssueNumber []string `json:"issue_number" binding:"required"`
}
