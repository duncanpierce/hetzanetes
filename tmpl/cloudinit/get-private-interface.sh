$(ip -j route list {{.PrivateIpRange}} | jq -r .[0].dev)