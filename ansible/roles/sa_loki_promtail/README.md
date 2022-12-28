sa_loki_promtail
================

[![Build Status](https://travis-ci.com/softasap/sa_loki_promtail.svg?branch=master)](https://travis-ci.com/softasap/sa_loki_promtail)
[![Build Status](https://github.com/softasap/sa_loki_promtail/workflows/CI/badge.svg?event=push)](https://github.com/softasap/sa_loki_promtail/actions?query=workflow%3ACI)

Example of usage:

Simple

```YAML

     - {
         role: "sa_loki_promtail",
         loki_promtail_version: "1.5.0"
       }
```

Advanced

```YAML

roles:

     - {
         role: "sa_loki_promtail",
         loki_version: "1.5.0",
         loki_user:   loki,
         loki_group:  loki,
         loki_base_dir: /opt/loki
       }
```

Configuring logging
-------------------

scrape_configs:
  - job_name: system
    entry_parser: raw
    static_configs:
    - targets:
        - localhost
      labels:
        job: varlogs
        __path__: /var/log/*log
  - job_name: nginx
    entry_parser: raw
    static_configs:
    - targets:
        - localhost
      labels:
        job: nginx
        __path__: /var/log/nginx/*log


Usage with ansible galaxy workflow
----------------------------------

If you installed the `sa_loki_promtail` role using the command


`
   ansible-galaxy install softasap.sa_loki_promtail
`

the role will be available in the folder `library/softasap.sa_loki_promtail`
Please adjust the path accordingly.

```YAML

     - {
         role: "softasap.sa_loki_promtail"
       }

```




Copyright and license
---------------------

Code is dual licensed under the [BSD 3 clause] (https://opensource.org/licenses/BSD-3-Clause) and the [MIT License] (http://opensource.org/licenses/MIT). Choose the one that suits you best.

Reach us:

Subscribe for roles updates at [FB] (https://www.facebook.com/SoftAsap/)

Join gitter discussion channel at [Gitter](https://gitter.im/softasap)

Discover other roles at  http://www.softasap.com/roles/registry_generated.html

visit our blog at http://www.softasap.com/blog/archive.html
