# SES - DNS Check

Check whether the DNS entries necessary for SES are done. 

## Input 

```bash
./checkdns my-domain.com eu-central-1
```

## Example
```bash
./checkdns contact.letsbuild-aws.com eu-central-1
=== SES DKIM Configuration ===
Status: SUCCESS

DKIM Tokens (selectors):
  n6x*********6._domainkey.contact.letsbuild-aws.com
  7ti*********pcjtm._domainkey.contact.letsbuild-aws.com
  wps5*********7sstrcv._domainkey.contact.letsbuild-aws.com

=== DNS Records ===
_amazonses.contact.letsbuild-aws.com                         Fail
_dmarc.contact.letsbuild-aws.com                             Pass
  v=DMARC1; p=none;
n6xanyluw2oml4h5my5test47gxpqqw6._domainkey.contact.letsbuild-aws.com Pass
  p=MII*********8wIDAQAB
7tij*********tay6lpcjtm._domainkey.contact.letsbuild-aws.com Pass

wps5fke*********k7sstrcv._domainkey.contact.letsbuild-aws.com Pass
```

## Getting started

User binary or do:

```bash
task build
```

With [installed go](https://go.dev/) before.
