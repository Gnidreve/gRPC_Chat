# TODO / Code-Review-Ergebnis

Stand: vollständige Analyse von Client (`application/`) und Server (`server/`) am 20.08.2026.

## Muss man sich nochmal angucken (kritisch)

- [ ] **Keine Authentifizierung — Identitäts-Spoofing möglich.** `Store.Join()` glaubt dem Client jede ID, jeden Nickname, jede Farbe blind. `User.id` wird in jeder Nachricht an alle Clients mitgeschickt (auch wenn die UI sie nicht anzeigt), daher kann sich jeder, der den Port erreicht, eine fremde ID schnappen und `SendMessage` damit aufrufen.
- [ ] **`broadcast()` in `server/internal/store/store.go` kann den ganzen Server blockieren.** Sendet synchron an alle Subscriber-Channels nacheinander (`ch <- ev`). Hängt ein Client (schlechte Verbindung, App eingefroren, Channel-Buffer bei 32 Events voll), blockiert `AddMessage()` — und damit *jeder* `SendMessage`-Call von *jedem* Nutzer — bis der langsame Reader wieder Platz macht.
- [ ] **Unbegrenztes Speicherwachstum.** `s.messages` wächst für immer, kein Cap, kein Trimming. Kein Limit auf Nachrichtenlänge.
- [ ] **Server läuft unverschlüsselt (Plaintext-TCP, kein TLS).** Bewusste Übergangslösung, um das Coolify/Traefik-h2c-Problem zu umgehen — aktuell ist aller Chat-Traffic im Klartext unterwegs.

## Sollte man sich angucken

- [ ] Kein Graceful Shutdown in `server/cmd/server/main.go` — jeder Redeploy kappt alle Verbindungen abrupt (Reconnect fängt's ab, sauberer wär's trotzdem).
- [ ] Kein Logging über eine Startzeile hinaus — im Fehlerfall komplett blind, keine Metriken, keine Sicht auf aktive Verbindungen.
- [ ] `downloadAndInstall` (`update_checker.dart`) prüft die APK-Integrität nicht (kein Hash-Vergleich gegen den GitHub-Release).
- [ ] Debug-Keystore-Signing ist nicht nur Bequemlichkeit, sondern eine Sicherheitsfrage — der Key ist öffentlich bekannt (Standard jeder Android-Dev-Umgebung).
- [ ] Fixe 8er-Farbpalette ohne Server-seitige Eindeutigkeit — ab dem 9. gleichzeitigen Nutzer sind Farbdopplungen garantiert.
- [ ] `ChatInputBar` leert das Textfeld sofort, auch wenn `sendMessage` später fehlschlägt — getippter Text ist dann weg.
- [ ] Reconnect-Fehler und Sende-Fehler teilen sich eine `_error`-Variable in `chat_screen.dart` — der jeweils letzte gewinnt, der andere verschwindet.

## Sicherheitslücken (Zusammenfassung)

1. Keine Auth → Impersonation möglich
2. Kein TLS zum Server → Traffic im Klartext mitlesbar
3. Kein Rate-Limiting → jeder kann fluten (Nachrichten, Join-Calls, Subscribe-Verbindungen)
4. Keine Eingabevalidierung (Nachrichtenlänge, Nickname-Länge, Farbformat)
5. Debug-Signing-Key ist öffentlich bekannt
6. Kein Integritäts-Check beim Update-Download

Für einen privaten Chat unter Freunden überschaubares Risiko — bei Wachstum oder fremden Nutzern wären 1–4 die ersten Baustellen.

## Code-Qualität

Insgesamt sauber und konsistent. Go-Seite: gute Trennung (`store`/`chatserver`/`cmd`), durchdachtes Concurrency-Handling (bis auf den Broadcast-Bug oben), gute Testabdeckung für die kniffligen Teile (Identität, Presence, History-Replay). Dart-Seite: konsistentes State-Machine-Muster in `AppGate`, klare Widget-Trennung, Kommentare erklären das Warum statt das Was.

- [ ] **Dart-Seite hat quasi keine automatisierten Tests** (nur ein Widget-Test). `UpdateChecker`, `UserProfileStore`, Farb-Parsing, Reconnect-Backoff sind ungetestet, während die Go-Seite gut abgedeckt ist.

## Verbesserungsvorschläge für zukünftige Versionen

- [ ] Redis-Migration (ohnehin geplant, jetzt dringlicher wegen Speicherwachstum + History-Verlust bei jedem Redeploy)
- [ ] Echtes TLS zum Server (Traefik-Label fixen oder TLS direkt im Go-Server via autocert)
- [ ] Leichtgewichtige Auth (z.B. ein bei Join ausgestelltes Secret, bei folgenden Calls geprüft)
- [ ] Rate-Limiting + Nachrichtenlängen-Limit serverseitig
- [ ] Push-Benachrichtigungen (aktuell: komplett geschlossene App bekommt nichts mit)
- [ ] History-Pagination statt immer alles zu replayen
- [ ] Nickname/Farbe nachträglich änderbar
- [ ] Server-Metriken/Logging
- [ ] Automatisierte Tests für die Dart-Service-Schicht
