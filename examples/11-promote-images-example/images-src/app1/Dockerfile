ARG image_from_name
ARG image_from_registry_host
ARG image_from_registry_namespace
ARG image_from_tag

FROM ${image_from_registry_host}/${image_from_registry_namespace}/${image_from_name}:${image_from_tag}

COPY ./${app_name}/app.sh /app.sh

CMD ["/app.sh"]