sa_loki
=======

[![Build Status](https://travis-ci.com/softasap/sa_loki.svg?branch=master)](https://travis-ci.com/softasap/sa_loki)
[![Build Status](https://github.com/softasap/sa_loki/workflows/CI/badge.svg?event=push)](https://github.com/softasap/sa_loki/actions?query=workflow%3ACI)

Example of usage:

Simple

```YAML

     - {
         role: "sa_loki",
         loki_version: "1.5.0"
       }
```

Advanced

```YAML

roles:

     - {
         role: "sa_loki",
         loki_version: "1.5.0",
         loki_user:   loki,
         loki_group:  loki,
         loki_base_dir: /opt/loki
       }
```

Connecting clients
------------------

In difference from prometheus, Loki uses push concept. I.e. promtail clients are pushing payload to loki cluster.
It works like a charm in AWS or similar cloud, when you are using private addresses, but if you do not have private network,
most likely you will need to limit access a bit...

```
sudo ufw allow from 192.168.1.215 proto tcp to any port 3100
```

Usage with ansible galaxy workflow
----------------------------------

If you installed the `sa_loki` role using the command


`
   ansible-galaxy install softasap.sa_loki
`

the role will be available in the folder `library/softasap.sa_loki`
Please adjust the path accordingly.

```YAML

     - {
         role: "softasap.sa_loki"
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
