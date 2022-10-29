import React from "react";

const Notice = () => {
  return (
    <ul>
      <li>
        1度の操作で変換できるのは50件までです。遅くとも15分後にこの制限は解除されるでしょう。(50
        /15min.)
      </li>
      <li>
        変換・合計が0件の場合はAPI利用制限がかかっている可能性があります。しばらく時間をおいてからご利用ください。
      </li>
    </ul>
  );
};

export default Notice;
