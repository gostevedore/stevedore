ARG parent_name
ARG parent_version

FROM ${parent_name}:${parent_version}

ARG whoami=unknown

RUN echo "Hey there, I'm ${whoami}!" > /whoami.txt

CMD ["cat","/whoami.txt"]
