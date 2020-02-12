#!/bin/bash

echo
echo -----------
date
echo -----------
cp html.temp html.go
sed -i '/9e31d7ffe4b03a35f666ae495f72964c/ r www/index.html' html.go 
sed -i '/9e31d7ffe4b03a35f666ae495f72964c/d' html.go
go build
ssh w3b screen -X -S camera quit
scp camera w3b:camera/
ssh w3b screen -d -S camera -m ./camera/camera
scp www/index.html w3b:camera/
scp www/drag.html w3b:camera/
