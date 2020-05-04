<h1 align="center">Welcome to mask-gae 👋</h1>
<p>
  <a href="https://github.com/gdg-twhk/mask-gae/blob/master/LICENSE" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
  <a href="https://twitter.com/CageChung" target="_blank">
    <img alt="Twitter: CageChung" src="https://img.shields.io/twitter/follow/CageChung.svg?style=social" />
  </a>
</p>

> mask.goodideas-studio.com backend API

## Usage

_dispatch.yaml_

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

## Author

👤 **Kai-Chu Chung**

* Website: https://kaichu.io/
* Github: [@cage1016](https://github.com/cage1016)

## Show your support

Give a ⭐️ if this project helped you!

## 📝 License

Copyright © 2020 [gdg-twhk](https://github.com/gdg-twhk).<br />
This project is [MIT](https://github.com/gdg-twhk/mask-gae/blob/master/LICENSE) licensed.

***
_This README was generated with ❤️ by [readme-md-generator](https://github.com/kefranabg/readme-md-generator)_