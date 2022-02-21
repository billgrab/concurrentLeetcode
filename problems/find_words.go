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

type workItem struct {
	r int
	c int
}

const (
	numWorkers = 16
)

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

	/*
	Even with the update of using a fixed pool of worker, for the 12x12 big-data-set case, submission on leetcode.com still got TimeLimitExceeded.
	Seems like the highest % of completion is around (~80, ~120) out of the total 12x12 items. 
	Of course, stats varies depends on: 
		1/.test cases exection order, and
		2/.the traffic load on leetcode.com.
	
	Here are some of the captured  stdout from leetcode.com portal during repeated submissions:
	
	(numWorkers = 16)
		...
		[worker:11]	Done 123: 	(r=9, c=4)
		[worker:10]	Done 124: 	(r=9, c=9)
		[worker:4]	Done 125: 	(r=10, c=10)
	(numWorkers = 8)
		...
		[worker:6]	Done 86: 	(r=7, c=1)
		[worker:7]	Done 87: 	(r=7, c=3)
		[worker:3]	Done 88: 	(r=7, c=2)
	*/

	/* Well, finally got a pass! with numWorkers = 8
		[02/20/2022 16:18] 62 / 62 test cases passed.
		Status: Accepted
		Runtime: 2408 ms
		Memory Usage: 8.2 MB

	And another one with numWorkers = 16
		[02/20/2022 16:24] 62 / 62 test cases passed.
		Status: Accepted
		Runtime: 2376 ms
		Memory Usage: 8.3 MB
	*/
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

	dispatcherBuf := rows*cols
	workDispatcher := make(chan *workItem, dispatcherBuf)

	bigBuf := rows*cols // might need to set higher
	var result chan string = make(chan string, bigBuf)
	var wg sync.WaitGroup
	var explored int64

	work := func(items <-chan *workItem, id int) {
		fmt.Printf("[worker:%d] started\n", id)
		defer wg.Done()

		for item := range items {
			used := createEmpty(rows, cols)
			explore(ctx, board, item.r, item.c, root, used, result)

			atomic.AddInt64(&explored, 1)
			fmt.Printf("[worker:%d]\tDone %d: \t(r=%d, c=%d)\n", id, atomic.LoadInt64(&explored), item.r, item.c)
		}
	}
	for i:=0; i<numWorkers; i++ {
		wg.Add(1)
		go work(workDispatcher, i)
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			workDispatcher <- &workItem{i, j}
		}
	}
	close(workDispatcher)
	
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
