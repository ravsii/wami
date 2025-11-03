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
   --recursive, -r  enables recursive walking for ALL paths. If disabled, only paths ending with '...' are treated as recursive
   --ignore-blank   ignore blank imports (e.g., '_ fmt')
   --ignore-dot     ignore dot imports (e.g., '. fmt')
   --ignore-same    ignore imports using the same alias as the original package (e.g., 'fmt fmt')
   --min uint       minimal amount of usages to appear in the output (inclusive) (default: 0)
   --max uint       maximum amount of usages to appear in the output (inclusive) (default: 0)
   --help, -h       show help
```

## todo

- group by
- strip prefix
- strip suffix
- aliases only
- outputs:
  - [x] text
  - [ ] colored
  - [ ] json
  - [ ] csv
- add benchmarks using hyperfine for big repos
  - k8s
  - docker
  - hugo?
  - lazyXXX?
  - go itself

## example output

```text
"fmt": 14 total usages
"github.com/ravsii/elgo": 14 total usages
"time": 9 total usages
"errors": 7 total usages
"log": 7 total usages
"testing": 6 total usages
"context": 5 total usages
   └ 2 usages as context
"net": 5 total usages
"sync": 5 total usages
   └ 1 usages as sync
"github.com/ravsii/elgo/examples/player": 4 total usages
"github.com/ravsii/elgo/socket": 4 total usages
"google.golang.org/grpc": 4 total usages
   └ 1 usages as grpc
"math/rand": 4 total usages
"strings": 4 total usages
"github.com/ravsii/elgo/grpc/pb": 3 total usages
   └ 1 usages as pb
"google.golang.org/grpc/codes": 3 total usages
   └ 1 usages as codes
"google.golang.org/grpc/status": 3 total usages
   └ 1 usages as status
"math": 3 total usages
"github.com/ravsii/elgo/grpc": 2 total usages
   └ 1 usages as elgo_grpc
"io": 2 total usages
"strconv": 2 total usages
"bufio": 1 total usage
"bytes": 1 total usage
"github.com/integrii/flaggy": 1 total usage
"google.golang.org/grpc/credentials/insecure": 1 total usage
"google.golang.org/protobuf/reflect/protoreflect": 1 total usage
   └ 1 usages as protoreflect
"google.golang.org/protobuf/runtime/protoimpl": 1 total usage
   └ 1 usages as protoimpl
"os": 1 total usage
"os/signal": 1 total usage
"reflect": 1 total usage
   └ 1 usages as reflect
"sort": 1 total usage
"syscall": 1 total usage
```
