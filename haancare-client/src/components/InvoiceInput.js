import React, { useState } from "react";

const InvoiceInput = ({ orderNum, onRegister }) => {
  const [inputValue, setInputValue] = useState("");

  const onclickRegister = () => {
    if (!inputValue.trim()) {
      alert("송장 번호를 입력하세요.");
      return;
    }
    onRegister(orderNum, inputValue);
  };

  return (
    <div className="invoice_input_wrapper">
      <input
        type="text"
        value={inputValue}
        onChange={(e) => setInputValue(e.target.value)}
        placeholder="송장 번호 입력"
      ></input>
      <button onClick={onclickRegister}>등록</button>
    </div>
  );
};

export default InvoiceInput;
