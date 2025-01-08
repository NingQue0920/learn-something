
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/select.h>

#define PORT 8080
#define MAX_CLIENTS 10

int main() {
    int server_fd, client_socket, max_sd, activity;
    struct sockaddr_in address;
    int addrlen = sizeof(address);
    char buffer[1024] = {0};
    fd_set readfds;
    int clients[MAX_CLIENTS] = {0};

    server_fd = socket(AF_INET, SOCK_STREAM, 0);
    address.sin_family = AF_INET;
    address.sin_addr.s_addr = INADDR_ANY;
    address.sin_port = htons(PORT);

    bind(server_fd, (struct sockaddr*)&address, sizeof(address));
    listen(server_fd, 3);

    // 设置服务器套接字为非阻塞模式
    fcntl(server_fd, F_SETFL, O_NONBLOCK);

    printf("Server listening on port %d (using select)\n", PORT);

    while (1) {
        FD_ZERO(&readfds);  // 清空文件描述符集合
        FD_SET(server_fd, &readfds);  // 将监听套接字添加到集合中
        max_sd = server_fd;

        // 将所有客户端套接字添加到集合中
        for (int i = 0; i < MAX_CLIENTS; i++) {
            int sd = clients[i];
            if (sd > 0) {
                FD_SET(sd, &readfds);  // 添加到集合中
            }
            if (sd > max_sd) {
                max_sd = sd;  // 更新最大文件描述符
            }
        }

        // 使用 select 检查哪些套接字有事件
        activity = select(max_sd + 1, &readfds, NULL, NULL, NULL);

        if (activity < 0) {
            perror("select error");
            continue;
        }

        // 检查是否有新的连接
        if (FD_ISSET(server_fd, &readfds)) {
            client_socket = accept(server_fd, (struct sockaddr*)&address, (socklen_t*)&addrlen);
            if (client_socket < 0) {
                perror("accept failed");
                continue;
            }

            printf("New connection accepted\n");

            // 将新的客户端添加到客户端数组中
            for (int i = 0; i < MAX_CLIENTS; i++) {
                if (clients[i] == 0) {
                    clients[i] = client_socket;
                    break;
                }
            }
        }

        // 检查客户端套接字是否有数据可读
        for (int i = 0; i < MAX_CLIENTS; i++) {
            int sd = clients[i];
            if (FD_ISSET(sd, &readfds)) {
                int bytes_read = read(sd, buffer, sizeof(buffer));
                if (bytes_read > 0) {
                    printf("Received: %s\n", buffer);
                    send(sd, buffer, bytes_read, 0);  // 回显数据
                } else if (bytes_read == 0) {
                    printf("Client disconnected\n");
                    close(sd);
                    clients[i] = 0;
                } else {
                    perror("Error reading from socket");
                    close(sd);
                    clients[i] = 0;
                }
            }
        }
    }

    close(server_fd);
    return 0;
}
