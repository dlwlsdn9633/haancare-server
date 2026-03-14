import React, { useState, useEffect } from "react";
import InvoiceInput from "./InvoiceInput";
import "./DeliveryList.css";

const DeliveryList = () => {
  const [deliveries, setDeliveries] = useState([]);
  const [version, setVersion] = useState('Loading...');
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    Promise.all([
      fetch("http://localhost:8080/version"),
      fetch("http://localhost:8080/deliveries?orderStat=1"),
    ])
      .then(async ([versionRes, deliveriesRes]) => {
        if (!versionRes.ok) {
          throw new Error(`Failed to fetch version: ${versionRes.statusText}`);
        }
        if (!deliveriesRes.ok) {
          throw new Error(
            `Failed to fetch deliveries: ${deliveriesRes.statusText}`,
          );
        }

        const versionData = await versionRes.json();
        const deliveriesData = await deliveriesRes.json();

        return [versionData, deliveriesData];
      })
      .then(([versionData, deliveriesData]) => {
        setVersion(versionData.version || "Unknown");
        setDeliveries(deliveriesData.deliveries || []);
      })
      .catch((err) => {
        console.error("API error", err);
        setError(err.message);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, []);

  const handleInvoiceUpdate = (orderNumToUpdate, newInvoice) => {
    const updateDeliveries = deliveries.map((item) => {
      if (item.orderNum === orderNumToUpdate) {
        return { ...item, invoice: newInvoice };
      }
      return item;
    });
    setDeliveries(updateDeliveries);
  };

  const getStatusText = (state) => {
    switch (state) {
      case 1:
        return "주문접수";
      case 5:
        return "배송완료";
      default:
        return `상태 ${state}`;
    }
  };

  if (isLoading) {
    return (
      <div className="list_container">
        <h1>데이터를 불러오는 중...</h1>
      </div>
    );
  }

  if (error) {
    return (
      <div className="list_container">
        <h1 style={{ color: "#ff0000" }}>에러 발생: {error}</h1>
      </div>
    );
  }

  return (
    <div className="list_container">
      <h1>배송 목록 (Version: {version})</h1>

      <table className="delivery_table">
        <thead>
          <tr>
            <th>주문 번호</th>
            <th>이름</th>
            <th>주문 상태</th>
            <th>송장 번호</th>
            <th>주문 날짜</th>
          </tr>
        </thead>
        <tbody>
          {deliveries.length > 0 ? (
            deliveries.map(delivery => {
              return (
              <tr key={delivery.orderNum}>
                <td data-label="주문 번호">{delivery.orderNum}</td>
                <td data-label="이름">{delivery.name}</td>
                <td data-label="주문 상태">
                  <span className={`status_badge status_${delivery.orderState}`}>
                    {getStatusText(delivery.orderState)}
                  </span>
                </td>
                <td data-label="송장 번호">
                  {(delivery.orderState === 1 && !delivery.invoice) ? (
                    <InvoiceInput
                      orderNum={delivery.orderNum}
                      onRegister={handleInvoiceUpdate}
                    />
                  ) : (
                    delivery.invoice
                  )}
                </td>
                <td data-label="주문 생성일">{new Date(delivery.createdAt).toLocaleString()}</td>
              </tr>
              );
            })
          ) : (
            <tr>
              <td colSpan="5" style={{ textAlign: "center" }}>
                표시할 데이터가 없습니다.
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
};

export default DeliveryList;
