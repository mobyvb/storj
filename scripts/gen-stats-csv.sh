#!/bin/bash

# osx: $ ./print-capt-ids.sh ~/Library/Application\ Support/Storj/Capt/
basepath=$1
i=0
while [ $i -le 99 ]
do
  nid=`identity ca id --ca.cert-path "${basepath}/f${i}/ca.cert"`
  if ((i % 2 == 0))
  then
    # good stats
    echo $nid,20,20,20,20
  else
    # bad stats
    echo $nid,20,5,20,5
  fi

  ((i++))
done


