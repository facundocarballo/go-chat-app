import { SERVER_URL } from "@/src/handlers/server";
import { Group } from "../group";
import axios, { AxiosRequestConfig } from "axios";
import { Friend } from "../friend";
import { Message } from "../message";

export class User {
  id: number;
  name: string;
  email: string;
  created_at: string;
  jwt?: string;
  groups?: Group[];
  friends?: User[];

  constructor(
    _id: number,
    _name: string,
    _email: string,
    _created_at: string,
    _jwt?: string
  ) {
    this.id = _id;
    this.name = _name;
    this.email = _email;
    this.created_at = _created_at;
    this.jwt = _jwt;
  }

  static async CreateUserWithData(
    _id: number,
    _name: string,
    _email: string,
    _created_at: string,
    _jwt?: string
  ): Promise<User> {
    const user = new User(_id, _name, _email, _created_at, _jwt);
    await user.GetFriends();
    await user.GetGroups();
    return user;
  }

  static async CreateUserWithMessages(
    owner: User,
    _id: number,
    _name: string,
    _email: string,
    _created_at: string,
    _jwt?: string
  ): Promise<User> {
    const user = new User(_id, _name, _email, _created_at, _jwt);

    return user;
  }

  static async Login(
    email: string,
    password: string
  ): Promise<User | undefined> {
    try {
      const {
        data: { user, jwt },
      } = await axios.post(SERVER_URL + "login", {
        email,
        password,
      });
      return await this.CreateUserWithData(
        user.id,
        user.name,
        user.email,
        user.created_at,
        jwt
      );
    } catch (err) {
      return;
    }
  }

  async GetMessages(to_id: string, isGroup: boolean): Promise<Message[] | undefined> {
    try {
      const axiosConfig = this._GetAxiosConfig();
      const query = isGroup ? "group-message?group_id=" : "user-message?friend_id="
      const { data } = await axios.get(
        SERVER_URL + `${query}${to_id}`,
        axiosConfig
      );
      const messages = Message.GetMessagesOf(data);
      return messages;
    } catch (err) {
      console.error(
        `Error getting the messages of user_a (${this.id}) and to_id (${to_id}).`
      );
      return;
    }
  }

  async GetFriends() {
    try {
      const axiosConfig = this._GetAxiosConfig();
      const { data } = await axios.get(SERVER_URL + "friends", axiosConfig);
      const f = await Friend.GetFriendsOf(this, data);
      this.friends = f;
    } catch (err) {
      console.error("Error getting all the friends of this user. " + err);
    }
  }

  async GetGroups() {
    try {
      const axiosConfig = this._GetAxiosConfig();
      const { data } = await axios.get(SERVER_URL + "group", axiosConfig);
      const groups = await Group.GetGroupsBy(data);
      this.groups = groups;
    } catch (err) {
      console.error("Error getting all the groups of this user. " + err);
    }
  }

  private _GetAxiosConfig(): AxiosRequestConfig {
    return {
      headers: {
        Authorization: `Bearer ${this.jwt}`,
        "Content-Type": "application/json",
      },
    };
  }
}
