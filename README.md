# dcard-rate-limiter

dcard web backend developer home test.

## Require

* Redis
* Zookeeper
* Go 1.16

## Install

可以使用 docker-compose 來啟動本服務，您可以前往 deployments/environment
底下中查看 [docker-compose.deploy](https://github.com/karta0898098/dcard-rate-limiter/blob/master/deployments/environment/docker-compose.deploy.yml "link")
的設定檔案。

 ```
 make docker.deploy
 ```

``推薦``也可使用 docker-compose 啟動所需要的環境配置，並且直接執行 go run 的方式執行速度較快，但是依賴 go 1.16 的 sdk。

```
make ratelimiter.dev.env
make ratelimiter.local
```

## Testing

啟動單元測試

```
make uint.testing
```

啟動整合測試

```
make integration.testing
```

手動測試入口點

```
GET http://localhost:18080/api/v1/protected
```

## 設計理念

本服務選擇使用 Redis slide window algorithm 的做法來實現 rate limit。 </br>
為何採用 Redis 來實作的原因如下 </br>
優點：

* 微服務中的最終一致性
* 存取速度速度快
* 實作快速

為了避免 race condition 的問題，所以選擇使用 [zookeeper](https://zookeeper.apache.org/ "link") 來做 lock 的動作，主要是覺得使用起來像原生的 mutex
可以直覺的開發因時間問題沒有實作 redis setnx lock 的方式，當然關於 redis lock 還有許多解法。 </br> e.g. :

* [redis setnx lock](https://redis.io/commands/setnx "link")
* [redisLock](https://redis.io/topics/distlock "link")

將 Redis RateLimit 設計成三層式的結構，讓他不單只是 middleware ，也可以快速拆成微服務的形式。 </br>
並且設計的接口提供了 URL 加上 IP 的形式讓每一隻 API 的資源可以分別控制，亦可全局管理。

關於選擇使用 [go-redis](https://github.com/go-redis/redis "link") 只是筆者較為熟悉，且封裝良好，社群的數量活躍並且在 redis 的官網[推薦清單](https://redis.io/clients#go "link")中 。


## TODO
* Zookeeper 移植到 infra package 裡面
* 實作 Redis setnx lock 並比較速度差異

## License
dcard-rate-limiter source code is available under an
MIT [License](https://github.com/karta0898098/dcard-rate-limiter/blob/master/LICENSE "link").

