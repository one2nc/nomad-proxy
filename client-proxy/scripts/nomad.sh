#!/usr/bin/env bash
sudo apt-get update
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y unzip curl vim \
    apt-transport-https \
    ca-certificates \
    software-properties-common
# Download Nomad
NOMAD_VERSION=0.8.6

echo "Fetching Nomad..."
cd /tmp/
curl -sSL https://releases.hashicorp.com/nomad/${NOMAD_VERSION}/nomad_${NOMAD_VERSION}_linux_amd64.zip -o nomad.zip
unzip -oq nomad.zip
sudo install nomad /usr/bin/nomad

sudo mkdir -p /etc/nomad.d
sudo chmod a+w /etc/nomad.d

# Move the certs to /tmp

(
cat <<-EOF
    bind_addr = "0.0.0.0"
    log_level = "INFO"
    enable_syslog = true
    # Specify the Nomad client data directory
    data_dir = "/opt/nomad/data"

    region = "global"
    datacenter = "dc1"

    leave_on_terminate = true
    leave_on_interrupt = false

    # Set Nomad to server mode
    server {
        enabled = true
    }

    client {
        enabled = true
    }
    tls {
	    rpc  = true
	    http = true
	    ca_file   = "/tmp/cert-chain.pem"
	    cert_file = "/tmp/server.pem"
	    key_file  = "/tmp/server-key.pem"

	    verify_https_client    = true
    }
EOF
) | sudo tee /etc/nomad.d/server.hcl

sudo nomad agent -config=/etc/nomad.d/server.hcl -dev 2>&1 &
