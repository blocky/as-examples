package as

import (
	"errors"
	"runtime"
)

type HostHTTPRequestInput struct {
	Method  string              `json:"method"`
	URL     string              `json:"url"`
	Headers map[string][]string `json:"headers"`
	Body    []byte              `json:"body"`
}

type HostHTTPRequestResult struct {
	IsOk  bool                  `json:"ok"`
	Value HostHTTPRequestOutput `json:"value"`
	Error string                `json:"error"`
}

type HostHTTPRequestOutput struct {
	StatusCode int                 `json:"status_code"`
	Body       []byte              `json:"body"`
	Headers    map[string][]string `json:"headers"`
}

//go:wasmimport env httpRequest
func _hostFuncHTTPRequest(offset, size uint32) uint64

func HostFuncHTTPRequest(
	input HostHTTPRequestInput,
) (
	HostHTTPRequestOutput,
	error,
) {
	inputData, err := Marshal(input)
	if err != nil {
		msg := "marshaling input data: " + err.Error()
		return HostHTTPRequestOutput{}, errors.New(msg)
	}

	inOffset, inSize := bytesToOffsetSize(inputData)
	outPtr := _hostFuncHTTPRequest(inOffset, inSize)
	runtime.KeepAlive(inputData)
	outputData := bytesFromFatPtr(outPtr)

	var result HostHTTPRequestResult
	err = Unmarshal(outputData, &result)
	switch {
	case err != nil:
		msg := "unmarshaling output data: " + err.Error()
		return HostHTTPRequestOutput{}, errors.New(msg)
	case !result.IsOk:
		msg := "host fn returned error: " + result.Error
		return HostHTTPRequestOutput{}, errors.New(msg)
	}
	return result.Value, nil
}
