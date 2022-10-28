import { FormEvent, useEffect, useRef, useState } from "react";

import { NextPage } from "next";

import axios from "../lib/axios";
import { useRouter } from "next/router";
import useFetchSession from "../hooks/useFetchSession";

const LoginPage: NextPage = () => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL_BASE;
  const didEffect = useRef(false);
  const { csrfToken, isLoading } = useFetchSession(apiUrl);

  const onClick = async (event: FormEvent) => {
    event.preventDefault();

    const config = {
      headers: { "X-CSRF-Token": csrfToken },
    };
    const { data } = await axios.post(`${apiUrl}/auth/auth`, {}, config);
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
