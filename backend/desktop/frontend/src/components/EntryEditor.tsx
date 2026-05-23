import { FormEvent } from "react";
import { EntryFormState } from "../types";

type Props = {
  title: string;
  form: EntryFormState;
  onChange: (form: EntryFormState) => void;
  onSubmit: () => void;
  onCancel: () => void;
  loading: boolean;
};

export function EntryEditor({
  title,
  form,
  onChange,
  onSubmit,
  onCancel,
  loading,
}: Props) {
  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    onSubmit();
  }

  return (
    <div className="modal-backdrop">
      <div className="modal panel">
        <h2>{title}</h2>
        <form onSubmit={handleSubmit} className="form-stack">
          <label>
            Title
            <input
              value={form.title}
              onChange={(e) => onChange({ ...form, title: e.target.value })}
              required
            />
          </label>
          <label>
            Username
            <input
              value={form.username}
              onChange={(e) => onChange({ ...form, username: e.target.value })}
              required
            />
          </label>
          <label>
            Password
            <input
              type="password"
              value={form.password}
              onChange={(e) => onChange({ ...form, password: e.target.value })}
              required
            />
          </label>
          <label>
            URL
            <input
              value={form.url}
              onChange={(e) => onChange({ ...form, url: e.target.value })}
              placeholder="optional"
            />
          </label>
          <div className="row-actions">
            <button type="button" className="btn" onClick={onCancel} disabled={loading}>
              Отмена
            </button>
            <button type="submit" className="btn primary" disabled={loading}>
              {loading ? "Сохранение…" : "Сохранить"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
