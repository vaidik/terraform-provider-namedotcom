# terraform-provider-namedotcom

[Terraform][1] provider for manage DNS records in [Name.com][2].

[![Build Status](https://travis-ci.org/vaidik/terraform-provider-namedotcom.svg?branch=master)](https://travis-ci.org/vaidik/terraform-provider-namedotcom)

## Usage

### Quick Example

```
provider "namedotcom" {
    user    = "vaidik"
    token   = "my-secret-token"
}

# Apex record: vaidik.in
resource "namedotcom_record" "vaidik-in" {
    domain_name =   "vaidik.in"
    answer      =   "127.0.0.1"
    type        =   "A"
}

# Sub-domain: blog.vaidik.in
resource "namedotcom_record" "blog-vaidik-in" {
    domain_name =   "vaidik.in"
    host        =   "blog"
    answer      =   "google.com"
    type        =   "CNAME"
    ttl         =   "600"
}
```

### Provider Configuration

Read more about how to obtain API credentials (user and token) from [API
docs][3].

After you have generated your API credentials, instantiate the provider like so:

```
provider "namedotcom" {
    user    = "<username>"
    token   = "<token>"
}
```

### DNS Records Resource Configuration

**Usage:**
```
resource "namedotcom_dns_record" "vaidik-in" {
    domain_name    = "vaidik.in"
    host           = "blog"
    answer         = "1.1.1.1"
    type           = "A"
    ttl            = "400"
}
```

**Arguments:**

* `domain_name`: the zone that the record belongs to: e.g. for a record for
  blog.example.org, domain would be "example.org".
* `host` (optional): hostname relative to the zone: e.g. for a record for
  blog.example.org, host would be "blog". 
* `type`: one of the following: A, AAAA, ANAME, CNAME, MX, NS, SRV, or TXT.
* `answer`: either the IP address for A or AAAA records; the target for ANAME,
  CNAME, MX, or NS records; the text for TXT records. For SRV records, answer
  has the following format: "{weight} {port} {target}" e.g.
  "1 5061 sip.example.org".
* `ttl` (optional): the time this record can be cached for in seconds.
* `priority` (optional): only required for MX and SRV records, it is ignored
  for all others.

Read more about each argument in the [API docs][4].

[1]: https://terraform.io/
[2]: https://www.name.com/
[3]: https://www.name.com/api-docs/
[4]: https://www.name.com/api-docs/DNS#CreateRecord
