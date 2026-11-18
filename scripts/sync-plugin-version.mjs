// Plugin installs only update when plugin.json's version changes, and changesets only
// bumps package.json, so this copies the version across after `changeset version`.

import { readFileSync, writeFileSync } from "node:fs";

const { version } = JSON.parse(readFileSync("package.json", "utf8"));
const pluginPath = ".claude-plugin/plugin.json";
const plugin = JSON.parse(readFileSync(pluginPath, "utf8"));

plugin.version = version;
writeFileSync(pluginPath, JSON.stringify(plugin, null, 2) + "\n");
