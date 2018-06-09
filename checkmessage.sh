#!/bin/bash

logfile='/var/log/messages'
tmpfile='/tmp/.dmessageinfo'
if test -e $tmpfile
then 
    md5=`head -1  $logfile |md5sum |awk '{print $1}'`
    fsize=`du -b $logfile | awk '{print $1}'`
    oldmd5=`head -1 $tmpfile| awk '{print $1}'`
    oldsize=`head -1 $tmpfile| awk '{print $2}'`
    if ([ $md5 = $oldmd5 ])
	then
	    size=$(($fsize-$oldsize))
	else
	    size=$fsize
	fi
	echo "$md5    $fsize" >$tmpfile
else
    md5=`head -1  $logfile |md5sum |awk '{print $1}'`
	fsize=`du -b $logfile | awk '{print $1}'`
	echo "$md5    $fsize" >$tmpfile
	exit
fi

if test -e $logfile
then

tail -c $size $logfile |perl -e '
while(<>){ 
	       if ( $_ =~ /.*(nf_conntrack|segfault|bucket table overflow|page allocation failure|Out of memory|file-max limit).*/i)
		   {
           print "ERROR ",$1,"in messages \n";
           exit 2;
		   }
		}
'
fi
