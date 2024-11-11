# client-registration-api
API for registering application clients for users of OpenEPI

## Generating mocks for unit testing
We use the tool mockery for generating mocks for unit testing. To install mockery, run the following command:
```
brew install mockery
```

To specify which interfaces to generate mocks for, specify in .mockery.yml file. Example:
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

