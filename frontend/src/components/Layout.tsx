// components/Layout.tsx

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
      <div className="relative flex h-screen w-full overflow-hidden bg-gradient-dark">
        <div className="fixed top-0 left-0 h-full z-50">
          <div
            className="absolute top-0 left-0 h-full overflow-hidden"
            style={{ width: isExpanded ? "240px" : "64px" }}
          >
            <div
              className="absolute top-0 left-0 h-full w-16"
              onMouseEnter={() => setIsExpanded(true)}
              onMouseLeave={() => setIsExpanded(false)}
            >
              <AppSidebar
                isExpanded={isExpanded}
                setIsExpanded={setIsExpanded}
                isFooterMenuOpen={isFooterMenuOpen}
                setIsFooterMenuOpen={setIsFooterMenuOpen}
              />
            </div>
          </div>
        </div>
        <main className="flex-1 overflow-auto pl-16">{children}</main>
      </div>
    </SidebarProvider>
  );
};

export default Layout;
