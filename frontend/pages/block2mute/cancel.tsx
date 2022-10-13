import { FormEvent, useRef, useState } from "react";

import { NextPage } from "next";

import axios from "axios";

const CancelPage: NextPage = () => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL_BASE;
  const [message, setMessage] = useState("処理をキャンセルしました。");
  const isLoggedout = useRef(false);

  const onClickLogout = async (event: FormEvent) => {
    event.preventDefault();
    const { data } = await axios.get(`${apiUrl}/auth/logout`);
    if (data.result == 1) {
      setMessage("ログアウトしました。ウィンドウを閉じてください。");
      isLoggedout.current = true;
    }
  };

  const onClickClose = async (event: FormEvent) => {
    window.close();
  };

  const logoutButton = (isLoggedout: boolean) => {
    if (isLoggedout) {
      return "";
    } else {
      return (
        <p>
          <button type="button" onClick={onClickLogout}>
            {" "}
            Logout
          </button>
        </p>
      );
    }
  };

  return (
    <div>
      <h1>ブロックミュート変換</h1>
      {message}
      {logoutButton(isLoggedout.current)}
    </div>
  );
};

export default CancelPage;
