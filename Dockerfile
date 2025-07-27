FROM debian:latest

COPY terminaldefense /game/terminaldefense
COPY static/ /game/static/

ENV TD_PORT=8080

CMD [ "/game/terminaldefense" ]
