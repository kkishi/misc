#!/bin/bash

media=/media/keisuke/EOS_DIGITAL

lo=$(ls -lc $media/DCIM/*/* --time-style="+,%Y%m%d," | cut -d , -f 2 | sort | head -1)
hi=$(ls -lc $media/DCIM/*/* --time-style="+,%Y%m%d," | cut -d , -f 2 | sort | tail -1)
cmd="rsync -Pav $media/ /tank/photos/keisuke/Pictures/R6/"$lo"_"$hi"/"
echo "Command: $cmd"

read -p "Are you sure? " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
  $cmd
fi
