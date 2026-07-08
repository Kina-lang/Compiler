export class CompilerMetrics {
  private readonly _metrics: Map<string, number[]> = new Map<string, number[]>(
    [],
  );

  constructor() {}

  public capture(captureName: string) {
    const t = performance.now();

    if (!this._metrics.has(captureName)) this._metrics.set(captureName, []);
    this._metrics.get(captureName)!.push(t);
  }

  public calculateDelta(captureName: string): number {
    const captures = this._metrics.get(captureName);
    if (!captures || captures.length < 2) return 0;

    return captures[captures.length - 1]! - captures[0]!;
  }
}
