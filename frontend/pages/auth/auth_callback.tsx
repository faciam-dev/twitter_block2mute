import { useRouter } from "next/router";
import { NextPage } from "next";

const axios = require("axios").default;

import { useEffect, useRef } from "react";

const AuthCallbackPage: NextPage = () => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL_BASE;
  const router = useRouter();
  const didEffect = useRef(false);
  const { oauth_token, oauth_verifier } = router.query;

  const getCallback = async () => {
    if (oauth_token === undefined) return;
    let client = axios.create({
      withCredentials: true,
    });

    const { data, headers } = await client.get(
      `${apiUrl}/auth/auth_callback?oauth_token=${oauth_token}&oauth_verifier=${oauth_verifier}`
    );

    router.push({
      pathname: "/",
    });
  };

  useEffect(() => {
    if (!didEffect.current) {
      getCallback();
    }
  });
  return <div>Login...</div>;
};

export default AuthCallbackPage;
