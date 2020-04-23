FROM irisnet/irishub:latest

COPY . /scripts

RUN sh /scripts/setup.sh

EXPOSE 26657

CMD iris start