
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <poll.h>

#define PORT 8080
#define MAX_CLIENTS 10

int main() {
    int server_fd, client_socket;
    struct sockaddr_in address;
    int addrlen = sizeof(address);
    char buffer[1024] = {0};

    struct pollfd fds[MAX_CLIENTS + 1];
    int nfds = 1;  // 当前文件描述符数量，包含服务器套接字

    server_fd = socket(AF_INET, SOCK_STREAM, 0);
    address.sin_family = AF_INET;
    address.sin_addr.s_addr = INADDR_ANY;
    address.sin_port = htons(PORT);

    bind(server_fd, (struct sockaddr*)&address, sizeof(address));
    listen(server_fd, 3);

    // 设置服务器套接字为非阻塞
    fcntl(server_fd, F_SETFL, O_NONBLOCK);

    printf("Server listening on port %d (using poll)\n", PORT);

    fds[0].fd = server_fd;
    fds[0].events = POLLIN;

    while (1) {
        int poll_count = poll(fds, nfds, -1);  // 阻塞直到某些套接字就绪
        if (poll_count < 0) {
            perror("poll error");
            continue;
        }

        // 处理新的连接
        if (fds[0].revents & POLLIN) {
            client_socket = accept(server_fd, (struct sockaddr*)&address, (socklen_t*)&addrlen);
            if (client_socket < 0) {
                perror("accept failed");
                continue;
            }

            printf("New connection accepted\n");

            // 将新的客户端套接字添加到 fds 中
            fds[nfds].fd = client_socket;
            fds[nfds].events = POLLIN;
            nfds++;
        }

        // 处理已连接客户端的数据
        for (int i = 1; i < nfds; i++) {
            if (fds[i].revents & POLLIN) {
                int bytes_read = read(fds[i].fd, buffer, sizeof(buffer));
                if (bytes_read > 0) {
                    printf("Received: %s\n", buffer);
                    send(fds[i].fd, buffer, bytes_read, 0);  // 回显数据
                } else if (bytes_read == 0) {
                    printf("Client disconnected\n");
                    close(fds[i].fd);
                    fds[i] = fds[nfds - 1];  // 用最后一个客户端替换当前客户端
                    nfds--;  // 减少文件描述符数量
                } else {
                    perror("Error reading from socket");
                    close(fds[i].fd);
                    fds[i] = fds[nfds - 1];
                    nfds--;
                }
            }
        }
    }

    close(server_fd);
    return 0;
}
