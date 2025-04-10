# Golang Gotchas

## Generic GO

### 1. Deferred Function Arguments Are Evaluated Immediately
```go
defer fmt.Println("result:", compute())
```

➡️ compute() runs immediately, not when fmt.Println is executed.

🔹 Gotcha: Deferred functions' arguments are evaluated when the defer is declared, not when it runs.

### 2. Loop Variable Capture in Goroutines
```go
for i := 0; i < 3; i++ {
    go func() { fmt.Println(i) }()
}
```

➡️ May print 3 3 3 instead of 0 1 2

🔹 Fix: Capture the variable explicitly:

```go
for i := 0; i < 3; i++ {
    go func(i int) { fmt.Println(i) }(i)
}
```

### 3. Nil Interfaces Are Not Always Nil

```go
var err error = (*MyError)(nil)
fmt.Println(err == nil) // false!
```

➡️ Interface is non-nil because its type is set, even if its value is nil.

🔹 Gotcha: A non-nil interface with a nil underlying value is not equal to nil.

### 4. Slice Capacity May Lead to Unexpected Sharing
```go
a := []int{1, 2, 3, 4}
b := a[:2]
c := a[2:]
b[1] = 99
fmt.Println(a) // [1 99 3 4]
```

➡️ Slices share backing arrays.

🔹 Gotcha: Slicing does not copy; it creates a view over the same array — changes affect all views.

### 5. Recover Only Works in the Same Goroutine
```go
go func() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
        }
    }()
    panic("fail") // will not be recovered!
}()
```

➡️ Recover is goroutine-local.

🔹 Gotcha: recover() only works within the same goroutine where panic occurred.

### 6. Map Is Not Safe for Concurrent Writes
```go
m := make(map[string]int)
go func() { m["key"] = 1 }()
go func() { m["key2"] = 2 }()
```

➡️ Runtime panic: "concurrent map writes"

🔹 Gotcha: Use sync.Map or mutex for concurrent access.

### 7. Unexported Struct Fields Don't Marshal to JSON
```go
type Data struct {
    name string // unexported
    Age  int
}
```

➡️ name is ignored by json.Marshal.

🔹 Gotcha: Only exported fields (capitalised) are marshalled.

### 8. Select with No Cases Blocks Forever
```go
select {} // blocks forever
```

➡️ Useful for waiting or simulating deadlock

🔹 Gotcha: A select {} with no case statements will never return — good or bad depending on intent.

### 9. Appending to a Slice Can Break Aliased Data
```go
s := []int{1, 2}
s2 := s[:1]
s = append(s, 3) // underlying array might change
s2[0] = 99
```

➡️ s2[0] may not affect s anymore if capacity is exceeded and new memory is allocated.

🔹 Gotcha: Appending can cause the slice to point to a new backing array.

### 10. Shadowing Variables in if, for, and range
```go
x := 10
if x := someFunc(); x != 0 {
    fmt.Println(x)
}
fmt.Println(x) // this x is still 10
```

➡️ The x inside the if is a different variable.

🔹 Gotcha: Shadowing can cause confusion and bugs — be careful with reused variable names.

### 11. Interface Nil vs Concrete Nil in JSON Marshalling
```go
type Person struct {
    Details interface{}
}
p := Person{}
b, _ := json.Marshal(p)
fmt.Println(string(b)) // {"Details":null}
```

➡️ Even though Details is nil, it's still set in the struct and marshals as null.

🔹 Gotcha: Interfaces are encoded even when holding nil, unless omitted explicitly (omitempty won't work unless the interface itself is nil).

### 12. Closing a Closed Channel Panics
```go
ch := make(chan int)
close(ch)
close(ch) // panic: close of closed channel
```

➡️ No safe way to check if a channel is closed before closing it.

🔹 Gotcha: Use design patterns (e.g., single closer goroutine) to control closure responsibility.

### 13. Iterating a Map is Random
```go
m := map[string]int{"a": 1, "b": 2, "c": 3}
for k := range m {
    fmt.Print(k, " ")
}
```

➡️ Output order is non-deterministic.

🔹 Gotcha: Since Go 1.12, map iteration is randomized to prevent reliance on ordering.

### 14. Go Routine Leaks from Blocking Channels
```go
func leaky() {
    ch := make(chan int)
    go func() {
        ch <- 42 // blocks forever if never received
    }()
}
```

➡️ The goroutine is stuck forever.

🔹 Gotcha: Always design channel sends/receives with backpressure or cancellation awareness.

### 15. Pointer Receiver vs Value Receiver Affects Interface Satisfaction
```go
type Doer interface { Do() }

type Thing struct{}
func (t Thing) Do() {} // value receiver

var _ Doer = Thing{}     // OK
var _ Doer = &Thing{}    // also OK

type Thing2 struct{}
func (t *Thing2) Do() {} // pointer receiver

var _ Doer = &Thing2{}   // OK
var _ Doer = Thing2{}    // compile error
```

➡️ Interfaces only match if the receiver method set matches.

🔹 Gotcha: Pointer vs value receivers influence whether a type satisfies an interface.

### 16. Struct Embedding Can Override Promoted Methods
```go
type A struct{}
func (A) Hello() { fmt.Println("A") }

type B struct{ A }
func (B) Hello() { fmt.Println("B") }

func main() {
    b := B{}
    b.Hello() // prints "B", not "A"
}
```
➡️ Embedded method is overridden by explicit method.

🔹 Gotcha: Promoted methods can be shadowed — intentional or not.

### 17. Defer in Loops Can Cause Resource Leaks
```go
files := []*os.File{f1, f2, f3}
for _, f := range files {
    defer f.Close() // runs at end of *main*, not end of loop
}
```

➡️ All files stay open until the end — not ideal for large batches.

🔹 Gotcha: Avoid deferring inside tight loops if managing resources.

### 18. Unbuffered Channels Can Deadlock on Send or Receive
```go
ch := make(chan int)
ch <- 1 // deadlocks if there's no concurrent receiver
```

➡️ Channels block without a peer.

🔹 Gotcha: Use buffered channels or coordinate goroutines.

### 19. JSON Numbers Are Floats by Default
```go
var m map[string]interface{}
json.Unmarshal([]byte(`{"value": 42}`), &m)
fmt.Printf("%T\n", m["value"]) // float64
```

➡️ Even whole numbers are float64.

🔹 Gotcha: Explicitly decode into typed structs for numeric accuracy.

### 20. Map Keys Are Compared by Value
```go
type Point struct{ X, Y int }
m := map[Point]string{{1, 2}: "A"}

p := Point{1, 2}
fmt.Println(m[p]) // OK

p2 := Point{1, 2}
fmt.Println(m[p2]) // Also OK

// But if you embed a slice, boom:
type Invalid struct{ V []int }
// m := map[Invalid]string{{[]int{1,2}}: "A"} // compile error!
```

➡️ Only comparable types (no slices, maps, functions) can be map keys.

🔹 Gotcha: Know the comparability rules — and be careful when composing types.

### 21. Panics in Goroutines Without Recovery Crash the Whole Program
```go
go func() {
    panic("oh no!") // crash!
}()
```
➡️ If not recovered inside the goroutine, the panic terminates the program.

🔹 Gotcha: Always recover() within the same goroutine to safely isolate failures.

### 22. defer Does Not Respect os.Exit()
```go
defer fmt.Println("won't run")
os.Exit(1)
```

➡️ Deferred functions are skipped when os.Exit() is called.

🔹 Gotcha: Cleanup logic in defer won’t run if you exit the process directly.

### 23. time.After Leaks if Not Used Properly
```go
select {
case <-ch:
case <-time.After(time.Second):
}
```
➡️ time.After creates a timer that isn't garbage-collected until it fires.

🔹 Gotcha: Use time.NewTimer() + Stop() when timing out in tight loops or large systems.

### 24. context.WithCancel Must Be Called With defer cancel()
```go
ctx, cancel := context.WithCancel(ctx)
// forget defer cancel() → leak
```

➡️ Not calling cancel() leads to context leaks.

🔹 Gotcha: Always clean up contexts you create — even if it feels unnecessary.

### 25. recover() Only Works in Deferred Functions
```go
if r := recover(); r != nil {
    fmt.Println("nope") // this won't recover
}
```
➡️ You must call recover() inside a deferred function.

🔹 Gotcha: recover() outside defer is a no-op.

### 26. Struct Tags Are Strings — No Type Safety
```go
type Foo struct {
    ID string `json:"id" wrongTag`
}
```
➡️ Invalid tag values silently ignored; reflect tools may misbehave.

🔹 Gotcha: Typos in tags don’t throw errors — be vigilant or use linters.

### 27. new(T) and &T{} Are Not Always the Same
```go
a := new(T)   // zero-allocated
b := &T{}     // literal initialisation
// subtle if T has methods with pointer/value receivers
```
➡️ Both return a pointer, but &T{} lets you initialise fields.

🔹 Gotcha: Mixing them in factory patterns may lead to subtle inconsistencies.

### 28. Generics: Constraint vs Implementation Confusion
```go
func Sum[T int | float64](a, b T) T {
    return a + b // won't compile: no guarantee that + is supported
}
```
➡️ You need to use constraints.Ordered or define your own.

🔹 Gotcha: The type constraint must guarantee the operation — Go doesn’t assume operators exist.

### 29. Reflection Requires Exported Fields
```go
type secret struct {
    hidden string
}

val := reflect.ValueOf(secret{"shh"})
fmt.Println(val.Field(0).Interface()) // panic
```
➡️ Accessing unexported fields via reflect.Value.Interface() panics.

🔹 Gotcha: Reflection only works fully on exported fields unless you use unsafe.

### 30. Goroutines + Loop + Shared Variables — Redux
Even if you pass the loop var correctly, shared results may cause issues:

```go
var wg sync.WaitGroup
results := make([]int, 5)
for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(i int) {
        defer wg.Done()
        results[i] = compute(i) // races with others!
    }(i)
}
```
➡️ Even with correct indexing, without locks the write to results[i] is not synchronized.

🔹 Gotcha: Writes to shared slice/map elements require synchronisation, even with unique indexes.

### 31. Memory Not Released on Sliced Arrays
```go
big := make([]byte, 1<<20) // 1MB
view := big[:1]            // keeps the whole array in memory!
```
🔹 Gotcha: Slicing a large array and keeping the small slice holds the full memory — consider copying (copy()).

### 32. go run Builds In-Memory, Missing Debug Symbols

If you `go run` vs `go build`, stack traces may be less useful or paths less clear.

### 33. `go test` Runs in Temp Directory

`go test` compiles and runs your code in a temporary location, but it executes with the same working directory — so relative paths work, but compiled artifacts and stack traces may refer to temp locations.

### 34. unsafe.Pointer Corruption

Even though not idiomatic, when using unsafe, you break all type safety. A classic mistake:

```go
i := 123
p := unsafe.Pointer(&i)
f := *(*float64)(p) // totally invalid
```

➡️ Produces garbage or panics at runtime.

🔹 Gotcha: unsafe bypasses Go's memory safety — use only when absolutely required and with full awareness.

## 🌐 net/http
### 34. http.Request.Body Can Only Be Read Once
```go
body, _ := io.ReadAll(r.Body)
r.Body.Close()
bodyAgain, _ := io.ReadAll(r.Body) // nil, already consumed!
```
🔹 Gotcha: You need to buffer and replace r.Body if you want to reuse it.

✅ Solution:

```go
b, _ := io.ReadAll(r.Body)
r.Body = io.NopCloser(bytes.NewBuffer(b))
```

### 35. http.ResponseWriter Must Be Written to
```go
func handler(w http.ResponseWriter, r *http.Request) {
    http.Error(w, "error", http.StatusBadRequest)
    return
    fmt.Fprintln(w, "success") // still writes! too late
}
```
🔹 Gotcha: Writing after headers are sent leads to mixed/undefined response states. Prefer early returns.

### 36. http.Server Timeouts Not Set By Default
```go
srv := &http.Server{Addr: ":8080"} // no Read/Write/IdleTimeout
```
🔹 Gotcha: No timeout means connections can hang forever. Always set these in production.

✅ Solution:

```go
srv.ReadTimeout = 5 * time.Second
srv.WriteTimeout = 10 * time.Second
srv.IdleTimeout = 120 * time.Second
```

## 🧾 encoding/json

### 37. omitempty Skips Zero Values, But Not nil Interfaces
```go
type Foo struct {
    Value interface{} `json:"value,omitempty"`
}
f := Foo{Value: nil}
json.Marshal(f) // {"value":null} ← NOT omitted
```
🔹 Gotcha: omitempty only omits true nil, not nil wrapped in an interface.

### 38. Embedded Anonymous Fields Must Be Exported to Marshal
```go
type embedded struct {
    Secret string
}
type parent struct {
    embedded // exported
}
```
🔹 Gotcha: in Go, if an embedded field is unexported, it is not treated as an embedded/promoted field, but it is still a regular field of the struct, and its own exported fields are still accessible for marshaling.

## ⏰ time

### 39. time.AfterFunc Must Be Stopped to Prevent Leaks
```go
t := time.AfterFunc(10*time.Second, callback)
t.Stop() // important!
```
🔹 Gotcha: Timers must be stopped or they’ll leak goroutines if never triggered.

### 40. Time Zones Are OS Dependent
```go
loc, _ := time.LoadLocation("Europe/London") // may fail in containers
```
🔹 Gotcha: Time zone data comes from the OS — Alpine images or minimal containers might not include it.

✅ Solution: Use static time zones (time.FixedZone) or bundle tzdata.

## 📦 filepath & os
### 41. filepath.Join() Removes Empty Elements
```go
fmt.Println(filepath.Join("foo", "", "bar")) // "foo/bar"
```
🔹 Gotcha: Join() removes empty segments, unlike some path libraries in other languages.

### 42. os.RemoveAll() is Recursive
```go
os.RemoveAll("/my/data") // dangerous
🔹 Gotcha: It deletes the entire tree — always double check what you're deleting.
```
## ⏳ context
### 43. Passing context.TODO() Silently Suppresses Intent
```go 
db.QueryContext(context.TODO(), "...") // runs forever if context is unused
```
🔹 Gotcha: context.TODO() is a placeholder — using it in production disables cancellation and timeout propagation.

✅ Use context.WithTimeout, WithDeadline, etc.

### 44. Forgetting defer cancel() Leaks Goroutines
```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// forgot: defer cancel()
```
🔹 Gotcha: Resources associated with the context (timers, goroutines) won’t be cleaned up.

🔀 sync
### 45. sync.WaitGroup.Add() Must Be Called Before Wait()
```go
var wg sync.WaitGroup
go func() {
    wg.Add(1) // RACE with wg.Wait()
    defer wg.Done()
    work()
}()
wg.Wait() // may see 0 before Add
```
🔹 Gotcha: Calling Add() after Wait() starts can cause a race.

✅ Always Add() in the same goroutine before starting new goroutines.

### 46. sync.Map.Range is Not Ordered
```go
m := sync.Map{}
m.Store("a", 1)
m.Store("b", 2)
m.Range(func(k, v any) bool {
    fmt.Println(k, v) // unpredictable order
    return true
})
```
🔹 Gotcha: If you rely on ordering, sort the keys manually after collecting them.

🔁 Generics
### 47. Type Inference Can Fail Unexpectedly
```go
func Sum[T Number](a, b T) T { return a + b }
Sum(1, 2) // ok
Sum(1, int64(2)) // compile error!
```
🔹 Gotcha: Generics require matching types exactly — mixing int and int64 breaks inference.

