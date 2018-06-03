#!/usr/bin/env python2

from coapthon.client.helperclient import HelperClient

host = "fdf0:a23f:8cae:5b97::2"
port = 5683
path = "basic"

with open('payload.jpg', 'r') as content_file:
    client = HelperClient(server=(host, port))
    content = content_file.read()
    response = client.post(path, "test")
    print response.pretty_print()
    client.stop()
