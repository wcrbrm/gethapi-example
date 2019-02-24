FROM ethereum/client-go:stable

# no need in genesis, we wil use --dev flag to start test enviroment instead
# COPY genesis.json /root/

COPY init.js /root/
COPY start.sh /root/
ENTRYPOINT ["sh", "/root/start.sh"]