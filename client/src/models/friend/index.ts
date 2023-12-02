import { GetParticularUser } from "@/src/handlers/server";
import { Message } from "../message";
import { User } from "../user";

export class Friend {
  id: number;
  user_a: User;
  user_b: User;
  start_friendship: string;
  messages?: Message[]

  constructor(
    _id: number,
    _user_a: User,
    _user_b: User,
    _strart_friendship: string
  ) {
    this.id = _id;
    this.user_a = _user_a;
    this.user_b = _user_b;
    this.start_friendship = _strart_friendship;
  }

  static async GetFriendsOf(owner: User, data: any[]): Promise<User[]> {
    let friends: User[] = [];
    for (let i = 0; i < data.length; i++) {
      let user: User|undefined = undefined;
      if (data[i].user_a == owner.id) {
        user = await GetParticularUser(data[i].user_b);
      } else {
        user = await GetParticularUser(data[i].user_a);
      }
      friends.push(user!)
    }
    return friends;
  }
}
