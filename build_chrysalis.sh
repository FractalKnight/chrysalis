#!/bin/bash
#

function ProgressBar {
# Process data
	let _progress=(${1}*100/${2}*100)/100
	let _done=(${_progress}*4)/10
	let _left=40-$_done
# Build progressbar string lengths
	_done=$(printf "%${_done}s")
	_left=$(printf "%${_left}s")

# 1.2 Build progressbar strings and print the ProgressBar line
# 1.2.1 Output example:
# 1.2.1.1 Progress : [########################################] 100%
printf "\rProgress : [${_done// /#}${_left// /-}] ${_progress}%%"

}

_start=1
_end=300

printf "Starting the Chrysalis process...\n"
sleep 10
nohup go run chrysalis.go & > /dev/null 2>&1
printf "Finding files for backup\n"
sleep 30
printf "Adding files to Chrysalis\n"
sleep 5
for number in $(seq ${_start} ${_end})
do 
	sleep 0.1
	ProgressBar ${number} ${_end}
done
printf "\nDone creating Chrysalis"
