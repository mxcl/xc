# xc - Markdown defined task runner. 


[![test](https://github.com/joe-davidson1802/xc/actions/workflows/test.yaml/badge.svg)](https://github.com/joe-davidson1802/xc/actions/workflows/test.yaml) 
[![docs](https://github.com/joe-davidson1802/xc/actions/workflows/docs.yml/badge.svg)](https://github.com/joe-davidson1802/xc/actions/workflows/docs.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/joe-davidson1802/xc.svg)](https://pkg.go.dev/github.com/joe-davidson1802/xc)

![xc](https://user-images.githubusercontent.com/19927761/156772881-10065864-ff4d-4225-ab2b-5adbbe628845.png)


[Docs](https://xcfile.dev/) | [Getting Started Guide](https://xcfile.dev/getting-started/) | [Github](https://github.com/joe-davidson1802/xc)

`xc` is a task runner designed to maximise convenience, and minimise complexity.

Each `xc` task is defined in simple, human-readable Markdown. Meaning that for people without the `xc` tool installed there is a clear source of useful commands in the README.md file.

# Example


Take the `tag` task in this repository:

````
## tag

Deploys a new tag for the repo.

Requires: test

```
export VERSION=`git rev-list --count HEAD`
echo Adding git tag with version v0.0.${VERSION}
git tag v0.0.${VERSION}
git push origin v0.0.${VERSION}
```
````

The task could be run simply with `xc tag`, but a side-effect is that the steps for pushing a tag without the use of `xc` are clearly documented too.

```
$ xc tag
go test ./...
?       github.com/joe-davidson1802/xc/cmd/xc   [no test files]
?       github.com/joe-davidson1802/xc/models   [no test files]
ok      github.com/joe-davidson1802/xc/parser   (cached)
ok      github.com/joe-davidson1802/xc/run      (cached)
export VERSION=78
echo Adding git tag with version v0.0.78
Adding git tag with version v0.0.78
git tag v0.0.78
git push origin v0.0.78
Total 0 (delta 0), reused 0 (delta 0), pack-reused 0
To github.com:joe-davidson1802/xc
 * [new tag]         v0.0.78 -> v0.0.78
```


# Tasks for this project:

## Tasks

### test

Test the project.

```
go test ./...
```

### build

Builds the `xc` binary.

```
go build ./cmd/xc
```

### tag
Deploys a new tag for the repo.

Requires: test
```
export VERSION=`git rev-list --count HEAD`
echo Adding git tag with version v0.0.${VERSION}
git tag v0.0.${VERSION}
git push origin v0.0.${VERSION}
```
