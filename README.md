# Flow

A collection of async/flow control functions

## Usage

### Eventually

Get the result of a function call later on:

```go
res := flow.Eventually(context.Background(), func(ctx context.Context) (int, error) {
    return 5, nil
})

...

<-res.Done()
if res.Err() != nil {
    panic(res.Err())
}

fmt.Println(res.Out()) // prints 5
```

#### Groups

If you need to wait for multiple results to resolve:

```go

fast := flow.Eventually(context.Background(), func(ctx context.Context) (int, error) {
    return 5, nil
})
slow := flow.Eventually(context.Background(), func(ctx context.Context) (string, error) {
    time.Sleep(time.Second * 5)
    return "bongo", nil
})

group := &flow.ResultGroup{}
group.Add(fast)
group.Add(slow)

// The wait function blocks until all results have resolved
group.Wait()

fmt.Println(fast.Out()) // prints 5
fmt.Println(slow.Out()) // prints bongo
```

### Retry

To retry a function a 3 times:

```go
i := 0
f := flow.Retry(func(ctx context.Context) (int, error) {
    defer func() { i++ }()
    if i < 2 {
        return 0, errors.New("demo error")
    }
    return 5, nil
}, 3)

// The above will call the function three times, and will return 5, nil on the third call
out, err := f(context.Background()) // err is nil
fmt.Println(out) // prints 5
```

To retyr a function 3 times with a millisecond delay between each try:

```go
i := 0
f := flow.RetryDelay(func(ctx context.Context) (int, error) {
    defer func() { i++ }()
    if i < 2 {
        return 0, errors.New("demo error")
    }
    return 5, nil
}, 3, time.Millisecond)

out, err := f(context.Background()) // err is nil
fmt.Println(out) // prints 5
```

### Throttle

To throttle a function call so that it runs once per second:

```go
i := 1
f := flow.Throttle(func(ctx context.Context) (int, error) {
    defer func() { i++ }()
    return i, nil
}, time.Second)

for range 3 {
    fmt.Println(f(context.Background()))
}
```

Which will print:

```
1 <nil>
0 throttled
0 throttled
```

To throttle a function call so that it runs once per second, and returns the first value without a throttled error:

```go
i := 1
f := flow.SilentThrottle(func(ctx context.Context) (int, error) {
    defer func() { i++ }()
    return i, nil
}, time.Millisecond)

for i := range 3 {
    if i == 2 {
        time.Sleep(time.Millisecond * 2)
    }
    fmt.Println(f(context.Background()))
    i++
}
```

Which will print:

```
1 <nil>
1 <nil>
2 <nil>
```
