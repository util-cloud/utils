#!/bin/sh

for i in $*
do
	mkdir $i ;
	cd $i ;
#	man $i | unexpand -a --tabs=7 > $i.txt ;
	man $i | col -b > $i.txt ;
	cd .. ;
done
