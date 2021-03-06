package degrade

import (
	"Asura/conf"
	blade "Asura/src"
	"Asura/src/ecode"
	store2 "Asura/app/middleware/degrade/store"

	"github.com/pkg/errors"
)

// This example create a cache middleware instance and two cache policy,
// then attach them to the specified path.
//
// The `PageCache` policy will attempt to cache the whole response by URI.
// It usually used to cache the common response.
//
// The `Degrader` policy usually used to prevent the API totaly unavailable if any disaster is happen.
// A succeeded response will be cached per 600s.
// The cache key is generated by specified args and its values.
// You can using file or memcache as cache backend for degradation currently.
//
// The `Cache` policy is used to work with multilevel HTTP caching architecture.
// This will cause client side response caching.
// We only support weak validator with `ETag` header currently.
func Example() {
	mc := store2.NewMemcache(&conf.Memcache{
		Host: "127.0.0.1",
		Port: "11211",
	})
	ca := New(mc)
	deg := NewDegrader(10)
	pc := NewPage(10)
	ctl := NewControl(10)
	filter := func(ctx *blade.Context) bool {
		if ctx.Request.Form.Get("cache") == "false" {
			return false
		}
		return true
	}

	engine := blade.Default()
	engine.GET("/users/profile", ca.Cache(deg.Args("name", "age"), nil), func(c *blade.Context) {
		values := c.Request.URL.Query()
		name := values.Get("name")
		age := values.Get("age")

		err := errors.New("error from others") // error from other call
		if err != nil {
			// mark this response should be degraded
			c.JSON(nil, ecode.Degrade)
			return
		}
		c.JSON(map[string]string{"name": name, "age": age}, nil)
	})
	engine.GET("/users/index", ca.Cache(pc, nil), func(c *blade.Context) {
		c.String(200, "%s", "Title: User")
	})
	engine.GET("/users/list", ca.Cache(ctl, filter), func(c *blade.Context) {
		c.JSON([]string{"user1", "user2", "user3"}, nil)
	})
	engine.Run(":18080")
}
