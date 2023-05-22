ARG image_from_name
ARG image_from_tag
ARG image_from_registry_host

FROM ${image_from_registry_host}/${image_from_name}:${image_from_tag}


COPY ./${app_name}/app.sh /app/run.sh

CMD ["/app/run.sh"]