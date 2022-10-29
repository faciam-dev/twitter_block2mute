import type { FormEvent } from "react";

import { NextPage } from "next";
import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/router";
import axios from "axios";
import LogoutButton from "../../components/logoutButton";
import useFetchSession from "../../hooks/useFetchSession";
import Notice from "../../components/notice";

const AllPage: NextPage = () => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL_BASE;
  const { csrfToken, isLoading } = useFetchSession(apiUrl);

  const router = useRouter();
  const didEffect = useRef(false);
  const [message, setMessage] = useState("現在、変換中です。");
  const isLoggedout = useRef(false);

  useEffect(() => {
    if (isLoading) {
      return;
    }
    if (!didEffect.current) {
      didEffect.current = true;
      const getIsAuth = async () => {
        const config = {
          headers: { "X-CSRF-Token": csrfToken },
        };
        try {
          const { data } = await axios.post(
            `${apiUrl}/block2mute/all`,
            {},
            config
          );

          if (data.num_success != undefined) {
            setMessage(
              data.num_success + "件ブロックからミュートに変換しました。"
            );
          }
        } catch (error) {
          let errorMessage = "変換中にエラーが発生しました。";
          if (error.error != undefined && error.error.detail != undefined) {
            setMessage(errorMessage + error.error.detail);
          }
        }
      };
      getIsAuth();
    }
  });

  const onClickLogout = async (event: FormEvent) => {
    event.preventDefault();
    const config = {
      headers: { "X-CSRF-Token": csrfToken },
    };
    const { data } = await axios.post(`${apiUrl}/auth/logout`, {}, config);
    if (data.result == 1) {
      setMessage("ログアウトしました。ウィンドウを閉じてください。");
      isLoggedout.current = true;
    }
  };

  return (
    <div>
      <h1>ブロックミュート変換</h1>
      <p>
        {message}
        <br></br>
      </p>
      <Notice></Notice>
      <LogoutButton
        clickLogout={onClickLogout}
        isLoggedout={isLoggedout.current}
      ></LogoutButton>
    </div>
  );
};

export default AllPage;
