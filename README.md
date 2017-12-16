% Elite:Dangerous Rescue Summary tool
% Igaguri

## 概要

Elite:Dangerousでの救助活動の統計をジャーナルファイルから生成するコマンドラインツールです。
ダブルクリックで実行されます。
64ビットWindows専用です


## 開発者向け情報

Makefileを使うには、Mingw版GNU Makeが入っているbash互換環境が必要です。
[Git for Windows](https://git-scm.com/download/win)のGit Bashを導入し、Makeを別途インストールした環境で動作確認しています。

もしくは、手動で`go build -o rescueSummary.exe main.go`を実行し、pack-resoyurceファイル内のファイルと一緒のディレクトリに入れてください。
