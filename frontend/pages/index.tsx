import useSWR from "swr";
import { useState } from "react";

const fetcher = (url: string) => fetch(url).then((res) => res.json());

export default function IndexPage() {
  const { data, error } = useSWR("http://localhost:8080/auth/is_auth", fetcher);

  if (!data) return <div>Auth....</div>;

  return <div>{data.result} </div>;
}
