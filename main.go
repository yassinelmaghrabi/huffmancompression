package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// struct for the nodes of the huffman tree
type huffmanNode struct {
	runecount runeCountPair
	left      *huffmanNode
	right     *huffmanNode
}

type runeCountPair struct {
	charecter rune
	count     int
}

// counter function for counting rune probabilities
func counter(text string) []runeCountPair {
	resultMap := make(map[rune]int)
	for _, char := range text {
		resultMap[char]++
	}
	characterCounts := make([]runeCountPair, 0, len(resultMap))
	for key, value := range resultMap {
		characterCounts = append(characterCounts, runeCountPair{
			charecter: key,
			count:     value,
		})
	}
	return characterCounts
}

// turns a list of huffman nodes into heap
func heapify(pairSlice []huffmanNode) []huffmanNode {
	heap := make([]huffmanNode, len(pairSlice))
	copy(heap, pairSlice)
	for i := len(heap)/2 - 1; i >= 0; i-- {
		siftDown(heap, i, len(heap))
	}
	return heap
}

// ensures the heap is still a heap after pop or add operations using a bottom up approach
func siftDown(heap []huffmanNode, start, end int) {
	root := start
	for {
		child := 2*root + 1
		if child >= end {
			return
		}
		if child+1 < end && heap[child+1].runecount.count < heap[child].runecount.count {
			child++
		}
		if heap[root].runecount.count <= heap[child].runecount.count {
			return
		}
		swap(heap, root, child)
		root = child
	}
}

func swap(slice []huffmanNode, i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// turns a list of runecountpairs to a huffman tree and returns the root of the tree for traversal
func buildHuffmanTree(pairs []runeCountPair) *huffmanNode {
	nodes := make([]huffmanNode, len(pairs))
	for i, pair := range pairs {
		nodes[i] = huffmanNode{runecount: pair}
	}

	heap := heapify(nodes)
	for len(heap) > 1 {
		left := heap[0]
		heap = heap[1:]
		siftDown(heap, 0, len(heap))

		right := heap[0]
		heap = heap[1:]
		siftDown(heap, 0, len(heap))

		newNode := huffmanNode{
			runecount: runeCountPair{
				count: left.runecount.count + right.runecount.count,
			},
			left:  &left,
			right: &right,
		}
		heap = append(heap, newNode)
		siftDown(heap, 0, len(heap))
	}
	return &heap[0]
}

// recurssive function that generates the codes with left being a 0 and right being a 1
func generateCodes(node *huffmanNode, prefix string, codes map[rune]string) {
	if node == nil {
		return
	}
	if node.runecount.charecter != 0 {
		codes[node.runecount.charecter] = prefix
	} else {
		generateCodes(node.left, prefix+"0", codes)
		generateCodes(node.right, prefix+"1", codes)
	}
}

func encodeText(text string, codes map[rune]string) string {
	var encoded strings.Builder
	for _, char := range text {
		encoded.WriteString(codes[char])
	}
	return encoded.String()
}

func readFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content.WriteString(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return content.String(), nil
}

func calculateEntropy(text string, pairs []runeCountPair) float64 {
	totalChars := len(text)
	entropy := 0.0
	for _, pair := range pairs {
		prob := float64(pair.count) / float64(totalChars)
		entropy -= prob * math.Log2(prob)
	}
	return entropy
}

func calculateCompressionEfficiency(text string, encodedText string, entropy float64) (float64, float64) {
	originalSize := len(text) * 8
	encodedSize := len(encodedText)

	compressionRatio := float64(originalSize) / float64(encodedSize)
	avgBitsPerChar := float64(encodedSize) / float64(len(text))
	efficiency := entropy / avgBitsPerChar

	return compressionRatio, efficiency
}

func main() {
	filePath := "./input.txt"

	text, err := readFile(filePath)
	if err != nil {
		color.Red("Error reading file: %v", err)
		return
	}
	color.Cyan("Original Text:\n%s\n\n", text)

	pairs := counter(text)

	root := buildHuffmanTree(pairs)

	codes := make(map[rune]string)
	generateCodes(root, "", codes)

	encodedText := encodeText(text, codes)

	entropy := calculateEntropy(text, pairs)
	compressionRatio, efficiency := calculateCompressionEfficiency(text, encodedText, entropy)

	color.Yellow("Character Codes:")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Character", "Frequency", "Huffman Code"})
	table.SetAlignment(tablewriter.ALIGN_CENTER)

	for _, pair := range pairs {
		character := string(pair.charecter)
		if pair.charecter == ' ' {
			character = "SPACE"
		}
		table.Append([]string{character, fmt.Sprintf("%d", pair.count), codes[pair.charecter]})
	}
	table.Render()

	color.Yellow("\nEncoded Text:")
	fmt.Println(encodedText)

	color.Green("\nCompression Results:")
	fmt.Printf("Entropy: %.4f bits per character\n", entropy)
	fmt.Printf("Compression Ratio: %.4f\n", compressionRatio)
	fmt.Printf("Compression Efficiency: %.4f\n", efficiency)
}
