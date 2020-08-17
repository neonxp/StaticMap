# StaticMap

Генератор статической картинки-карты для заданных географических координат и масштаба.

Может использоваться как сервер, так и отдельная библиотека.

```go
import "github.com/neonxp/StaticMap/pkg/static"
...
img, err := static.GetMapImage(lat, lon, zoom, width, height)
```
