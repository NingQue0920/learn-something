
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <arpa/inet.h>
#include <errno.h>

#define PORT 8080
#define MAX_CLIENTS 10

int set_nonblocking(int fd) {
    int flags = fcntl(fd, F_GETFL, 0);
    return fcntl(fd, F_SETFL, flags | O_NONBLOCK);
}

int main() {
    int server_fd, client_socket;
    struct sockaddr_in address;
    int addrlen = sizeof(address);
    char buffer[1024] = {0};
    int clients[MAX_CLIENTS] = {0};

    server_fd = socket(AF_INET, SOCK_STREAM, 0);
    address.sin_family = AF_INET;
    address.sin_addr.s_addr = INADDR_ANY;
    address.sin_port = htons(PORT);

    bind(server_fd, (struct sockaddr*)&address, sizeof(address));
    listen(server_fd, MAX_CLIENTS);

    set_nonblocking(server_fd);  // 设置服务器套接字为非阻塞模式
    printf("Server listening on port %d (non-blocking mode)\n", PORT);

    while (1) {
        // 尝试接受新连接
        client_socket = accept(server_fd, (struct sockaddr*)&address, (socklen_t*)&addrlen);
        if (client_socket >= 0) {
            // 设置为非阻塞，client_socket没有事件时，直接返回
            set_nonblocking(client_socket);
            for (int i = 0; i < MAX_CLIENTS; i++) {
                if (clients[i] == 0) {
                    clients[i] = client_socket;
                    printf("New client connected: %d\n", client_socket);
                    break;
                }
            }
        } else if (errno != EAGAIN && errno != EWOULDBLOCK) {
            perror("Accept failed");
        }

        // 轮询检查已连接的客户端是否有数据可读
        for (int i = 0; i < MAX_CLIENTS; i++) {
            if (clients[i] > 0) {
                // 设置为非阻塞之后，如果没数据，read会返回-1 
                int bytes_read = read(clients[i], buffer, sizeof(buffer));
                if (bytes_read > 0) {
                    printf("Received from client %d: %s\n", clients[i], buffer);
                    send(clients[i], buffer, bytes_read, 0);  // 回显数据
                } else if (bytes_read == 0) {
                    printf("Client %d disconnected\n", clients[i]);
                    close(clients[i]);
                    clients[i] = 0;
                }
            }
        }
        sleep(1);  // 控制 CPU 使用率
    }
    close(server_fd);
    return 0;
}
