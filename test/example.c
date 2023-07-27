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

    hostnm = gethostbyname("askdjbv;oasd");
    port = 40066;

    server.sin_family      = AF_INET;
    server.sin_port        = htons(port);
    server.sin_addr.s_addr = *((unsigned long *)hostnm->h_addr);
//
    if ((s = socket(AF_INET, SOCK_DGRAM, IPPROTO_UDP)) < 0)
    {
        perror("Socket()");
        exit(3);
    }

    if (connect(s, (struct sockaddr *)&server, sizeof(server)) < 0)
    {
        perror("Connect()");
        exit(4);
    }

    write(s, "GET lib.ru", 10);

    char buf1[1000];
    int res = read(s, buf1, 1000);
    write(2, buf1, res);


//    int sockfd = socket(AF_INET, SOCK_STREAM, 0);
//    write(sockfd, "sdf\n", 5);
//    char buf2[100] = "askjdfnlksadf";
//    int asd = read(0, buf2, 100);
//     write(1, buf2, asd);

//    s = socket(AF_INET, SOCK_STREAM, 0);
//    write(s, "hellosber\n", 11);
//
//    read(s, buf2, 9);
//    write(2, buf2, 6);
//    write(2, "\n", 2);
}