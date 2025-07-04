# Weather CLI

Go言語で作られたシンプルな天気予報CLIツールです。

## 機能

- 指定した都市の現在の天気情報を取得
- 気温、湿度、風速、天気状況を表示
- 日本気象庁（JMA）のデータを使用
- **APIキー不要**

## 使用方法

### インストール

```bash
go build -o weather-cli
```

### 基本的な使用方法

```bash
# 東京の天気を取得
./weather-cli -city tokyo

# 大阪の天気を取得
./weather-cli -city osaka

# 札幌の天気を取得
./weather-cli -city sapporo
```

### オプション

- `-city`: 都市名を指定（デフォルト: Tokyo）

### 対応都市

- tokyo（東京）
- osaka（大阪）
- kyoto（京都）
- yokohama（横浜）
- nagoya（名古屋）
- sapporo（札幌）
- fukuoka（福岡）
- sendai（仙台）
- hiroshima（広島）
- naha（那覇）

## データソース

本アプリケーションは[Open-Meteo](https://open-meteo.com/)のJMA APIを使用して、日本気象庁のデータを取得しています。APIキーは不要で、無料で利用できます。

## 注意事項

- 現在は日本の主要都市のみ対応しています
- 天気情報は日本気象庁の高解像度データ（5km）を使用しています