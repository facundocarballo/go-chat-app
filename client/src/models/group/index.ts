import { Message } from "../message";
import { User } from "../user";

export class Group {
    id: number;
    name: string;
    description: string;
    owner: User;
    sent: string;
    users: User[];
    messages?: Message[]
  
    constructor(
      _id: number,
      _name: string,
      _description: string,
      _owner: User,
      _sent: string
    ) {
      this.id = _id;
      this.name = _name;
      this.description = _description;
      this.owner = _owner;
      this.sent = _sent;
      this.users = [];
    }
}