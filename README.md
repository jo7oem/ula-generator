# ula-generator
[![Go](https://github.com/jo7oem/ula-generator/actions/workflows/go.yml/badge.svg)](https://github.com/jo7oem/ula-generator/actions/workflows/go.yml)
ula-generator は IPv6 ネットワークのための Unique Local Address (ULA) のプレフィックスを生成するツールです。
[RFC 4193](https://tools.ietf.org/html/rfc4193) に基づいています。


ula-generator is a tool to generate a unique local address (ULA) prefix for IPv6 networks. 
It is based on the [RFC 4193](https://tools.ietf.org/html/rfc4193) .

## 使い方 / Usage

```bash
$ ula-generator
fd63:1463:47c8:3dbf:8f00::/56
```

### 詳細 / Detail
#### ULA生成方法 / How to generate ULA
ULAの生成方法は RFC 4193 に基づいています。
1. 現在時刻を取得する
2. 1.で取得した現在時刻を NTP Timestamp Format([RFC 5905](https://datatracker.ietf.org/doc/html/rfc5905)) 形式に変換する
3. EUI-64を生成するためのもととなるインターフェースを選出する
   - 以下の優先度で選出する
     1. IPv6 のグローバルアドレスを持つインターフェース
     2. IPv6 のULAアドレスを持つインターフェース
     3. IPv6 のリンクローカルアドレスを持つインターフェース
     4. MACアドレスを持つインターフェース
   - 同じ優先度のインターフェースが複数ある場合は、その中から最も若番のインターフェースを選出する
4. 3.で選出したインターフェースよりEUI-64を生成する
   1. リンクローカルアドレスを持つ場合は、そのアドレスの下位64bitを使用する
   2. 64bitのMACアドレスを持つ場合は、そのMACアドレスを使用する
   3. それ以外の場合は、インターフェースのMACアドレスを取得し、そのMACアドレスからEUI-64を生成する
5. 2.と 4.で生成した値を結合する
6. 5.で生成した値をSHA-1でハッシュ化する
7. fd::/8 と 6. で生成した値の末尾40ビットを結合しULAのプレフィックスとして採用する
