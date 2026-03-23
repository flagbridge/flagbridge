import { FlagBridge, type FlagBridgeConfig } from "@flagbridge/sdk-node";

// OpenFeature Provider interface (simplified)
// Full implementation will depend on @openfeature/server-sdk
export class FlagBridgeProvider {
  readonly metadata = { name: "flagbridge" };
  private client: FlagBridge;

  constructor(config: FlagBridgeConfig) {
    this.client = new FlagBridge(config);
  }

  async resolveBooleanEvaluation(flagKey: string, defaultValue: boolean) {
    const value = await this.client.getBooleanValue(flagKey, defaultValue);
    return { value };
  }

  async resolveStringEvaluation(flagKey: string, defaultValue: string) {
    const value = await this.client.getStringValue(flagKey, defaultValue);
    return { value };
  }

  async resolveNumberEvaluation(flagKey: string, defaultValue: number) {
    const value = await this.client.getNumberValue(flagKey, defaultValue);
    return { value };
  }
}
