FROM golang:1.17 as Builder

WORKDIR /root

COPY . .

RUN ./make.bash build

FROM golang:1.17 as Configure

WORKDIR /

RUN apt-get update && apt-get install xz-utils bzip2 && \
    wget http://pttapp.cc/data-archives/bbs_backup_lastest.tar.xz && \
    tar -Jxvf bbs_backup_lastest.tar.xz && \
    wget http://pttapp.cc/data-archives/dump.shm.lastest.tar.bz2 && \
    tar -jxvf dump.shm.lastest.tar.bz2

FROM golang:1.17

WORKDIR /

RUN groupadd -g 99 bbs && \
    useradd -d /home/bbs \
            --gid 99 \
            --uid 9999 bbs

COPY --from=Configure /home/bbs/ /home/bbs/
COPY --from=Configure dump.shm .
COPY conf/ /conf/
COPY --from=Builder /root/Ptt-backend .

RUN chown bbs:bbs -R dump.shm conf

EXPOSE 8081/tcp

USER bbs

CMD [ "/Ptt-backend" ]
