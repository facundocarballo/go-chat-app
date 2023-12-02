import { User } from "../models/user";

export interface IContext {
  user?: User;

  setUser: (_user: User | undefined) => void;
}
