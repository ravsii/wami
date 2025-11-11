# wami

[![Build](https://github.com/ravsii/wami/actions/workflows/build.yml/badge.svg)](https://github.com/ravsii/wami/actions/workflows/build.yml) [![Test](https://github.com/ravsii/wami/actions/workflows/test.yml/badge.svg)](https://github.com/ravsii/wami/actions/workflows/test.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/ravsii/wami)](https://goreportcard.com/report/github.com/ravsii/wami) [![codecov](https://codecov.io/github/ravsii/wami/graph/badge.svg?token=3sQsbMyLJX)](https://codecov.io/github/ravsii/wami) [![Go Reference](https://pkg.go.dev/badge/github.com/ravsii/wami.svg)](https://pkg.go.dev/github.com/ravsii/wami)

**wami** ("What Are My Imports?") is a tool for Go projects that analyzes
imports and their aliases. It offers a wide range of options, filters, and
output formats, making it easy to integrate with other tools.

## Table of Contents

<!--toc:start-->
- [wami](#wami)
  - [Table of Contents](#table-of-contents)
  - [Install](#install)
  - [Usage](#usage)
  - [Outputs](#outputs)
    - [Text output](#text-output)
    - [JSON output](#json-output)
    - [CSV output](#csv-output)
<!--toc:end-->

## Install

```sh
go install github.com/ravsii/wami@latest
```

## Usage

```sh
NAME:
   wami - What are my imports? (wami) is a cli for import analisys for go apps.

USAGE:
   wami [global options] [arguments...]

GLOBAL OPTIONS:
   --aliases-only, -a          only output imports that have aliases. Note: all imports will be parsed anyways, for a total amount of usages
   --format string, -f string  output format (text, text-colored, json, csv) (default: "text-colored")
   --ignore regexp             regexp to ignore import paths
   --ignore-alias regexp       regexp to ignore import aliases
   --ignore-blank              ignore blank imports (e.g., '_ fmt')
   --ignore-dot                ignore dot imports (e.g., '. fmt')
   --ignore-same               ignore imports using the same alias as the original package (e.g., 'fmt fmt')
   --include regexp            regexp to include import paths
   --include-alias regexp      regexp to include import aliases
   --max uint                  maximum amount of usages to appear in the output (inclusive) (default: 0)
   --min uint                  minimal amount of usages to appear in the output (inclusive) (default: 0)
   --recursive, -r             enables recursive walking for ALL paths. If disabled, only paths ending with '...' are treated as recursive
   --help, -h                  show help
```

## Outputs

Here’s an example output in multiple formats, generated from the
[`kubernetes`](https://github.com/kubernetes/kubernetes) repository - one of the
largest Go projects:

```sh
> wami <path> --min 300 --max 350
````

### Text output

<details>
<summary>Show example</summary>

```sh
k8s.io/client-go/kubernetes/scheme: 349 total usages
 ├ 210 usages as scheme
 ├ 13 usages as clientscheme
 ├ 6 usages as clientsetscheme
 ├ 3 usages as k8sscheme
 ├ 2 usages as clientgoscheme
 ├ 1 usages as cgoscheme
 ├ 1 usages as clientgokubescheme
 ├ 1 usages as kubernetesscheme
 └ 1 usages as typedscheme
syscall: 345 total usages
regexp: 342 total usages
 └ 1 usages as re
k8s.io/apimachinery/pkg/api/equality: 320 total usages
 ├ 215 usages as apiequality
 └ 62 usages as equality
github.com/onsi/gomega: 317 total usages
 ├ 8 usages as .
 └ 2 usages as o
```

</details>

### JSON output

<details>
<summary>Show example</summary>

```json
[
  {
    "path": "k8s.io/client-go/kubernetes/scheme",
    "count": 349,
    "aliases": [
      {
        "count": 210,
        "name": "scheme"
      },
      {
        "count": 13,
        "name": "clientscheme"
      },
      {
        "count": 6,
        "name": "clientsetscheme"
      },
      {
        "count": 3,
        "name": "k8sscheme"
      },
      {
        "count": 2,
        "name": "clientgoscheme"
      },
      {
        "count": 1,
        "name": "cgoscheme"
      },
      {
        "count": 1,
        "name": "clientgokubescheme"
      },
      {
        "count": 1,
        "name": "kubernetesscheme"
      },
      {
        "count": 1,
        "name": "typedscheme"
      }
    ]
  },
  {
    "path": "syscall",
    "count": 345
  },
  {
    "path": "regexp",
    "count": 342,
    "aliases": [
      {
        "count": 1,
        "name": "re"
      }
    ]
  },
  {
    "path": "k8s.io/apimachinery/pkg/api/equality",
    "count": 320,
    "aliases": [
      {
        "count": 215,
        "name": "apiequality"
      },
      {
        "count": 62,
        "name": "equality"
      }
    ]
  },
  {
    "path": "github.com/onsi/gomega",
    "count": 317,
    "aliases": [
      {
        "count": 8,
        "name": "."
      },
      {
        "count": 2,
        "name": "o"
      }
    ]
  }
]
```

</details>

### CSV output

Alias column has the following format:

```csv
"<count1>,<alias1>;<count2>,<alias2>;..."
```

<details>
<summary>Show example</summary>

```csv
import,count,aliases
k8s.io/client-go/kubernetes/scheme,349,"210,scheme;13,clientscheme;6,clientsetscheme;3,k8sscheme;2,clientgoscheme;1,cgoscheme;1,clientgokubescheme;1,kubernetesscheme;1,typedscheme"
syscall,345,
regexp,342,"1,re"
k8s.io/apimachinery/pkg/api/equality,320,"215,apiequality;62,equality"
github.com/onsi/gomega,317,"8,.;2,o"
```

</details>
