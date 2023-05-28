ARG image_from_name
ARG image_from_tag
ARG image_from_registry_host

FROM ${image_from_registry_host}/${image_from_name}:${image_from_tag}
ARG app_name
COPY ./${app_name}/app.sh /app.sh

CMD ["/app.sh"]
