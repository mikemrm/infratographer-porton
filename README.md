# Portón (Gate in spanish)

Portón is an opinionated API Gateway based on [krakend](https://www.krakend.io/)
which includes a set of plugins to integrate with the Infratographer services.
The plugins are to be added in the near future.

Currently, this is based on krakend enterprise, so make sure you have a license file
before running this.

# Defaults

Portón comes with Krakend's flexible configuration enabled by default
[[1](https://www.krakend.io/docs/enterprise/configuration/flexible-config/)].

The following directories are empty by default and should be populated
by an orchestration engine if you want to take advantage of flexible configuration:

* `/etc/krakend/settings`
* `/etc/krakend/partials`
* `/etc/krakend/templates`

**NOTE**: These settings are **NOT** to be included into the container image,
as they'll be dependent on the deployment and the services fronted by the gateway.

Plugins are built and copied into the resulting container image. The directory
containing the plugins is `/usr/lib/krakend/plugins`.

# References

- [1] https://www.krakend.io/docs/enterprise/configuration/flexible-config/