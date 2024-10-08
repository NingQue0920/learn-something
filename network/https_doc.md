### Https加密过程
    第一次握手，用于协商并交换信息，以及两个随机数
    1. CLIENT_HELLO: 客户端发送支持的加密算法、支持的压缩算法、支持的扩展等信息给服务器
    2. SERVER_HELLO: 服务器选择加密算法、压缩算法、扩展等信息返回给客户端
---    
    第二次握手， 用于验证服务器身份，传输公钥，并完成唯一一次的非对称加密
    3. CERTIFICATE: 服务器发送证书给客户端，证书中包含很多信息，但最重要的是**公钥**、**证书签名**和**服务器地址**
        公钥用于后续的加密，证书签名用于验证公钥是否被篡改，服务器地址用于验证是否是目标机器。
        验签：客户端用CA的公钥解密证书签名，得到摘要，再用摘要和证书中的公钥对比，如果一致则证明公钥没有被篡改。
    4. CLIENT_KEY_EXCHANGE: 客户端验证证书，用服务器的公钥加密一个随机数，生成（Per-master Secret 预主密钥），发送给服务器
---
    对称密钥生成
    5. SERVER_KEY_EXCHANGE: 服务器用私钥解密客户端发送的随机数，结合第一步中双方交换的两个随机数生成（Master Secret 主密钥）。这个主密钥是对称加密的密钥，用于后续的对称加密
        同理，客户端也会根据三个随机数生成密钥，用于后续的对称加密。
---    
    客户端通知服务器，后续通信将使用对称加密
    6. CHANGE_CIPHER_SPEC: 客户端发送改变加密算法的通知给服务器，表示后续的通信将使用对称加密
    7. FINISHED: 客户端发送握手结束通知给服务器，表示握手结束。
---
    服务器通知客户端，后续通信将使用对称加密
    8. CHANGE_CIPHER_SPEC: 服务器发送改变加密算法的通知给客户端，表示后续的通信将使用对称加密
    9. FINISHED: 服务器发送握手结束通知给客户端，表示握手结束。
