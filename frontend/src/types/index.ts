// src/types/index.ts
export type LinkType = {
  shortURL: string;
  longURL: string;
  clickCount: number;
  lastClicked: string | null;
  createdAt: string;
};

export type User = {
  id: number;
  email: string;
  createdAt: string;
};

export type AuthMode = 'login' | 'register';
