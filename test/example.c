#include <stdlib.h>
#include <sys/socket.h>
#include <netdb.h>
#include <stdio.h>
#include <unistd.h>


int main() {
    unsigned short port;       /* port client will connect to         */
    char buf[12];              /* data buffer for sending & receiving */
    struct hostent *hostnm;    /* server host name information        */
    struct sockaddr_in server; /* server address                      */
    int s;                     /* client socket                       */

    hostnm = gethostbyname("lib.ru");
    port = 80;

    server.sin_family      = AF_INET;
    server.sin_port        = htons(port);
    server.sin_addr.s_addr = *((unsigned long *)hostnm->h_addr);
//
    if ((s = socket(AF_INET, SOCK_STREAM, 0)) < 0)
    {
        perror("Socket()");
        exit(3);
    }

    if (connect(s, (struct sockaddr *)&server, sizeof(server)) < 0)
    {
        perror("Connect()");
        exit(4);
    }
//    int sockfd = socket(AF_INET, SOCK_STREAM, 0);
//    write(sockfd, "sdf\n", 5);
}