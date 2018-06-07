#!/bin/bash
for i in `ls`; do
	echo "compressing directory... : $i"
	tar czf $i.tgz $i
	echo "mv directory... : $i"
	cp $i.tgz /mnt/root/bak/vay/
	rm -f $i.tgz
	echo "$i done"
	echo ""
done
