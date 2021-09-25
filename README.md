# awsl

A WebSocket Linker

## 支持的协议

### server端

* http proxy (http)
* socks4/5 proxy (socks,socks5,socks4)
* websocket over tls (awsl)
* quic (quic)
* tcp (tcp)

### client端

* websocket over tls (awsl)
* quic (quic)
* tcp (tcp)

## awsl.service demo

```
[Unit]
Description=awsl
After=network-online.target network.target
Wants=network-online.target systemd-networkd-wait-online.service network.target

[Service]
Restart=on-failure
User=user
Group=users
ExecStart=awsl -m 80
ExecReload=/bin/kill $MAINPID
KillMode=mixed
KillSignal=SIGTERM
TimeoutStopSec=5s
StandardOutput=file:awsl.log
StandardError=file:awsl.error.log

[Install]
WantedBy=multi-user.target
```
