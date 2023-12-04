export class Message {
  id: number;
  sender_id: number;
  recipient_id: number;
  message: string;
  sent: string;

  constructor(
    _id: number,
    _sender: number,
    _recipient: number,
    _message: string,
    _sent: string
  ) {
    this.id = _id;
    this.sender_id = _sender;
    this.recipient_id = _recipient;
    this.message = _message;
    this.sent = _sent;
  }

  static GetMessagesOf(data: any[]): Message[] {
    if (data == null) return [];
    let messages: Message[] = [];
    for (const m of data) {
      messages.push(new Message(m.id, m.user_id, m.to_id, m.message, m.sent));
    }
    return messages;
  }
}
