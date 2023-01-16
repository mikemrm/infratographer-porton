// package plugin implements the porton krakend.io plugin
package plugin

/*
The portonRegisterer object implements the Registerer interface
for Lura's server and client plugins.

We must be careful not to change the `RegisterClients` and `RegisterHandlers`
methods' signatures, as they are called by the plugin loader.

It is not recommended to import the packages that define these interfaces.
That is:

* "github.com/luraproject/lura/v2/transport/http/client/plugin"
* "github.com/luraproject/lura/v2/transport/http/server/plugin"

This is to avoid any dependency issues that may arise from the plugin loader
importing the same packages as the main application. We instead stick to just
relying on the interfaces' signatures.

Also note that the plugin loader interface does not work if the `*Registerer`
variable instances are pointers. We need to use values instead. We decided to
build on top of `string` since that's what's working and documented in the
official plugin loader example.
*/
type portonRegisterer string

func NewPortonRegisterer(name string) portonRegisterer {
	return portonRegisterer(name)
}
