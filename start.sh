#!/usr/bin/env bash

mkdir -p /var/run/mitmproxy
rm -f /var/run/mitmproxy/*.sock

nkfproxy &
sedproxy &
awkproxy &
miniproxy &

sed -i "s/\$hostproxy_addr/$hostproxy_addr/g" /etc/nginx/conf.d/default.conf
/usr/local/openresty/bin/openresty -g 'daemon off;'
