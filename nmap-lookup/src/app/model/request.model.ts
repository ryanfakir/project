export class Host {
    public querykeys: string[] = [];
    constructor(list ?: string[]) {
        this.querykeys = list;
    }
}