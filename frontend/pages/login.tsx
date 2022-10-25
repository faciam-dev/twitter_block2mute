import type { FormEvent } from "react";

import { NextPage } from "next";

import axios from "axios";

const LoginPage: NextPage = () => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL_BASE;

  const onClick = async (event: FormEvent) => {
    event.preventDefault();
    const { data } = await axios.get(`${apiUrl}/auth/auth`);
    window.location.href = data.to_url;
  };

  return (
    <div>
      <h1>ブロックミュート変換</h1>

      <form>
        <button type="submit" onClick={onClick}>
          login
        </button>
      </form>
    </div>
  );
};

export default LoginPage;
