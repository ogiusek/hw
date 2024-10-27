# http wraps

httpwrap allows us to handle methods and parse them to http when allowing for customization.

we guarantee that httpwrap wont change

## example handle

```go

import "github.com/ogiusek/hw/src/hw"
type PingCommand struct {
  Count int
}

func Ping(command PingCommand) any {
  command.Count += 1
  return command
}

func AddPing(){
  http.handle("/ping", hw.Wrap(hw.Run(Ping), nil))
}
```

here we have example command which maps Ping command edits it and returns.

## responses

### nil
```go
func Ping(command CommandType) any {
  // ...
  return nil 
}
```

this command would return:
- 204

### custom response

```go
func Ping(command CommandType) any {
  w := hw.NewResponse()
  w.WriteHeader(418)
  return w
}
```
custom response is reccomended to use only in wraps (not reccomended in commands)


this command would return:
- 418

### struct
```go
func Ping(command CommandType) any {
  return command
}
```

this command would return:
- 200
- command in body as json

this is default encoding.
everything what is not nil and custom response will be parsed to json with default encoder.
default encoder can be found in:
- /src/hw/encoder.go

## wraps

wraps can be used as middleware.

```go
func Catch[T any](fn hw.Wrapper[T]) hw.Wrapper[T] {
	return func(args T, r *http.Request) any {
		res := fn(args, r)
		if err, ok := res.(error); ok {
			w := hw.NewResponse()
			// parse error to hw
			return w
		}
		return res
	}
}
```

this example wrap can be used like this to catch errors:

```go
func AddPing(){
  http.handle("/ping", hw.Wrap(Catch(hw.Run(Ping)), nil))
}
```

if we become overwhelmed by amount of our wraps we can have them in one place:

```go
// executable is our method type (func[T any](T) any)
func Run[T any](fn hw.Executable[T]) hw.Wrapper[T] {
	return Catch(Validate(hw.Run(fn)))
}

func AddPing(){
  // i use here hw.Wrap to be able to put here endpoint specific wraps (like authorization)
  http.handle("/ping", hw.Wrap(Run(Ping), nil))
}
```

## request decoding

we can have custom decoding of commnds.

```go
type ParserOptions[T any] struct {
	Decoder Decoder[T]
}
func AddPing(){
  // i use here hw.Wrap to be able to put here endpoint specific wraps (like authorization)
  http.handle("/ping", hw.Wrap(Run(Ping), &hw.ParserOptions[CommandType]{
    Decoder: func(r *http.Request) any { // should return CommandType, Resp or nil
      var args CommandType
      // decode from request
      return args
    },
  }))
}
```

eventualy we can have universal decoding

```go
func DefaultDecoder[T any](r *http.Request) any { // should return CommandType, Resp or nil
  var args T
  // create args
  return args
}
```

when decoder would return unexpected type app returns 500 and runs
```go
log.Printf("decoder returns unexpected type returned: %v", returned)
```

## response encoding

we have to use wraps for custom encoding and parse to `Resp`.
encoding allowable outputs are in 