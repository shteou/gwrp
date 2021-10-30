# GitHub WebHook Reverse Proxy (gwrp)

gwrp (pronounced `gɜ:rp`) is a reverse proxy, or router, for your Github WebHooks.

gwrp validates any incoming webhooks, and routes them to other service endpoints.  
Routing rules are defined by a triple of the GitHub event, a target endpoint (POST)
and a `jq` query which must successfully evaluate for the route to match.

gwrp is useful where you want to route large numbers of GitHub webhooks to different components. This
avoids the limitation on the number of webhooks for each

## Example configuration

The following environment variables define two rules, `GWRP_PRS` and `ALL_PINGS`.

Each rule is of the form `github,event,list|url|jq query`.

```
RULE_GWRP_PRS="pull_request|https://my-receiver:8080/gwrp-webhook|.repository.full_name=='shteou/gwrp'"
RULE_ALL_PINGS="ping|https://my-receiver:8080/pings|."
```

GWRP_PRS is defined for `pull_request` events only, where the repository's full name is `shteou/gwrp`.  
ALL_PINGS is defined for `ping` events, and matches all json payloads, via the identity operator (`.`).

## GitHub setup

A GitHub Webhook should be setup pointing to your gwrp instance on the `/webhook` endpoint. It must specify
the json payload option. You can choose to send all events to gwrp, or just a subset.

You can also setup an organization webhook to forward to gwrp.

## Helm Installations

A fairly standard helm chart is supplied. `gitHub.secretKey` must be supplied, and an arbitrary number
of rules can be configured as in the following example::

```yaml
gitHub:
  secretKey: "your-secret-key"

rules:
  GWRP_PRS:
    events: pull_requests
    route: https://my-receiver:8080/gwrp-webhook
    jq: .repository.full_name == \"shteou/gwrp\"
  ALL_PINGS:
    events: ping
    route: https://my-receiver:8080/pings
    jq: .
```


