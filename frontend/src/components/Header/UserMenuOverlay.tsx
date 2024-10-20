import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/providers/AuthContext";
import { useUserProfile } from "@/hooks/useUserProfile";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import {
  LogIn,
  LogOut,
  UserPlus,
  BadgeCheck,
  CreditCard,
  Bell,
  Sparkles,
  ChevronsUpDown,
} from "lucide-react";

interface UserMenuOverlayProps {
  isExpanded: boolean;
}

const UserMenuOverlay: React.FC<UserMenuOverlayProps> = ({ isExpanded }) => {
  const router = useRouter();
  const { isAuthenticated, logout } = useAuth();
  const { profile } = useUserProfile();
  const [isOpen, setIsOpen] = useState(false);

  return (
    <Popover open={isOpen} onOpenChange={setIsOpen}>
      <PopoverTrigger asChild>
        <button className="absolute bottom-4 left-0 w-full px-4 py-2 flex items-center ">
          <Avatar className="h-8 w-8 rounded-lg">
            <AvatarImage src="/path-to-user-avatar.jpg" alt="User" />
            <AvatarFallback>U</AvatarFallback>
          </Avatar>
          {isExpanded && (
            <div className="ml-2 text-left flex items-center justify-between w-full">
              <div className="flex flex-col">
                <div className="text-sm font-semibold">
                  {isAuthenticated ? profile?.username : "Guest"}
                </div>
                <div className="text-xs">
                  {isAuthenticated ? profile?.email : "Not logged in"}
                </div>
              </div>
              <ChevronsUpDown className="ml-auto size-4" />
            </div>
          )}
        </button>
      </PopoverTrigger>
      <PopoverContent className="w-56 " align="start" side="right">
        {isAuthenticated ? (
          <>
            <div className="flex flex-col space-y-1 p-2">
              <p className="text-sm font-medium">
                {isAuthenticated ? profile?.username : "Guest"}
              </p>
              <p className="text-xs text-muted-foreground">
                {isAuthenticated ? profile?.email : "Not logged in"}
              </p>
            </div>
            <div className="h-px bg-border my-2" />
            <button className="w-full text-left px-2 py-1.5 text-sm hover:bg-accent hover:text-accent-foreground">
              <Sparkles className="mr-2 h-4 w-4 inline" />
              Upgrade to Pro
            </button>
            <button className="w-full text-left px-2 py-1.5 text-sm hover:bg-accent hover:text-accent-foreground">
              <BadgeCheck className="mr-2 h-4 w-4 inline" />
              Account
            </button>
            <button className="w-full text-left px-2 py-1.5 text-sm hover:bg-accent hover:text-accent-foreground">
              <CreditCard className="mr-2 h-4 w-4 inline" />
              Billing
            </button>
            <button className="w-full text-left px-2 py-1.5 text-sm hover:bg-accent hover:text-accent-foreground">
              <Bell className="mr-2 h-4 w-4 inline" />
              Notifications
            </button>
            <div className="h-px bg-border my-2" />
            <button
              onClick={logout}
              className="w-full text-left px-2 py-1.5 text-sm hover:bg-accent hover:text-accent-foreground"
            >
              <LogOut className="mr-2 h-4 w-4 inline" />
              Log out
            </button>
          </>
        ) : (
          <>
            <button
              onClick={() => router.push("/login")}
              className="w-full text-left px-2 py-1.5 text-sm hover:bg-accent hover:text-accent-foreground"
            >
              <LogIn className="mr-2 h-4 w-4 inline" />
              Login
            </button>
            <button
              onClick={() => router.push("/signup")}
              className="w-full text-left px-2 py-1.5 text-sm hover:bg-accent hover:text-accent-foreground"
            >
              <UserPlus className="mr-2 h-4 w-4 inline" />
              Register
            </button>
          </>
        )}
      </PopoverContent>
    </Popover>
  );
};

export default UserMenuOverlay;
