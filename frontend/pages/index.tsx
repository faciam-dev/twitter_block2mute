import type { NextPage } from "next";

import axios from "axios";
import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/router";
import { assertIsDefined } from "../helpers/assert";

const IndexPage: NextPage = () => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL_BASE;
  const router = useRouter();
  const didEffect = useRef(false);
  const [message, setMessage] = useState("is auth....");

  useEffect(() => {
    if (!didEffect.current) {
      didEffect.current = true;
      const getIsAuth = async () => {
        const config = {
          withCredentials: true,
        };
        const { data } = await axios.get(`${apiUrl}/auth/is_auth`, config);

        if (data.result != 1) {
          router.push("/login");
        } else {
          setMessage("logged in");
          router.push("/block2mute/");
        }
      };
      getIsAuth();
    }
  });

  return <div>{message}</div>;
};

export default IndexPage;
