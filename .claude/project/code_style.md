# sq-ops-backend コードスタイル

## 全般
- ファイルの末尾には必ず改行を入れる

## OpenAPI
- operationId は、`{HTTPメソッド}{リソース名}` 形式で命名する
  ex) `GetUsers`, `PostUsers`
- required などの配列値はインライン形式で記述する ["item1", "item2"]
  ex) `required: [ "id", "name" ]`

## UseCase, Repository インターフェースの関数名
- 単数取得
  - `Get` 形式で命名する
    - 例: `Get`, `GetByName`
- 複数取得
  - `List` 形式で命名する
    - 例: `List`, `ListByStatus`
- 作成
  - `Create` 形式で命名する
    - 例: `Create`
- 更新
  - `Update` 形式で命名する
    - 例: `Update`, `UpdateStatus`
- 削除
  - `Delete` 形式で命名する
    - 例: `Delete`, `DeleteByIDs`

- **!Important 既存のソースコードを参考にする**
