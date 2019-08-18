# 「任意の自然数の各桁を、一桁になるまで掛け算する回数の最大回数とその数を示せ」

## 概要
Keiko ToriiさんのTweet https://twitter.com/KeikoUTorii/status/1161322092165337088 の問題を解いてみよう、と言う話。

 ref. "Smallest number of multiplicative persistence n"
  https://oeis.org/A003001

## アイデア
1. 少し変形して、n桁の素数の積（2,3,5,7）の形で解く

2. 数字に 2 と 5 が混ざると0が出て必ず次で終了するため、(2,3,7) と (3,5,7) の２つの文字プールに分けて、数字列を作る

3. 累乗の事前計算辞書を作る。（1^n ～ 9^n の計算済み辞書を作っておく。nは計算予定の最大桁数。例えば 2222452425 を 2^6 * 4^2 * 5^2 に変形した上で、計算済みの 2^6 と 4^2 と 5^2 を取り出して、2回の乗算だけで数字列の乗算を完了させる）

4. コア数-1のマルチスレッドで並列計算

## 実行方法
引数なしで、go run multiply_persist.go とやれば、

   試行情報 ... 桁数:最大乗算数:試行回数
   最大回数 ... Found(最大乗算数): 数字列

の２つの情報が出力されます（前者は、一つの桁が終わるごとに１つの出力が出ます）


## 性能
手元の i5 8600K (6core) では、5スレッドで 300桁で3秒、500桁で 25秒、といったところ。

## 考察
実際に実行してみると、n=29桁で11回以降、n=95で4、n=374で3 を最後に、それ以降は全て最大でも2回で計算が終わる。 https://twitter.com/shirouzu/status/1162944937308016646

これは例えば、1000桁の場合、全て2でも最初の乗算結果は300桁を越えており、この中に0を含まない確率は非常に小さい。また偶然 0 が無くとも、2 と 5 が同時に含まれると、次の計算で終了してしまう。

…という状況を11回以上も乗り越えるのは、むしろ桁が多い方が厳しくなっていく、という感じに見える。

## おまけ
私にとっては、最初の golang用習作プログラムなので、間違いがあればツッコミ歓迎。


オリジナル URL: https://github.com/shirouzu/multiply_persist


