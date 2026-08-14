package ingest

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"strings"

	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func decodeBody(raw []byte, contentEncoding string) ([]byte, error) {
	switch strings.ToLower(contentEncoding) {
	case "", "identity":
		return raw, nil
	case "gzip":
		reader, err := gzip.NewReader(bytes.NewReader(raw))
		if err != nil {
			return nil, fmt.Errorf("invalid gzip body: %w", err)
		}
		defer reader.Close()
		var output bytes.Buffer
		if _, err := output.ReadFrom(reader); err != nil {
			return nil, fmt.Errorf("invalid gzip body: %w", err)
		}
		return output.Bytes(), nil
	case "deflate":
		reader, err := zlib.NewReader(bytes.NewReader(raw))
		if err != nil {
			return nil, fmt.Errorf("invalid deflate body: %w", err)
		}
		defer reader.Close()
		var output bytes.Buffer
		if _, err := output.ReadFrom(reader); err != nil {
			return nil, fmt.Errorf("invalid deflate body: %w", err)
		}
		return output.Bytes(), nil
	default:
		return nil, fmt.Errorf("unsupported content-encoding %q", contentEncoding)
	}
}

type bodyEncoding string

const (
	bodyEncodingProtobuf bodyEncoding = "protobuf"
	bodyEncodingJSON     bodyEncoding = "json"
)

func unmarshalOTLP(decoded []byte, contentType string, message proto.Message) (bodyEncoding, error) {
	lowerType := strings.ToLower(contentType)

	if strings.Contains(lowerType, "protobuf") || strings.HasSuffix(lowerType, "+proto") {
		if err := proto.Unmarshal(decoded, message); err != nil {
			return bodyEncodingProtobuf, fmt.Errorf("invalid OTLP protobuf body: %w", err)
		}
		return bodyEncodingProtobuf, nil
	}

	if strings.Contains(lowerType, "json") || bytes.HasPrefix(bytes.TrimSpace(decoded), []byte("{")) {
		options := protojson.UnmarshalOptions{DiscardUnknown: true}
		if err := options.Unmarshal(decoded, message); err != nil {
			return bodyEncodingJSON, fmt.Errorf("invalid OTLP JSON body: %w", err)
		}
		return bodyEncodingJSON, nil
	}

	if err := proto.Unmarshal(decoded, message); err != nil {
		return bodyEncodingProtobuf, fmt.Errorf("invalid OTLP protobuf body: %w", err)
	}
	return bodyEncodingProtobuf, nil
}

func parseTraceRequest(decoded []byte, contentType string) (*tracepb.TracesData, bodyEncoding, error) {
	request := &tracepb.TracesData{}
	encoding, err := unmarshalOTLP(decoded, contentType, request)
	if err != nil {
		return nil, encoding, err
	}
	return request, encoding, nil
}

func parseLogsRequest(decoded []byte, contentType string) (*logspb.LogsData, bodyEncoding, error) {
	request := &logspb.LogsData{}
	encoding, err := unmarshalOTLP(decoded, contentType, request)
	if err != nil {
		return nil, encoding, err
	}
	return request, encoding, nil
}
