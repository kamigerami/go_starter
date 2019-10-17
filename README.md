### GO_STARTER

A starter project used as a personal bootstrap for golang projects.

Contains things like:

- Makefile
- Gitlab-CI file
- Dockerfile
- Config struct
- CLI Arguments
- Project structure
- Http client

### How to

- Clone project.
- Rename go_starter to your project name in all files
- Rename any variables and configurations
- Get to work.

### CLI Flags

```
  -dry-run
    	Dry run flag
  -version
    	Logs current version
```

#### Env variables

```
GO_STARTER_DUMMY_URL_PATH                string
GO_STARTER_DUMMY_BASE_URL                string
GO_STARTER_DUMMY_USERNAME                string
GO_STARTER_DUMMY_PASSWORD                string
GO_STARTER_DUMMY_USER_AGENT              string
GO_STARTER_DUMMY_DEBUG                   bool
```

### Todo

Change http client to use go-resty

https://github.com/go-resty/resty
