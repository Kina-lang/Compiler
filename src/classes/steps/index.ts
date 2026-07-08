import { BuildastStep } from "./BuildastStep";
import { ReadfileStep } from "./ReadfileStep";
import { TokenizeStep } from "./TokenizeStep";
import type { KinaCompiler } from "../KinaCompiler";
import { SemanticAnalysisStep } from "./SemanticAnalysisStep";
import { BuildIRStep } from "./BuildIRStep";
import { OptimizeIRStep } from "./OptimizeIRStep";
import { PrepareFSStep } from "./PrepareFSStep";
import { CompileStep } from "./CompileStep";

export const getSteps = (compiler: KinaCompiler) => ({
  PrepareFS: new PrepareFSStep(compiler),
  Readfile: new ReadfileStep(compiler),
  Tokenize: new TokenizeStep(compiler),
  Buildast: new BuildastStep(compiler),
  SemanticAnalysis: new SemanticAnalysisStep(compiler),
  BuildIR: new BuildIRStep(compiler),
  OptimizeIR: new OptimizeIRStep(compiler),
  Compile: new CompileStep(compiler),
});
