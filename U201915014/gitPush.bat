#!/bin/bash
echo "Start submitting code to the local repository"
echo "The current directory is：%cd%"
git add *
echo;



echo "Commit the changes to the local repository"
echo "please enter the commit info...."
set /p message=
set now=%date:~0,10% %time:~0,8%
git commit -m "%message% by FishInMars(Windows) %now%"
echo;
 
echo "Commit the changes to the remote git server"
git push
echo;
 
echo "Batch execution complete!"
echo;

pause