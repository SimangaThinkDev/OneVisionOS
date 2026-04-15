#!/bin/bash

FLAG_FILE="$HOME/.onevision_welcome_done"

if [ ! -f "$FLAG_FILE" ]; then
    firefox-esr --new-window "file:///opt/onevision/welcome/index.html"
    touch "$FLAG_FILE"
fi
