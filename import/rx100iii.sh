#!/bin/bash

for d in "3239-3438" "disk"; do
  if [[ -d "/media/keisuke/${d}" ]]; then
    media="/media/keisuke/${d}";
  fi
done

if [[ $media = "" ]]; then
  echo "SD card not found"
  exit
fi

dir=$media/DCIM
echo "Make sure that everything you want to copy is under $dir:"
echo
find $media
echo

lo=$(ls -lc $dir/*/* --time-style="+,%Y%m%d," | cut -d , -f 2 | sort | head -1)
hi=$(ls -lc $dir/*/* --time-style="+,%Y%m%d," | cut -d , -f 2 | sort | tail -1)
cmd="rsync -Pav $dir/ /tank/photos/ayumi/Pictures/RX100III_"$lo"_"$hi"/"
echo "Command: $cmd"

read -p "Are you sure? " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
  $cmd
fi
