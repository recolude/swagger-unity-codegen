# Swagger Unity Codegen

[![Build Status](https://travis-ci.com/recolude/swagger-unity-codegen.svg?branch=master)](https://travis-ci.com/recolude/swagger-unity-codegen) [![Coverage](https://codecov.io/gh/recolude/swagger-unity-codegen/branch/master/graph/badge.svg)](https://codecov.io/gh/recolude/swagger-unity-codegen)

Generate valid networking code for Unity3D that **WILL NOT REQUIRE** any external dependencies/DLLs. Created and used for [Recolude](https://app.recolude.com)'s Unity Plugin. Lots of cute unity things you can do here that wouldn't make sense sitting in original swagger codegen repo.

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
swag3d --file api/openapi-spec/swagger.json generate \
	--config-name="RecoludeConfig" \
	--config-menu="Recolude/Config" \
	--tags "RecordingService" \
	--namespace Recolude.API \
	--out "Scripts/Recolude/API"
```

## Features

### Generate Scriptable Object For Configuring Services

The swagger tool will generate a [Scriptable Object](https://docs.unity3d.com/Manual/class-ScriptableObject.html) that can be used to store different keys found in your security definitions.

![Imgur](https://i.imgur.com/WHI9XV2.png)

### A Library You Can Use To Generate Your Own Code

You don't need a swagger file to generate your own unity code! This allows you to generate c# code as part of something like custom build pipelines that use in-house API definitions.

```go
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/recolude/swagger-unity-codegen/unitygen"
	"github.com/recolude/swagger-unity-codegen/unitygen/definition"
	"github.com/recolude/swagger-unity-codegen/unitygen/path"
	"github.com/recolude/swagger-unity-codegen/unitygen/property"
)

func main() {
	// Create a new file
	f, err := os.Create("api.cs")
	if err != nil {
		panic(err)
	}

	// What we want our succesful response to be
	responseDefinition := definition.NewObject(
		"EchoResponse",
		[]property.Property{
			property.NewString("response", ""),
			property.NewString("serverTime", "date-time"),
		},
	)

	service := unitygen.NewService(
		"ExampleService",
		[]path.Path{
			path.NewPath(
				"api/echo",     // The URL endpoint, ex: example.com/api/echo
				"Echo",         // Name of the function that gets generated
				http.MethodGet, // Get Rquest
				nil,            // Different tags
				nil,            // Security Definitions (Like API Keys)

				// Different responses that can be sent to the function
				map[string]path.Response{
					"200": path.NewResponse("", responseDefinition),
				},

				// Paramters to call the function
				[]path.Parameter{
					path.NewParameter(
						path.QueryParameterLocation,   // Where the parameter should be located
						"whatToEcho",                  // name of the query param
						true,                          // require the parameter
						property.NewString("val", ""), // name of the variable in c#
					),
				},
			),
		},
	)

	// Write the service C# code to api.cs
	fmt.Fprint(f, responseDefinition.ToCSharp())
	fmt.Fprint(f, service.ToCSharp(nil, "ServiceConfig"))
}

```

When the above code executes, it will generate the following C# code:

```c#
[System.Serializable]
public class EchoResponse
{
    public string response;

    public System.DateTime serverTime;
}

public class ExampleService
{
    public ServiceConfig Config { get; }

    public ExampleService(ServiceConfig Config)
    {
        this.Config = Config;
    }

    public class EchoUnityWebRequest
    {

        public EchoResponse success;

        public UnityWebRequest UnderlyingRequest { get; }

        public EchoUnityWebRequest(UnityWebRequest req)
        {
            this.UnderlyingRequest = req;
        }

        public IEnumerator Run()
        {
            yield return this.UnderlyingRequest.SendWebRequest();
            if (UnderlyingRequest.responseCode == 200)
            {
                success = JsonUtility.FromJson<EchoResponse>(UnderlyingRequest.downloadHandler.text);
            }
        }

    }
    public EchoUnityWebRequest Echo(string val)
    {
        var unityNetworkReq = new UnityWebRequest(string.Format("{0}api/echo?whatToEcho={1}", this.Config.BasePath, UnityWebRequest.EscapeURL(whatToEcho)), UnityWebRequest.kHttpVerbGET);
        unityNetworkReq.downloadHandler = new DownloadHandlerBuffer();
        return new EchoUnityWebRequest(unityNetworkReq);
    }
}
```

## TODO

Ordered by priority (to me)!

- [X] Scriptable Object For Configuration.
- [ ] Easier Request Building
- [ ] Don't include definitions that are never used
- [ ] Support Searilizing Bodies
- [ ] Optional Parameters In Request Body
- [ ] Embedded object definitions
- [ ] Embedded array object definitions
- [ ] Generate a scriptable object for any definition found in the swagger file
- [ ] Ability to generate `*.unitypackage`
- [ ] YAML support
- [ ] OpenAPI support
- [ ] Oauth security definition
- [ ] XML consumption/produce Support (but why tho)
