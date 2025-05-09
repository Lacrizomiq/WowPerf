// components/Layout.tsx
"use client";

import React, { useState, useEffect } from "react";
import AppSidebar from "@/components/Header/Sidebar";
import MobileSidebar from "@/components/Header/Mobile/MobileSidebar";
import { Button } from "@/components/ui/button";
import { Menu } from "lucide-react";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const [isDesktopSidebarExpanded, setIsDesktopSidebarExpanded] =
    useState(true);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [isMobileView, setIsMobileView] = useState(false);

  useEffect(() => {
    const handleResize = () => {
      const mobile = window.innerWidth < 768; // md breakpoint
      setIsMobileView(mobile);

      if (mobile) {
        setIsDesktopSidebarExpanded(false);
      }
    };

    handleResize(); // Execute on mount
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  const toggleDesktopSidebar = () => {
    if (!isMobileView) {
      setIsDesktopSidebarExpanded(!isDesktopSidebarExpanded);
    }
  };

  // Define sidebar widths for main content padding
  const sidebarExpandedWidth = "w-64"; // 16rem
  const sidebarCollapsedWidth = "w-16"; // 4rem

  const mainContentPaddingLeft = isMobileView
    ? "pl-0"
    : isDesktopSidebarExpanded
    ? "md:pl-64"
    : "md:pl-16";

  return (
    <div className="relative flex h-screen w-full overflow-hidden bg-[#1A1D21] text-slate-200">
      {/* Desktop Sidebar */}
      {!isMobileView && (
        <div
          className={`
            fixed left-0 top-0 z-30 h-full flex-shrink-0 
            transition-all duration-300 ease-in-out
            ${
              isDesktopSidebarExpanded
                ? sidebarExpandedWidth
                : sidebarCollapsedWidth
            }
          `}
        >
          <AppSidebar
            isExpanded={isDesktopSidebarExpanded}
            onToggleSidebar={toggleDesktopSidebar}
          />
        </div>
      )}

      {/* Mobile Burger Button & Sidebar */}
      {isMobileView && (
        <>
          <Button
            variant="ghost"
            size="icon"
            className="fixed top-4 left-4 z-50 text-white bg-slate-800/80 backdrop-blur-sm md:hidden"
            onClick={() => setIsMobileMenuOpen(true)}
            aria-label="Open mobile menu"
          >
            <Menu className="h-6 w-6" />
          </Button>

          <MobileSidebar
            isOpen={isMobileMenuOpen}
            onClose={() => setIsMobileMenuOpen(false)}
          />
        </>
      )}

      {/* Main Content */}
      <main
        className={`flex-1 overflow-y-auto transition-all duration-300 ease-in-out 
                  ${mainContentPaddingLeft}`}
      >
        {/* Add internal padding to content for spacing */}
        <div className="p-4 sm:p-6 lg:p-8">{children}</div>
      </main>
    </div>
  );
};

export default Layout;
