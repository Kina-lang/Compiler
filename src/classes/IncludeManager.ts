export class IncludeManager {
  private readonly _includes: Set<string> = new Set();

  constructor() {}

  public add(includePath: string) {
    this._includes.add(includePath);
  }

  public getAll(): Set<string> {
    return this._includes;
  }
}
