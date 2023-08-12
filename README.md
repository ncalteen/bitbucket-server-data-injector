## About

A simple go program to inject empty projects and repos to BBS for migration
testing

### Usage

The program uses Basic auth and the root URL of the BBS instance to peform
operations. These are expected as ENV vars: `BBS_API_URL`, `BBS_API_USER`, and
`BBS_API_PASS`

```
export BBS_API_URL="https://yourBBSinstance.io/"
```

```
export BBS_API_USER="admin"
```

```
export BBS_API_PASS="securepassword"
```

Additionally two command line flags can be used to create a set amount of
project and repo in each project: `--projects-count` and `--repos-count`. These
both default to 1.

```
go run main.go --projects-count 10 --repos-count 10
```
