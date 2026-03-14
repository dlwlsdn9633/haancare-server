import React, { useState } from "react";
import { FaTruck } from "react-icons/fa";
import "./Sidebar.css";
import logo from "../assets/images/h_logo.jpg";

const Sidebar = () => {
  const [collapsed, setCollapsed] = useState(false);

  return (
    <aside className={`sidebar_container ${collapsed ? "collapsed" : ""}`}>
      <button
        className="sidebar_toggle"
        onClick={() => setCollapsed(!collapsed)}
      >
        ☰
      </button>

      <div className="sidebar_logo">
        <img src={logo} alt="Haancare Logo" />
      </div>

      <nav className="sidebar_nav">
        <ul>
          <li className="active">
            <a href="/">
              <FaTruck className="menu_icon" />
              <span className="menu_text">배송 목록</span>
            </a>
          </li>
        </ul>
      </nav>
    </aside>
  );
};

export default Sidebar;
