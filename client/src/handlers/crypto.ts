import { createHash } from "crypto";

export const GetHash = (txt: string): string => {
    const hash = createHash('sha256');
    hash.update(txt);
    return hash.digest('hex');
}