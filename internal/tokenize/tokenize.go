package tokenize

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

//go:embed vocab/o200k_base.tiktoken
var o200kVocab []byte

func HeuristicCount(text string) int64 {
	return int64(len(text)/4) + 1
}

var (
	o200kOnce   sync.Once
	o200kErr    error
	o200kLoaded *Encoder
)

func O200K() *Encoder {
	o200kOnce.Do(func() {
		o200kLoaded, o200kErr = newEncoder("o200k_base", o200kVocab)
	})
	if o200kErr != nil {
		return nil
	}
	return o200kLoaded
}

func newEncoder(name string, vocab []byte) (*Encoder, error) {
	ranks := make(map[string]int, 200_000)
	scanner := bufio.NewScanner(bytes.NewReader(vocab))
	scanner.Buffer(make([]byte, 1024*1024), 4*1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid vocab line %q", line)
		}
		tokenBytes, err := base64.StdEncoding.DecodeString(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid vocab token %q: %w", parts[0], err)
		}
		rank, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid vocab rank %q: %w", parts[1], err)
		}
		ranks[string(tokenBytes)] = rank
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read vocab: %w", err)
	}

	regex, err := compilePretokenizer()
	if err != nil {
		return nil, err
	}

	return &Encoder{name: name, ranks: ranks, regex: regex}, nil
}
