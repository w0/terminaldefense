FROM debian:latest

COPY terminaldefense terminaldefense
COPY static/ static/

CMD [ "/terminaldefense" ]
