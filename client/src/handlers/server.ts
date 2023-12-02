import { User } from "../models/user";
import axios from 'axios'

export const SERVER_URL = "http://localhost:3690/";

export const GetParticularUser = async (id: number): Promise<User|undefined> => {
    try {
        const {data} = await axios.get(SERVER_URL + `user?id=${id}`);
        return new User(data.id, data.name, data.email, data.created_at);
    } catch (err) {
        console.error(`Error getting this user with id: ${id}. ` + err);
        return
    }
}