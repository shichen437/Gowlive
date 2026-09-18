#!/bin/sh

restart_attempts=0
max_restart_attempts=99
restart_delay=2
max_restart_delay=180

nginx -g 'daemon off;' &
nginx_pid=$!

while true; do
    echo "[entrypoint] Starting gowlive backend..."
    /gowlive/main

    exit_code=$?
    restart_attempts=$((restart_attempts + 1))

    if [ $restart_attempts -gt $max_restart_attempts ]; then
        echo "[entrypoint] Reached max restart attempts ($max_restart_attempts). Exiting."
        exit 1
    fi

    echo "[entrypoint] Backend exited with code $exit_code (attempt $restart_attempts/$max_restart_attempts). Restarting in ${restart_delay}s..."
    sleep $restart_delay

    restart_delay=$((restart_delay * 2))
    if [ $restart_delay -gt $max_restart_delay ]; then
        restart_delay=$max_restart_delay
    fi
done

kill $nginx_pid 2>/dev/null
exit $?
