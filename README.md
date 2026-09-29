# gotodo-lesson
参考記事[https://qiita.com/hrk_ym/items/c73c5ad41c92688c3b94]

- リポジトリ層の引数
1件取得 (GetTodoByID)

表現方法: nil（ポインタ）

意味: 「対象のデータが存在しない」

一覧取得 (GetTodos)

表現方法: []（空スライス）

意味: 「条件に該当するデータが0件」

## test

### service層 
- GetTodos 正常系    
- GetTodoByID 正常系  
- GetTodoByID エラー系
- CreateTodo 正常系   
- UpdateTodo 正常系   
- UpdateTodo エラー系 
- DeleteTodo 正常系   

### repository層
- GetTodoByIDの異常系はnodetsの方ではエラー型を作成してそれに変換して変えさせている
- テストの方でその仕様を先に書き(その時点ではテスト失敗)その後repositoryでNotFoundDataErrorを実装して再テスト(TDDの練習目的)
- nodetsの時と違うポイントとしてNotFoundDataErrorというクラスを作るのではなくエラー型の値を作り､それを変数(ErrNotFound)に入れる｡ 判定の際はerror.Is(err, ErrNotFound)を使う
- Goでもエラー型のストラクトなら作れんじゃねとは思う
- エラー値を使って判定 このエラー値ってなんぞや
- nodets の new NotFoundDataError()こいつ型として宣言したものをインスタンス化したやつ?
- Goの場合はちゃうんやと
```Go
var ErrNotFound = errors.New("todo not found")
```
- errors.New()はGo標準ライブラリが用意してるエラー型から値が作られる｡今回はerrors.New()にtodo not foundを渡し､todo not foundというメッセを持つエラー値を作成し､ErrNotFound(変数)に代入している?
- 最後にerrors.Is(err, ErrNotFound)で返却されたエラーはErrNotFoundとして使えるのかチェックする
- Goでもエラー型のストラクトなら作れんじゃねとは思う←作れるけど今回のID取得のエラーに特筆して情報をもたせる必要がないため､型まで作る必要なしと判断
- errors.New()などで変数を定義するエラーをセンチネルエラーという頭にErrという文字をつけるのがしきたり
- NotFoundDataErrorのような名前にしないのは型と勘違いしないため