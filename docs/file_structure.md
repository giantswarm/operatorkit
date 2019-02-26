# File Structure

- operatorkit resources should have a consistent structure to make code
naviagation easier.
- Each method of the resource should be saved to a separate file.
e.g. `create.go` should contain the `EnsureCreated` method.
- Additional files can be added as needed.
- Unit test files should have the usual Go `_test` suffix.

## Simple interface

- When using the simple resource interface at least the following files should
exist.

```
resource
└── example
    ├── create.go
    ├── delete.go
    ├── error.go
    └── resource.go
```

## CRUD interface

- When using the CRUD resource interface at least the following files should
exist.

```
resource
└── example
    ├── create.go
    ├── current.go
    ├── delete.go
    ├── desired.go
    ├── error.go
    ├── resource.go
    └── update.go
```
