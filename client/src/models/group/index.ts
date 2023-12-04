import axios from "axios";
import { Message } from "../message";
import { User } from "../user";
import { SERVER_URL } from "@/src/handlers/server";

export class Group {
  id: number;
  name: string;
  description: string;
  owner: User;
  sent: string;
  users: User[];
  messages?: Message[];

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

  static async GetGroupsBy(data: any[]): Promise<Group[]> {
    let groups: Group[] = [];
    for (let i = 0; i < data.length; i++) {
      let group: Group | undefined = undefined;
      group = new Group(
        data[i].id,
        data[i].name,
        data[i].description,
        data[i].owner,
        data[i].sent
      );
      group.users = await this.GetGroupUsers(data[i].id);
      groups.push(group!);
    }
    return groups;
  }

  static async GetGroupUsers(id: number): Promise<User[]> {
    let users: User[] = [];
    try {
      const { data } = await axios.get(
        SERVER_URL + `group-users?group_id${id}`
      );
      for (const u of data) {
        users.push(
          new User(data.id, data.name, data.emial, new Date().toISOString())
        );
      }
    } catch (err) {
      console.error("Error reading the users of this group: ", id);
    }
    return users;
  }
}
