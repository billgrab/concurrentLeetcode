package problems

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// TrieNode ...
type TrieNode struct {
	children [26]*TrieNode
	word string
}

// findWords is the required method for https://leetcode.com/problems/word-search-ii/
// 
func findWords(board [][]byte, words []string) []string {
	// Running on leetcode.com successfully
	result, _ := FindWordsInSequence(board, words)
	_ = result

	// Submitting on leetcode.com got TimeLimitExceeded error, but runs fine on local machine -- see find_words.test.go for parallel + sequential flavors.
	// likely due to limitation of test server on leetcode.com
	/*
	=== RUN   TestFindWords
	=== RUN   TestFindWords/samll_data_set
	[2022-02-20 03:06:52.944264 -0800 PST m=+0.001021784] starting test case [samll data set]
	explored num of cells: 16
	Parallel        == pal_uSec: 172, total_uSec: 322, result: [eat oath]
	Sequential      == seq_uSec: 36, total_uSec: 410, result: [oath eat]
	=== RUN   TestFindWords/big_data_set
	[2022-02-20 03:06:52.944731 -0800 PST m=+0.001488737] starting test case [big data set]
	explored num of cells: 144
	Parallel        == pal_uSec: 1105006, total_uSec: 1105110, result: [a aa aaa aaaa aaaaa aaaaaa aaaaaaa aaaaaaaa aaaaaaaaa aaaaaaaaaa]
	Sequential      == seq_uSec: 1275218, total_uSec: 2380384, result: [a aa aaa aaaa aaaaa aaaaaa aaaaaaa aaaaaaaa aaaaaaaaa aaaaaaaaaa]
	--- PASS: TestFindWords (2.38s)
    --- PASS: TestFindWords/samll_data_set (0.00s)
    --- PASS: TestFindWords/big_data_set (2.38s)
	*/
	result, _ = FindWordsInParallel(board, words)
	return result
}

// FindWordsInParallel ...
func FindWordsInParallel(board [][]byte, words []string) (strs []string, lapse time.Duration) {
	root := BuildTrie(words)
	rows, cols := len(board), len(board[0])

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	createEmpty := func(r, c int) [][]bool {
		var res [][]bool
		for i:=0; i < r; i++ {
			res = append(res, make([]bool, c))
		}
		return res
	}
	t0 := time.Now()
	defer func() {
		lapse = time.Since(t0)
	}()

	bigBuf := rows*cols // might need to set higher
	var result chan string = make(chan string, bigBuf)
	var wg sync.WaitGroup
	var explored int64
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			wg.Add(1)
			go func(r, c int) {
				defer func() {
					atomic.AddInt64(&explored, 1)
					// fmt.Printf("Done %d: (r=%d, c=%d)\n", atomic.LoadInt64(&explored), r, c)
					wg.Done()
				}()
				used := createEmpty(rows, cols)
				explore(ctx, board, r, c, root, used, result)
			}(i, j)
		}
	}
	
	agrDone := make(chan interface{})
	go func() {
		defer func() { 
			agrDone <- true
		}()
		for str := range result {
			if !Contains(strs, str) {
				strs = append(strs, str)
			}
		}
	}()
	
	wg.Wait()
	fmt.Printf("explored num of cells: %d\n", explored)
	close(result)
	
	<-agrDone
	return
}

// FindWordsInSequence ...
func FindWordsInSequence(board [][]byte, words []string) (strs []string, lapse time.Duration) {
	root := BuildTrie(words)
	rows, cols := len(board), len(board[0])

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	createEmpty := func(r, c int) [][]bool {
		var res [][]bool
		for i:=0; i < r; i++ {
			res = append(res, make([]bool, c))
		}
		return res
	}
	t0 := time.Now()
	defer func() {
		lapse = time.Since(t0)
	}()

	bigBuf := rows*cols // might need to set higher
	var result chan string = make(chan string, bigBuf)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for str := range result {
			if !Contains(strs, str) {
				strs = append(strs, str)
			}
		}
	}()

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			used := createEmpty(rows, cols)
			explore(ctx, board, i, j, root, used, result)
		}
	}
	close(result)
	
	wg.Wait()
	return
}

func explore(ctx context.Context, board [][]byte, i, j int, node *TrieNode, visited [][]bool, result chan string)  {
	select {
	case <-ctx.Done():
		fmt.Println("ctx cancelled due to timeout")
		return
	default:
	}

    rows, cols := len(board), len(board[0])
    if i<0 || j<0 || i>=rows || j>=cols || visited[i][j] {
        return
    }

    c := board[i][j]
    if node.children[c-'a'] == nil {
        return
    }
    
    visited[i][j] = true
    if w := node.children[c-'a'].word; w != "" {
        result <- w
    }
    nextNode := node.children[c-'a']
	childCtx, childCancel := context.WithCancel(ctx)
	defer childCancel()
    explore(childCtx, board, i-1, j, nextNode, visited, result)
    explore(childCtx, board, i+1, j, nextNode, visited, result)
    explore(childCtx, board, i, j-1, nextNode, visited, result)
    explore(childCtx, board, i, j+1, nextNode, visited, result)
    visited[i][j] = false
}

// Contains ...
func Contains(objs []string, target string ) bool {
	for _, o := range objs {
		if strings.EqualFold(o, target) {
			return true
		}
	}
	return false
}

// BuildTrie ...
func BuildTrie(words []string) *TrieNode {
	root := &TrieNode{}
	for _, w := range words {
		current := root
		for _, r := range w {
			if current.children[r-'a'] == nil {
				current.children[r-'a'] = &TrieNode{}
			}
			current = current.children[r-'a']
		}
		current.word = w
	}
	return root
}
