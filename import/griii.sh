#!/bin/bash

media="/media/keisuke/RICOH GR"
dcim="$media"/DCIM/
echo "The following files will not be copied. Make sure that it's OK."
echo
find "$media" | grep -v "$dcim"
echo

lo=$(ls -lc "$dcim"/*/* --time-style="+,%Y%m%d," | cut -d , -f 2 | sort | head -1)
hi=$(ls -lc "$dcim"/*/* --time-style="+,%Y%m%d," | cut -d , -f 2 | sort | tail -1)
dst=/tank/photos/keisuke/Pictures/GRIII/"$lo"_"$hi"/
echo "Command: rsync -Pav $dcim $dst"

read -p "Are you sure? " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
  rsync -Pav "$dcim" $dst
fi
