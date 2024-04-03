#!/bin/bash

media=$(ls -d /run/user/1000/gvfs/mtp\:host\=Google_Pixel_5_11041FDD4000U7/内部共有ストレージ)
dir=$media/DCIM
echo "Make sure that everything you want to copy is under $dir:"
echo
#find $media
#echo

lo=$(ls -lc $dir/Camera/PXL_* --time-style="+,%Y%m%d," | cut -d , -f 2 | sort | head -1)
hi=$(ls -lc $dir/Camera/PXL_* --time-style="+,%Y%m%d," | cut -d , -f 2 | sort | tail -1)
cmd="rsync -Pav $dir/ /tank/photos/ayumi/Pictures/PIXEL5_"$lo"_"$hi"/"
echo "Command: $cmd"

read -p "Are you sure? " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
  $cmd
fi
