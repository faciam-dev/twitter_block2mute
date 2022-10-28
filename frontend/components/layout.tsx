import { useState } from "react";

const Layout = ({ children }) => {
  const [csrfToken, setCsrfToken] = useState("");

  return (
    <>
      <main>{children}</main>
    </>
  );
};

export default Layout;
