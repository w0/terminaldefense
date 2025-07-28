FROM debian:latest

COPY terminaldefense terminaldefense
COPY static/ static/

ENV TD_PORT=8080

CMD [ "/terminaldefense" ]
