"use client";

import React, { useState } from "react";
import Sidebar from "@/components/Header/Sidebar";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const [mainMargin, setMainMargin] = useState(64);

  return (
    <div className="flex min-h-screen bg-[#002440]">
      <Sidebar setMainMargin={setMainMargin} />
      <main
        className="flex-1 transition-all duration-300 overflow-y-auto"
        style={{ marginLeft: `${mainMargin}px` }}
      >
        {children}
      </main>
    </div>
  );
};

export default Layout;
