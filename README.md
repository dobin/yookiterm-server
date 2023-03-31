# Yookiterm-server

The main backend part of yookiterm. Provides challenges, and hints to where the actual 
yookiterm-lxdserver are. The UI will mostly query this server.


## What is yookiterm

Yookiterm provides per-user Linux root containers via JavasScript
terminal, and accompagning tutorials and writeups of
certain topics. It is currently used as a plattform
teaching exploit development at an university.

## What is yookiterm-server

yookiterm-server provides the following functionality: 
* authentication (user based and SSO)
* Provide the HTML UI
* Deliver challenges (from `yookiterm-challenges`)

It does not: 
* Interact directly with containers (thats what `yookiterm-lxdserver` is for)


## Install

Make sure you have something like the following directory structure: 
* /home/yookiterm/
  * yookiterm-server/ (this)
  * yookiterm/ (Angular UI)
  * yookiterm-slides/ (slides PDF)
  * yookiterm-challenges/ (challenges markdown)


Install yookiterm-server:
```
$ cd /home/yookiterm
$ git clone https://github.com/dobin/yookiterm-server.git
$ cd yookiterm-server

# build
$ go get
$ go build

# configure
$ cp yookierm-server.yml.sample yookiterm-server.yml
$ vi yookiterm-server.yml

# create base container
$ lxd.lxc init images:debian/11/amd64 Debian64
$ lxd.lxc init images:debian/11/i386  Debian32
```

Other:
```
# get challenges
$ cd /home/yookiterm
$ git clone https://github.com/dobin/yookiterm-challenges.git

# provide UI (optional)
$ git clone https://github.com/dobin/yookiterm.git
```

## configure reverse proxy

`Caddyfile`:
```
exploit.courses {
        reverse_proxy http://10.10.10.100:8080
}

vmaslr.yookiterm.ch {
        reverse_proxy http://10.10.10.101:8000
}

vmnoaslr.yookiterm.ch {
        reverse_proxy http://10.10.10.102:8000
}
```

`10.10.10.100` will run `yookiterm-server` (so this) in a container.

`10.10.10.101` and `10.10.10.102` are VM's providing `yookiterm-lxdserver`.


## Config file

Things to update:
* jwtsecret: A unique random string, keep it secret. **Use the same for yookiterm-lxdserver**!
* admin_password
* user_password
* container_hosts
* base_containers

`yookiterm-server.yml`:
```yml
jwtsecret: "supersecret"
server_addr: "[::]:80"
server_banned_ips:
server_url: "https://my.website"  # used for SSO

challenges_dir: "../yookiterm-challenges"
slides_dir: "../yookiterm-slides/"
frontend_dir: "../yookiterm/app/"

admin_password: "<pw>"  # admin access
user_password: "<pw>"  # login without SSO, any username with this password

googleId: ""
googleSecret: ""
azureId: ""
azureSecret: ""

# hostname is the public hostname/port of the VM hosting the relevant yookiterm-lxdserver
container_hosts:
- hostnamealias: ubuntuaslr
  hostname: container.my.website:41443
  aslr: true
  arch: intel
  sshbaseport: 51000
- hostnamealias: ubuntunoaslr
  hostname: container.my.website:42443
  aslr: false
  arch: intel
  sshbaseport: 52000
#- hostnamealias: ubuntuarm
#  hostname: container.my.website:43443
#  aslr: true
#  arch: arm
#  sshbaseport: 53000

# these containers are copied for each user on request
base_containers:
- id: "1"
  name: "Debian32"
  bits: "32"
- id: "2"
  name: "Debian64"
  bits: "64"

```


# Systemd Service 

```
cp yookiterm.service /etc/systemd/system
systemctl enable yookiterm
systemctl start yookiterm
```
