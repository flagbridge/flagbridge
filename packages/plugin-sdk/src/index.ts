export interface PluginManifest {
  name: string;
  version: string;
  description: string;
  author: string;
  hooks: string[];
}

export interface PluginContext {
  projectId: string;
  config: Record<string, unknown>;
}

export interface FlagBridgePlugin {
  manifest: PluginManifest;
  onActivate?(ctx: PluginContext): Promise<void>;
  onDeactivate?(ctx: PluginContext): Promise<void>;
  onFlagEvaluated?(ctx: PluginContext, flagKey: string, value: unknown): Promise<void>;
}

export function definePlugin(plugin: FlagBridgePlugin): FlagBridgePlugin {
  return plugin;
}
