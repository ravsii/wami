# wami - What Are My Imports?

wami, or What Are My Imports? - is a tool for go projects to analyze imports and
their aliases. It has a lot of options, filters and output formats so it could
be integrated with other tools very easily.

## Table of Contents

<!--toc:start-->
- [wami](#wami)
  - [Table of Contents](#table-of-contents)
  - [Usage](#usage)
  - [todo](#todo)
  - [example output](#example-output)
<!--toc:end-->

## Usage

```sh
NAME:
   wami - What are my imports? (wami) is a cli for import analisys for go apps.

USAGE:
   wami [global options] [arguments...]

GLOBAL OPTIONS:
   --format string, -f string  output format (text, json) (default: text)
   --recursive, -r             enables recursive walking for ALL paths. If disabled, only paths ending with '...' are treated as recursive
   --ignore-blank              ignore blank imports (e.g., '_ fmt')
   --ignore-dot                ignore dot imports (e.g., '. fmt')
   --ignore-same               ignore imports using the same alias as the original package (e.g., 'fmt fmt')
   --min uint                  minimal amount of usages to appear in the output (inclusive) (default: 0)
   --max uint                  maximum amount of usages to appear in the output (inclusive) (default: 0)
   --help, -h                  show help
```

## todo

- group by
- strip prefix
- strip suffix
- aliases only
- outputs:
  - [x] text
  - [ ] colored
  - [x] json
  - [ ] csv

## Example

Example output running on kubernetes repo:

```sh
> go run . <path> --min 100 --max 110 --ignore-same --ignore-blank
```

```sh
"k8s.io/component-base/featuregate": 110 total usages
"k8s.io/kubectl/pkg/scheme": 110 total usages
"github.com/google/cel-go/common/types": 109 total usages
   └ 6 usages as "celtypes"
"k8s.io/cli-runtime/pkg/genericclioptions": 108 total usages
"sigs.k8s.io/yaml": 105 total usages
   ├ 3 usages as "k8syaml"
   ├ 1 usages as "sigsyaml"
   └ 1 usages as "sigyaml"
"k8s.io/apimachinery/pkg/util/net": 103 total usages
   ├ 72 usages as "utilnet"
   ├ 7 usages as "netutil"
   ├ 1 usages as "apiutil"
   └ 1 usages as "machineryutilnet"
"k8s.io/kubernetes/test/e2e/framework/node": 103 total usages
   └ 103 usages as "e2enode"
"k8s.io/kubectl/pkg/util/templates": 102 total usages
"k8s.io/utils/clock/testing": 102 total usages
   ├ 77 usages as "testingclock"
   ├ 10 usages as "clocktesting"
   ├ 9 usages as "testclock"
   ├ 4 usages as "testclocks"
   ├ 1 usages as "baseclocktest"
   └ 1 usages as "clock"
"sigs.k8s.io/structured-merge-diff/v6/fieldpath": 102 total usages
"github.com/google/cadvisor/info/v1": 101 total usages
   ├ 40 usages as "cadvisorapi"
   ├ 39 usages as "info"
   ├ 7 usages as "cadvisorapiv1"
   ├ 1 usages as "cadvisorv1"
   └ 1 usages as "v10"
```
