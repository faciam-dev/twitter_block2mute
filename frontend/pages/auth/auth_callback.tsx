import useSWR from "swr";
import React from "react";
import { useRouter } from "next/router";

const fetcher = (url: string) => fetch(url).then((res) => res.json());

export default function AuthCallbackPage() {
  const router = useRouter();
  const { oauth_token, oauth_verifier } = router.query;

  const { data, error } = useSWR(
    "http://localhost:8080/auth/auth_callback?oauth_token=" +
      oauth_token +
      "&oauth_verifier=" +
      oauth_verifier,
    fetcher
  );

  if (data) {
    window.close();
  }
}
