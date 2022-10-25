import { FormEvent } from "react";

type Props = {
  clickLogout;
  isLoggedout: boolean;
};

const LogoutButton = ({ clickLogout, isLoggedout }: Props) => {
  const onClickLogout = async (event: FormEvent) => {
    clickLogout(event);
  };

  if (isLoggedout) {
    return <p></p>;
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

export default LogoutButton;
