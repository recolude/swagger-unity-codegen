# Swagger Unity Codegen

[![Build Status](https://travis-ci.com/recolude/swagger-unity-codegen.svg?branch=master)](https://travis-ci.com/recolude/swagger-unity-codegen) [![Coverage](https://codecov.io/gh/recolude/swagger-unity-codegen/branch/master/graph/badge.svg)](https://codecov.io/gh/recolude/swagger-unity-codegen) [![Go Report Card](https://goreportcard.com/badge/github.com/recolude/swagger-unity-codegen)](https://goreportcard.com/report/github.com/recolude/swagger-unity-codegen)

**[Currently In Beta: Only Supports Swagger 2.0 JSON ATM]**

Generate valid networking code for Unity3D that takes advantage of [Unity's Web Request](https://docs.unity3d.com/ScriptReference/Networking.UnityWebRequest.html) object instead of something like RestSharp. This project was both made for and is used by [Recolude](https://app.recolude.com)'s Unity Plugin. Lots of cute unity things you can do here that wouldn't make sense sitting in original swagger codegen repo.

PRs + Issues Welcome.

## Installation

Built and tested with golang 1.14, but should probably work on versions as far back as 1.10 (when string builder was introduced).

```bash
git clone https://github.com/recolude/swagger-unity-codegen.git
cd swagger-unity-codegen
go install ./cmd/swag3d
```

## Dependencies

The code produced from this tool will depend on two external DLLs. Trust me I tried my best to use just the builtin Unity Serializer at first, but it's just not powerful enough to take into account all the different types of definitions a swagger file can have 😓.

The two DLLs are: 

1. [Newtonsoft.Json](https://www.nuget.org/packages/Newtonsoft.Json/)
2. [JSONSubTypes](https://www.nuget.org/packages/JsonSubTypes/)

Make sure that the DLLs you grab from the packages both target the same version of .NET
If you don't want to download them through nuget, you can checkout the dlls folder in this repository. I haven't tested them with most versions of Unity though.

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
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

### Example

Command used to generate Recolude's code.

```bash
swag3d generate \
	--file api/openapi-spec/swagger.json \
	--config-name="RecoludeConfig" \
	--config-menu="Recolude/Config" \
	--tags "RecordingService" \
	--namespace Recolude.API \
	--out "Scripts" \
	--scriptable-object-config=false
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
	"github.com/recolude/swagger-unity-codegen/unitygen/model"
	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
	"github.com/recolude/swagger-unity-codegen/unitygen/path"
)

func main() {
	// Create a new file
	f, err := os.Create("api.cs")
	if err != nil {
		panic(err)
	}

	// What we want our succesful response to be
	responseDefinition := model.NewObject(
		"EchoResponse",
		[]model.Property{
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
					"200": path.NewDefinitionResponse("", responseDefinition),
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

    [SerializeField]
    private string serverTime;

    public System.DateTime ServerTime { get => System.DateTime.Parse(serverTime); }
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
            Interpret(this.UnderlyingRequest);
        }

        public void Interpret(UnityWebRequest req)
        {
            if (req.responseCode == 200)
            {
                success = JsonConvert.DeserializeObject<EchoResponse>(req.downloadHandler.text);
            }
        }
    }

    public class EchoRequestParams
    {
        private bool whatToEchoSet = false;
        private string whatToEcho;
        public string WhatToEcho { get { return whatToEcho; } set { whatToEchoSet = true; whatToEcho = value; } }
        public void UnsetWhatToEcho() { whatToEcho = null; whatToEchoSet = false; }

        public UnityWebRequest BuildUnityWebRequest(string baseURL)
        {
            var finalPath = baseURL + "api/echo";
            var queryAdded = false;

            if (whatToEchoSet)
            {
                finalPath += (queryAdded ? "&" : "?") + "whatToEcho=";
                queryAdded = true;
                finalPath += UnityWebRequest.EscapeURL(whatToEcho.ToString());
            }

            return new UnityWebRequest(finalPath, UnityWebRequest.kHttpVerbGET);
        }
    }
	
    public EchoUnityWebRequest Echo(EchoRequestParams requestParams)
    {
        var unityNetworkReq = requestParams.BuildUnityWebRequest(this.Config.BasePath);
        unityNetworkReq.downloadHandler = new DownloadHandlerBuffer();
        return new EchoUnityWebRequest(unityNetworkReq);
    }

    public EchoUnityWebRequest Echo(string whatToEcho)
    {
        return Echo(new EchoRequestParams()
        {
            WhatToEcho = whatToEcho,
        });
    }
}
```

## Progress Towards V1

Ordered by priority (to me)!

- [X] Scriptable Object For Configuration.
- [X] Easier Request Building.
- [X] Don't include definitions that are never used.
- [X] Support for System.DateTime.
- [X] Support Searilizing Bodies.
- [ ] Polymorphism
- [ ] Implement [Fluent Interface Pattern](https://en.wikipedia.org/wiki/Fluent_interface) For Creating Requests.
- [ ] Optional Parameters In Request Body.
- [ ] Required Fields
- [X] Embedded object definitions.
- [ ] Embedded array object definitions.
- [ ] Generate a scriptable object for any definition found in the swagger file.
- [ ] Ability to generate `*.unitypackage`.
- [ ] YAML support.
- [ ] Oauth security definition.
