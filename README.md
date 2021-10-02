# GitHub WebHook Reverse Proxy (gwrp)

gwrp (pronounced `gɜ:rp`) is a reverse proxy, or router, for your Github WebHooks.

gwrp validates any incoming webhooks, and routes them to other service endpoints.  
Routing rules are defined by `jq` queries.


## Example configuration

The following environment variables define a rule called `GWRP_PRS`.
The payload is passed on to 

```
RULE_GWRP_PRS=".action=='PR' and .repository.full_name=='shteou/gwrp'"
RULE_GWRP_PRS_URL="https://my-receiver:8080/gwrp-webhook"
```