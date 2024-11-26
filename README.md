# client-registration-api
API for registering application clients for users of OpenEPI

*NB!* This service needs to run in the same cluster as our authentication solution.

## Generating mocks for unit testing
We use the tool mockery for generating mocks for unit testing. To install mockery, see installation instructions:
https://vektra.github.io/mockery/latest/installation/

To specify which interfaces to generate mocks for, specify in .mockery.yaml file. Example:
```
with-expecter: True
dir: tests/mocks/{{.InterfaceDirRelative}}
mockname: "Mock{{.InterfaceName}}"
outpkg: "{{.PackageName}}"
filename: "mock_{{.InterfaceName}}.go"
all: True
packages:
  github.com/openearthplatforminitiative/client-registration-api/keycloak:
```

The above config will create mocks for all interfaces in the `keycloak` package.

To generate mocks for unit testing, run the following command:
```
mockery
```

The mocks should be checked in to the repository.

