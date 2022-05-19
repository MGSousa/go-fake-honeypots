# Fake Honeypots

Simple honeypots for SSH and Telnet with a fake Shell (can be extended)

## Install manually

> Clone repo
```sh
git clone https://github.com/MGSousa/go-fake-honeypots.git && cd go-fake-honeypots
```

> Build Go binaries
```sh
./build.sh
```

> Execute the Honeypot(s): <br>
```sh
cd ./dist/

./sshd
./telnetd
```

## Install (via Docker)
> sshd
```sh
docker build -t fh-sshd -f Dockerfile.ssh .
docker run --rm -p 22:22 -it fh-sshd
```
> telnetd
```sh
docker build -t fh-telnetd -f Dockerfile.telnetd .
docker run --rm -p 23:23 -it fh-telnetd
```

<br>
- All users trying to connecting via TELNET to port 23 will be shown a fake CISCO router login (Any input will lead to telnet shell) <br>
- All users trying to connect via SSH to port 22 will login into a fake shell (Password is: password.. it's possible to also remove pass auth &/or use key auth)<br>
- All the actions executed by malicious users will be saved into fh-telnet.log / fh-ssh.log <br><br>
If you want it to run 24/7, you can setup a systemd unit/supervisord running in background keeping the script up or just lunch the command with screen <br>

