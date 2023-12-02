import { SERVER_URL } from "@/src/handlers/server";
import { Group } from "../group";
import axios, { AxiosRequestConfig } from "axios";
import { Friend } from "../friend";

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
    await user._GetFriends()
    return user;
  };

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
      return await this.CreateUserWithData(user.id, user.name, user.email, user.created_at, jwt);
    } catch (err) {
      return;
    }
  }

  async _GetFriends() {
    try {
      const axiosConfig = this._GetAxiosConfig();
      const { data } = await axios.get(SERVER_URL + "friends", axiosConfig);
      const f = await Friend.GetFriendsOf(this, data);
      this.friends = f;
    } catch (err) {
      console.error("Error getting all the friends of this user. " + err);
    }
  }

  private async _GetGroups() {}

  private _GetAxiosConfig(): AxiosRequestConfig {
    return {
      headers: {
        Authorization: `Bearer ${this.jwt}`,
        "Content-Type": "application/json",
      },
    };
  }
}
