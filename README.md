# Rate limiter

This is a naive implementation of [token bucket](https://en.wikipedia.org/wiki/Token_bucket) rate limiter

# Installation

`go get github.com/AvoidMe/ratelimiter`

# Example usage

```go
import "github.com/AvoidMe/ratelimiter"

// Init new rate limiter
limiter := ratelimiter.NewLimiter(10, 1*time.Second)
go limiter.Start()
defer limiter.Stop()

// Check if client can make request
if limiter.Get() {
    // client allowed to make request
}

// Check if client can make 10 requests (for batch usage)
if limiter.GetN(10) {
    // client allowed to make 10 requests
}
```
