#!/usr/bin/env python2

from coapthon.server.coap import CoAP
from basic_resource import BasicResource
from coapthon.resources.resource import Resource

class WellKnownResource(Resource):
    def __init__(self, name="WellKnownResource", coap_server=None):
        super(WellKnownResource, self).__init__(name, coap_server, visible=True,
                                            observable=True, allow_children=True)
        self.payload = """
</.well-known/core>;ct=40,</actuators/leds>;title="LEDs: ?color=r|g|b, POST/PUT mode=on|off";rt="Control",</sensors/sht21>;title="Temperature and Humidity";rt="Sht21",</sensors/max44009>;title="Light";rt="max44009"
        """

    def render_GET(self, request):
        return self

class CoAPServer(CoAP):
    def __init__(self, host, port):
        CoAP.__init__(self, (host, port))
        self.add_resource('basic/', BasicResource())
        self.add_resource('.well-known/core', WellKnownResource())

def main():
    server = CoAPServer("::", 5683)
    try:
        server.listen(10)
    except KeyboardInterrupt:
        print "Server Shutdown"
        server.close()
        print "Exiting..."

if __name__ == '__main__':
    main()
