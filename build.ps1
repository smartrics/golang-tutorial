# build.ps1

param (
    [string]$Task = "all"
)

function Test {
    go test -v ./...
}

function Lint {
    golangci-lint run
}

function Vet {
    go vet ./...
}

function Fmt {
    go fmt ./...
}

function Bench {
    go test -bench=. ./...
}

switch ($Task) {
    "test" { Test }
    "lint" { Lint }
    "vet"  { Vet }
    "fmt"  { Fmt }
    "bench" { Bench }
    "all" {
        Test
        Vet
        Lint
    }
    default { Write-Output "Unknown task: $Task" }
}
