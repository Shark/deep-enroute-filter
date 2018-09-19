Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/xenial64"

  config.vm.define "router" do |router|
    router.vm.hostname = "router"
    router.vm.network "forwarded_port", guest: 8080, host: 8080, host_ip: '127.0.0.1'
    router.vm.network "private_network", auto_config: false, ip: "fdf0:a23f:8cae:5b97::1"
    router.vm.synced_folder ".", "/vagrant"

    # configure the private network interface
    # workaround for a bug in Vagrant which fails to bring up the interface with an IPv6 address
    router.vm.provision "shell", inline: <<~SHELL
      set -euxo pipefail

      # IFACE=$(ip link show | grep DOWN | cut -d: -f2 | tr -d '[:space:]')
      IFACE=enp0s8
      [[ -z $IFACE ]] && exit 0

      cat > /etc/network/interfaces.d/private.cfg <<EOF
auto $IFACE
iface $IFACE inet6 static
address fdf0:a23f:8cae:5b97::1/64
pre-up echo 0 > /proc/sys/net/ipv6/conf/$IFACE/accept_dad
pre-up echo 1 > /proc/sys/net/ipv6/conf/all/forwarding
up ip -6 addr add fd1b:eff6:e45:e487::1/64 dev $IFACE
down ip -6 addr del fd1b:eff6:e45:e487::1/64 dev $IFACE
EOF
      ifup $IFACE
    SHELL

    router.vm.provision "shell", inline: <<~SHELL
      set -euxo pipefail

      sudo add-apt-repository -y ppa:longsleep/golang-backports
      sudo apt-get update
      sudo apt-get -y --no-install-recommends install golang-go build-essential libnetfilter-queue-dev

      mkdir -p ~/go/bin ~/go/src/gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering ~/app
      export GOPATH=~/go
      export PATH=$PATH:~/go/bin
      curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
      cp -r /vagrant/filter ~/go/src/gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering
      cd ~/go/src/gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter
      dep ensure

      cd vendor/github.com/zubairhamed/canopus/
      mkdir openssl
      cd openssl
      git init
      git remote add origin https://github.com/openssl/openssl.git
      git fetch origin master
      git checkout b9b5181dd2f52ff0560a33b116396cdae5e48048
      ./config
      make -j4
      sudo cp libssl.so.1.1 libcrypto.so.1.1 /usr/lib

      cd ~/go/src/gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter
      go build .
      sudo cp filter /usr/bin/enroute-filter

      sudo cp enroute-filter.service /etc/systemd/system
      sudo systemctl daemon-reload
      sudo systemctl enable enroute-filter
      sudo systemctl start enroute-filter
    SHELL
  end

  config.vm.define "device" do |device|
    device.vm.hostname = "device"
    device.vm.network "private_network", auto_config: false, ip: "fdf0:a23f:8cae:5b97::2"

    # configure the private network interface
    # workaround for a bug in Vagrant which fails to bring up the interface with an IPv6 address
    device.vm.provision "shell", inline: <<~SHELL
      set -euxo pipefail

      # IFACE=$(ip link show | grep DOWN | cut -d: -f2 | tr -d '[:space:]')
      IFACE=enp0s8
      [[ -z $IFACE ]] && exit 0

      cat > /etc/network/interfaces.d/private.cfg <<EOF
auto $IFACE
iface $IFACE inet6 static
address fdf0:a23f:8cae:5b97::2/64
gateway fdf0:a23f:8cae:5b97::1
pre-up echo 0 > /proc/sys/net/ipv6/conf/$IFACE/accept_dad
EOF
      ifup $IFACE
    SHELL

    device.vm.provision "shell", inline: <<~SHELL
      set -euxo pipefail

      sudo apt-get update
      sudo apt-get -y --no-install-recommends install python python-pip python-setuptools
      mkdir -p /usr/src/app
      chown vagrant:vagrant /usr/src/app
    SHELL

    device.vm.provision "file", source: "examples/coap-server", destination: "/usr/src/app"

    device.vm.provision "shell", inline: <<~SHELL
      set -euxo pipefail

      cd /usr/src/app
      pip install --no-cache-dir -r requirements.txt
      cp coap-server.service /etc/systemd/system
      systemctl daemon-reload
      systemctl enable coap-server
      systemctl start coap-server
    SHELL
  end

  config.vm.define "client" do |client|
    client.vm.hostname = "client"
    client.vm.network "private_network", auto_config: false, ip: "fdf0:a23f:8cae:5b97::3"

    # configure the private network interface
    # workaround for a bug in Vagrant which fails to bring up the interface with an IPv6 address
    client.vm.provision "shell", inline: <<~SHELL
      set -euxo pipefail

      # IFACE=$(ip link show | grep DOWN | cut -d: -f2 | tr -d '[:space:]')
      IFACE=enp0s8
      [[ -z $IFACE ]] && exit 0

      cat > /etc/network/interfaces.d/private.cfg <<EOF
auto $IFACE
iface $IFACE inet6 static
address fd1b:eff6:e45:e487::2/64
gateway fd1b:eff6:e45:e487::1
pre-up echo 0 > /proc/sys/net/ipv6/conf/$IFACE/accept_dad
EOF
      ifup $IFACE
    SHELL

    client.vm.provision "shell", inline: <<~SHELL
      set -euxo pipefail

      curl -sL https://deb.nodesource.com/setup_10.x | sudo -E bash -
      sudo apt-get -y --no-install-recommends install nodejs
      sudo npm install -g coap-cli
    SHELL
  end

  # Share an additional folder to the guest VM. The first argument is
  # the path on the host to the actual folder. The second argument is
  # the path on the guest to mount the folder. And the optional third
  # argument is a set of non-required options.

  config.vm.provider "virtualbox" do |vb|
    # Display the VirtualBox GUI when booting the machine
    vb.gui = true

    # Customize the amount of memory on the VM:
    vb.memory = "512"
  end
end
