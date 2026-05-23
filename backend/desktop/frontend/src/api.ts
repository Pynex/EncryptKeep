import {
  AddEntry as AddEntryRaw,
  DeleteEntry as DeleteEntryRaw,
  GetEntry as GetEntryRaw,
  GetVaultStats as GetVaultStatsRaw,
  HasStoredKeys as HasStoredKeysRaw,
  InitializeNewKeys as InitializeNewKeysRaw,
  IsUnlocked as IsUnlockedRaw,
  ListEntries as ListEntriesRaw,
  Lock as LockRaw,
  Sync as SyncRaw,
  Unlock as UnlockRaw,
  UpdateEntry as UpdateEntryRaw,
} from "../wailsjs/go/main/App";
import { main, vault } from "../wailsjs/go/models";

function ensureWailsBridge() {
  const bridge = (window as any)?.go?.main?.App;
  if (!bridge) {
    throw new Error(
      "Wails bridge is unavailable. Open the desktop window via `wails dev` instead of browser preview."
    );
  }
}

export async function HasStoredKeys(): Promise<boolean> {
  ensureWailsBridge();
  return HasStoredKeysRaw();
}

export async function IsUnlocked(): Promise<boolean> {
  ensureWailsBridge();
  return IsUnlockedRaw();
}

export async function Unlock(masterPassword: string): Promise<void> {
  ensureWailsBridge();
  return UnlockRaw(masterPassword);
}

export async function InitializeNewKeys(
  privateKeyHex: string,
  masterPassword: string
): Promise<void> {
  ensureWailsBridge();
  return InitializeNewKeysRaw(privateKeyHex, masterPassword);
}

export async function Lock(): Promise<void> {
  ensureWailsBridge();
  return LockRaw();
}

export async function ResetStoredKeys(): Promise<void> {
  ensureWailsBridge();
  const fn = (window as any)?.go?.main?.App?.ResetStoredKeys;
  if (typeof fn !== "function") {
    throw new Error(
      "ResetStoredKeys is unavailable. Restart `wails dev` to regenerate bindings."
    );
  }
  return fn();
}

export async function ListEntries(): Promise<Array<vault.PasswordEntry>> {
  ensureWailsBridge();
  return ListEntriesRaw();
}

export async function GetEntry(id: string): Promise<vault.PasswordEntry> {
  ensureWailsBridge();
  return GetEntryRaw(id);
}

export async function AddEntry(
  title: string,
  username: string,
  password: string,
  url: string
): Promise<void> {
  ensureWailsBridge();
  return AddEntryRaw(title, username, password, url);
}

export async function UpdateEntry(entry: vault.PasswordEntry): Promise<void> {
  ensureWailsBridge();
  return UpdateEntryRaw(entry);
}

export async function DeleteEntry(id: string): Promise<void> {
  ensureWailsBridge();
  return DeleteEntryRaw(id);
}

export async function Sync(): Promise<void> {
  ensureWailsBridge();
  return SyncRaw();
}

export async function GetVaultStats(): Promise<main.VaultStats> {
  ensureWailsBridge();
  return GetVaultStatsRaw();
}
