Caddy
=========
[![Galaxy](https://img.shields.io/badge/galaxy-samdoran.caddy-blue.svg?style=flat)](https://galaxy.ansible.com/samdoran/caddy)
[![Build Status](https://travis-ci.com/samdoran/ansible-role-caddy.svg?branch=master)](https://travis-ci.com/samdoran/ansible-role-caddy)

Install [Caddy](https://caddyserver.com) with a basic config set to include `/etc/caddy/conf.d/*`. Roles that wish to use Caddy should place their config file in that directory.

This role now supports Caddy v2 installed from a repository. There is a tasks file that can be used to clean up the previous Caddy v1 installation. It should be run **before** running the role to update to Caddy v2.

If you need custom plugins, define them in `caddy_plugins`. This will build and download a custom Caddy binary and replace the one installed by the system package manager.

Requirements
------------

None.

Role Variables
--------------

| Name              | Default Value       | Description          |
|-------------------|---------------------|----------------------|
| `caddy_user` | `caddy` | Caddy user |
| `caddy_group` | `caddy` | Caddy group |
| `caddy_service_name` | `caddy` | Name of the service for starting/stopping/enabling. |
| `caddy_default_port` | `80` | Default port Caddy will bind to. |
| `caddy_config_path` | `/etc/caddy` | Path to Caddy config. |
| `caddy_config_file` | `caddy.conf` | Name of the Caddy config file. |
| `caddy_root` | `/usr/share/caddy` | Path to the default root served by Caddy. |
| `caddy_global_config_options` | `[]` | List of global Caddy config options. Syntax must be correct. |
| `caddy_plugins` | `[]` | List of [plugins](https://caddyserver.com/download) to be added to Caddy. This will download a custom Caddy binary and replace the one installed from the repository. |


Dependencies
------------

- `community.general` collection

Example Playbook
----------------

    - name: Install Caddy
      hosts: all
      roles:
         - samdoran.caddy

    - name: Remove Caddy v1 and install Caddy v2
      hosts: all
      tasks:
        - name: Cleanup Caddy v1
          import_role:
            name: samdoran.caddy
            tasks_from: cleanup-v1.yml

        - name: Install Caddy v2
          import_role:
            name: samdoran.caddy

License
-------

Apache 2.0
