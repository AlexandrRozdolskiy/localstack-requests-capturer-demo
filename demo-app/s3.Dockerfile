FROM ubuntu

WORKDIR /app

RUN apt update -y
RUN apt install pip -y
RUN pip install awscli

COPY demo-app/local.sh ./

CMD [ "/app/local.sh" ]