# Portón (Gate in spanish)

Portón is an opinionated API Gateway based on [krakend](https://www.krakend.io/)
which includes a set of plugins to integrate with the Infratographer services.
The plugins are to be added in the near future.

Currently, this is based on krakend enterprise, so make sure you have a license file
before running this.

# Defaults

Portón comes with Krakend's flexible configuration enabled by default
[[1](https://www.krakend.io/docs/enterprise/configuration/flexible-config/)].

Plugins are built and copied into the resulting container image. The directory
containing the plugins is `/opt/krakend/plugins`.

The base krakend configuration and directories for the flexible configuration are not included
in the image. They are expected to be mounted into the container at runtime or
built into an extended all-in-one image.

# porton plugin

The porton plugin is a krakend plugin that provides authorization services
for the Infratographer services. It is built from the source in the `plugin`
directory.

It takes the following configuration options:

- `authz_service`: The URL of the auth service.
- `authz_service.endpoint`: The endpoint of the auth service.
- `authz_service.timeout`: The timeout for the auth service call in milliseconds. (default: `1000`)
- `action`: The action to validate against the auth service.
- `tenant_source`: The source of the tenant ID. Can be `header` or `path`. (default: `path`)

# References

- [1] https://www.krakend.io/docs/enterprise/configuration/flexible-config/