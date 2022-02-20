package main

import (
	"fmt"
	"testing"
	"time"

	"concurrent.leetcode.com/problems"
	"github.com/stretchr/testify/assert"
)

func TestFindWords(t *testing.T) {
	// Case1 -- small
	expectedLenResult1 := 2
	board1:= [][]byte{
		{'o', 'a', 'a', 'n'},
		{'e', 't', 'a', 'e'},
		{'i', 'h', 'k', 'r'},
		{'i', 'f', 'l', 'v'},
	}
	words1 := []string{"oath", "pea", "eat", "rain"}

	//Case2 -- big
	expectedLenResult2 := 10
	board2 := [][]byte{
		{'a','a','a','a','a','a','a','a','a','a','a','a'},
		{'a','a','a','a','a','a','a','a','a','a','a','a'},
		{'a','a','a','a','a','a','a','a','a','a','a','a'},
		{'a','a','a','a','a','a','a','a','a','a','a','a'},
		{'a','a','a','a','a','a','a','a','a','a','a','a'},
		{'a','a','a','a','a','a','a','a','a','a','a','a'},
		{'a','a','a','a','a','a','a','a','a','a','a','a'},
		{'a','a','a','a','a','a','a','a','a','a','a','a'},
		{'a','a','a','a','a','a','a','a','a','a','a','a'},
		{'a','a','a','a','a','a','a','a','a','a','a','a'},
		{'a','a','a','a','a','a','a','a','a','a','a','a'},
		{'a','a','a','a','a','a','a','a','a','a','a','a'},
	}
	words2 := []string{"lllllll","fffffff","ssss","s","rr","xxxx","ttt","eee","ppppppp","iiiiiiiii","xxxxxxxxxx","pppppp",
	"xxxxxx","yy","jj","ccc","zzz","ffffffff","r","mmmmmmmmm","tttttttt","mm","ttttt","qqqqqqqqqq","z","aaaaaaaa","nnnnnnnnn",
	"v","g","ddddddd","eeeeeeeee","aaaaaaa","ee","n","kkkkkkkkk","ff","qq","vvvvv","kkkk","e","nnn","ooo","kkkkk","o",
	"ooooooo","jjj","lll","ssssssss","mmmm","qqqqq","gggggg","rrrrrrrrrr","iiii","bbbbbbbbb","aaaaaa","hhhh","qqq","zzzzzzzzz",
	"xxxxxxxxx","ww","iiiiiii","pp","vvvvvvvvvv","eeeee","nnnnnnn","nnnnnn","nn","nnnnnnnn","wwwwwwww","vvvvvvvv","fffffffff",
	"aaa","p","ddd","ppppppppp","fffff","aaaaaaaaa","oooooooo","jjjj","xxx","zz","hhhhh","uuuuu","f","ddddddddd","zzzzzz",
	"cccccc","kkkkkk","bbbbbbbb","hhhhhhhhhh","uuuuuuu","cccccccccc","jjjjj","gg","ppp","ccccccccc","rrrrrr","c","cccccccc",
	"yyyyy","uuuu","jjjjjjjj","bb","hhh","l","u","yyyyyy","vvv","mmm","ffffff","eeeeeee","qqqqqqq","zzzzzzzzzz","ggg",
	"zzzzzzz","dddddddddd","jjjjjjj","bbbbb","ttttttt","dddddddd","wwwwwww","vvvvvv","iii","ttttttttt","ggggggg","xx",
	"oooooo","cc","rrrr","qqqq","sssssss","oooo","lllllllll","ii","tttttttttt","uuuuuu","kkkkkkkk","wwwwwwwwww","pppppppppp",
	"uuuuuuuu","yyyyyyy","cccc","ggggg","ddddd","llllllllll","tttt","pppppppp","rrrrrrr","nnnn","x","yyy","iiiiiiiiii",
	"iiiiii","llll","nnnnnnnnnn","aaaaaaaaaa","eeeeeeeeee","m","uuu","rrrrrrrr","h","b","vvvvvvv","ll","vv","mmmmmmm","zzzzz",
	"uu","ccccccc","xxxxxxx","ss","eeeeeeee","llllllll","eeee","y","ppppp","qqqqqq","mmmmmm","gggg","yyyyyyyyy","jjjjjj",
	"rrrrr","a","bbbb","ssssss","sss","ooooo","ffffffffff","kkk","xxxxxxxx","wwwwwwwww","w","iiiiiiii","ffff","dddddd",
	"bbbbbb","uuuuuuuuu","kkkkkkk","gggggggggg","qqqqqqqq","vvvvvvvvv","bbbbbbbbbb","nnnnn","tt","wwww","iiiii","hhhhhhh",
	"zzzzzzzz","ssssssssss","j","fff","bbbbbbb","aaaa","mmmmmmmmmm","jjjjjjjjjj","sssss","yyyyyyyy","hh","q","rrrrrrrrr",
	"mmmmmmmm","wwwww","www","rrr","lllll","uuuuuuuuuu","oo","jjjjjjjjj","dddd","pppp","hhhhhhhhh","kk","gggggggg","xxxxx",
	"vvvv","d","qqqqqqqqq","dd","ggggggggg","t","yyyy","bbb","yyyyyyyyyy","tttttt","ccccc","aa","eeeeee","llllll",
	"kkkkkkkkkk","sssssssss","i","hhhhhh","oooooooooo","wwwwww","ooooooooo","zzzz","k","hhhhhhhh","aaaaa","mmmmm"}

	testCases := []struct {
		desc string
		board [][]byte
		words []string
		expectedLenResult int
	} {
		{desc: "samll data set",
			board: board1, words: words1, expectedLenResult: expectedLenResult1},
		{desc: "big data set",
			board: board2, words: words2, expectedLenResult: expectedLenResult2},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t0 := time.Now()
			fmt.Printf("[%s] starting test case [%s]\n", t0, tc.desc)
	
			var result []string
			result, palDur := problems.FindWordsInParallel(tc.board, tc.words)
			fmt.Printf(" Parallel\t == pal_uSec: %v, total_uSec: %v, result: %v\n", palDur.Microseconds(), time.Since(t0).Microseconds(), result)
			assert.Equal(t, tc.expectedLenResult, len(result), tc.desc)
	
			result, seqDur := problems.FindWordsInSequence(tc.board, tc.words)
			fmt.Printf(" Sequential\t == seq_uSec: %v, total_uSec: %v, result: %v\n", seqDur.Microseconds(), time.Since(t0).Microseconds(), result)
			assert.Equal(t, tc.expectedLenResult, len(result), tc.desc)
		})
	}
}
