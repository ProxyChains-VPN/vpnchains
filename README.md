# vpnchains
## About
Allows to intercept child proccesses' connect() calls (see LD_PRELOAD) and de facto connect via wireguard tunnel.
Implemented with the use of gvisor netstack and wireguard-go (that is implemented with the former library).

## How to install
First of all, install libbson (mongo-db-driver). Most likely, it exists in official repos of your Linux distribution.  
Secondly, run `install.sh` script; this one will check if libbson includes exist in `/usr/include` directory and compile and move library in `/usr/lib/libvpnchains_inject.so`. The executable will be located at `build/vpnchains`.

## How to use
**!!! Command args have priority over environment variables, and environment variables have priority over default values!**
- If the port for the IPC server is neither specified in flag `` nor in the environment variable `VPNCHAINS_IPC_SERVER_PORT` port 45454 will be used. Has to be specified explicitly for more than one vpnchains instance.
- If the path to the intercepting library is neither specified in flag `-lib-path` nor in the environment variable `VPNCHAINS_INJECT_LIB_PATH` path `/usr/lib/libvpnchains_inject.so` will be used.
- If the default size of the buffer used for reading from sockets is neither specified in flag `-buf` nor in the environment variable `VPNCHAINS_BUF_SIZE` an amount of 65536 will be used.
- If the default mtu for the wireguard tunnel is neither specified in flag `-mtu` nor in the environment variable `VPNCHAINS_MTU` an amount of 1420 will be used.
- If the default path to the wireguard config is not specified in flag `-config` than relative path `wg0.conf` will be used.

## Supported OS
Arch Linux, kernel 6.4.7-arch1-1; probably Ubuntu 22.04; everything else (Linux and *unix) may have some installation issues, but should work

## Apps that are guaranteed to work
- wget
- curl
- traceroute –tcp
- pacman
- apt
- snap
- git
- **firefox**
- **chromium**
- nslookup -vc
- dig +tcp (???)
- probably all other apps using TCP connections

## Apps that won't work
- quic (because of UDP)
- nslookup (because of UDP)
- anything else non-TCP
