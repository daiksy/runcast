# Weather CLI

Go言語で作られたシンプルな天気予報CLIツールです。

## 機能

- 指定した都市の現在の天気情報を取得
- 気温、湿度、風速、天気状況を表示
- 温度単位の選択（摂氏、華氏、ケルビン）

## 使用方法

### インストール

```bash
go build -o weather-cli
```

### 基本的な使用方法

```bash
# 東京の天気を取得
./weather-cli -city Tokyo

# 大阪の天気を取得
./weather-cli -city Osaka

# 華氏で表示
./weather-cli -city "New York" -units imperial
```

### オプション

- `-city`: 都市名を指定（デフォルト: Tokyo）
- `-units`: 温度単位を指定（metric, imperial, kelvin）

## API Key の設定

実際に使用する場合は、OpenWeatherMap のAPI Keyが必要です：

1. https://openweathermap.org/api でアカウントを作成
2. API Keyを取得
3. 環境変数 `OPENWEATHER_API_KEY` に設定

## 注意事項

- 現在はデモ用のAPIキーを使用しているため、実際の天気データは取得できません
- 本格的に使用する場合は、有効なAPI Keyを設定してください