package main

import (
        "fmt"
        "log"
        "strconv"

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
                            Type:     schema.TypeString,
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
    client, _ := makeClient(d)
    domainName := d.Get("domain_name").(string)
    host := ""
    hostValue, hostOk := d.GetOk("host"); if hostOk {
        host = hostValue.(string)
    }
    record_type := d.Get("type").(string)
    answer := d.Get("answer").(string)

    record, err := client.CreateRecord(&namecom.Record{
        DomainName: domainName,
        Host: host,
        Type: record_type,
        Answer: answer,
    })

    if err != nil {
        return fmt.Errorf("Failed to create record: %s", err)
    }

    d.SetId(record.Fqdn)
    d.Set("record_id", fmt.Sprintf("%v", record.ID))

    return resourceNamedotcomRecordRead(d, meta)
}

func resourceNamedotcomRecordRead(d *schema.ResourceData, m interface{}) error {
    client, _ := makeClient(d)
    domainName := d.Get("domain_name").(string)
    id, _ := strconv.ParseInt(d.Get("record_id").(string), 10, 32) //, 10, 32)

    record, _ := client.GetRecord(&namecom.GetRecordRequest{
        DomainName: domainName,
        ID: int32(id),
    })

    d.SetId(record.Fqdn)
    d.Set("domain_name", record.DomainName)
    d.Set("host", record.Host)
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
