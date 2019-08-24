## Minimal DIY UDID Retriever
* Green Verified (signed) `.mobileconfig`
* HTTPS web-server
* Go for cross-platform single-binary deployment
* Open-source and private
  * I created this to protect my Ad Hoc testers' privacy
* [Live demo](https://udid.aliask.co)

## Rough guide for Ubuntu 18.04 LTS AMD64 VPS
* No other guide/support will be provided
* Non-bug issues will be closed
* Sorry, I don't have the bandwidth
* Forks/maintainers welcome

1. `git clone https://github.com/EmperorEarth/udid && cd udid`
1. Create VPS  
  a. Digitalocean Droplet, AWS EC2, etc  
  b. Recommend 18.04 LTS amd64 if using Ubuntu
1. Choose URL  
  a. Need this for Let's Encrypt certificate  
  b. Can be `domain.tld` or `subdomain.domain.tld`
1. Point URL to VPS  
1. Replace `subdomain.domain.tld` in listed files with your chosen URL  
  a. `main.go`:`23`  
  b. `udid_unsigned.mobileconfig`:`8` (keep `/upload`)
1. Generate a UUID  
  a. [Online UUID generator](https://www.uuidgenerator.net/version4)
1. Replace `12345678-1234-1234-1234-1234567890ab` with your generated UUID  
  a. `udid_unsigned.mobileconfig`:`21`
1. `GOARCH=amd64 GOOS=linux go build -o udid`  
  a. `go build` [documentation](https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies)  
  b. [valid `GOARCH` + `GOOS` combinations](https://golang.org/doc/install/source#environment)
1. `sftp`/`psftp` to VPS  
1. `put`/`mput` `udid` binary and `udid_unsigned.mobileconfig`  
1. `sudo setcap 'cap_net_bind_service=+ep' ./udid && chmod 500 ./udid`  
  a. Allows server binary to bind to 80 & 443 (and other privileged ports < 1024)  
  b. Changes server binary permissions to read/execute by logged in user
1. `./udid`  
  a. Starts server to generate TLS Certificate from Let's Encrypt
1. Navigate to `subdomain.domain.tld/foo` adjust URL (keep `/foo`)  
  a. Encourages server to generate TLS certificate faster
1. Refresh until browser receives valid TLS certificate (locked lock icon left of URL in Chrome)  
1. `ls certificates/` should show a new file `subdomain.domain.tld`  
1. Using whichever text editor you prefer, copy parts of `certiticates/subdomain.domain.tld` into various files  
  a. Copy from `-----BEGIN EC PRIVATE KEY-----` until `-----END EC PRIVATE KEY-----` into `private-key.pem`  
  b. Copy from the first `-----BEGIN CERTIFICATE-----` until the first `-----END CERTIFICATE-----` into `certificate.pem`  
  c. Copy from the second `-----BEGIN CERTIFICATE-----` until the second `-----END CERTIFICATE-----` into `certificate-authority.pem`
1. `openssl smime -sign -signer ./certificate.pem -inkey ./private-key.pem -certfile ./certificate-authority.pem -nodetach -outform der -in ./udid_unsigned.mobileconfig -out ./udid.mobileconfig`  
  a. Signs your `.mobileconfig` file so users will see a green `Valid`
1. `sudo vi /etc/systemd/system/udid.service`  
  a. Creates config file for a `systemd` service that will start on startup/crash  
  b. See `Sample SystemD config file` section for sample config
1. `sudo systemctl enable udid`  
  a. Links and enables systemd on each startup/crash
1. `sudo reboot`  
  a. `sudo systemctl start udid` won't log properly until restart  
  b. Something wonky with `journalctl`. if you have a fix, please file an issue
1. Navigate to `subdomain.domain.tld` on an iPhone  

## Sample SystemD config file
* Replace `username` with VPS username
```
[Unit]
Description=UDID service
After=network.target
After=systemd-user-sessions.service
After=network-online.target

[Service]
WorkingDirectory=/home/username
ExecStart=/home/username/udid
ExecStop=/usr/bin/pkill udid
Restart=on-failure
RestartSec=30

[Install]
WantedBy=multi-user.target
```

## Assorted notes
* Lines `15`, `17`, `23`, `25` are customizable in `udid_unsigned.mobileconfig`
* Go cross compiles, so I recommend installing it locally
* Installing Go [documentation](https://golang.org/doc/install#install)
* If using `ufw`, make sure `http`&`https` are allowed
  * `sudo ufw allow http`
  * `sudo ufw allow https`
