#!/bin/bash
echo "Start submitting code to the local repository"
git add *
echo;



echo "Commit the changes to the local repository"
echo "please enter the commit info...."
read info
date=`date +%Y/%m/%d`
time=`date +%r`
info="$info by FishInMars(Ubuntu) $date $time"
git commit -m "$info"
echo;
 
echo "Commit the changes to the remote git server"
git push origin master
echo;
 
echo "Batch execution complete!"
echo;
