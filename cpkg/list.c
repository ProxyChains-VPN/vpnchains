#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>

typedef struct UdpPacket_s {
    struct sockaddr_in from;
    char *msg;
    ssize_t msg_len;
} UdpPacket;

typedef struct Node_s {
    struct Node_s *prev, *next;
    void *value; // передаёт адрес значения
} Node;

typedef Node List;

void init_list(List *list) {
    list->prev = list;
    list->next = list;
    list->value = NULL;
}

Node *init_node(Node *node, void* ptr, Node *new) {
    new->value = ptr;

    new->next->prev = new;
    new->prev->next = new;

    return new;
}

Node *add_after(Node *node, void *ptr) {
    Node *new = malloc(sizeof(Node));

    new->prev = node;
    new->next = node->next;

    return init_node(node, ptr, new);
}

Node *add_before(Node *node, void *ptr) {
    Node *new = malloc(sizeof(Node));

    new->next = node;
    new->prev = node->prev;

    return init_node(node, ptr, new);
}

void *erase(Node *node) {
    node->prev->next = node->next;
    node->next->prev = node->prev;

    void* value = node->value;

    free(node->value);
    node->value = NULL;

    free(node);
    node = NULL;

    return value;
}

