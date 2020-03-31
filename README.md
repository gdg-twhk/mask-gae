# mask.goodideas-studio.com 後台

```yaml
dispatch:
  - url: "*/favicon.ico"
    module: default

  - url: "*/api/pharmacies*"
    module: pharmacy

  - url: "*/docs*"
    module: docs

  - url: "*/api/feedback*"
    module: feedback

  - url: "*/googleae8f4bcce8bec00c.html"
    module: ownership

  - url: "mask-9999.appspot.com/"
    module: default
```

```shell script
$ make
Usage:

  deploy_pharmacy [v=version-name]   deploy pharmacy service
  deploy_docs [v=version-name]       deploy docs service
  deploy_feedback [v=version-name]   deploy feedback service
  deploy_dispatch                    deploy disptach
  build_swagger                      generate swagger docs
```