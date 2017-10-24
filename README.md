LBSENGINE
-----

## What

LBSENGINE实现了一个通用的GEO索引引擎,[2017Go基金会中国黑客马拉松](http://gohack2017.golangfoundation.org/)参赛项目，获得一等奖。

## How

项目成于2017Go基金会中国黑客马拉松，演示PPT见[lbsengine.pdf](https://github.com/sillydong/lbsengine/blob/master/lbsengine.pdf)

开发可参考`example`中的代码实现

## Who

- [sillydong](https://github.com/sillydong)

- [3zhen](https://github.com/3zheng)

## Todo

1. 维护一份内存索引，启动时从永久存储初始化内存索引
2. 优化序列化/反序列化方法
3. 参考Redis优化geohash算法
4. ...

## Thanks

1. [redis](http://github.com/go-redis/redis)
2. [geohash](http://github.com/mmcloughlin/geohash)
3. [go-cache](http://github.com/patrickmn/go-cache)
4. [murmur hash](https://github.com/huichen/murmur)
5. [gorp](https://github.com/go-gorp/gorp)
6. [echo](https://github.com/labstack/echo)
7. [go-bindata-assetfs](https://github.com/elazarl/go-bindata-assetfs)
8. [echo-static](https://github.com/Code-Hex/echo-static)
9. [msgp](https://github.com/tinylib/msgp)
10. [tidb](https://github.com/pingcap/tidb)
