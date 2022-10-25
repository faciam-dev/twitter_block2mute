import { FormEvent, useRef, useState } from "react";

import LogoutButton from "../../components/LogoutButton";

import { NextPage } from "next";

import axios from "axios";

const CancelPage: NextPage = () => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL_BASE;
  const [message, setMessage] = useState("処理をキャンセルしました。");
  const isLoggedout = useRef(false);

  const onClickLogout = async (event: FormEvent) => {
    event.preventDefault();
    const config = {
      withCredentials: true,
    };

    try {
      const { data } = await axios.get(`${apiUrl}/auth/logout`, config);
      if (data.result == 1) {
        setMessage("ログアウトしました。ウィンドウを閉じてください。");
        isLoggedout.current = true;
      }
    } catch (error) {}
  };

  return (
    <div>
      <h1>ブロックミュート変換</h1>
      {message}
      <LogoutButton
        clickLogout={onClickLogout}
        isLoggedout={isLoggedout.current}
      ></LogoutButton>
    </div>
  );
};

export default CancelPage;
