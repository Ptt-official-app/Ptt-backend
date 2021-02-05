#!/bin/bash

# usage: 
#    ./localtest.bash -m auto  --> Ignore interactions
#    ./localtest.bash          --> interactive mode
if [ $# -eq 0 ]; then
  # No argument provided
  mode="interactive"
else
  while [[ "$#" -gt 0 ]]; do
      case $1 in
          -m|--mode) mode="$2"; shift ;;
          *) echo "Unknown parameter passed: $1"; exit 1 ;;
      esac
      shift
  done
fi
case "$mode" in
  auto|interactive)
  	echo "$mode mode"
    ;;
  *)
  	echo "$mode mode not supported"
    exit 1
    ;;
esac

home="./home"
zipHome="bbs_backup_lastest.tar.xz"

if [ "$mode" == "auto" ]; then
  response="Y"
else
  echo "Do you want to update home? (Y/N)"
  read response
fi
if [ "$response" == "Y" ]; then
  updateDbResult=$(curl https://pttapp.cc/update-static-file -H "Authorization: pttapp")
  
  if [[ $updateDbResult == *"finish"* ]]; then
	echo "Update complete!"
    if [ -d "$home" ]; then
      if [ "$mode" == "auto" ]; then
	    response="Y"
	  else
		echo "$home exist, do you want to backup (Y/N)"
  	  	read response
	  fi
      if [ "$response" == "Y" ]; then
  	    timeStamp=`date +%s%N | cut -b1-13`
	    oldHome=${home}_${timeStamp}
  	    mv $home $oldHome
      else
  	    rm -rf $home
  	  fi
	fi
    if [ -f "$zipHome" ]; then
      rm $zipHome
	fi
    wget http://pttapp.cc/data-archives/bbs_backup_lastest.tar.xz
    tar -Jxvf bbs_backup_lastest.tar.xz
  else
    echo "Update failed: $updateDbResult"
  fi
fi

if [ "$mode" == "auto" ]; then
  response="Y"
else
  echo "Do you want to rebuild source? (Y/N)"
  read response
fi
if [ "$response" == "Y" ]; then
  cd cmd
  go build
  cd ..
fi

if [ -d "$home" ]; then
  ./cmd/cmd
else
  echo "Something went wrong, can't find your $home folder"
fi
