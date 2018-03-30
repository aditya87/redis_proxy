#! /bin/bash

redis-server --port 7777 --daemonize yes
/app/redis_proxy &
sleep 2
/app/integration
