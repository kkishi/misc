#!/bin/bash

media=/media/keisuke/7000-8000
dir=$media/DCIM/100GOPRO
echo "Make sure that everything you want to copy is under $dir:"
echo
find $media
echo

lo=$(ls -lc $dir/* --time-style="+,%Y%m%d," | cut -d , -f 2 | sort | head -1)
hi=$(ls -lc $dir/* --time-style="+,%Y%m%d," | cut -d , -f 2 | sort | tail -1)
cmd="rsync -Pav $dir/ /tank/photos/keisuke/Pictures/GOPRO/"$lo"_"$hi"/"
echo "Command: $cmd"

read -p "Are you sure? " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
  $cmd
fi
