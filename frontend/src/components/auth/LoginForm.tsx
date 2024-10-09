import React, { useState } from "react";

interface LoginFormProps {
  onLogin: (username: string, password: string) => Promise<void>;
}

export function LoginForm({ onLogin }: LoginFormProps) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await onLogin(username, password);
    } catch (error) {
      console.error("Login error:", error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        placeholder="Username"
        className="w-full px-2 py-2 mb-2 text-white border-2 rounded-md bg-deep-blue"
        required
      />
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="Password"
        className="w-full px-2 py-2 mb-2 text-white border-2 rounded-md bg-deep-blue"
        required
      />
      <button
        type="submit"
        className="w-full px-2 py-2 mb-2 text-white border-2 rounded-md bg-deep-blue"
      >
        Log In
      </button>
    </form>
  );
}
