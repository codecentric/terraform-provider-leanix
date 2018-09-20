package leanix

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
)

func apiTokenAndLeanixBasicAuthHeader() (string, string) {
	const apiToken = "aVQEzWKwE2sSp3rhVKWwaVQEzWKwE2sSp3rhVKWw"
	return apiToken, "Basic " + base64.StdEncoding.EncodeToString([]byte("apitoken:"+apiToken))
}

func TestGetAuthorizationHeader(t *testing.T) {
	apiToken, leanixBasicAuthHeader := apiTokenAndLeanixBasicAuthHeader()
	authRoute, authHeader := NewAuthRouteDefinition(t, apiToken)

	testServer := NewTestServer(
		t,
		TestRoute{
			TestResourceAndMethod{Resource: "/services/mtm/v1/oauth2/token", Method: "POST"}: authRoute,
		},
	)

	defer testServer.Close()
	client := NewLeanixClient(
		testServer.URL,
		leanixBasicAuthHeader,
	)
	header, err := client.getAuthorizationHeader()
	if err != nil {
		t.Errorf("LeanixClient.getAuthorizationHeader() returned an error: %s", err)
	}
	assertEqual(t, header, authHeader)
}

func TestCreateWebhookSubscription(t *testing.T) {
	apiToken, leanixBasicAuthHeader := apiTokenAndLeanixBasicAuthHeader()
	authRoute, authHeader := NewAuthRouteDefinition(t, apiToken)

	subscription := &WebhookSubscription{
		Identifier:   "hook",
		DeliveryType: "PUSH",
		TagSets: [][]string{
			{"pathfinder", "FACT_SHEET_CREATED"},
			{"pathfinder", "FACT_SHEET_UPDATED"},
		},
		WorkspaceId:         "8751abbf-8093-410d-a090-10c7735952cf",
		TargetUrl:           "https://bjir4u9ata.execute-api.eu-central-1.amazonaws.com/test/events",
		TargetMethod:        "POST",
		AuthorizationHeader: "Basic bGVhbml4Omt2Y1l2djVuVEJUQTNXcGQK",
		Callback:            "delivery.payload = {\"lol\" : \"lel\"}",
		IgnoreError:         true,
		WorkspaceConstraint: "ANY",
		Active:              false,
	}
	subscriptionId := "id"
	expectedBody, err := json.Marshal(*subscription)
	if err != nil {
		t.Fatal(err)
	}

	createRoute := &TestRouteDefinition{
		ExpectedHeader: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": authHeader,
		},
		ExpectedBody: expectedBody,
		ResponseStatus: func(header http.Header, body []byte) int {
			return http.StatusOK
		},
		ResponseBody: func(header http.Header, body []byte) []byte {
			subscriptionWithId := *subscription
			subscriptionWithId.Id = &subscriptionId
			response := &WebhookSubscriptionResponse{
				Status:       "Ok",
				Subscription: &subscriptionWithId,
			}
			responseMarshal, err := json.Marshal(response)
			if err != nil {
				t.Fatal(err)
			}
			return responseMarshal
		},
	}

	testServer := NewTestServer(
		t,
		TestRoute{
			TestResourceAndMethod{Resource: "/services/mtm/v1/oauth2/token", Method: "POST"}:       authRoute,
			TestResourceAndMethod{Resource: "/services/webhooks/v1/subscriptions", Method: "POST"}: createRoute,
		},
	)

	defer testServer.Close()
	client := NewLeanixClient(
		testServer.URL,
		leanixBasicAuthHeader,
	)
	subscriptionResponse, err := client.CreateWebhookSubscription(*subscription)
	if err != nil {
		t.Fatalf("LeanixClient.CreateWebhookSubscription() returned an error: %s", err)
	}
	assertEqual(t, subscriptionResponse.Id, &subscriptionId)
	subscriptionResponse.Id = nil // we remove the ID and check if the rest of the struct is also equal
	assertEqual(t, subscriptionResponse, subscription)
}

func TestReadWebhookSubscription(t *testing.T) {
	apiToken, leanixBasicAuthHeader := apiTokenAndLeanixBasicAuthHeader()
	authRoute, authHeader := NewAuthRouteDefinition(t, apiToken)

	subscriptionId := "id"
	subscription := &WebhookSubscription{
		Id:           &subscriptionId,
		Identifier:   "hook",
		DeliveryType: "PUSH",
		TagSets: [][]string{
			{"pathfinder", "FACT_SHEET_CREATED"},
			{"pathfinder", "FACT_SHEET_UPDATED"},
		},
		WorkspaceId:         "8751abbf-8093-410d-a090-10c7735952cf",
		TargetUrl:           "https://bjir4u9ata.execute-api.eu-central-1.amazonaws.com/test/events",
		TargetMethod:        "POST",
		AuthorizationHeader: "Basic bGVhbml4Omt2Y1l2djVuVEJUQTNXcGQK",
		Callback:            "delivery.payload = {\"lol\" : \"lel\"}",
		IgnoreError:         true,
		WorkspaceConstraint: "ANY",
		Active:              false,
	}

	getRoute := &TestRouteDefinition{
		ExpectedHeader: map[string]string{
			"Authorization": authHeader,
		},
		ExpectedBody: []byte{},
		ResponseStatus: func(header http.Header, body []byte) int {
			return http.StatusOK
		},
		ResponseBody: func(header http.Header, body []byte) []byte {
			response := &WebhookSubscriptionResponse{
				Status:       "Ok",
				Subscription: subscription,
			}
			responseMarshal, err := json.Marshal(response)
			if err != nil {
				t.Fatal(err)
			}
			return responseMarshal
		},
	}

	testServer := NewTestServer(
		t,
		TestRoute{
			TestResourceAndMethod{Resource: "/services/mtm/v1/oauth2/token", Method: "POST"}:                        authRoute,
			TestResourceAndMethod{Resource: "/services/webhooks/v1/subscriptions/" + subscriptionId, Method: "GET"}: getRoute,
		},
	)

	defer testServer.Close()
	client := NewLeanixClient(
		testServer.URL,
		leanixBasicAuthHeader,
	)
	subscriptionResponse, err := client.ReadWebhookSubscription(subscriptionId)
	if err != nil {
		t.Fatalf("LeanixClient.ReadWebhookSubscription() returned an error: %s", err)
	}
	assertEqual(t, subscriptionResponse, subscription)
}

func TestUpdateWebhookSubscription(t *testing.T) {
	apiToken, leanixBasicAuthHeader := apiTokenAndLeanixBasicAuthHeader()
	authRoute, authHeader := NewAuthRouteDefinition(t, apiToken)

	subscriptionId := "id"
	subscription := &WebhookSubscription{
		Id:           &subscriptionId,
		Identifier:   "hook",
		DeliveryType: "PUSH",
		TagSets: [][]string{
			{"pathfinder", "FACT_SHEET_CREATED"},
			{"pathfinder", "FACT_SHEET_UPDATED"},
		},
		WorkspaceId:         "8751abbf-8093-410d-a090-10c7735952cf",
		TargetUrl:           "https://bjir4u9ata.execute-api.eu-central-1.amazonaws.com/test/events",
		TargetMethod:        "POST",
		AuthorizationHeader: "Basic bGVhbml4Omt2Y1l2djVuVEJUQTNXcGQK",
		Callback:            "delivery.payload = {\"lol\" : \"lel\"}",
		IgnoreError:         true,
		WorkspaceConstraint: "ANY",
		Active:              false,
	}

	expectedBody, err := json.Marshal(*subscription)
	if err != nil {
		t.Fatal(err)
	}

	updateRoute := &TestRouteDefinition{
		ExpectedHeader: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": authHeader,
		},
		ExpectedBody: expectedBody,
		ResponseStatus: func(header http.Header, body []byte) int {
			return http.StatusOK
		},
		ResponseBody: func(header http.Header, body []byte) []byte {
			response := &WebhookSubscriptionResponse{
				Status:       "Ok",
				Subscription: subscription,
			}
			responseMarshal, err := json.Marshal(response)
			if err != nil {
				t.Fatal(err)
			}
			return responseMarshal
		},
	}

	testServer := NewTestServer(
		t,
		TestRoute{
			TestResourceAndMethod{Resource: "/services/mtm/v1/oauth2/token", Method: "POST"}:                        authRoute,
			TestResourceAndMethod{Resource: "/services/webhooks/v1/subscriptions/" + subscriptionId, Method: "PUT"}: updateRoute,
		},
	)

	defer testServer.Close()
	client := NewLeanixClient(
		testServer.URL,
		leanixBasicAuthHeader,
	)
	subscriptionResponse, err := client.UpdateWebhookSubscription(*subscription)
	if err != nil {
		t.Fatalf("LeanixClient.UpdateWebhookSubscription() returned an error: %s", err)
	}
	assertEqual(t, subscriptionResponse, subscription)
}

func TestDeleteWebhookSubscription(t *testing.T) {
	apiToken, leanixBasicAuthHeader := apiTokenAndLeanixBasicAuthHeader()
	authRoute, authHeader := NewAuthRouteDefinition(t, apiToken)

	subscriptionId := "id"
	subscription := &WebhookSubscription{
		Id:           &subscriptionId,
		Identifier:   "hook",
		DeliveryType: "PUSH",
		TagSets: [][]string{
			{"pathfinder", "FACT_SHEET_CREATED"},
			{"pathfinder", "FACT_SHEET_UPDATED"},
		},
		WorkspaceId:         "8751abbf-8093-410d-a090-10c7735952cf",
		TargetUrl:           "https://bjir4u9ata.execute-api.eu-central-1.amazonaws.com/test/events",
		TargetMethod:        "POST",
		AuthorizationHeader: "Basic bGVhbml4Omt2Y1l2djVuVEJUQTNXcGQK",
		Callback:            "delivery.payload = {\"lol\" : \"lel\"}",
		IgnoreError:         true,
		WorkspaceConstraint: "ANY",
		Active:              false,
	}

	deleteRoute := &TestRouteDefinition{
		ExpectedHeader: map[string]string{
			"Authorization": authHeader,
		},
		ExpectedBody: []byte{},
		ResponseStatus: func(header http.Header, body []byte) int {
			return http.StatusOK
		},
		ResponseBody: func(header http.Header, body []byte) []byte {
			response := &WebhookSubscriptionResponse{
				Status:       "Ok",
				Subscription: subscription,
			}
			responseMarshal, err := json.Marshal(response)
			if err != nil {
				t.Fatal(err)
			}
			return responseMarshal
		},
	}

	testServer := NewTestServer(
		t,
		TestRoute{
			TestResourceAndMethod{Resource: "/services/mtm/v1/oauth2/token", Method: "POST"}:                           authRoute,
			TestResourceAndMethod{Resource: "/services/webhooks/v1/subscriptions/" + subscriptionId, Method: "DELETE"}: deleteRoute,
		},
	)

	defer testServer.Close()
	client := NewLeanixClient(
		testServer.URL,
		leanixBasicAuthHeader,
	)
	subscriptionResponse, err := client.DeleteWebhookSubscription(subscriptionId)
	if err != nil {
		t.Fatalf("LeanixClient.DeleteWebhookSubscription() returned an error: %s", err)
	}
	assertEqual(t, subscriptionResponse, subscription)
}

func NewAuthRouteDefinition(t *testing.T, apiToken string) (*TestRouteDefinition, string) {
	tokenType := "Bearer"
	accessToken := "this_is_my_awesome_oauth_token"
	return &TestRouteDefinition{
		ExpectedHeader: map[string]string{
			"Content-Type":  "application/x-www-form-urlencoded",
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte("apitoken:"+apiToken)),
		},
		ExpectedBody: []byte(url.Values{"grant_type": {"client_credentials"}}.Encode()),
		ResponseStatus: func(header http.Header, body []byte) int {
			return http.StatusOK
		},
		ResponseBody: func(header http.Header, body []byte) []byte {
			response := &AuthResponse{
				AccessToken: accessToken,
				TokenType:   tokenType,
			}
			responseMarshal, err := json.Marshal(response)
			if err != nil {
				t.Fatal(err)
			}
			return responseMarshal
		},
	}, tokenType + " " + accessToken
}
