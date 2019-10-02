package main

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
                            Type:             schema.TypeInt,
                            Optional:         true,
                        },
                },
        }
}

func makeClient (d *schema.ResourceData) (*namecom.NameCom, error) {
    user := d.Get("user")
    token := d.Get("token")

    client := namecom.New(user.(string), token.(string))
    _, err := client.HelloFunc(&namecom.HelloRequest{})
    if err != nil {
        log.Printf("Error connecting with Name.com API")
    }

    return client, err
}

func resourceNamedotcomRecordCreate(d *schema.ResourceData, meta interface{}) error {
    return resourceNamedotcomRecordRead(d, meta)
}

func resourceNamedotcomRecordRead(d *schema.ResourceData, m interface{}) error {
    client, _ := makeClient(d)
    domainName := d.Get("domain_name").(string)
    id := d.Get("id").(int32)

    record, _ := client.GetRecord(&namecom.GetRecordRequest{
        DomainName: domainName,
        ID: id,
    })

    d.SetId(fmt.Sprintf("%v", record.ID))
    d.Set("domain_name", record.DomainName)
    d.Set("host", record.DomainName)
    d.Set("fqdn", record.Fqdn)
    d.Set("type", record.Type)
    d.Set("answer", record.Answer)
    d.Set("ttl", record.TTL)
    d.Set("priority", record.Priority)

    return nil
}

func resourceNamedotcomRecordUpdate(d *schema.ResourceData, m interface{}) error {
        return resourceNamedotcomRecordRead(d, m)
}

func resourceNamedotcomRecordDelete(d *schema.ResourceData, m interface{}) error {
        return nil
}
