#!/bin/bash
i=0
while [ $i -le 1800000 ]
do
    echo `date` >> metrics1.txt
    curl -XGET -s $IP:$PORT/metrics | grep --extended-regexp '^starrocks_be_tablet_meta_mem_bytes' >> metrics1.txt
    sleep 10
    let i++
done
