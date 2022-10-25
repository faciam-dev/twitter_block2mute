import type { FormEvent } from "react";

import { NextPage } from "next";
import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/router";
import axios from "axios";

const IndexPage: NextPage = () => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL_BASE;
  const router = useRouter();
  const didEffect = useRef(false);
  const [totalBlock, setTotalBlock] = useState(0);

  useEffect(() => {
    if (!didEffect.current) {
      didEffect.current = true;
      const getIsAuth = async () => {
        const config = {
          withCredentials: true,
        };

        try {
          const { data } = await axios.get(`${apiUrl}/blocks/ids`, config);

          if (data.total != undefined) {
            setTotalBlock(data.total);
          }
        } catch (error) {
          router.push("/login");
        }
      };
      getIsAuth();
    }
  });

  const onClickOk = async (event: FormEvent) => {
    event.preventDefault();
    router.push("/block2mute/all");
  };

  const onClickCancel = async (event: FormEvent) => {
    event.preventDefault();
    router.push("/block2mute/cancel");
  };

  return (
    <div>
      <h1>ブロックミュート変換</h1>
      <p>
        現在あなたは{totalBlock}
        件ブロックしています。ブロックを全てミュートに変換しますか。 <br />
      </p>

      <form>
        <p>
          <button type="button" onClick={onClickOk}>
            OK
          </button>
        </p>
        <p>
          <button type="button" onClick={onClickCancel}>
            Cancel
          </button>
        </p>
      </form>
    </div>
  );
};

export default IndexPage;
