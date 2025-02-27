# Hello World - Attesting a Function Call

This example shows you how to use the Blocky Attestation Service (Blocky AS) to
attest a simple function calls.

You'll learn how to:

- Create a function that returns a `"Hello, World!"` message
- Log messages from the function
- Invoke functions in the Blocky AS using its `bky-as` CLI
- Extract function output from the Blocky AS attestation

## Setup

- Install the Blocky AS CLI by following the
  [setup instructions](https://blocky-docs.redocly.app/attestation-service/setup)
  in the Blocky AS documentation.
- Make sure you also have
  [Docker](https://www.docker.com/) and [jq](https://jqlang.org/) installed on
  your system.

## Quick Start

To run this example, call:

```bash
make run
```

You will see the following output extracted from a Blocky AS response:

```
Log:
Writing "Hello, World!" to host

Output:
Hello, World!
```

## Walkthrough

### Step 1: Create a function that returns a "Hello, World!" message.

Our first goal is to create a simple function that returns a `"Hello, World!"`
message. We will write this function in Go and compile it to WebAssembly (WASM)
to run on the Blocky AS server. If you open [`main.go`](./main.go) you'll see
our function there:

```go
//export helloWorld
func helloWorld(inputPtr, secretPtr uint64) uint64 {
	msg := "Hello, World!"

	as.Log(fmt.Sprintf("Writing \"%s\" to host\n", msg))

	return as.WriteToHost([]byte(msg))
}
```

You will notice a few things:

- The `helloWorld` function is exported so that it can be invoked by the
  Blocky AS server in a WASM runtime.
- The function takes two `uint64` arguments and returns a `uint64`. These are
  fat pointers to shared memory managed by the Blocky AS server, where the first
  32 bits are a memory address and the second 32 bits are the size of the data.
  The memory space is sandboxed and shared between the TEE host program (Blocky
  AS server) and the WASM runtime (your function). The `inputPtr` and
  `secretPtr` arguments carry user-defined function input and secrets,
  though we don't make use of them in this example. The output of the function
  is also a memory pointer, whose value will be returned to the user.
- The function uses `as.Log` to write a message to the Blocky AS server log, 
  maintained separately for each function invocation. You can log messages
  to debug or monitor your function's behavior.
- The function calls `as.WriteToHost` to write a byte array (serialized from
  `msg`) to shared memory. The host (Blocky AS server) will create an 
  attestation over that array as a part of its response.

### Step 2: Compile the function to WebAssembly (WASM)

To invoke our function in the Blocky AS server, we first need to compile
it into a WASM file. If you inspect the `build` target in the
[`Makefile`](./Makefile), you'll see the build command:

```bash
@docker run --rm \
  -v .:/src \
  -w /src \
  tinygo/tinygo:0.31.2 \
  tinygo build -o tmp/x.wasm -target=wasi main.go
```

where we use `docker` to run [`tinygo`](https://tinygo.org/) to compile 
[`main.go`](./main.go) to WASM and save it to `tmp/x.wasm`. You can build our
function by calling:

```bash
make build
```

### Step 3: Invoke the function on the Blocky AS server

To invoke our function, we first need to define an invocation template.
We have one set up already in [`fn-call.json`](./fn-call.json) that looks like:

```json
[
  {
    "code_file": "tmp/x.wasm",
    "function": "helloWorld"
  }
]
```

where `code_file` is the path to the WASM file we compiled earlier, and
`function` is the name of the exported function we want to call.

To invoke our function, we need to pass the invocation template to `bky-as`.
If you inspect the `run` target in the [`Makefile`](./Makefile), you'll see the
command:

```bash
cat fn-call.json | bky-as attest-fn-call > tmp/out.json
```

where we use `cat` to read the [`fn-call.json`](./fn-call.json), pipe it to
`bky-as attest-fn-call`, and save the output to `tmp/out.json`.

So then to run our function, you can call:

```bash
make run
```

### Step 4: Extract function output from the Blocky AS attestation

The `run` target will extract the log and the attested output of the function 
call. 
After running:

```bash
make run
```

you should see:

```
Log:
Writing "Hello, World!" to host

Output:
Hello, World!
```

To dive deeper, let's again look at the `run` target in the 
[`Makefile`](./Makefile). There you will see that we save the output of the
`bky-as attest-fn-call` command to `tmp/out.json`, which contains:

```json
{
  "enclave_attested_application_public_key": {
    "enclave_attestation": {
      "Platform": "plain",
      "PlAttests": [
        "eyJEYXRhIjoiZXlKamRYSjJaVj...", 
        "eyJEYXRhIjoiUzIxMWNuRjRTWE...", 
        "eyJEYXRhIjoiVTJ0MlFubDZha2...",
        "eyJEYXRhIjoiVEc5TFUwNU5Sal...",
        "eyJEYXRhIjoidHgyUmlmeEVITy..."
      ]
    },
    "public_key": {
      "curve_type": "p256k1",
      "data": "BKmurqxIrdHeTwJN0YCV/4xbOv1iCA5jdSkvByzjH6UccaRDrB8KM295IkeihMQJOLoKSNMF5/mKypRbUp7Lkcs="
    }
  },
  "function_calls": [
    {
      "transitive_attestation": "WyJXeUpOTWtab1QxUlNhRTV...",
      "claims": {
        "hash_of_code": "3aa94a482d4c37fb86a913f499ddcd22c316cd26293285bf063d015c160121e1f8821019d4e141ac1eb17030f556368a7edbd3d4cc24f159107b2bb07fb27a05",
        "function": "helloWorld",
        "hash_of_input": "a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26",
        "output": "SGVsbG8sIFdvcmxkIQ==",
        "hash_of_secrets": "9375447cd5307bf7473b8200f039b60a3be491282f852df9f42ce31a8a43f6f8e916c4f8264e7d233add48746a40166eec588be8b7b9b16a5eb698d4c3b06e00"
      },
      "logs": "V3JpdGluZyBvdXQgIkhlbGxvLCBXb3JsZCEiCg=="
    }
  ]
}
```

The `enclave_attested_application_public_key` contains the `enclave_attestation`
over the Blocky AS server public key. The `function_calls` section contains the
`transitive_attestation` over the function call. The `bky-as` CLI verifies the
`enclave_attestation` by making sure that it has been signed by the TEE hardware
manufacturer's private key and that the code measurement of the Blocky AS server
running inside our TEE matches one in the `acceptable_measurements` list in 
[`config.toml`](../config.toml). The `bky-as` CLI then extracts the enclave
attested application public key, generated by the Blocky AS server on startup,
and uses it to verify the signature of the `transitive_attestation` and extract
its `claims`. You can learn more about this process in
the [Attestations in the Blocky Attestation Service](https://blocky-docs.redocly.app/attestation-service/concepts#attestations-in-the-blocky-attestation-service)
section in our documentation.

The `claims` section contains attested information about the execution of
the function.
You can see:

- `hash_of_code`: The hash of the WASM code executed by the Blocky AS server
- `function`: The name of the function executed by the Blocky AS server
- `hash_of_input`: The hash of the input data used by the function. In this
  example this is the hash of the empty string, since we didn't specify any
  input.
- `hash_of_secrets`: The hash of the secrets used by the function. In this
  example this is the hash of the empty string, since we didn't specify any
  secrets.
- `output`: The output of the function encoded in base64.

Finally, the `logs` field contains the logs generated by the function execution
and encoded in base64. Notice that while `output` is a part of the attested
`claims`, the `logs` are not attested and are only a part of the server response. 

If you look at the `run` target in the [Makefile](./Makefile) again, you will
see that we use `jq` to extract the `output` and the `logs` fields from
`tmp/out.json` and decode them from base64.

## Next steps

Now that you have successfully run the example, you can start modifying it to
fit your own needs. Check out other examples in this repository, to learn what
else you can do with Blocky AS, and in particular the 
[Hello World - Bringing A Blocky AS Function Call Attestation On Chain](../hello_world_on_chain)
example.


