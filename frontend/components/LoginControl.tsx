import React, { useEffect, useState } from "react";
import useSWR from "swr";

type State = {
  isLoggedIn: boolean;
  isLoadingAuthApi: boolean;
};

const fetcher = (url: string) => fetch(url).then((res) => res.json());

export function LoginControl() {
  const [state, setState] = useState<State>({
    isLoggedIn: false,
    isLoadingAuthApi: false,
  });

  const { data, error } = useSWR(
    state.isLoadingAuthApi ? `http://localhost:8080/auth/auth` : null,
    fetcher
  );

  useEffect(() => {
    console.log(data);

    if (!error && data && state.isLoadingAuthApi) {
      window.open(data.to_url);
      state.isLoadingAuthApi = false;
    }
  });

  function handleLoginClick() {
    setState({ isLoggedIn: true, isLoadingAuthApi: true });
  }

  function handleLogoutClick() {
    setState({ isLoggedIn: false, isLoadingAuthApi: false });
  }

  const isLoggedIn: boolean = state.isLoggedIn;
  let button;

  if (state.isLoggedIn) {
    button = "";
  } else {
    button = <LoginButton onClick={handleLoginClick} />;
  }

  return (
    <div>
      {button}
      <Greeting isLoggedIn={isLoggedIn} />
    </div>
  );
}

function GetAuthUrl() {
  const { data, error } = useSWR(`http://localhost:8080/auth/auth`, fetcher);
  return {
    toUrl: data.to_url,
    isLoading: !error && !data,
    isError: error,
  };
}

function LoginButton(props) {
  return <button onClick={props.onClick}>Login</button>;
}

function LogoutButton(props) {
  return <button onClick={props.onClick}>Logout</button>;
}

function UserGreeting(props) {
  return <h1></h1>;
}

function GuestGreeting(props) {
  return <h1>ツイッターアカウントでログインしてください。</h1>;
}

function Greeting(props) {
  const isLoggedIn = props.isLoggedIn;
  if (isLoggedIn) {
    return <UserGreeting />;
  }
  return <GuestGreeting />;
}
