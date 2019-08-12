# gotemplate

CLI app that helps you to kickstart new projects by generating a predetermined file structure with parameters.

```
go run main.go
```
### Parameters

- **Name**: 
- **Namespace**: 
- **Path**: [$PWD] (hit enter/change it)

### File structure
```
{namespace}
├── cmd
│   └── service
│       └── main.go
├── configs
│   ├── deployment.yml
│   ├── ingress.yml
│   └── service.yml
├── Dockerfile
├── docs
│   ├── api
│   │   └── readme.md
│   └── redoc.go
├── Jenkinsfile
├── lint.sh
├── pkg
├── pushSonarqube.sh
├── readme.md
├── runTests.sh
└── test
```

Project {name} has been created in $PWD/{namespace} 🍺