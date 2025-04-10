package gotchas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func compute() int { return 42 }

// Gotcha 1: Deferred Function Arguments Are Evaluated Immediately
func TestDeferArgs(t *testing.T) {
	var b bytes.Buffer
	func() {
		defer fmt.Fprintf(&b, "result: %d", compute())
	}()
	if b.String() != "result: 42" {
		t.Errorf("got %s", b.String())
	}
}

// Gotcha 2: Loop Variable Capture in Goroutines
func TestLoopVariableCapture(t *testing.T) {
	got := make(chan int, 3)
	for i := range 3 {
		go func(i int) {
			got <- i
		}(i)
	}
	m := map[int]bool{}
	for range 3 {
		m[<-got] = true
	}
	if len(m) != 3 {
		t.Error("loop variable captured incorrectly")
	}
}

// Gotcha 3: Nil Interfaces Are Not Always Nil
type MyError struct{}

func (e *MyError) Error() string { return "fail" }

func TestNilInterface(t *testing.T) {
	var err error = (*MyError)(nil)
	if err == nil {
		t.Error("expected non-nil interface")
	}
}

// Gotcha 4: Slice Capacity Sharing
func TestSliceSharing(t *testing.T) {
	a := []int{1, 2, 3, 4}
	b := a[:2]
	b[1] = 99
	if a[1] != 99 {
		t.Error("expected slice to share backing array")
	}
}

// Gotcha 5: Recover Only Works in Same Goroutine
func TestRecoverGoroutine(t *testing.T) {
	done := make(chan bool)
	go func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic to be recovered inside same goroutine")
			}
			done <- true
		}()
		panic("fail")
	}()
	<-done
}

// Gotcha 6: Map Concurrent Write Panic
func TestMapConcurrency(t *testing.T) {
	// intentionally skipped: would crash the test suite
}

// Gotcha 7: Unexported Fields Not Marshalled
type Data struct {
	name string
	Age  int
}

func TestUnexportedMarshal(t *testing.T) {
	d := Data{name: "hidden", Age: 30}
	b, _ := json.Marshal(d)
	if strings.Contains(string(b), "hidden") {
		t.Error("unexported field should not be marshalled")
	}
}

// Gotcha 8: select{} Blocks Forever
func TestSelectBlock(t *testing.T) {
	// intentionally skipped: would hang the test suite
}

// Gotcha 9: Slice Append May Reallocate
func TestSliceAliasing_Reallocation(t *testing.T) {
	s := make([]int, 2, 4) // capacity is 4
	s[0] = 1
	s[1] = 2

	s2 := s[:1]      // alias of same backing array
	s = append(s, 3) // no reallocation yet
	s2[0] = 99       // should affect s

	if s[0] != 99 {
		t.Error("expected s2[0] to alias s[0] before reallocation")
	}

	s = append(s, 4) // still ok, capacity not exceeded
	s = append(s, 5) // now reallocation happens
	s2[0] = 100      // won't affect s anymore

	if s[0] == 100 {
		t.Error("expected reallocation to break aliasing")
	}
}

// Gotcha 10: Variable Shadowing
func someFunc() int { return 1 }

func TestShadowing(t *testing.T) {
	x := 10
	if x := someFunc(); x != 1 {
		t.Error("shadowed variable incorrect")
	}
	if x != 10 {
		t.Error("outer variable should not change")
	}
}

// Gotcha 11: Interface Nil vs Concrete Nil in JSON Marshalling
type Person struct {
	Details interface{} `json:"details,omitempty"`
}

func TestInterfaceNilMarshal(t *testing.T) {
	p := Person{}
	b, _ := json.Marshal(p)
	if !strings.Contains(string(b), "details") {
		t.Log("omitempty excluded interface field")
	}
}

// Gotcha 12: Closing a Closed Channel Panics
func TestCloseClosedChannel(t *testing.T) {
	ch := make(chan int)
	close(ch)
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on double close")
		}
	}()
	close(ch)
}

// Gotcha 13: Iterating a Map is Random
func TestMapIterationRandom(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	// no assert: iteration order is intentionally not stable
	t.Log("map iteration order is randomized:", keys)
}

// Gotcha 14: Goroutine Leak from Blocking Channel
func TestBlockingSendLeak(t *testing.T) {
	ch := make(chan int)
	go func() {
		ch <- 42 // blocked forever if no receiver
	}()
	select {
	case <-time.After(10 * time.Millisecond):
		t.Log("goroutine likely blocked (no receiver)")
	}
}

// Gotcha 15: Pointer Receiver vs Value Receiver Interface Satisfaction
type Doer interface {
	Do()
}

type Thing struct{}

func (t Thing) Do() {}

type Thing2 struct{}

func (t *Thing2) Do() {}

func TestReceiverMismatch(t *testing.T) {
	var _ Doer = Thing{}
	var _ Doer = &Thing2{} // OK
	// var _ Doer = Thing2{} // would not compile
}

// Gotcha 16: Struct Embedding Method Shadowing
type A struct{}

func (A) Hello() string { return "A" }

type B struct{ A }

func (B) Hello() string { return "B" }

func TestEmbeddedMethodOverride(t *testing.T) {
	b := B{}
	if b.Hello() != "B" {
		t.Error("method should be overridden")
	}
}

// Gotcha 17: defer in Loop Causes Resource Leak
func TestDeferInLoop(t *testing.T) {
	files := []*os.File{}
	for i := 0; i < 3; i++ {
		f, err := os.CreateTemp("", "test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name()) // executed at end, not per loop
		files = append(files, f)
	}
	for _, f := range files {
		if _, err := f.WriteString("ok"); err != nil {
			t.Error(err)
		}
	}
}

// Gotcha 18: Unbuffered Channels Can Deadlock
func TestUnbufferedChannelDeadlock(t *testing.T) {
	ch := make(chan int)
	go func() {
		ch <- 1
	}()
	<-ch // ok only because of paired goroutine
}

// Gotcha 19: JSON Numbers Are Floats by Default
func TestJSONNumbersAsFloat64(t *testing.T) {
	var m map[string]interface{}
	json.Unmarshal([]byte(`{"value": 42}`), &m)
	if _, ok := m["value"].(float64); !ok {
		t.Error("expected float64 from json number")
	}
}

// Gotcha 20: Map Key Comparability
type Point struct{ X, Y int }

func TestMapKeyComparability(t *testing.T) {
	m := map[Point]string{{1, 2}: "ok"}
	if m[Point{1, 2}] != "ok" {
		t.Error("comparable struct should work as map key")
	}
	// Uncommenting below would fail to compile
	// type Invalid struct { V []int }
	// _ = map[Invalid]string{}
}

// Gotcha 21: Panics in Goroutines Without Recovery Crash the Whole Program
func TestPanicInGoroutine(t *testing.T) {
	done := make(chan bool)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- true
			}
		}()
		panic("crash")
	}()
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Error("goroutine panic not recovered")
	}
}

// Gotcha 22: defer Does Not Run After os.Exit
func TestDeferWithOsExit(t *testing.T) {
	if os.Getenv("TEST_EXIT") == "1" {
		defer t.Log("this will not be printed")
		os.Exit(0)
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestDeferWithOsExit", "-test.v")
	cmd.Env = append(os.Environ(), "TEST_EXIT=1") // ✅ assign the env

	out, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 0 {
			t.Errorf("unexpected error: %v", err)
		}
	}

	if bytes.Contains(out, []byte("this will not be printed")) {
		t.Error("defer should not run after os.Exit")
	}
}

// Gotcha 23: time.After Leaks if Used in Loops
func TestTimeAfterLeak(t *testing.T) {
	for range 10 {
		timer := time.NewTimer(time.Millisecond)

		// Clean usage: either Stop() before it fires, or safely read
		select {
		case <-timer.C:
			// timer fired
		case <-time.After(10 * time.Millisecond):
			// timer didn't fire yet, stop and drain if needed
			if !timer.Stop() {
				<-timer.C
			}
		}
	}
	t.Log("timers used and cleaned safely")
}

// Gotcha 24: context.WithCancel Must Be Deferred
func TestContextCancelLeak(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	select {
	case <-ctx.Done():
	case <-time.After(10 * time.Millisecond):
		t.Log("context active")
	}
}

// Gotcha 25: recover() Only Works Inside defer
func TestRecoverOutsideDefer(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic recovery inside defer")
		}
	}()
	panic("fail")
}

// Gotcha 26: Struct Tags are Strings and Not Checked
type TagStruct struct {
	ID string `json:"id" wrongTag`
}

func TestStructTag(t *testing.T) {
	tag := reflect.TypeOf(TagStruct{}).Field(0).Tag
	if tag.Get("wrongTag") != "" {
		t.Log("got wrong tag value:", tag)
	}
}

type Action interface {
	Do() string
}

type ThingOne struct{}        // implements with value receiver
func (t ThingOne) Do() string { return "ThingOne" }

type ThingTwo struct{}         // implements with pointer receiver
func (t *ThingTwo) Do() string { return "ThingTwo" }

func TestNewVsLiteralPointer(t *testing.T) {
	var a Action = new(ThingOne) // OK: *Thing1 has Do() via value receiver
	if a.Do() != "ThingOne" {
		t.Error("unexpected Do() on *ThingOne")
	}

	// Literal form
	var b Action = &ThingTwo{} // OK: pointer receiver needed
	if b.Do() != "ThingTwo" {
		t.Error("unexpected Do() on &ThingTwo")
	}

	// THIS would fail to compile:
	// var c Action = ThingTwo{} // ❌ doesn't implement Action (needs pointer)
}

// Gotcha 28: Generics Constraint Mismatch (Go 1.18+)
type Number interface {
	~int | ~float64
}

func Add[T Number](a, b T) T {
	return a + b
}

func TestGenericsConstraintMismatch(t *testing.T) {
	i := Add(1, 2)     // OK: both int
	f := Add(1.5, 2.5) // OK: both float64

	if i != 3 || f != 4.0 {
		t.Error("unexpected result from Add")
	}

	// This would fail to compile due to type mismatch:
	// _ = Add(1, 2.0)

	// Why? Because type inference can't resolve a common type between int and float64,
	// and the constraint doesn't allow mixed-type parameters.
}

// Gotcha 29: Reflection Requires Exported Fields
func TestReflectUnexported(t *testing.T) {
	type secret struct {
		value string
	}
	s := secret{"hidden"}
	v := reflect.ValueOf(s)
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when accessing unexported field via Interface()")
		}
	}()
	_ = v.Field(0).Interface()
}

// Gotcha 30: Goroutine + Shared Slice Write = Race
func TestSharedSliceWriteRace(t *testing.T) {
	s := make([]int, 5)
	var wg sync.WaitGroup
	for i := range s {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s[i] = i
		}(i)
	}
	wg.Wait()
	for i, v := range s {
		if v != i {
			t.Logf("potential race at index %d: got %d", i, v)
		}
	}
}

// Gotcha 31: Slicing Keeps Underlying Array in Memory
func TestSliceMemoryRetention(t *testing.T) {
	big := make([]byte, 1<<20)
	small := big[:1]
	_ = small
	t.Log("small slice retains large backing array")
}

// Gotcha 32 (Revisited): go run builds in-memory and may omit debug symbols
func TestGoRunDebugSymbolNote(t *testing.T) {
	// This is informational only. When using `go run`, binary is built in-memory.
	// Stack traces or runtime paths may be missing debug symbols or file paths.
	t.Log("Note: `go run` builds a temp binary in-memory. Use `go build` for full symbol output.")
}

// Gotcha 33: go test Runs in Temp Directory
func TestGoTestWorkingDirectory(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Current working directory: %s", cwd)

	tmp := os.TempDir()
	if strings.Contains(cwd, tmp) {
		t.Log("This test is running from a temp directory (rare)")
	} else {
		t.Log("This test is running from source directory (expected)")
	}
}

// Gotcha 34: http.Request.Body Can Only Be Read Once
func TestRequestBodyReadOnce(t *testing.T) {
	body := io.NopCloser(strings.NewReader("test"))
	req := &http.Request{Body: body}
	data1, _ := io.ReadAll(req.Body)
	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewBuffer(data1))
	data2, _ := io.ReadAll(req.Body)
	if string(data1) != string(data2) {
		t.Error("expected re-readable body")
	}
}

// Gotcha 35: http.ResponseWriter Shouldn't Be Written After Headers
func TestHttpWriteAfterHeader(t *testing.T) {
	rr := httptest.NewRecorder()
	http.Error(rr, "fail", http.StatusBadRequest)
	fmt.Fprintln(rr, "extra") // won't be included
	if !strings.Contains(rr.Body.String(), "fail") {
		t.Error("expected error message")
	}
}

// Gotcha 36: http.Server Timeouts Not Set
func TestHttpServerTimeouts(t *testing.T) {
	srv := &http.Server{}
	if srv.ReadTimeout == 0 || srv.WriteTimeout == 0 {
		t.Log("timeouts not set by default")
	}
}

// Gotcha 37: JSON omitempty on interface holding nil
func TestJSONOmitemptyNilInterface(t *testing.T) {
	type Wrapper struct {
		Value interface{} `json:"value,omitempty"`
	}
	w := Wrapper{Value: nil}
	data, _ := json.Marshal(w)
	if string(data) != "{}" {
		t.Error("expected omitempty to omit nil interface")
	}
}

// Gotcha 38: Embedded Anonymous Fields Must Be Exported
type inner struct {
	private string // unexported
	Public  string
}
type outer struct {
	inner
}

func TestUnexportedFieldInsideEmbeddedStruct(t *testing.T) {
	o := outer{
		inner: inner{
			private: "secret",
			Public:  "shown",
		},
	}
	data, _ := json.Marshal(o)
	jsonStr := string(data)

	if strings.Contains(jsonStr, "secret") {
		t.Error("unexported fields should not be marshalled")
	}
	if !strings.Contains(jsonStr, "shown") {
		t.Error("exported fields should be marshalled")
	}
}

// Gotcha 39: time.AfterFunc Must Be Stopped
func TestAfterFuncStop(t *testing.T) {
	ran := false
	timer := time.AfterFunc(10*time.Millisecond, func() {
		ran = true
	})
	timer.Stop()
	time.Sleep(20 * time.Millisecond)
	if ran {
		t.Error("timer should not have executed")
	}
}

// Gotcha 40: Time Zone Load Depends on OS
func TestTimeZoneLoading(t *testing.T) {
	_, err := time.LoadLocation("Europe/London")
	if err != nil {
		t.Log("timezone might be missing in OS image")
	}
}

// Gotcha 41: filepath.Join Drops Empty
func TestFilepathJoinEmpty(t *testing.T) {
	result := filepath.Join("a", "", "b")
	expected := filepath.Join("a", "b")

	if result != expected {
		t.Errorf("filepath.Join should collapse empty: got %v, want %v", result, expected)
	}
}

// Gotcha 42: os.RemoveAll is Recursive
func TestRemoveAllSafety(t *testing.T) {
	dir := t.TempDir()
	nested := filepath.Join(dir, "child")
	os.Mkdir(nested, 0755)
	if err := os.RemoveAll(dir); err != nil {
		t.Error("failed to remove directory")
	}
}

// Gotcha 43: context.TODO Silently Omits Cancellation
func TestContextTODO(t *testing.T) {
	ctx := context.TODO()
	select {
	case <-ctx.Done():
		t.Error("context.TODO should not be cancellable")
	case <-time.After(10 * time.Millisecond):
	}
}

// Gotcha 44: context.WithCancel Leak Without Cancel
func TestContextCancelLeakAvoided(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // proper cleanup
	<-ctx.Done()
}

// Gotcha 45: sync.WaitGroup.Add After Wait
func TestWaitGroupAddAfterWait(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	done := make(chan struct{})
	go func() {
		defer wg.Done()
		<-done
	}()
	go func() {
		wg.Wait()
		// wg.Add(1) would race here
		close(done)
	}()
	time.Sleep(20 * time.Millisecond)
}

// Gotcha 46: sync.Map.Range is Unordered
func TestSyncMapRange(t *testing.T) {
	var m sync.Map
	m.Store("a", 1)
	m.Store("b", 2)
	found := map[string]bool{}
	m.Range(func(k, v any) bool {
		found[k.(string)] = true
		return true
	})
	if len(found) != 2 {
		t.Error("unexpected sync.Map range behavior")
	}
}

// Gotcha 47: Type Inference Fails with Mismatched Types (Go 1.18+ Generics)
func Sum[T int | int64](a, b T) T {
	return a + b
}

func TestGenericSumInference(t *testing.T) {
	x := Sum[int](1, 2)
	if x != 3 {
		t.Error("expected sum to be 3")
	}

	// Mixing types (int and int64) would fail to compile:
	// _ = Sum(1, int64(2)) // Uncommenting will cause compile error due to type mismatch

	// Instead, both types must match
	y := Sum[int64](1, 2)
	if y != 3 {
		t.Error("expected int64 sum to be 3")
	}
}
