package leanix

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLeanixWebhookSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceLeanixWebhookSubscriptionCreate,
		Read:   resourceLeanixWebhookSubscriptionRead,
		Update: resourceLeanixWebhookSubscriptionUpdate,
		Delete: resourceLeanixWebhookSubscriptionDelete,

		Schema: map[string]*schema.Schema{
			"identifier": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"target_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"target_method": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"workspace_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"authorization_header": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"callback": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace_constraint": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ANY",
			},
			"payload_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "DEFAULT",
			},
			"active": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"ignore_error": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"tag_set": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tags": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceLeanixWebhookSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	leanixClient := meta.(*LeanixClient)

	subscription := WebhookSubscription{
		Identifier:          d.Get("identifier").(string),
		DeliveryType:        "PUSH",
		TagSets:             extractTagSets(d.Get("tag_set")),
		WorkspaceId:         d.Get("workspace_id").(string),
		TargetUrl:           d.Get("target_url").(string),
		TargetMethod:        d.Get("target_method").(string),
		AuthorizationHeader: d.Get("authorization_header").(string),
		Callback:            d.Get("callback").(string),
		IgnoreError:         d.Get("ignore_error").(bool),
		WorkspaceConstraint: d.Get("workspace_constraint").(string),
		PayloadMode:         d.Get("payload_mode").(string),
		Active:              d.Get("active").(bool),
	}
	created, err := leanixClient.CreateWebhookSubscription(subscription)
	if err != nil {
		return err
	}

	d.SetId(*created.Id)
	return nil
}

func resourceLeanixWebhookSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	leanixClient := meta.(*LeanixClient)

	subscriptionId := d.Id()
	if subscriptionId == "" {
		return errors.New("Terraform internal resource ID not set. Cannot delete resource!")
	}

	subscription, err := leanixClient.ReadWebhookSubscription(subscriptionId)
	if err != nil {
		return err
	}
	if subscription == nil {
		d.SetId("")
		return nil
	}

	d.Set("identifier", subscription.Identifier)
	d.Set("tag_set", packageTagSets(subscription.TagSets))
	d.Set("workspace_id", subscription.WorkspaceId)
	d.Set("target_url", subscription.TargetUrl)
	d.Set("target_method", subscription.TargetMethod)
	d.Set("authorization_header", subscription.AuthorizationHeader)
	d.Set("callback", subscription.Callback)
	d.Set("ignore_error", subscription.IgnoreError)
	d.Set("workspace_constraint", subscription.WorkspaceConstraint)
	d.Set("payload_mode", subscription.PayloadMode)
	d.Set("active", subscription.Active)

	return nil
}

func resourceLeanixWebhookSubscriptionUpdate(d *schema.ResourceData, meta interface{}) error {
	leanixClient := meta.(*LeanixClient)

	subscriptionId := d.Id()
	if subscriptionId == "" {
		return errors.New("Terraform internal resource ID not set. Cannot delete resource!")
	}

	subscription := WebhookSubscription{
		Id:                  &subscriptionId,
		Identifier:          d.Get("identifier").(string),
		DeliveryType:        "PUSH",
		TagSets:             extractTagSets(d.Get("tag_set")),
		WorkspaceId:         d.Get("workspace_id").(string),
		TargetUrl:           d.Get("target_url").(string),
		TargetMethod:        d.Get("target_method").(string),
		AuthorizationHeader: d.Get("authorization_header").(string),
		Callback:            d.Get("callback").(string),
		IgnoreError:         d.Get("ignore_error").(bool),
		WorkspaceConstraint: d.Get("workspace_constraint").(string),
		PayloadMode:         d.Get("payload_mode").(string),
		Active:              d.Get("active").(bool),
	}
	updated, err := leanixClient.UpdateWebhookSubscription(subscription)
	if err != nil {
		return err
	}

	d.SetId(*updated.Id)
	return nil
}

func resourceLeanixWebhookSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	leanixClient := meta.(*LeanixClient)

	subscriptionId := d.Id()
	if subscriptionId == "" {
		return errors.New("Terraform internal resource ID not set. Cannot delete resource!")
	}

	_, err := leanixClient.DeleteWebhookSubscription(subscriptionId)
	if err != nil {
		return err
	}

	return nil
}

func extractTagSets(value interface{}) [][]string {
	var extractedTagSets [][]string
	for setIndex, rawTagSet := range value.(*schema.Set).List() {
		extractedTagSets = append(extractedTagSets, []string{})
		castedTagSet := rawTagSet.(map[string]interface{})
		for _, tag := range castedTagSet["tags"].([]interface{}) {
			castedTag := tag.(string)
			extractedTagSets[setIndex] = append(extractedTagSets[setIndex], castedTag)
		}
	}
	return extractedTagSets
}

func packageTagSets(tagSets [][]string) []map[string][]string {
	var packagedTagSets []map[string][]string
	for _, tagSet := range tagSets {
		tagsMap := map[string][]string{"tags": tagSet}
		packagedTagSets = append(packagedTagSets, tagsMap)
	}
	return packagedTagSets
}
