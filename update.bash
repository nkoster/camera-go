#!/bin/bash

cp html.temp html.go
sed -i '/9e31d7ffe4b03a35f666ae495f72964c/ r www/index.html' html.go 
sed -i '/9e31d7ffe4b03a35f666ae495f72964c/d' html.go
go build
scp camera-go w3b:camera/
#scp www/index.html w3b:camera/
#scp www/drag.html w3b:camera/
date
