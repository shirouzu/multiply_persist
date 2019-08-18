package main

//  2019/08/18 H.Shirouzu https://github.com/shirouzu/multiply_persist
//
// 「任意の自然数の各桁を、一桁になるまで掛け算する回数の最大回数とその数を示せ」
// という問題を解く。
// （少し変形して、n桁の素数の積の形で解く）
// 
// https://twitter.com/KeikoUTorii/status/1161322092165337088
//
//  ref. "Smallest number of multiplicative persistence n"
//   https://oeis.org/A003001

import (
	"fmt"
	"time"
	"sync"
	"math/big"
	"runtime"
)

// 出力
//  試行情報 ... 桁数:最大乗算数:試行回数
//  最大回数 ... Found(最大乗算数): 数字列

var NMAX int = 0      // 本当は排他制御必要
var NUMS [][]int      // 試行する素数グループ
var EDIC [][]*big.Int // 1-9累乗辞書（EDIC[2][4] = 2^4）
var SUMMARY_ONLY bool // trueで試行出力のみに

type ANS struct {
	kind int
	idx [3]int
}

func main() {
	start_tick := time.Now()

	thr_num := runtime.NumCPU()-1 // スレッド数（マルチスレッド時、序盤は出力順が乱れます）
	col_start := 1        // 開始桁数
	col_max := 500        // 最大桁数
	SUMMARY_ONLY = false

	if thr_num == 0 {
		thr_num = 1
	}
	fmt.Printf("Start (%d threads)\n", thr_num)

	NUMS = [][]int{{7,3,2}, {7,5,3}} // 2 と 5 の組み合わせは次で必ず消えるため除外

	setup_edic(col_max) // 1-9について、最大桁数の累乗辞書を作る
	wg := new(sync.WaitGroup)

	for i:=0; i < thr_num; i++ {
		wg.Add(1)
		go thread_proc(i, col_start, col_max, thr_num, wg)
	}
	wg.Wait()
	fmt.Printf("\nfin (%.1f sec)", float64(time.Now().Sub(start_tick))/1000000000)
}

// 1-9の累乗辞書作成
func setup_edic(col_max int) {
	EDIC = make([][]*big.Int, 10)
	m := big.NewInt(int64(0))

	for i:=1; i < 10; i++ {
		EDIC[i] = make([]*big.Int, col_max)
		for j:=0; j < col_max; j++ {
			v := big.NewInt(int64(i))
			e := big.NewInt(int64(j))
			v.Exp(v, e, m)
			EDIC[i][j] = v
		}
	}
}

// スレッドメイン処理
func thread_proc(idx int, start int, col_max int, thr_num int, wg *sync.WaitGroup) {
	for cols:=start+idx; cols < col_max; cols+=thr_num {
		ans := ANS{}
		r, cnt := combi(cols, 1, &ans)
		fmt.Printf("%4d:%2d:%-8d", cols, r, cnt)
	}
	wg.Done()
}

// 素数累乗のcols桁の全組み合わせを試行
func combi(cols int, dep int, ans *ANS) (int, int) {
	sum := big.NewInt(1)
	max_dep := 0
	cnt := 0

	for i:=0; i <= cols; i++ {
		ans.idx[0]=i
		for j:=0; i+j <= cols; j++ {
			ans.idx[1]=j
			k := cols - (i+j)
			ans.idx[2]=k
			for n, nums:= range NUMS {
				ans.kind=n
				sum.Set(EDIC[nums[0]][i])
				sum.Mul(sum, EDIC[nums[1]][j])
				sum.Mul(sum, EDIC[nums[2]][k])
				r := mul_eva(dep, sum, ans)
				if r > max_dep {
					max_dep = r
				}
				cnt++
			}
		}
	}
	return	max_dep, cnt
}

// 数字文字列を乗算
func str2sum(s string) *big.Int {
	sum := big.NewInt(1)
	bs := []byte(s)

	for idx:=0; idx != -1; {
		idx = str2e(bs, idx, sum)
	}

//	fmt.Printf("%s -> %s\n", s, sum.String())
	return	sum
}

// 数字列から、１つの数字を選び、数字累乗に変換し、sumに掛ける
//  str2e('2411124', 0, 1) --> ('.4111.4', 1, 4)  // 2^2=4 をsumに乗算、次回位置1を返す
//  str2e('.4111.4', 1, 4) --> ('..111..', 2, 64) // 4^2=16をsumに乗算、次回位置2を返す
//  str2e('..111..', 2, 4) --> ('.......',-1, 64) // 1^3=1 をsumに乗算、完了(-1)を返す
func str2e(bs []byte, idx int, sum *big.Int) int {
	var v byte = 0
	vnum := 0
	next_idx := -1

	for ; idx < len(bs); idx++ {
		c := bs[idx]
		if c == 0 {
			continue
		}
		if c == '0' {
			sum.SetInt64(int64(0))
			return -1
		}
		if v == 0 {
			v = c
			vnum = 1
			bs[idx] = 0
		} else {
			if v == c {
				vnum++
				bs[idx] = 0
			} else {
				if next_idx == -1 {
					next_idx = idx
				}
			}
		}
	}
	if vnum > 0 {
		sum.Mul(sum, EDIC[int(v - '0')][vnum])
	}

	return next_idx
}

// 乗算結果の評価
func mul_eva(dep int, sum *big.Int, ans *ANS) int {
	s := sum.String()
	new_cols := len(s)

	if new_cols == 1 {
		if dep >= NMAX {
			NMAX = dep
			if !SUMMARY_ONLY {
				put_data("Found", dep, ans)
			}
		}
	} else {
		new_dep := dep + 1
		new_sum := str2sum(s)
		//put_data("Next ", dep, ans)

		dep = mul_eva(new_dep, new_sum, ans)
	}
	return dep
}

// （主に）解答出力
func put_data(tag string, dep int, ans *ANS) {
	s := fmt.Sprintf("\n%s(%d): ", tag, dep)
	for i, v := range ans.idx {
		for j:=0; j < v; j++ {
			s += fmt.Sprintf("%d", NUMS[ans.kind][i])
		}
	}
	fmt.Printf("%s\n", s)
}

