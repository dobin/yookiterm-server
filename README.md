# Yookiterm-server

The main backend part of yookiterm. Provides challenges, and hints to where the actual 
yookiterm-lxdserver are. The UI will mostly query this server.


## What is yookiterm

Yookiterm provides per-user Linux root containers via JavasScript
terminal, and accompagning tutorials and writeups of
certain topics. It is currently used as a plattform
teaching exploit development at an university.


## Install

Make sure you have something like the following directory structure: 
* /opt/exploit.courses/
  * yookiterm-server/ (this)
  * yookiterm/ (Angular UI)
  * yookiterm-slides/ (slides PDF)
  * yookiterm-challenges/ (challenges markdown)


```
# build
go get
go build

# configure
cp yookierm-server.yml.sample yookiterm-server.yml
vi yookiterm-server.yml

# get challenges
cd ../
git clone https://github.com/dobin/yookiterm-challenges.git

# provide UI (optional)
git clone https://github.com/dobin/yookiterm.git

# create base container0
lxc init ubuntu:16.04 hlUbuntu32
lxc init ubuntu:16.04 hlUbuntu64
```

## Config file

Things to update:
* jwtsecret: A unique random string, keep it secret. Use same for yookiterm-lxdserver!
* server_domain: Not used
* admin_password
* user_password
* container_hosts
* base_containers

```yml
jwtsecret: "<choose secret>"
server_addr: "[::]:8090"
server_banned_ips:
server_maintenance: false
server_domain: "container.exploit.courses"
challenges_dir: "../yookiterm-challenges"
admin_password: "<pw>"
user_password: "<pw>"

container_hosts:
- hostnamealias: ubuntuaslr
  hostname: container.exploit.courses:41443
  aslr: true
  arch: intel
  sshbaseport: 51000
- hostnamealias: ubuntunoaslr
  hostname: container.exploit.courses:42443
  aslr: false
  arch: intel
  sshbaseport: 52000
- hostnamealias: ubuntuarm
  hostname: container.exploit.courses:43443
  aslr: true
  arch: arm
  sshbaseport: 53000

base_containers:
- id: "1"
  name: "hlUbuntu32"
  bits: "32"
- id: "2"
  name: "hlUbuntu64"
  bits: "64"
```
