ARG image_from_name
ARG image_from_tag

FROM ${image_from_name}:${image_from_tag}

ARG whoami=unknown

RUN echo "Hey there, I'm ${whoami}!" > /whoami.txt

CMD ["cat","/whoami.txt"]
