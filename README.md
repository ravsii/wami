# wami

What are my imports? is a tool for go projects to analyze import

## todo

- group by
- strip prefix
- strip suffix
- ignore
- min amount
- max amount
- aliases only

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
