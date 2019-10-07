package namedotcom

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/namedotcom/go/namecom"
)

func resourceNamedotcomRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceNamedotcomRecordCreate,
		Read:   resourceNamedotcomRecordRead,
		Update: resourceNamedotcomRecordUpdate,
		Delete: resourceNamedotcomRecordDelete,

		Schema: map[string]*schema.Schema{
			"record_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fqdn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				// TODO: do we need ForceNew
			},
			"answer": {
				Type:     schema.TypeString,
				Required: true,
				// TODO: do we need ForceNew?
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func makeClient(d *schema.ResourceData, m interface{}) (*namecom.NameCom, error) {
	user := m.(*Config).User
	token := m.(*Config).Token

	log.Printf("user %s and token %s", user, token)
	client := namecom.New(user, token)
	_, err := client.HelloFunc(&namecom.HelloRequest{})
	if err != nil {
		log.Printf("Error connecting with Name.com API")
	}

	return client, err
}

func makeRecord(d *schema.ResourceData) *namecom.Record {
	domainName := d.Get("domain_name").(string)
	host := ""
	hostValue, hostOk := d.GetOk("host")
	if hostOk {
		host = hostValue.(string)
	}
	record_type := d.Get("type").(string)
	answer := d.Get("answer").(string)

	record := namecom.Record{
		DomainName: domainName,
		Host:       host,
		Type:       record_type,
		Answer:     answer,
	}

	recordId, recordOk := d.GetOk("record_id")
	if recordOk {
		record.ID = int32(recordId.(int))
	}

	return &record
}

func resourceNamedotcomRecordCreate(d *schema.ResourceData, m interface{}) error {
	client, _ := makeClient(d, m)
	record, err := client.CreateRecord(makeRecord(d))

	if err != nil {
		return fmt.Errorf("Failed to create record: %s", err)
	}

	d.SetId(record.Fqdn)
	d.Set("record_id", record.ID)

	return resourceNamedotcomRecordRead(d, m)
}

func resourceNamedotcomRecordRead(d *schema.ResourceData, m interface{}) error {
	client, _ := makeClient(d, m)
	domainName := d.Get("domain_name").(string)
	id := d.Get("record_id").(int)

	record, err := client.GetRecord(&namecom.GetRecordRequest{
		DomainName: domainName,
		ID:         int32(id),
	})

	if err != nil {
		d.SetId("")
		return nil
	}

	d.SetId(record.Fqdn)
	d.Set("record_id", record.ID)
	d.Set("domain_name", record.DomainName)
	d.Set("host", record.Host)
	d.Set("type", record.Type)
	d.Set("answer", record.Answer)
	d.Set("ttl", record.TTL)
	d.Set("priority", record.Priority)

	return nil
}

func resourceNamedotcomRecordUpdate(d *schema.ResourceData, m interface{}) error {
	client, _ := makeClient(d, m)
	record, err := client.UpdateRecord(makeRecord(d))

	if err != nil {
		return fmt.Errorf("Failed to update record: %s", err)
	}

	d.SetId(record.Fqdn)
	d.Set("record_id", record.ID)

	return resourceNamedotcomRecordRead(d, m)
}

func resourceNamedotcomRecordDelete(d *schema.ResourceData, m interface{}) error {
	client, _ := makeClient(d, m)
	domainName := d.Get("domain_name").(string)
	id := d.Get("record_id").(int)

	_, err := client.DeleteRecord(&namecom.DeleteRecordRequest{
		DomainName: domainName,
		ID:         int32(id),
	})

	if err != nil {
		return fmt.Errorf("Failed to destroy the resource: %s", err)
	}

	d.SetId("") // Explicitly destroy the resource
	return nil
}
