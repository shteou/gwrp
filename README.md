# GitHub WebHook Reverse Proxy (gwrp)

gwrp (pronounced `gɜ:rp`) is a reverse proxy, or router, for your Github WebHooks.

gwrp validates any incoming webhooks, and routes them to other service endpoints.  
Routing rules are defined by `jq` queries.


## Example configuration

The following environment variables define two rules, GWRP_PRS and GWRP_BRANCHES.

Each rule is of the form `github,event,list|url|jq query`.


```
RULE_GWRP_PRS="pull_request|https://my-receiver:8080/gwrp-webhook|.action=='PR' and .repository.full_name=='shteou/gwrp'"
RULE_GWRP_BRANCHES="|https://my-receiver:8080/gwrp-webhook|.action=='PR' and .repository.full_name=='shteou/gwrp'"
```