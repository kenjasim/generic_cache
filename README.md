# generic_cache
Implementaton of an In memory cache using Go Generics

## Example

```go
cache := NewGenericCache[*TestCacheValue](10 * time.Second)
cv := TestCacheValue{}
cache.Set("test", cv)
cachedValue, err := cache.Get("test")
if err != nil{
    log.Fatal(err)
}
```