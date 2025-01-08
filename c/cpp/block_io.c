
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/wait.h>

#define PORT 8080

void handle_client(int client_socket) {
    char buffer[1024] = {0};
    int bytes_read;

    while ((bytes_read = read(client_socket, buffer, sizeof(buffer))) > 0) {
        printf("Received: %s\n", buffer);
        send(client_socket, buffer, bytes_read, 0);  // 回显数据
        memset(buffer, 0, sizeof(buffer));
    }
    close(client_socket);
    printf("Client disconnected\n");
}

int main() {
    int server_fd, client_socket;
    struct sockaddr_in address;
    int addrlen = sizeof(address);

    server_fd = socket(AF_INET, SOCK_STREAM, 0);
    address.sin_family = AF_INET;
    address.sin_addr.s_addr = INADDR_ANY;
    address.sin_port = htons(PORT);

    bind(server_fd, (struct sockaddr*)&address, sizeof(address));
    listen(server_fd, 5);

    printf("Server listening on port %d\n", PORT);

    while (1) {
        client_socket = accept(server_fd, (struct sockaddr*)&address, (socklen_t*)&addrlen);
        if (client_socket < 0) {
            perror("Accept failed");
            continue;
        }

        printf("New connection accepted\n");

        // 使用 fork 为每个客户端创建一个新进程
        pid_t pid = fork();
        if (pid == 0) {  // 子进程
            close(server_fd);
            handle_client(client_socket);
            exit(0);
        } else if (pid > 0) {  // 父进程
            close(client_socket);
        } else {
            perror("Fork failed");
        }
    }
    close(server_fd);
    return 0;
}
