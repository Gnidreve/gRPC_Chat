# TODO / Code-Review-Ergebnis

Stand: vollständige Analyse von Client (`application/`) und Server (`server/`) am 20.08.2026.
Update 20.08.2026: Abgleich gegen die installierten Go- (`server/.agents/skills/`) und
Flutter/Dart-Skills (`application/.agents/skills/`) — Funde markiert mit „(Skill-Audit)“,
erledigte Punkte sind abgehakt.

## Muss man sich nochmal angucken (kritisch)

- [ ] **Keine Authentifizierung — Identitäts-Spoofing möglich.** `Store.Join()` glaubt dem Client jede ID, jeden Nickname, jede Farbe blind. `User.id` wird in jeder Nachricht an alle Clients mitgeschickt (auch wenn die UI sie nicht anzeigt), daher kann sich jeder, der den Port erreicht, eine fremde ID schnappen und `SendMessage` damit aufrufen.
- [x] ~~`broadcast()` in `server/internal/store/store.go` kann den ganzen Server blockieren.~~ Behoben (Skill-Audit, `golang-concurrency`): `broadcast()` sendet jetzt non-blocking (`select`/`default`) — ein hängender Subscriber verliert nur sein eigenes Event statt `AddMessage()`/`Subscribe()` für alle zu blockieren.
- [x] ~~Unbegrenztes Speicherwachstum.~~ Behoben: `s.messages` ist jetzt auf `maxHistory = 500` gedeckelt (älteste Nachrichten werden verworfen, verbundene Clients behalten ihre bereits empfangene History).
- [ ] **Server läuft unverschlüsselt (Plaintext-TCP, kein TLS).** Bewusste Übergangslösung, um das Coolify/Traefik-h2c-Problem zu umgehen — aktuell ist aller Chat-Traffic im Klartext unterwegs. `golang-grpc`-Skill: „TLS MUST be enabled in production" — weiterhin offen, keine sichere Lösung ohne Coolify-Proxy-Konfiguration in der Hand.

## Sollte man sich angucken

- [x] ~~Kein Graceful Shutdown in `server/cmd/server/main.go`.~~ Behoben (Skill-Audit, `golang-grpc`): `SIGINT`/`SIGTERM` werden jetzt abgefangen, `httpServer.Shutdown()` mit 15s-Timeout-Fallback auf `Close()`, Health-Status wird vor dem Draining auf `NOT_SERVING` gesetzt.
- [x] ~~Kein Logging über eine Startzeile hinaus.~~ Behoben: strukturiertes JSON-Logging via `log/slog`, plus Logging- und Panic-Recovery-Interceptors für alle Unary-/Stream-RPCs (`server/cmd/server/interceptors.go`).
- [ ] `downloadAndInstall` (`update_checker.dart`) prüft die APK-Integrität nicht (kein Hash-Vergleich gegen den GitHub-Release).
- [ ] Debug-Keystore-Signing ist nicht nur Bequemlichkeit, sondern eine Sicherheitsfrage — der Key ist öffentlich bekannt (Standard jeder Android-Dev-Umgebung).
- [ ] Fixe 8er-Farbpalette ohne Server-seitige Eindeutigkeit — ab dem 9. gleichzeitigen Nutzer sind Farbdopplungen garantiert.
- [ ] `ChatInputBar` leert das Textfeld sofort, auch wenn `sendMessage` später fehlschlägt — getippter Text ist dann weg.
- [ ] Reconnect-Fehler und Sende-Fehler teilen sich eine `_error`-Variable in `chat_screen.dart` — der jeweils letzte gewinnt, der andere verschwindet.
- [ ] **(Skill-Audit, `golang-grpc`)** Kein gRPC Health-Check-Service registriert — behoben: `grpc_health_v1` läuft jetzt mit (`server/cmd/server/main.go`), relevant falls Coolify/ein Orchestrator je Readiness-Probes nutzt.

## Neue Funde aus dem Skill-Audit (server, Go)

Alle unten als erledigt markierten Punkte sind bereits umgesetzt und über `go build`, `go vet`,
`go test -race -shuffle=on`, `golangci-lint run` (0 Findings) lokal verifiziert.

- [x] Kein `.golangci.yml`, kein Lint-Schritt in CI (`golang-lint`, `golang-continuous-integration`). → `server/.golangci.yml` (33 Linter, gosec/errcheck/paralleltest/etc.) + neuer Workflow `.github/workflows/test-server.yml` (build/vet/`go mod tidy`-Check/Test mit `-race -shuffle=on`/Lint/`govulncheck`, läuft auf jedem PR + Push nach `server/**`).
- [x] Kein Makefile (`golang-project-layout`). → `server/Makefile` mit `build`/`test`/`lint`/`lint-fix`/`fmt`/`run`.
- [x] `net.Listen` statt `(*net.ListenConfig).Listen` (noctx-Linter). → behoben in `main.go`.
- [x] `http.Server` ohne `ReadHeaderTimeout` (gosec G112, Slowloris-Risiko). → `ReadHeaderTimeout: 5s` gesetzt.
- [x] `int` → `int32`-Konvertierung des Online-Counts ohne Bounds-Check (gosec G115, `golang-safety`). → `clampToInt32()` mit explizitem Grenzwert-Check.
- [x] Veraltetes `golang.org/x/net/http2/h2c` (staticcheck SA1019 — als deprecated markiert). → durch die native `http.Server.Protocols`-API (Go 1.24+) ersetzt; die x/net-Abhängigkeit ist jetzt nur noch indirekt (Transitiv-Dependency von gRPC).
- [x] Fehlende Doc-Kommentare auf exportierten Typen/Funktionen (revive: `package-comments`, `exported`). → ergänzt in `main.go`, `chatserver.go`, `store.go`.
- [x] Tests: `err != ErrUnknownUser` statt `errors.Is` (errorlint), fehlendes `t.Parallel()` in allen Testfunktionen (paralleltest), ungeprüfte `lis.Close()`/`conn.Close()` in `t.Cleanup` (errcheck). → alle behoben.
- [x] `os.Exit()` nach `defer stop()` in `main()` — der Defer würde nie laufen (gocritic `exitAfterDefer`). → `main()` in `run() error` extrahiert, `os.Exit` erst nach Rückkehr aus `run()`.
- [x] Keine Interceptors für Logging/Panic-Recovery (`golang-grpc`). → `server/cmd/server/interceptors.go`: Unary- und Stream-Varianten für beides.
- [ ] **(offen, bewusst zurückgestellt)** Kein `gosec`/`govulncheck` als eigenständiger CI-Schritt außerhalb des Lint-Workflows — `govulncheck` läuft jetzt zwar in `test-server.yml`, aber ein separates CodeQL-/Security-Workflow (`golang-continuous-integration`-Skill) wurde nicht aufgesetzt — für ein Projekt dieser Größe aktuell nicht verhältnismäßig, aber vorgemerkt.
- [ ] **(offen, bewusst zurückgestellt)** Rate-Limiting (`golang.org/x/time/rate`, `golang-security`-Skill) — nicht umgesetzt, da das ohne Auth (siehe oben) wenig Wirkung hätte; gehört zusammen mit der Auth-Frage entschieden.
- [ ] **(offen, bewusst nicht umgesetzt)** Volle Observability-Suite (Prometheus-Metriken, OpenTelemetry-Tracing, Grafana-Dashboards laut `golang-observability`-Skill) — für einen privaten Chat mit einer Handvoll Nutzern deutlich über das gerechtfertigte Maß hinaus; strukturiertes Logging (siehe oben) deckt den aktuellen Bedarf.

## Neue Funde aus dem Skill-Audit (application, Flutter/Dart)

- [x] **Ressourcen-Leak:** `downloadAndInstall()` in `update_checker.dart` erzeugte pro Aufruf einen neuen `http.Client()`, der nie geschlossen wurde — bei jedem "Erneut versuchen" nach einem fehlgeschlagenen Download ein weiterer Leak. → behoben: `client.close()` in `finally`.
- [x] **Fehlender Status-Check:** derselbe Download-Pfad prüfte `response.statusCode` nie, bevor der Stream in eine `.apk`-Datei geschrieben und dem System-Installer übergeben wurde — ein GitHub-Ratelimit/Fehler (403/5xx) hätte eine HTML-Fehlerseite als "APK" installieren lassen. → behoben: wirft jetzt `HttpException` bei Status ≠ 200 (vom bestehenden `try`/`catch` in `update_screen.dart` bereits abgefangen).
- [x] `MediaQuery.of(context).size`/`.padding` statt `MediaQuery.sizeOf`/`.paddingOf` in `message_bubble.dart` (pro Nachricht gebaut!) und `chat_input_bar.dart` — unnötige Rebuilds bei jedem `MediaQuery`-Wechsel (u.a. Tastatur-Ein-/Ausblenden), nicht nur bei Breiten-/Padding-Änderungen. → behoben, gleiche Werte, engere Dependency.
- [x] Keine `Key`s auf interaktiven Widgets (Sende-Button, Nachrichtenfeld, Nickname-Feld, Submit-Button) — blockiert stabile Widget-/Integrationstests. → `ValueKey`s ergänzt (`chat_send_button`, `chat_message_field`, `profile_nickname_field`, `profile_submit_button`).
- [ ] **(offen, braucht Design-Entscheidung)** `UpdateChecker` injiziert `http.Client` nicht über den Konstruktor — macht die Klasse praktisch untestbar, ist aber der eigentliche Grund für die fehlende Testabdeckung (siehe unten). Kleiner API-Schnitt, betrifft aber Call-Sites.
- [ ] **(offen, niedrige Priorität)** GitHub-Release-JSON wird in `check()` ad-hoc geparst statt über ein `fromJson`-Modell — kein Bug, aber unnötig untestbar isoliert vom Netzwerk-Call.
- [ ] **(offen, braucht Abwägung)** `analysis_options.yaml` setzt keine `analyzer: language: strict-casts/strict-inference/strict-raw-types` — würde die manuellen `as String?`/`as Map<String, dynamic>`-Casts in `update_checker.dart` statisch absichern, dürfte aber neue Lint-Findings anderswo aufdecken (eigener Aufräum-Task).
- [ ] **(bewusst nicht umgesetzt)** Layered MVVM/Repository-Architektur laut `flutter-apply-architecture-best-practices`-Skill — für 7 Screens/Widgets und 3 Services aktuell nicht gerechtfertigt, explizit als bewusste Entscheidung vermerkt statt als Lücke.
- Geprüft und sauber: alle `TextEditingController`/`AnimationController`/`StreamSubscription`/`Timer`/`ScrollController`-Instanzen werden korrekt disposed (`chat_input_bar.dart`, `profile_setup_screen.dart`, `online_indicator.dart`, `chat_screen.dart`, `app_gate.dart`).

## Sicherheitslücken (Zusammenfassung)

1. Keine Auth → Impersonation möglich
2. Kein TLS zum Server → Traffic im Klartext mitlesbar
3. Kein Rate-Limiting → jeder kann fluten (Nachrichten, Join-Calls, Subscribe-Verbindungen)
4. Keine Eingabevalidierung (Nachrichtenlänge, Nickname-Länge, Farbformat)
5. Debug-Signing-Key ist öffentlich bekannt
6. Kein Integritäts-Check beim Update-Download

Für einen privaten Chat unter Freunden überschaubares Risiko — bei Wachstum oder fremden Nutzern wären 1–4 die ersten Baustellen.

## Code-Qualität

Insgesamt sauber und konsistent. Go-Seite: gute Trennung (`store`/`chatserver`/`cmd`), durchdachtes Concurrency-Handling (Broadcast-Bug jetzt behoben, siehe oben), gute Testabdeckung für die kniffligen Teile (Identität, Presence, History-Replay), jetzt zusätzlich mit CI-Absicherung (Lint, Race-Tests, `go mod tidy`-Check, `govulncheck`). Dart-Seite: konsistentes State-Machine-Muster in `AppGate`, klare Widget-Trennung, Kommentare erklären das Warum statt das Was.

- [ ] **Dart-Seite hat quasi keine automatisierten Tests** (nur ein Widget-Test). `UpdateChecker`, `UserProfileStore`, Farb-Parsing, Reconnect-Backoff sind ungetestet, während die Go-Seite gut abgedeckt ist. Die neu ergänzten `Key`s und die Injectable-Client-Frage oben sind Voraussetzungen, um das anzugehen.

## Verbesserungsvorschläge für zukünftige Versionen

- [ ] Redis-Migration (ohnehin geplant, jetzt weniger dringlich seit History gedeckelt ist, aber Speicherverlust bei jedem Redeploy bleibt)
- [ ] Echtes TLS zum Server (Traefik-Label fixen oder TLS direkt im Go-Server via autocert)
- [ ] Leichtgewichtige Auth (z.B. ein bei Join ausgestelltes Secret, bei folgenden Calls geprüft)
- [ ] Rate-Limiting + Nachrichtenlängen-Limit serverseitig
- [ ] Push-Benachrichtigungen (aktuell: komplett geschlossene App bekommt nichts mit)
- [ ] History-Pagination statt immer alles zu replayen
- [ ] Nickname/Farbe nachträglich änderbar
- [ ] Automatisierte Tests für die Dart-Service-Schicht (`UpdateChecker` mit injiziertem `http.Client`, `UserProfileStore`, Farb-Parsing, Reconnect-Backoff)
- [ ] `analysis_options.yaml` um `strict-casts`/`strict-inference`/`strict-raw-types` erweitern (eigener Cleanup-Task)
