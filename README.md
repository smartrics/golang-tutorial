# GOLANG Tutorial

This is a step-by-step tutorial designed with the help of ChatGPT.  
It walks through the design and implementation of a simple banking system in Go, with each part building incrementally on top of the previous one.

It is designed to guide mid-to-senior engineers through building a real-world backend system from scratch — a banking platform — using idiomatic Go, clean architecture principles, and test-driven development.

The purpose is to help experienced Go developers refresh their knowledge through practical, progressive exercises — with a focus on writing production-quality code, structuring real-world projects, and revisiting concurrency, interfaces, and clean design.

## 🎯 What You'll Build

  * Core banking logic (accounts, transfers)
  * Clean interfaces and structs
  * Transfer pipelines with decorators and mixins
  * An async job processor (fan-in/fan-out, concurrency patterns)
  * A testable, composable transfer engine
  * A full HTTP API (routes, JSON I/O, handlers)
  * End-to-end integration tests
  * Monitoring/debug endpoints for observability

Each **part is implemented in a dedicated branch** named `partX` (e.g. `part1`, `part2`, ... `part8`).

---

## 🔁 Tutorial Structure

| Part | Branch   | Focus                                         |
|------|----------|-----------------------------------------------|
| 1    | [part1](https://github.com/smartrics/golang-tutorial/tree/part1) | Structs, interfaces, basic testing |
| 2    | [part2](https://github.com/smartrics/golang-tutorial/tree/part2) | Table-driven tests, test coverage |
| 3    | [part3](https://github.com/smartrics/golang-tutorial/tree/part3) | Mocking, interfaces, error handling |
| 4    | [part4](https://github.com/smartrics/golang-tutorial/tree/part4) | Functional pipelines, decorators |
| 5    | [part5](https://github.com/smartrics/golang-tutorial/tree/part5) | Transfer engine + integration tests |
| 6    | [part6](https://github.com/smartrics/golang-tutorial/tree/part6) | Async processor with concurrency |
| 7    | [part7](https://github.com/smartrics/golang-tutorial/tree/part7) | Transfer engine + processor coordination |
| 8    | [part8](https://github.com/smartrics/golang-tutorial/tree/part8) | HTTP server, monitoring, E2E testing |

> 💡 You can check out any part using `git checkout partX`.

---

## 🚀 Getting Started

Clone the repo and initialise your Go module:

```bash
git clone https://github.com/smartrics/golang-tutorial.git
cd golang-tutorial
git checkout part1   # or any part you want to start from

go mod tidy
```

*Note: Use your own module path if you're adapting the tutorial.*

## 🧱 Build Philosophy
  * Test-first (TDD) development
  * Realistic banking domain use case
  * Clean architecture + Hexagonal layering
  * Async processing, observability, and end-to-end testing

## Part 1

### 🎯 Goal:

 * Refresh Go syntax and control flow
 * Implement a practical Bank Account Simulator
 * Write unit tests using idiomatic Go

### ✅ Requirements

*Objective*: Build a basic in-memory bank account system in Go, focusing on fundamental language constructs and idiomatic style.

#### 💡 Functional Requirements

 * Create accounts with an ID and an initial balance.
 * Deposit funds into an account.
 * Withdraw funds, with an error returned if funds are insufficient.
 * View balance and represent an account as a formatted string.

#### 🧠 Non-Functional Requirements

 * Use Go structs, methods, and error handling idioms.
 * Follow Go's naming conventions and receiver patterns.
 * Make Account immutable.
 * Encapsulate account behaviour behind methods (not just raw field access).
 * Add unit tests to validate logic using the built-in testing package.

### Golang gotchas

|     | Gotcha| Description |
| --- |  ---  |    ---      |
| ❌  | Pointer vs Value Receiver   | Value receivers won't mutate original struct             |
| ❌  | String-based error matching | Use `errors.Is()` instead                                |
| ❌  | Range copies values         | Use `&slice[i]` or index access                          |
| ❌  | Defer order & scope         | Defer runs LIFO after return                             |
| ❌  | Ignoring errors             | Always check returned errors                             |
| ❌  | Interface `nil != nil`      | Use `errors.Is()` or `errors.As()`                       |
| ❌  | Lowercase fields            | Not exported outside package                             |
| ✅  | Zero values work            | e.g. `Account{}` is safe unless you enforce constructors |

