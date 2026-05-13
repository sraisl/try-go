# tries

`tries` ist ein kleines Terminal-Werkzeug zum schnellen Wechseln in Projekt- oder Experiment-Verzeichnisse unter `~/src/tries`.

Die Anwendung zeigt vorhandene Verzeichnisse an, filtert sie während der Eingabe und kann bei fehlenden Treffern direkt ein neues Verzeichnis mit Datumspräfix anlegen.

## Funktionen

- Verzeichnisse unter `~/src/tries` durchsuchen
- Treffer live nach eingegebenen Suchbegriffen filtern
- mit `Enter` in das gewählte Verzeichnis wechseln
- bei keinen Treffern ein neues Verzeichnis anlegen
- neue Verzeichnisse bekommen ein Präfix im Format `YYYY-MM-DD-name`
- Verzeichnisse nach zweifacher Bestätigung mit `Backspace` löschen
- farbige Terminal-Oberfläche mit Bubble Tea und Lip Gloss

## Installation

Die Shell muss den Verzeichniswechsel selbst ausführen. Deshalb wird die Go-Anwendung über eine `zsh`-Funktion eingebunden.

Binary installieren:

```zsh
go install .
```

Das installiert `tries` standardmäßig nach `$(go env GOPATH)/bin`, also meist `~/go/bin`. Dieses Verzeichnis muss in deinem `PATH` liegen.

Alternativ systemweit nach `/usr/local/bin`:

```zsh
go build -o tries .
sudo install -m 0755 tries /usr/local/bin/tries
```

Einmalig in der aktuellen Shell:

`/path/to/try-go/tries.zsh` ist ein Platzhalter und muss auf deinen lokalen Repo-Pfad zeigen.

```zsh
source /path/to/try-go/tries.zsh
```

Dauerhaft in `~/.zshrc`:

```zsh
source /path/to/try-go/tries.zsh
```

Danach kann das Werkzeug mit folgendem Befehl gestartet werden:

```zsh
tries
```

## Bedienung

- Tippen: Suche eingeben
- `Backspace`: Suchtext löschen
- `↑` / `↓`: Auswahl bewegen
- `Enter`: Treffer öffnen oder neues Verzeichnis anlegen
- `Backspace` ohne Suchtext: ausgewähltes Verzeichnis zum Löschen markieren
- `Backspace` erneut: markiertes Verzeichnis löschen
- `Esc`: abbrechen oder Löschbestätigung zurücknehmen

## Entwicklung

Formatieren:

```zsh
gofmt -w main.go main_test.go
```

Prüfen:

```zsh
go vet ./...
```

Tests:

```zsh
go test ./...
```

Build:

```zsh
go build -o tries .
```

## Releases

GitHub Actions erstellt bei Tags automatisch einen Release mit downloadbaren Binaries.

Beispiel:

```zsh
git tag v0.1.0
git push origin v0.1.0
```

Gebaut werden Binaries für Linux, macOS und Windows, jeweils für `amd64` und `arm64`.
