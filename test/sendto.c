#include <stdio.h>
#include <sys/socket.h>
#include <string.h>
#include <arpa/inet.h>
#include <stdlib.h>

int main(){

    int fd = socket(AF_INET, SOCK_DGRAM, 0);
    struct sockaddr_in name;
    memset(&name, 0, sizeof(struct sockaddr_in));
    name.sin_family = AF_INET;
    name.sin_port = 44444;
    inet_pton(AF_INET, "127.0.0.1", &name.sin_addr);
    char msg[] = "Hello\n";
    int ip = name.sin_addr.s_addr;
    printf("sendto ip: %s, port: %d\n", inet_ntoa(name.sin_addr), name.sin_port);
    printf("sendto res: %d\n", sendto(fd, msg, strlen(msg), 0, (struct sockaddr*)&name, sizeof(name)));

    return 0;
}