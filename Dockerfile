FROM krakend/krakend-ee-plugin-builder:2.1.0 AS pluginbuilder

WORKDIR /go/src/porton

RUN mkdir -p /go/src/porton/lib

# Run the plugin builder
# ...

FROM krakend/krakend-ee:2.1.1

# TODO: Run krakend as a non-root user

# For more information see https://www.krakend.io/docs/enterprise/configuration/flexible-config/
ENV FC_ENABLE=1
ENV FC_SETTINGS="/etc/krakend/settings"
ENV FC_PARTIALS="/etc/krakend/partials"
ENV FC_TEMPLATES="/etc/krakend/templates" 

# Note that these are expected to be volumes when running
# via an orchestration engine
RUN mkdir -p /etc/krakend/settings
RUN mkdir -p /etc/krakend/partials
RUN mkdir -p /etc/krakend/templates
RUN mkdir -p /usr/lib/krakend/plugins 

# Copy plugins from the pluginbuilder
COPY --from=pluginbuilder /go/src/porton/lib /usr/lib/krakend/plugins