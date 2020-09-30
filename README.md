# Swagger Unity Codegen

[![Build Status](https://travis-ci.com/recolude/swagger-unity-codegen.svg?branch=master)](https://travis-ci.com/recolude/swagger-unity-codegen) [![Coverage](https://codecov.io/gh/recolude/swagger-unity-codegen/branch/master/graph/badge.svg)](https://codecov.io/gh/recolude/swagger-unity-codegen)

Generate valid networking code for Unity3D. Created and used for Recolude's Unity Plugin.

Lots of cute unity things you can do here that wouldn't make sense sitting in original swagger codegen repo.

PRs + Issues Welcome.

## Installation

Built and tested with golang 1.14, but should probably work on versions as far back as 1.10 (when string builder was introduced).

```bash
git clone https://github.com/recolude/swagger-unity-codegen.git
cd swagger-unity-codegen
go install ./cmd/swag3d
```

## Usage

```
NAME:
   swag3d - swagger and Unity3D meet

USAGE:
   swag3d [global options] command [command options] [arguments...]

VERSION:
   0.1.0

DESCRIPTION:
   Generate C# code specifically for Unity3D from a swagger file

AUTHOR:
   Elijah C Davis

COMMANDS:
   generate, g  Generate c# code from a swagger file
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --file value   where to load swagger from (default: "swagger.json")
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

### Example

```bash
swag3d --file api/openapi-spec/swagger.json generate --namespace Recolude.API > API.cs
```

## TODO

- [ ] Embedded object definitions
- [ ] Embedded array object definitions
- [ ] YAML support
- [ ] OpenAPI support
- [ ] Oauth security definition
