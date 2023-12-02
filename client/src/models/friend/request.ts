import { User } from "../user";

export class FriendRequest {
    id: number;
    user_a: User;
    user_b: User;
    sent: string;
  
    constructor(
      _id: number,
      _user_a: User,
      _user_b: User,
      _sent: string
    ) {
      this.id = _id;
      this.user_a = _user_a;
      this.user_b = _user_b;
      this.sent = _sent;
    }
}