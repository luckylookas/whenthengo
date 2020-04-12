FROM ubuntu:18.04
COPY ./executable /executable
RUN chmod 777 /executable
CMD /executable
