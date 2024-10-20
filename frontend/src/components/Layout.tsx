"use client";

import React from "react";
import AppSidebar from "@/components/Header/Sidebar";
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <SidebarProvider>
      <div className="flex h-screen w-full overflow-hidden bg-[#002440]">
        <AppSidebar />
        <div className="flex-1 overflow-auto">
          <main className="relative min-h-full">
            <SidebarTrigger className="absolute top-0 left-0 z-10 text-white" />
            {children}
          </main>
        </div>
      </div>
    </SidebarProvider>
  );
};

export default Layout;
