# Contributing

## Development workflow

1. Read the relevant scenarios in `cli/features/*.feature`.
2. Implement the behavior in the CLI packages.
3. Add or update Go tests that cover each scenario.
4. Run the local checks before opening a pull request:

```bash
go -C cli test -race ./...
go -C cli vet ./...
go -C cli build ./...
```

Every pull request and every push to `main` runs the same checks through GitHub Actions in `.github/workflows/ci.yml`.

## Feature scenarios and tests

Feature files describe the expected behavior in business-readable language. They are the project requirements and acceptance criteria; they are not executed directly yet because this repository does not use a Gherkin runner.

For each scenario, contributors must add automated Go coverage in the package that owns the behavior. Keep the scenario wording and the test intent aligned:

| Feature scenario                | Go test coverage                                |
| ------------------------------- | ----------------------------------------------- |
| Complete deployment environment | `internal/dependencies` and `internal/checkcmd` |
| Missing optional tools          | `internal/dependencies` and `internal/checkcmd` |
| Missing required tool           | `internal/dependencies` and `internal/checkcmd` |
| Installation plan               | `internal/installcmd` and `internal/packages`   |
| Applied installation            | `internal/installcmd` and `internal/packages`   |

If the project later adopts Godog or another BDD runner, the existing scenarios can become executable without changing their role as acceptance criteria.
