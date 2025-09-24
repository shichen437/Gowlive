#!/bin/sh

exec sh -c /gowlive/main & nginx -g 'daemon off;'
