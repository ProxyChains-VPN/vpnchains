#pragma once

#define IPC_SOCK_PATH "/tmp/vpnchains.socket"
#define IPC_PORT_ENV_VAR "VPNCHAINS_IPC_SERVER_PORT"
#define SO_EXPORT __attribute__((visibility("default")))