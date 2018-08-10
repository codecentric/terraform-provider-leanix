package leanix

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type LeanixClient struct {
	url                string
	apiToken           string
	http               *http.Client
	authorizationToken *string
	sync.Mutex
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type WebhookSubscriptionResponse struct {
	Status       string               `json:"status"`
	Subscription *WebhookSubscription `json:"data"`
}

func NewLeanixClient(url string, apiToken string) *LeanixClient {
	httpClient :=
		&http.Client{
			Timeout: time.Second * time.Duration(10),
		}
	return &LeanixClient{
		url:                url,
		apiToken:           apiToken,
		http:               httpClient,
		authorizationToken: nil,
	}
}

// Use the API token to get an OAuth2 token from LeanIX.
// The OAuth2 token can be used for further requests.
// This function returns the complete header, including the token type.
// To avoid requesting a new token while the old one is still valid,
// we synchronize calls towards this method and cache the token until it
// becomes invalid.
func (leanix *LeanixClient) getAuthorizationHeader() (string, error) {
	leanix.Lock()
	defer leanix.Unlock()
	if leanix.authorizationToken != nil {
		// TODO here we should check if the token is still valid
		return *leanix.authorizationToken, nil
	} else {
		postUrl := leanix.url + "/services/mtm/v1/oauth2/token"
		postBody := url.Values{"grant_type": {"client_credentials"}}
		req, err := http.NewRequest("POST", postUrl, strings.NewReader(postBody.Encode()))
		req.SetBasicAuth("apitoken", leanix.apiToken)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(postBody.Encode())))
		resp, err := leanix.http.Do(req)

		// Process response
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		authResponse := AuthResponse{}
		err = json.NewDecoder(resp.Body).Decode(&authResponse)
		if err != nil {
			return "", err
		}
		newToken := authResponse.TokenType + " " + authResponse.AccessToken
		leanix.authorizationToken = &newToken
		return *leanix.authorizationToken, nil
	}
}

// Create a new webhook subscription at LeanIX.
// This method needs a valid authorization header so it will attempt to get one.
func (leanix *LeanixClient) CreateWebhookSubscription(subscription WebhookSubscription) (*WebhookSubscription, error) {
	authorizationHeader, err := leanix.getAuthorizationHeader()
	if err != nil {
		return nil, err
	}

	postUrl := leanix.url + "/services/webhooks/v1/subscriptions"
	postBody, err := json.Marshal(subscription)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", authorizationHeader)

	resp, err := leanix.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	subscriptionResponse := WebhookSubscriptionResponse{}
	err = json.NewDecoder(resp.Body).Decode(&subscriptionResponse)
	if err != nil {
		return nil, err
	}
	if subscriptionResponse.Subscription == nil || subscriptionResponse.Subscription.Id == nil {
		return nil, errors.New("Failed to create subscription '" + subscription.Identifier + "'. Maybe it was already created outside of Terraform?")
	}
	newSubscription := subscriptionResponse.Subscription

	return newSubscription, nil
}

// Read a new webhook subscription from LeanIX.
// This method needs a valid authorization header so it will attempt to get one.
func (leanix *LeanixClient) ReadWebhookSubscription(subscriptionId string) (*WebhookSubscription, error) {
	authorizationHeader, err := leanix.getAuthorizationHeader()
	if err != nil {
		return nil, err
	}

	getUrl := leanix.url + "/services/webhooks/v1/subscriptions/" + subscriptionId

	req, err := http.NewRequest("GET", getUrl, nil)
	req.Header.Add("Authorization", authorizationHeader)

	resp, err := leanix.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	subscriptionResponse := WebhookSubscriptionResponse{}
	err = json.NewDecoder(resp.Body).Decode(&subscriptionResponse)
	if err != nil {
		return nil, err
	}
	if subscriptionResponse.Subscription == nil || subscriptionResponse.Subscription.Id == nil {
		return nil, errors.New("Failed to read subscription '" + subscriptionId + "'. Maybe it was already deleted outside of Terraform?")
	}
	subscription := subscriptionResponse.Subscription

	return subscription, nil
}

// Updates an existing webhook subscription at LeanIX.
// This method needs a valid authorization header so it will attempt to get one.
func (leanix *LeanixClient) UpdateWebhookSubscription(subscription WebhookSubscription) (*WebhookSubscription, error) {
	authorizationHeader, err := leanix.getAuthorizationHeader()
	if err != nil {
		return nil, err
	}

	putUrl := leanix.url + "/services/webhooks/v1/subscriptions/" + *subscription.Id
	putBody, err := json.Marshal(subscription)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", putUrl, bytes.NewBuffer(putBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", authorizationHeader)

	resp, err := leanix.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	subscriptionResponse := WebhookSubscriptionResponse{}
	err = json.NewDecoder(resp.Body).Decode(&subscriptionResponse)
	if err != nil {
		return nil, err
	}
	if subscriptionResponse.Subscription == nil || subscriptionResponse.Subscription.Id == nil {
		return nil, errors.New("Failed to update subscription '" + subscription.Identifier + "'. Maybe it was already created outside of Terraform?")
	}
	newSubscription := subscriptionResponse.Subscription

	return newSubscription, nil
}

// Delete a webhook subscription at LeanIX.
// This method needs a valid authorization header so it will attempt to get one.
func (leanix *LeanixClient) DeleteWebhookSubscription(subscriptionId string) (*WebhookSubscription, error) {
	authorizationHeader, err := leanix.getAuthorizationHeader()
	if err != nil {
		return nil, err
	}

	deleteUrl := leanix.url + "/services/webhooks/v1/subscriptions/" + subscriptionId

	req, err := http.NewRequest("DELETE", deleteUrl, nil)
	req.Header.Add("Authorization", authorizationHeader)

	resp, err := leanix.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	subscriptionResponse := WebhookSubscriptionResponse{}
	err = json.NewDecoder(resp.Body).Decode(&subscriptionResponse)
	if err != nil {
		return nil, err
	}
	if subscriptionResponse.Subscription == nil || subscriptionResponse.Subscription.Id == nil {
		return nil, errors.New("Failed to delete subscription with ID'" + subscriptionId + "'. Maybe it was already deleted outside of Terraform?")
	}
	deletedSubscription := subscriptionResponse.Subscription

	return deletedSubscription, nil
}
