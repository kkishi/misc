#!/bin/bash

media=/media/keisuke/EOS_DIGITAL
dir=$media/DCIM  # TODO: Copy the whole card
echo "Make sure that everything you want to copy is under $dir:"
echo
find $media
echo

lo=$(ls -lc $dir/*/* --time-style="+,%Y%m%d," | cut -d , -f 2 | sort | head -1)
hi=$(ls -lc $dir/*/* --time-style="+,%Y%m%d," | cut -d , -f 2 | sort | tail -1)
cmd="rsync -Pav $dir/ /tank/photos/keisuke/Pictures/5DII/"$lo"_"$hi"/"
echo "Command: $cmd"

read -p "Are you sure? " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
  $cmd
fi
