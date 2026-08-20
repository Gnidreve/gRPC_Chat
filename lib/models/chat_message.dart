/// Reine Datenhülle für eine Chat-Nachricht.
/// Keine Methoden, keine Logik – nur Felder fürs Templating.
class ChatMessage {
  final String? senderName; // null = eigene Nachricht
  final String text;
  final String time;
  final bool isOwn;

  const ChatMessage({
    this.senderName,
    required this.text,
    required this.time,
    required this.isOwn,
  });
}

/// Statische Beispiel-Daten, damit das Template ohne Backend aussieht
/// wie im Mockup. Später durch echte Daten ersetzen.
const List<ChatMessage> demoMessages = [
  ChatMessage(
    senderName: 'Mara',
    text: 'Hey, seid ihr schon am Set?',
    time: '09:14',
    isOwn: false,
  ),
  ChatMessage(
    text: 'Ja, sind gerade angekommen 👋',
    time: '09:15',
    isOwn: true,
  ),
  ChatMessage(
    senderName: 'Jonas',
    text: 'Perfekt, ich bring noch Kaffee mit',
    time: '09:16',
    isOwn: false,
  ),
  ChatMessage(
    text: 'Sehr gerne 🙏',
    time: '09:16',
    isOwn: true,
  ),
  ChatMessage(
    senderName: 'Mara',
    text: 'Bin in 5 Minuten da',
    time: '09:18',
    isOwn: false,
  ),
];
