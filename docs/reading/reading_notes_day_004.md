# Redis for Persistence

## Redis on Docker
- Start the redis server using the redis image
```
docker image pull redis

docker run -d --name redis -p 6379:6379 redis:<version>
```

- Run the redis-cli via docker to connect to the redis instance
```
docker exec -it redis redis-cli
```
- Run the cli is preinstalled via binary
```
redis-cli -h 127.0.0.1 -p 6379
```

- Docker uses internal config files for Redis.
    - can use redis with a local config file via:
    - add the redis.conf file to /data/ into the container via Dockerfile
    ```
    FROM redis
    COPY redis.conf /path/to/local/redis.conf
    CMD ["redis-server", "/usr/local/etc/redis/redis.conf"]
    ```
    - can also do this via cli options
    ```
    docker run -v /myredis/conf:/usr/local/etc/redis --name myredis redis redis-server /usr/local/etc/redis/redis.conf
    ```
    - this mounts a volume from local direcotry
    - may rewrite exisitng config files

## Go Redis

go get github.com/redis/go-redis/v9

## Connecting to Redis Server
- Create new connection via TLS
redis.NewClient(&redis.Options{

        Addr: "localhost:6379, 
        Password: "",
        DB: 0, //default db
        TLSConfig: &tls.Config{
            MinVersion:  tls.VersionTLS12,
            //Ceritficates: []tls.Certificate{cert},
        },
    })

- Connect via url
redis.ParseURL("redis://<user>:<pass>@localhost:6379/<db>")

- Connect via SSH
ssh.ClientConfig{
    User: "root",
    Auth: []ssh.AuthMethod{ssh.Password("password")},
    HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    Timeout: 15 * time.Second,
}

ssh.Dial("tcp", "remoteIP:22", sshConfig)

redis.NewClient(redis.Options{
    Addr: net.JoinHostPort("127.0.0.1", "6379")
    Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {

        return sshClient.Dial(network, addr)
    },

    ReadTimeout: -1,
    WriteTimeout: -1,
}

## Executing Commands
val, err := rdb.Get(ctx, "Get").Result()

- execute custom commands
rdb.Do(ctx, "get", "key").Result()

- redis.Nil to check if server responded with nil

- redis.Conn to represent a single connection to the redis server


## Redis HSET

- used to store objects in Redis
- HSET key field value [field value ...]
    - HSET myhash field1 "hello"
    - HGET myhash field1
        - "hello"

- In go
    - rdb.HSet(ctx, "myhash", "field1", "Hello").Result()
    - rdb.HGetAll(ctx, "myhash")
        - gives us a result of every single key val that was stored up until that point


