docker run -d \
    --name timber \
    --restart unless-stopped \
    --log-opt max-size=1m \
    -e PUID=$UID \
    -e PGID=$GID \
    -p 6800:6800 \
    -p 6888:6888 \
    -p 6888:6888/udp \
    -p 3000:3000 \
    -v /home/timber/Downloads:/data \
    -e TZ=Asia/Shanghai \
    timberzhang/go-file:V2.2