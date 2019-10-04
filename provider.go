package main

import (
        "github.com/hashicorp/terraform-plugin-sdk/terraform"
        "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Config struct {
    User string
    Token string
}

func Provider() terraform.ResourceProvider {
        return &schema.Provider{
                Schema: map[string]*schema.Schema{
                    "user": {
                        Type:        schema.TypeString,
                        Required:    true,
                        DefaultFunc: schema.EnvDefaultFunc("NAMEDOTCOM_USER", nil),
                        Description: "A registered Name.com username.",
                    },
                    "token": {
                        Type:        schema.TypeString,
                        Required:    true,
                        DefaultFunc: schema.EnvDefaultFunc("NAMEDOTCOM_TOKEN", nil),
                        Description: "A registered Name.com API token.",
                    },
                },

                ResourcesMap: map[string]*schema.Resource{
                    "namedotcom_record": resourceNamedotcomRecord(),
                },

                ConfigureFunc: providerConfigure,
        }
}

func providerConfigure(data *schema.ResourceData) (interface{}, error) {
    providerConfig := Config{
        User: data.Get("user").(string),
        Token: data.Get("token").(string),
    }

    return &providerConfig, nil
}
