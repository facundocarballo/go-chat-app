import { Group } from "../group";
import { User } from "../user";

export class Message {
  id: number;
  sender: User;
  recipient: User | Group;
  message: string;
  sent: string;

  constructor(
    _id: number,
    _sender: User,
    _recipient: User | Group,
    _message: string,
    _sent: string
  ) {
    this.id = _id;
    this.sender = _sender;
    this.recipient = _recipient;
    this.message = _message;
    this.sent = _sent;
  }
}
