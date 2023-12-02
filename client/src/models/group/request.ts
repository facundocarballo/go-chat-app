import { Group } from ".";
import { User } from "../user";

export class GroupRequest {
    id: number;
    user: User;
    group: Group;
    sent: string;
  
    constructor(
      _id: number,
      _user: User,
      _group: Group,
      _sent: string
    ) {
      this.id = _id;
      this.user = _user;
      this.group = _group;
      this.sent = _sent;
    }
}