# Since we are in the asset folder, move two
# folders up to get the program name

dir=$(pwd)/..
program=$(basename $(builtin cd $dir; pwd))

# The folders we are working with

codeFolder="/home/per/code"
binFolder="/home/per/bin/softteam"
softtubeFolder="/softtube/bin"

# Perform the copy

cp -rf $codeFolder/$program/build/* $binFolder/$program
cp -rf $codeFolder/$program/build/* $softtubeFolder/$program

# No longer needed after version 2.6.9
# cp -rf $codeFolder/$program/assets $binFolder/$program
# cp -rf $codeFolder/$program/assets $softtubeFolder/$program
