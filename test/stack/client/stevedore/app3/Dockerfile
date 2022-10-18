ARG image_from

FROM ${image_from}

ARG message=app3
RUN echo "Hey there! Welcome to ${message}" > /message.txt

CMD ["cat","/message.txt"]
