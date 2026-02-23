export type User = {
  id: string;
  email: string;
  username: string;
  avatarUrl: string;
  role: "user" | "admin";
  createdAt: string;
  updatedAt: string;
};

export type AuthResponse = {
  token: string;
  user: User;
};
