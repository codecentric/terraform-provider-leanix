package leanix

type WebhookSubscription struct {
	Id                  *string    `json:"id",omitempty`
	Identifier          string     `json:"identifier"`
	DeliveryType        string     `json:"deliveryType"`
	TagSets             [][]string `json:"tagSets"`
	WorkspaceId         string     `json:"workspaceId",omitempty`
	TargetUrl           string     `json:"targetUrl"`
	TargetMethod        string     `json:"targetMethod"`
	AuthorizationHeader string     `json:"authorizationHeader",omitempty`
	Callback            string     `json:"callback",omitempty`
	IgnoreError         bool       `json:"ignoreError"`
	WorkspaceConstraint string     `json:"workspaceConstraint"`
	Active              bool       `json:"active"`
}
