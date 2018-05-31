# Filter

## Dependencies

- `libnetfilter-queue-dev` package

## Setup

Add an iptables rule: `iptables -t raw -A PREROUTING -p udp --dport 5683 -j NFQUEUE --queue-num 0`
