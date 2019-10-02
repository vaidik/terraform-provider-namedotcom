package main

import (
        "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
        return &schema.Provider{
                Schema: map[string]*schema.Schema{
                    "user": {
                        Type:        schema.TypeString,
                        Optional:    true,
                        DefaultFunc: schema.EnvDefaultFunc("NAMEDOTCOM_USER", nil),
                        Description: "A registered Name.com username.",
                    },
                    "token": {
                        Type:        schema.TypeString,
                        Optional:    true,
                        DefaultFunc: schema.EnvDefaultFunc("NAMEDOTCOM_TOKEN", nil),
                        Description: "A registered Name.com API token.",
                    },
                },
                ResourcesMap: map[string]*schema.Resource{
                    "namedotcom_record": resourceNamedotcomRecord(),
                },
        }
}
