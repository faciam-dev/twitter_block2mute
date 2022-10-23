import type { FormEvent } from "react";

import { NextPage } from "next";
import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/router";
import axios from "axios";

const AllPage: NextPage = () => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL_BASE;
  const router = useRouter();
  const didEffect = useRef(false);
  const [message, setMessage] = useState("");
  const [totalSuccess, setTotalSuccess] = useState(0);
  const isLoggedout = useRef(false);

  useEffect(() => {
    if (!didEffect.current) {
      didEffect.current = true;
      const getIsAuth = async () => {
        const config = {
          withCredentials: true,
        };
        const { data } = await axios.get(`${apiUrl}/block2mute/all`, config);

        if (data.num_success != undefined) {
          setTotalSuccess(data.num_success);
        }
      };
      getIsAuth();
    }
  });

  const onClickLogout = async (event: FormEvent) => {
    event.preventDefault();
    const { data } = await axios.get(`${apiUrl}/auth/logout`);
    if (data.result == 1) {
      setMessage("ログアウトしました。ウィンドウを閉じてください。");
      isLoggedout.current = true;
    }
  };

  return (
    <div>
      <h1>ブロックミュート変換</h1>
      <p>
        {totalSuccess}
        件ブロックからミュートに変換しました。<br></br>
        ※0件変換の場合はAPI制限がかかっている可能性があります。
        {message}
      </p>

      <form>
        <p>
          <button type="button" onClick={onClickLogout}>
            {" "}
            Logout
          </button>
        </p>
      </form>
    </div>
  );
};

export default AllPage;
