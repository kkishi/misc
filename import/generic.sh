#!/bin/bash

media="$1"  # /media/keisuke/SDHC
dest="$2"  # /tank/photos/ayumi/Pictures/RX100III

dates=$(find "$media"/DCIM/*/* -printf "\"%p\" " | xargs -n 1 -P 32 exiftool -DateTimeOriginal -d '%Y%m%d' -s -S | sort)
lo=$(echo $dates | sed -e 's/ /\n/g' | head -1)
hi=$(echo $dates | sed -e 's/ /\n/g' | tail -1)

echo "Command: rsync -Pav $media/ $dest/"$lo"_"$hi"/"

read -p "Are you sure? " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
  rsync -Pav "$media"/ "$dest"/"$lo"_"$hi"/
fi
