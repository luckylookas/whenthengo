# for release workflow
FROM ubuntu:18.04
RUN mkdir /whenthen
COPY ./executable /whenthen/executable
RUN chmod 777 /whenthen/executable
CMD /whenthen/executable
