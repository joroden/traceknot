package tokenize

import (
	"container/heap"
	"regexp"
)

const o200kRegex = `[^\r\n\p{L}\p{N}]?[\p{Lu}\p{Lt}\p{Lm}\p{Lo}\p{M}]*[\p{Ll}\p{Lm}\p{Lo}\p{M}]+(?i:'s|'t|'re|'ve|'m|'ll|'d)?|[^\r\n\p{L}\p{N}]?[\p{Lu}\p{Lt}\p{Lm}\p{Lo}\p{M}]+[\p{Ll}\p{Lm}\p{Lo}\p{M}]*(?i:'s|'t|'re|'ve|'m|'ll|'d)?|\p{N}{1,3}| ?[^\s\p{L}\p{N}]+[\r\n/]*|\s*[\r\n]+|\s+`

func compilePretokenizer() (*regexp.Regexp, error) {
	return regexp.Compile(o200kRegex)
}

type Encoder struct {
	name  string
	ranks map[string]int
	regex *regexp.Regexp
}

func (e *Encoder) Name() string {
	return e.name
}

func (e *Encoder) Count(text string) int64 {
	indices := e.regex.FindAllStringIndex(text, -1)
	if len(indices) == 0 {
		if text == "" {
			return 0
		}
		return int64(len(bpeEncode(e.ranks, []byte(text))))
	}

	total := 0
	for _, span := range splitPieces(text, indices) {
		total += len(bpeEncode(e.ranks, []byte(text[span[0]:span[1]])))
	}
	return int64(total)
}

func splitPieces(text string, indices [][]int) [][]int {
	textEnd := len(text)
	output := make([][]int, 0, len(indices))
	for _, span := range indices {
		piece := text[span[0]:span[1]]
		isWhitespace := true
		for _, char := range piece {
			if char != ' ' && char != '\t' {
				isWhitespace = false
				break
			}
		}
		if !isWhitespace || span[1] == textEnd {
			output = append(output, span)
			continue
		}

		containsNewline := false
		for _, char := range piece {
			if char == '\r' || char == '\n' {
				containsNewline = true
				break
			}
		}
		if containsNewline {
			output = append(output, span)
			continue
		}
		for offset := span[0]; offset < span[1]; offset++ {
			output = append(output, []int{offset, offset + 1})
		}
	}
	return output
}

type bpeNode struct {
	id         int
	token      string
	prev, next *bpeNode
}

type pairEntry struct {
	rank    int
	leftID  int
	rightID int
	key     string
}

type pairHeap []pairEntry

func (h pairHeap) Len() int { return len(h) }
func (h pairHeap) Less(i, j int) bool {
	if h[i].rank != h[j].rank {
		return h[i].rank < h[j].rank
	}
	return h[i].leftID < h[j].leftID
}
func (h pairHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *pairHeap) Push(x any)   { *h = append(*h, x.(pairEntry)) }
func (h *pairHeap) Pop() any {
	old := *h
	last := old[len(old)-1]
	*h = old[:len(old)-1]
	return last
}

func buildBPEChain(data []byte) (head *bpeNode, nodes []*bpeNode) {
	nodes = make([]*bpeNode, len(data))
	var tail *bpeNode
	for index, b := range data {
		node := &bpeNode{id: index, token: string([]byte{b})}
		nodes[index] = node
		if tail != nil {
			tail.next = node
			node.prev = tail
		} else {
			head = node
		}
		tail = node
	}
	return head, nodes
}

func bpeEncode(ranks map[string]int, data []byte) []string {
	if len(data) == 0 {
		return nil
	}

	head, nodes := buildBPEChain(data)
	if head.next == nil {
		return []string{head.token}
	}

	h := &pairHeap{}
	heap.Init(h)
	for node := head; node.next != nil; node = node.next {
		pushPair(h, ranks, node)
	}

	for h.Len() > 0 {
		entry := heap.Pop(h).(pairEntry)
		node := nodes[entry.leftID]

		if node.next == nil || node.next.id != entry.rightID || node.token+node.next.token != entry.key {
			continue
		}

		right := node.next
		merged := node.token + right.token
		next := right.next
		if next != nil {
			next.prev = node
		}
		node.next = next
		node.token = merged
		right.prev = nil
		right.next = nil

		if node.prev != nil {
			pushPair(h, ranks, node.prev)
		}
		if node.next != nil {
			pushPair(h, ranks, node)
		}
	}

	var output []string
	for node := head; node != nil; node = node.next {
		output = append(output, node.token)
	}
	return output
}

func pushPair(h *pairHeap, ranks map[string]int, left *bpeNode) {
	if left.next == nil {
		return
	}
	key := left.token + left.next.token
	rank, ok := ranks[key]
	if !ok {
		return
	}
	heap.Push(h, pairEntry{rank: rank, leftID: left.id, rightID: left.next.id, key: key})
}
