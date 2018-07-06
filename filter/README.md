# Filter

## Dependencies

- `libnetfilter-queue-dev` package

## Setup

On the host, make sure that the containers can communicate:
```
ip6tables -A FORWARD -s fdf0:a23f:8cae:5b97::/64 -d fdba:cd7e:4c8e:a6fd::/64 -j ACCEPT
ip6tables -A FORWARD -s fdba:cd7e:4c8e:a6fd::/64 -d fdf0:a23f:8cae:5b97::/64 -j ACCEPT
```

Add the following iptables rule to intercept traffic:
```
ip6tables -t raw -I PREROUTING -p udp --dport 5683 -j NFQUEUE --queue-num 0
```
