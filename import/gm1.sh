#!/bin/bash

media=/media/keisuke/$1
dir=$media  # We copy everything because video files are stored outside of DCIM

lo=$(ls -lc $dir/DCIM/*/* --time-style="+,%Y%m%d," | cut -d , -f 2 | sort | head -1)
hi=$(ls -lc $dir/DCIM/*/* --time-style="+,%Y%m%d," | cut -d , -f 2 | sort | tail -1)
cmd="rsync -Pav $dir/ /tank/photos/keisuke/Pictures/GM1/"$lo"_"$hi"/"
echo "Command: $cmd"

read -p "Are you sure? " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
  $cmd
fi
