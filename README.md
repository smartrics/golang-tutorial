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
| 1    | [part1](https://github.com/smartrics/golang-tutorial/tree/part1) | Syntax & Language Fundamentals |
| 2    | [part2](https://github.com/smartrics/golang-tutorial/tree/part2) | Structs, Methods, Interfaces |
| 3    | [part3](https://github.com/smartrics/golang-tutorial/tree/part3) | Error Handling, Testing, Tooling |
| 4    | [part4](https://github.com/smartrics/golang-tutorial/tree/part4) | Application Structure, Project Organisation, and CI/CD Readiness |
| 5    | [part5](https://github.com/smartrics/golang-tutorial/tree/part5) | Functional Composition, OO Patterns, and Advanced Struct Techniques |
| 6    | [part6](https://github.com/smartrics/golang-tutorial/tree/part6) | Concurrency, Channels and Project Architecture |
| 7    | [part7](https://github.com/smartrics/golang-tutorial/tree/part7) | Interface Composition, APIs, and I/O Integration |
| 8    | [part8](https://github.com/smartrics/golang-tutorial/tree/part8) | Common Golang gotchas |

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