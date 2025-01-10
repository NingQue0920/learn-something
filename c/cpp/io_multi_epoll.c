// Epoll-based server implementation
void run_epoll_server() {
    int server_fd = create_and_bind();
    make_socket_non_blocking(server_fd);
    char buffer[BUFFER_SIZE];
    // every fd has a wait queue , when process is block , will be added into this queue .
    int epoll_fd = epoll_create1(0);
    if (epoll_fd == -1) {
        perror("epoll_create1");
        exit(EXIT_FAILURE);
    }
    
    struct epoll_event event, events[MAX_EVENTS];
    event.events = EPOLLIN;
    event.data.fd = server_fd;
    //register fd[add fd to the read-black-tree] , register callback 
    if (epoll_ctl(epoll_fd, EPOLL_CTL_ADD, server_fd, &event) == -1) {
        perror("epoll_ctl");
        exit(EXIT_FAILURE);
    }
    
    printf("Epoll server starting on port %d...\n", PORT);
    
    while (1) {
        // When the network card recevices data , it calls a hardware interrupt.
        // kernel: 
        //        - Handle data packet .
        //        - Call the callback function. 
        //        - The calback function adds the fd to rdllist.
        // epoll_wait : 
        //        - Check if the rdllist is empty , if not , return all fds in rdllist ,  else blocks . 
        //        - if return , main process will execute accept(for new connection) , read(handle new data) .
    //            - During main process execute other command , others socket may recevice data and added to rdllist.
        int n = epoll_wait(epoll_fd, events, MAX_EVENTS, -1);
        
        for (int i = 0; i < n; i++) {
            if (events[i].data.fd == server_fd) {
                // Handle new connection
                struct sockaddr_in client_addr;
                socklen_t client_len = sizeof(client_addr);
                int client_fd = accept(server_fd, (struct sockaddr *)&client_addr, &client_len);
                
                if (client_fd == -1) {
                    perror("accept");
                    continue;
                }
                
                make_socket_non_blocking(client_fd);
                event.events = EPOLLIN | EPOLLET;  // Edge-triggered
                event.data.fd = client_fd;
                
                if (epoll_ctl(epoll_fd, EPOLL_CTL_ADD, client_fd, &event) == -1) {
                    perror("epoll_ctl");
                    close(client_fd);
                    continue;
                }
                
                printf("New connection from %s on socket %d\n", 
                       inet_ntoa(client_addr.sin_addr), client_fd);
            } else {
                // Handle data from client
                int fd = events[i].data.fd;
                int bytes_read = read(fd, buffer, sizeof(buffer) - 1);
                
                if (bytes_read <= 0) {
                    if (bytes_read < 0) {
                        perror("read");
                    }
                    printf("Connection on socket %d closed\n", fd);
                    epoll_ctl(epoll_fd, EPOLL_CTL_DEL, fd, NULL);
                    close(fd);
                } else {
                    buffer[bytes_read] = '\0';
                    printf("Received from socket %d: %s", fd, buffer);
                    // Echo back
                    write(fd, buffer, bytes_read);
                }
            }
        }
    }
}
