export type EntryFormState = {
  title: string;
  username: string;
  password: string;
  url: string;
};

export const emptyEntryForm = (): EntryFormState => ({
  title: "",
  username: "",
  password: "",
  url: "",
});
