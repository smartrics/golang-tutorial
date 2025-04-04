# build.ps1

param (
    [string]$Task = "all"
)

function Test {
    go test -v ./tests/...
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
    go test -bench=. ./tests/...
}

function Mock {
    moq -out .\mocks\bank\bank_account_moq.go -pkg bank_mocks internal\bank BankAccount 
}

switch ($Task) {
    "test" { Test }
    "lint" { Lint }
    "vet"  { Vet }
    "fmt"  { Fmt }
    "bench" { Bench }
    "mock" { Mock }
    "all" {
        Mock
        Test
        Vet
        Lint
        Fmt
    }
    default { Write-Output "Unknown task: $Task" }
}
