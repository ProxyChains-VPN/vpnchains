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
    bind(fd, (struct sockaddr*)&name, sizeof(name));
    char msg[10];
    struct sockaddr_in from;
    int from_len = sizeof(from);
    while(1){
	recvfrom(fd, msg, 10, 0, (struct sockaddr*)&from, &from_len);
	int ip = from.sin_addr.s_addr;
	printf("received from %s: %s\n", inet_ntoa(from.sin_addr), msg);
    }

    return 0;
}