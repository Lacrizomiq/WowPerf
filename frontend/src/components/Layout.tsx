"use client";

import React, { useState } from "react";
import AppSidebar from "@/components/Header/Sidebar";
import { SidebarProvider } from "@/components/ui/sidebar";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const [isFooterMenuOpen, setIsFooterMenuOpen] = useState(false);

  return (
    <SidebarProvider>
      <div className="flex h-screen w-full overflow-hidden bg-gradient-dark">
        <div
          onMouseEnter={() => setIsExpanded(true)}
          onMouseLeave={() => setIsExpanded(false)}
          className="flex-shrink-0 transition-all duration-300"
          style={{ width: isExpanded ? "240px" : "64px" }}
        >
          <AppSidebar
            isExpanded={isExpanded}
            setIsExpanded={setIsExpanded}
            isFooterMenuOpen={isFooterMenuOpen}
            setIsFooterMenuOpen={setIsFooterMenuOpen}
          />
        </div>
        <main className="flex-1 overflow-auto">{children}</main>
      </div>
    </SidebarProvider>
  );
};

export default Layout;
