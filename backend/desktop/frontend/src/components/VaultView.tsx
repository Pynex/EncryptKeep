import { useCallback, useEffect, useState } from "react";
import {
  AddEntry,
  DeleteEntry,
  GetVaultStats,
  ListEntries,
  Lock,
  Sync,
  UpdateEntry,
} from "../api";
import { vault } from "../../wailsjs/go/models";
import { emptyEntryForm, EntryFormState } from "../types";
import { EntryEditor } from "./EntryEditor";

type Props = {
  onLocked: () => void;
};

function formatTime(value: unknown) {
  if (value == null || value === "") return "—";
  try {
    return new Date(value as string | number | Date).toLocaleString();
  } catch {
    return String(value);
  }
}

export function VaultView({ onLocked }: Props) {
  const [entries, setEntries] = useState<vault.PasswordEntry[]>([]);
  const [stats, setStats] = useState({ entryCount: 0, lastSync: "" });
  const [loading, setLoading] = useState(true);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const [editorOpen, setEditorOpen] = useState(false);
  const [editing, setEditing] = useState<vault.PasswordEntry | null>(null);
  const [form, setForm] = useState<EntryFormState>(emptyEntryForm());
  const [revealedId, setRevealedId] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    setError("");
    const [list, vaultStats] = await Promise.all([ListEntries(), GetVaultStats()]);
    setEntries(list ?? []);
    setStats(vaultStats);
  }, []);

  useEffect(() => {
    (async () => {
      setLoading(true);
      try {
        await refresh();
      } catch (e) {
        setError(String(e));
      } finally {
        setLoading(false);
      }
    })();
  }, [refresh]);

  async function handleSync() {
    setBusy(true);
    setError("");
    try {
      await Sync();
      await refresh();
    } catch (e) {
      setError(String(e));
    } finally {
      setBusy(false);
    }
  }

  function handleLock() {
    Lock();
    onLocked();
  }

  function openCreate() {
    setEditing(null);
    setForm(emptyEntryForm());
    setEditorOpen(true);
  }

  function openEdit(entry: vault.PasswordEntry) {
    setEditing(entry);
    setForm({
      title: entry.title,
      username: entry.username,
      password: entry.password,
      url: entry.url ?? "",
    });
    setEditorOpen(true);
  }

  async function saveEntry() {
    setBusy(true);
    setError("");
    try {
      if (editing) {
        await UpdateEntry(
          vault.PasswordEntry.createFrom({
            ...editing,
            title: form.title,
            username: form.username,
            password: form.password,
            url: form.url,
          })
        );
      } else {
        await AddEntry(form.title, form.username, form.password, form.url);
      }
      setEditorOpen(false);
      await refresh();
    } catch (e) {
      setError(String(e));
    } finally {
      setBusy(false);
    }
  }

  async function removeEntry(id: string) {
    if (!confirm("Удалить запись?")) return;
    setBusy(true);
    setError("");
    try {
      await DeleteEntry(id);
      await refresh();
    } catch (e) {
      setError(String(e));
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="vault-layout">
      <header className="vault-header">
        <div>
          <h1>Vault</h1>
          <p className="muted">
            Записей: {stats.entryCount} · Синхронизация: {formatTime(stats.lastSync)}
          </p>
        </div>
        <div className="row-actions">
          <button className="btn" onClick={handleSync} disabled={busy || loading}>
            Синхронизировать
          </button>
          <button className="btn primary" onClick={openCreate} disabled={busy}>
            Добавить
          </button>
          <button className="btn danger" onClick={handleLock} disabled={busy}>
            Заблокировать
          </button>
        </div>
      </header>

      {error && <p className="error banner">{error}</p>}
      {loading && <p className="muted">Загрузка…</p>}

      {!loading && entries.length === 0 && (
        <p className="muted empty">Нет записей. Нажмите «Добавить».</p>
      )}

      <div className="entry-list">
        {entries.map((entry) => (
          <article key={entry.id} className="entry-card panel">
            <div className="entry-main">
              <h3>{entry.title || "—"}</h3>
              <p className="muted">{entry.username}</p>
              {entry.url && (
                <a href={entry.url} target="_blank" rel="noreferrer" className="link">
                  {entry.url}
                </a>
              )}
              <p className="muted small">Обновлено: {formatTime(entry.updated_at)}</p>
            </div>
            <div className="entry-side">
              <code className="secret">
                {revealedId === entry.id ? entry.password : "••••••••"}
              </code>
              <div className="row-actions">
                <button
                  className="btn small"
                  onClick={() =>
                    setRevealedId(revealedId === entry.id ? null : entry.id)
                  }
                >
                  {revealedId === entry.id ? "Скрыть" : "Показать"}
                </button>
                <button className="btn small" onClick={() => openEdit(entry)}>
                  Изменить
                </button>
                <button
                  className="btn small danger"
                  onClick={() => removeEntry(entry.id)}
                  disabled={busy}
                >
                  Удалить
                </button>
              </div>
            </div>
          </article>
        ))}
      </div>

      {editorOpen && (
        <EntryEditor
          title={editing ? "Редактировать запись" : "Новая запись"}
          form={form}
          onChange={setForm}
          onSubmit={saveEntry}
          onCancel={() => setEditorOpen(false)}
          loading={busy}
        />
      )}
    </div>
  );
}
