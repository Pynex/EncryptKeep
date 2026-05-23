import { FormEvent, useEffect, useState } from "react";
import {
  HasStoredKeys,
  InitializeNewKeys,
  IsUnlocked,
  ResetStoredKeys,
  Unlock,
} from "../api";

type Props = {
  onUnlocked: () => void;
};

export function UnlockView({ onUnlocked }: Props) {
  const [hasKeys, setHasKeys] = useState(true);
  const [masterPassword, setMasterPassword] = useState("");
  const [privateKeyHex, setPrivateKeyHex] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    (async () => {
      try {
        if (await IsUnlocked()) {
          onUnlocked();
          return;
        }
        setHasKeys(await HasStoredKeys());
      } catch (e) {
        setError(String(e));
      }
    })();
  }, [onUnlocked]);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      if (hasKeys) {
        await Unlock(masterPassword);
      } else {
        await InitializeNewKeys(privateKeyHex.trim(), masterPassword);
      }
      onUnlocked();
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setLoading(false);
    }
  }

  async function resetStoredKeys() {
    if (
      !confirm(
        "Удалить локально сохраненный аккаунт и ввести новый приватный ключ?"
      )
    ) {
      return;
    }
    setError("");
    setLoading(true);
    try {
      await ResetStoredKeys();
      setHasKeys(false);
      setMasterPassword("");
      setPrivateKeyHex("");
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="panel unlock-panel">
      <h1>EncryptKeep</h1>
      <p className="muted">
        {hasKeys
          ? "Введите мастер-пароль для разблокировки."
          : "Первый запуск: укажите приватный ключ (64 hex) и мастер-пароль."}
      </p>
      <form onSubmit={submit} className="form-stack">
        {!hasKeys && (
          <label>
            Private key (hex)
            <input
              type="password"
              value={privateKeyHex}
              onChange={(e) => setPrivateKeyHex(e.target.value)}
              autoComplete="off"
              placeholder="64 hex characters"
            />
          </label>
        )}
        <label>
          Master password
          <input
            type="password"
            value={masterPassword}
            onChange={(e) => setMasterPassword(e.target.value)}
            autoComplete="current-password"
            minLength={8}
          />
        </label>
        {error && <p className="error">{error}</p>}
        <button type="submit" className="btn primary" disabled={loading}>
          {loading ? "Подключение…" : hasKeys ? "Разблокировать" : "Создать и войти"}
        </button>
        {hasKeys && (
          <button type="button" className="btn" onClick={resetStoredKeys} disabled={loading}>
            Добавить новый аккаунт
          </button>
        )}
      </form>
    </div>
  );
}
