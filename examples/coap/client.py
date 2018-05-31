#!/usr/bin/env python2

from coapthon.client.helperclient import HelperClient

host = "::1"
port = 5683
path = "basic"

with open('payload.jpg', 'r') as content_file:
    client = HelperClient(server=(host, port))
    content = content_file.read()
    response = client.post(path, "test")
    print response.pretty_print()
    client.stop()
