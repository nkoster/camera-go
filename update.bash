#!/bin/bash

while inotifywait -e close_write www/* -e close_write *.go
do
    echo
    echo -----------
    date
    echo -----------
    scp www/index.html w3b:camera/
    scp www/drag.html w3b:camera/
    cp html.temp html.go
    sed -i '/9e31d7ffe4b03a35f666ae495f72964c/ r www/index.html' html.go 
    sed -i '/9e31d7ffe4b03a35f666ae495f72964c/d' html.go
    go build
    ssh w3b screen -X -S camera quit
    scp camera w3b:camera/
    ssh w3b screen -d -S camera -m ./camera/camera
done
