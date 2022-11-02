import type { NextPage } from "next";

import axios from "../lib/axios";
import useFetchSession from "../hooks/useFetchSession";
import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/router";

const IndexPage: NextPage = (props) => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL_BASE;
  const router = useRouter();
  const didEffect = useRef(false);
  const [message, setMessage] = useState("is auth....");

  const { csrfToken, isLoading } = useFetchSession(apiUrl);

  useEffect(() => {
    if (isLoading) {
      return;
    }
    if (!didEffect.current) {
      didEffect.current = true;
      const getIsAuth = async (csrfToken: string) => {
        const config = {
          headers: { "X-CSRF-Token": csrfToken },
        };

        try {
          const { data } = await axios.post(
            `${apiUrl}/auth/is_auth`,
            {},
            config
          );

          if (data.result != 1) {
            router.push("/login");
          } else {
            setMessage("logged in");
            router.push("/block2mute/");
          }
        } catch (error) {
          setMessage("エラーによりログイン状態判定が出来ませんでした。");
        }
      };
      getIsAuth(csrfToken);
    }
  });

  return <div>{message}</div>;
};

export default IndexPage;
