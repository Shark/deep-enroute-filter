# En-Route Filter

## Use with Vagrant

Run `vagrant up`.

Run the following snippet to insert the filtering rule:

```bash
vagrant ssh router
sudo ip6tables -t raw -I PREROUTING -p udp --dport 5683 -j NFQUEUE --queue-num 0
```

The en-route filter web interface should be accessible at [http://localhost:8080].

Connect to the client and make a CoAP request:

```bash
$ vagrant ssh client
$ coap get coap://[fdf0:a23f:8cae:5b97::2]/basic
(2.05)  Basic Resource
```

You should now see that the request has been processed in the web interface.