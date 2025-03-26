// components/Layout.tsx

"use client";

import React, { useState, useEffect } from "react";
import AppSidebar from "@/components/Header/Sidebar";
import MobileSidebar from "@/components/Header/Mobile/MobileSidebar";
import { SidebarProvider } from "@/components/ui/sidebar";
import { Button } from "@/components/ui/button";
import { Menu } from "lucide-react";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    const handleResize = () => {
      setIsMobile(window.innerWidth < 768);
    };

    handleResize();
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  return (
    <SidebarProvider>
      <div className="relative flex h-screen w-full overflow-hidden bg-gradient-dark">
        {/* Desktop Sidebar */}
        <div className="hidden md:block">
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
                />
              </div>
            </div>
          </div>
        </div>

        {/* Mobile Sidebar & Burger Button */}
        <div className="md:hidden">
          <Button
            variant="ghost"
            size="icon"
            className="fixed top-4 left-4 z-50"
            onClick={() => setIsMobileMenuOpen(true)}
          >
            <Menu className="h-6 w-6" />
          </Button>

          <MobileSidebar
            isOpen={isMobileMenuOpen}
            onClose={() => setIsMobileMenuOpen(false)}
          />
        </div>

        {/* Main Content */}
        <main
          className={`flex-1 overflow-auto ${!isMobile ? "pl-16" : "pl-0"}`}
        >
          {children}
        </main>
      </div>
    </SidebarProvider>
  );
};

export default Layout;
