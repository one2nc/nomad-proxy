# nomad-proxy

Sample Nomad job:
```hcl
job "nomad_proxy_%(dc)s" {
  datacenters = ["%(dc)s"]
  type        = "service"

  constraint {
    attribute = "${node.class}"
    value     = "service"
  }

  update {
    max_parallel = 1
    min_healthy_time = "10s"
    healthy_deadline = "2m"
    progress_deadline = "5m"
    auto_revert = true
  }

  group "nomad_proxy" {
    count = %(service_count)s

    constraint {
      distinct_hosts = true
    }

    restart {
      attempts = 10
      interval = "5m"
      delay    = "25s"
      mode     = "delay"
    }

    task "nomad-proxy" {
      driver = "docker"

      template {
          data = <<EOH
      {
  "rules": {
    "*": {
      "PUT": [
        "{{env "meta.twofa_token"}}"
      ],
      "POST": [
        "{{env "meta.twofa_token"}}"
      ],
      "DELETE": [
        "{{env "meta.twofa_token"}}"
      ]
    }
  }
}
  EOH
        destination = "/tmp/2fa.json"
        perms = "644"
      }

      env {
        JOB_PREFIX = "%(prefix)s_%(dc)s"
        SERVER_ADDR = "http://${NOMAD_IP_np}:4646"
        DC          = "%(dc)s"
      }

      config {
        image = "tsl8/nomad-client-proxy:%(version)s"
        command = "./proxy"
        args = [%(token_args)s]

        volumes = [
          "tmp/2fa.json:/etc/opt/nomad_proxy_2fa.json"
        ]

        port_map {
          np = 9988
        }

        logging {
          type = "syslog"

          config {
            syslog-format  = "rfc3164"
            syslog-address = "%(log_dest)s"
            tag            = "%(dc)s-nomad-proxy"
          }
        }
      }

      resources {
        cpu    = 50
        memory = 256

        network {
          mbits = 1
          port  "np" {
            static = "9988"
          }
        }
      }

      service {
        name = "nomad-proxy"
        port = "np"

        check {
          type     = "tcp"
          interval = "60s"
          timeout  = "10s"
        }
      }
    }
  }
}

```